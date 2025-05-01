package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Omotolani98/cloudxp/cmd"
	"github.com/Omotolani98/cloudxp/internal/ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cloudxp <command> [arguments]")
		fmt.Println("Commands: add <mission_id> <xp>, log, status, complete <mission_id> <objective_index>, show <mission_id>, reset")
		return
	}

	switch os.Args[1] {
	case "ui":
		ui.RunTUI()
	case "add":
		if len(os.Args) != 4 {
			fmt.Println("Usage: cloudxp add <mission_id> <xp>")
			return
		}
		missionID := os.Args[2]
		xp, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Invalid XP value.")
			return
		}
		cmd.AddXP(missionID, xp)
	case "log":
		cmd.LogMissions()
	case "complete":
		if len(os.Args) != 4 {
			fmt.Println("Usage: cloudxp complete <mission_id> <objective_index>")
			return
		}
		missionID := os.Args[2]
		index, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Invalid objective index.")
			return
		}
		cmd.CompleteObjective(missionID, index)
	case "status":
		cmd.Status()
	case "show":
		if len(os.Args) != 3 {
			fmt.Println("Usage: cloudxp show <mission_id>")
			return
		}
		missionID := os.Args[2]
		cmd.ShowMissionDetails(missionID)
	case "reset":
		fmt.Print("Are you sure you want to reset all progress? (y/N): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm == "y" || confirm == "Y" {
			cmd.ResetProgress()
		} else {
			fmt.Println("Reset canceled.")
		}
	default:
		fmt.Println("Unknown command. Available: add, log, status, complete, show, reset")
	}
}
