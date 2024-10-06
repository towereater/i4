package model

type Interval struct {
	Start  string   `json:"start"`
	End    string   `json:"end"`
	Key    string   `json:"key"`
	Params []string `json:"params"`
}
