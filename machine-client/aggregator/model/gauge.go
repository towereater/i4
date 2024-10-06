package model

type StringGauge struct {
	Timestamp string   `json:"timestamp"`
	Key       string   `json:"key"`
	Value     string   `json:"value"`
	Params    []string `json:"params"`
}

type FloatGauge struct {
	Timestamp string   `json:"timestamp"`
	Key       string   `json:"key"`
	Value     float32  `json:"value"`
	Params    []string `json:"params"`
}
