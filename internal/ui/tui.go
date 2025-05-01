package ui

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Omotolani98/cloudxp/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

type viewState int

const (
	listView viewState = iota
	detailView
)

type savedMessage struct{}

type model struct {
	missions        []models.Mission
	missionCursor   int
	objectiveCursor int
	state           viewState
	showSaved       bool
}

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("219"))
	borderStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1, 2)
)

func statusSymbol(done bool) string {
	if done {
		return "✓"
	}
	return " "
}

func renderMissionList(m model) string {
	var out string
	out += titleStyle.Render("📘 Missions") + "\n\n"
	for i, mission := range m.missions {
		line := fmt.Sprintf("  [%s] %s (%d/%d XP)",
			statusSymbol(mission.Completed),
			mission.Title,
			mission.XP,
			mission.TotalXP,
		)
		if i == m.missionCursor {
			line = selectedStyle.Render("> " + line)
		}
		out += line + "\n"
	}
	return out
}

func renderObjectivesPanel(m model) string {
	if len(m.missions) == 0 || m.missionCursor >= len(m.missions) {
		return ""
	}
	mission := m.missions[m.missionCursor]

	// Count completed
	completed := 0
	for _, o := range mission.Objectives {
		if o.Completed {
			completed++
		}
	}

	title := titleStyle.Render(fmt.Sprintf("📘 Objectives for: %s (%d/%d Complete)",
		mission.Title, completed, len(mission.Objectives)))

	// Create a container for all objectives
	var objectiveLines []string

	// Add all objectives to the container
	for i, obj := range mission.Objectives {
		status := fmt.Sprintf("[%s]", statusSymbol(obj.Completed))
		text := fmt.Sprintf("  %s %s", status, obj.Title)

		if i == m.objectiveCursor && m.state == detailView {
			objectiveLines = append(objectiveLines, selectedStyle.Render("> "+text))
		} else if obj.Completed {
			objectiveLines = append(objectiveLines, lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(text))
		} else {
			objectiveLines = append(objectiveLines, text)
		}
	}

	// Combine title with objectives, ensuring they're all in one column
	return title + "\n\n" + lipgloss.JoinVertical(lipgloss.Left, objectiveLines...)
}

func renderStatusBar(m model) string {
	if m.showSaved {
		return "💾 Saved!"
	}
	totalXP, earnedXP, completed := 0, 0, 0
	for _, mission := range m.missions {
		totalXP += mission.TotalXP
		earnedXP += mission.XP
		if mission.Completed {
			completed++
		}
	}
	level := 0
	if earnedXP > 0 {
		level = earnedXP / 500
	}
	return fmt.Sprintf("XP: %d / %d   •   Level: %d   •   Completed: %d/%d Missions",
		earnedXP, totalXP, level, completed, len(m.missions))
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.state == listView && m.missionCursor > 0 {
				m.missionCursor--
			} else if m.state == detailView && m.objectiveCursor > 0 {
				m.objectiveCursor--
			}
		case "down":
			if m.state == listView && m.missionCursor < len(m.missions)-1 {
				m.missionCursor++
			} else if m.state == detailView {
				objs := m.missions[m.missionCursor].Objectives
				if m.objectiveCursor < len(objs)-1 {
					m.objectiveCursor++
				}
			}
		case "enter":
			if m.state == listView {
				m.state = detailView
				m.objectiveCursor = 0
			} else {
				m.state = listView
				m.objectiveCursor = 0
			}
		case "c":
			if m.state == detailView && m.missionCursor < len(m.missions) {
				mission := &m.missions[m.missionCursor]
				if m.objectiveCursor < len(mission.Objectives) {
					obj := &mission.Objectives[m.objectiveCursor]
					obj.Completed = !obj.Completed

					// Recalculate XP
					completedCount := 0
					for _, o := range mission.Objectives {
						if o.Completed {
							completedCount++
						}
					}
					if len(mission.Objectives) > 0 {
						mission.XP = (mission.TotalXP / len(mission.Objectives)) * completedCount
					} else {
						mission.XP = 0
					}
					mission.Completed = completedCount == len(mission.Objectives) && len(mission.Objectives) > 0
				}
				// Save updated missions to YAML file
				updatedFile := models.MissionsFile{MainQuest: "", Missions: m.missions}
				out, err := yaml.Marshal(&updatedFile)
				if err != nil {
					fmt.Println("[!] Failed to serialize YAML:", err)
					os.Exit(1)
				}

				err = os.WriteFile("quests/missions.yaml", out, 0644)
				if err != nil {
					fmt.Println("[!] Failed to write missions.yaml:", err)
					os.Exit(1)
				}

				m.showSaved = true
				return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
					return savedMessage{}
				})
			}
		}
	case savedMessage:
		m.showSaved = false
	}
	return m, nil
}

func (m model) View() string {
	missions := renderMissionList(m)
	details := renderObjectivesPanel(m)
	status := renderStatusBar(m)

	return lipgloss.JoinVertical(lipgloss.Left,
		borderStyle.Render(missions),
		borderStyle.Render(details),
		status,
	)
}
func RunTUI() {
	file, err := os.Open("quests/missions.yaml")
	if err != nil {
		fmt.Println("[!] Failed to open missions.yaml:", err)
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("[!] Failed to read missions.yaml:", err)
		os.Exit(1)
	}

	var missionsFile models.MissionsFile
	err = yaml.Unmarshal(data, &missionsFile)
	if err != nil {
		fmt.Println("[!] Failed to parse YAML:", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model{missions: missionsFile.Missions})
	if err := p.Start(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
