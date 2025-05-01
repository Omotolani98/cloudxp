package cmd

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Objective struct {
	Title     string `yaml:"title"`
	Completed bool   `yaml:"completed"`
}

type Mission struct {
	ID         string      `yaml:"id"`
	Title      string      `yaml:"title"`
	TotalXP    int         `yaml:"total_xp"`
	XP         int         `yaml:"xp"`
	Completed  bool        `yaml:"completed"`
	Objectives []Objective `yaml:"objectives"`
}

type MissionsFile struct {
	MainQuest string    `yaml:"main_quest"`
	Missions  []Mission `yaml:"missions"`
}

func AddXP(missionID string, xp int) {
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

	var missionsFile MissionsFile
	err = yaml.Unmarshal(data, &missionsFile)
	if err != nil {
		fmt.Println("[!] Failed to parse YAML:", err)
		os.Exit(1)
	}

	updated := false
	for i, m := range missionsFile.Missions {
		if m.ID == missionID {
			missionsFile.Missions[i].XP += xp
			if missionsFile.Missions[i].XP >= missionsFile.Missions[i].TotalXP {
				missionsFile.Missions[i].XP = missionsFile.Missions[i].TotalXP
				missionsFile.Missions[i].Completed = true
				fmt.Printf("[✓] %s completed! Total XP: %d\n", m.Title, m.TotalXP)
			} else {
				fmt.Printf("[+] %d XP added to %s. Total: %d/%d XP\n",
					xp, m.Title, missionsFile.Missions[i].XP, missionsFile.Missions[i].TotalXP)
			}
			updated = true
			break
		}
	}

	if !updated {
		fmt.Println("[!] Mission not found.")
		os.Exit(1)
	}

	newData, err := yaml.Marshal(&missionsFile)
	if err != nil {
		fmt.Println("[!] Failed to serialize YAML:", err)
		os.Exit(1)
	}

	err = os.WriteFile("quests/missions.yaml", newData, 0644)
	if err != nil {
		fmt.Println("[!] Failed to write missions.yaml:", err)
		os.Exit(1)
	}
}

func LogMissions() {
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

	var missionsFile MissionsFile
	err = yaml.Unmarshal(data, &missionsFile)
	if err != nil {
		fmt.Println("[!] Failed to parse YAML:", err)
		os.Exit(1)
	}

	fmt.Printf("\nMain Quest: %s\n\n", missionsFile.MainQuest)
	for _, m := range missionsFile.Missions {
		status := "[ ]"
		if m.Completed {
			status = "[✓]"
		}
		fmt.Printf("%s [%s] %s - %d/%d XP\n", status, m.ID, m.Title, m.XP, m.TotalXP)
	}
}

func Status() {
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

	var missionsFile MissionsFile
	err = yaml.Unmarshal(data, &missionsFile)
	if err != nil {
		fmt.Println("[!] Failed to parse YAML:", err)
		os.Exit(1)
	}

	totalXP := 0
	earnedXP := 0
	completed := 0

	for _, m := range missionsFile.Missions {
		totalXP += m.TotalXP
		earnedXP += m.XP
		if m.Completed {
			completed++
		}
	}

	level := earnedXP / 500
	fmt.Printf("\nStatus Report:\n")
	fmt.Printf("Level: %d\n", level)
	fmt.Printf("XP: %d / %d\n", earnedXP, totalXP)
	fmt.Printf("Missions Completed: %d / %d\n\n", completed, len(missionsFile.Missions))
}

func CompleteObjective(missionID string, objectiveIndex int) {
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

	var missionsFile MissionsFile
	err = yaml.Unmarshal(data, &missionsFile)
	if err != nil {
		fmt.Println("[!] Failed to parse YAML:", err)
		os.Exit(1)
	}

	updated := false
	for i, m := range missionsFile.Missions {
		if m.ID == missionID {
			if objectiveIndex < 0 || objectiveIndex >= len(m.Objectives) {
				fmt.Println("[!] Invalid objective index")
				os.Exit(1)
			}
			if m.Objectives[objectiveIndex].Completed {
				fmt.Println("[-] Objective already completed")
				os.Exit(0)
			}

			missionsFile.Missions[i].Objectives[objectiveIndex].Completed = true

			xpPerObjective := m.TotalXP / len(m.Objectives)
			missionsFile.Missions[i].XP += xpPerObjective
			fmt.Printf("[+] Objective '%s' completed. %d XP awarded.\n",
				m.Objectives[objectiveIndex].Title, xpPerObjective)

			allDone := true
			for _, obj := range missionsFile.Missions[i].Objectives {
				if !obj.Completed {
					allDone = false
					break
				}
			}
			if allDone {
				missionsFile.Missions[i].Completed = true
				fmt.Printf("[✓] Mission '%s' completed!\n", m.Title)
			}

			updated = true
			break
		}
	}

	if !updated {
		fmt.Println("[!] Mission not found.")
		os.Exit(1)
	}

	newData, err := yaml.Marshal(&missionsFile)
	if err != nil {
		fmt.Println("[!] Failed to serialize YAML:", err)
		os.Exit(1)
	}

	err = os.WriteFile("quests/missions.yaml", newData, 0644)
	if err != nil {
		fmt.Println("[!] Failed to write missions.yaml:", err)
		os.Exit(1)
	}
}

func ShowMissionDetails(missionID string) {
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

	var missionsFile MissionsFile
	err = yaml.Unmarshal(data, &missionsFile)
	if err != nil {
		fmt.Println("[!] Failed to parse YAML:", err)
		os.Exit(1)
	}

	for _, m := range missionsFile.Missions {
		if m.ID == missionID {
			fmt.Printf("\nMission: %s\n", m.Title)
			fmt.Printf("ID: %s\n", m.ID)
			fmt.Printf("XP: %d / %d\n", m.XP, m.TotalXP)
			fmt.Printf("Status: %v\n", map[bool]string{true: "Completed", false: "In Progress"}[m.Completed])
			fmt.Println("Objectives:")
			for i, obj := range m.Objectives {
				status := "[ ]"
				if obj.Completed {
					status = "[✓]"
				}
				fmt.Printf("  %s (%d) %s\n", status, i, obj.Title)
			}
			return
		}
	}

	fmt.Println("[!] Mission not found.")
}

func ResetProgress() {
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

	var missionsFile MissionsFile
	err = yaml.Unmarshal(data, &missionsFile)
	if err != nil {
		fmt.Println("[!] Failed to parse YAML:", err)
		os.Exit(1)
	}

	for i := range missionsFile.Missions {
		missionsFile.Missions[i].XP = 0
		missionsFile.Missions[i].Completed = false
		for j := range missionsFile.Missions[i].Objectives {
			missionsFile.Missions[i].Objectives[j].Completed = false
		}
	}

	newData, err := yaml.Marshal(&missionsFile)
	if err != nil {
		fmt.Println("[!] Failed to serialize YAML:", err)
		os.Exit(1)
	}

	err = os.WriteFile("quests/missions.yaml", newData, 0644)
	if err != nil {
		fmt.Println("[!] Failed to write missions.yaml:", err)
		os.Exit(1)
	}

	fmt.Println("[*] All missions and objectives have been reset.")
}
