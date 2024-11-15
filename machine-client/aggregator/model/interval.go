package model

type DataInterval struct {
	Start  string   `json:"start"`
	End    string   `json:"end"`
	Key    string   `json:"key"`
	Params []string `json:"params,omitempty"`
}
