package model

type Target struct {
	Name    string `json:"name"`
	Machine string `json:"machine"`
	User    string `json:"user"`
	Pass    string `json:"pass"`
	Folder  string `json:"folder"`
	File    string `json:"file"`
}
