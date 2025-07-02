package model

type DataContent struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

type DataGauge struct {
	Machine   string `json:"machine"`
	Key       string `json:"key"`
	Value     any    `json:"value"`
	Timestamp string `json:"ts"`
	Params    []any  `json:"params,omitempty"`
}

type DataInterval struct {
	Machine string `json:"machine"`
	Key     string `json:"key"`
	Value   any    `json:"value"`
	Start   string `json:"start"`
	End     string `json:"end"`
	Params  []any  `json:"params,omitempty"`
}
