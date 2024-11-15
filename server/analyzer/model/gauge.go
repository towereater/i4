package model

type DataGauge struct {
	Timestamp string   `json:"timestamp"`
	Key       string   `json:"key"`
	Value     any      `json:"value"`
	Params    []string `json:"params,omitempty"`
}
