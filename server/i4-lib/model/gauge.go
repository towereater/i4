package model

type DataGauge struct {
	Machine   string `json:"machine" bson:"machine"`
	Key       string `json:"key" bson:"key"`
	Value     any    `json:"value" bson:"value"`
	Timestamp string `json:"ts" bson:"ts"`
	Params    []any  `json:"params,omitempty" bson:"params,omitempty"`
}

type DataGaugeSum struct {
	Machine string  `json:"machine" bson:"machine"`
	Key     string  `json:"key" bson:"key"`
	Value   any     `json:"value" bson:"value"`
	Count   int64   `json:"count" bson:"count"`
	Sum     float64 `json:"sum" bson:"sum"`
}
