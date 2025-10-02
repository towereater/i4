package model

type DataInterval struct {
	Machine string `json:"machine" bson:"machine"`
	Key     string `json:"key" bson:"key"`
	Value   any    `json:"value" bson:"value"`
	Start   string `json:"start" bson:"start"`
	End     string `json:"end" bson:"end"`
	Params  []any  `json:"params,omitempty" bson:"params,omitempty"`
}

type DataIntervalSum struct {
	Machine string  `json:"machine" bson:"machine"`
	Key     string  `json:"key" bson:"key"`
	Value   any     `json:"value" bson:"value"`
	Count   int64   `json:"count" bson:"count"`
	Sum     float64 `json:"sum" bson:"sum"`
}
