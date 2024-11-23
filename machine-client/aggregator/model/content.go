package model

type DataContent struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

type DataGauge struct {
	Key       string   `json:"key"`
	Value     any      `json:"value"`
	Timestamp string   `json:"ts"`
	Params    []string `json:"params,omitempty"`
}

type DataInterval struct {
	Key    string   `json:"key"`
	Start  string   `json:"start"`
	End    string   `json:"end"`
	Params []string `json:"params,omitempty"`
}
