package models

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
