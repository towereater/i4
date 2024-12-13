package model

type Target struct {
	Id      string `json:"id"`
	Machine string `json:"machine"`
	NetIp   string `json:"netip"`
	User    string `json:"user"`
	Pass    string `json:"pass"`
	File    string `json:"file"`
}
