package model

type StringDataGauge struct {
	Timestamp string   `json:"timestamp"`
	Key       string   `json:"key"`
	Value     string   `json:"value"`
	Params    []string `json:"params,omitempty"`
}

type IntDataGauge struct {
	Timestamp string   `json:"timestamp"`
	Key       string   `json:"key"`
	Value     int32    `json:"value"`
	Params    []string `json:"params,omitempty"`
}

type FloatDataGauge struct {
	Timestamp string   `json:"timestamp"`
	Key       string   `json:"key"`
	Value     float32  `json:"value"`
	Params    []string `json:"params,omitempty"`
}
