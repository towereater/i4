package model

type DataContent struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

type DataGauge struct {
	Machine   string `json:"machine" bson:"machine"`
	Key       string `json:"key" bson:"key"`
	Value     any    `json:"value" bson:"value"`
	Timestamp string `json:"ts" bson:"ts"`
	Params    []any  `json:"params,omitempty" bson:"params,omitempty"`
}

type DataGaugeSum struct {
	Id struct {
		Machine string `json:"machine" bson:"machine"`
		Key     string `json:"key" bson:"key"`
		Value   any    `json:"value" bson:"value"`
	} `json:"id" bson:"_id"`
	Count int64 `json:"count" bson:"count"`
}

type DataInterval struct {
	Machine string `json:"machine" bson:"machine"`
	Key     string `json:"key" bson:"key"`
	Value   any    `json:"value" bson:"value"`
	Start   string `json:"start" bson:"start"`
	End     string `json:"end" bson:"end"`
	Params  []any  `json:"params,omitempty" bson:"params,omitempty"`
}

type DataIntervalSum struct {
	Id struct {
		Machine string `json:"machine" bson:"machine"`
		Key     string `json:"key" bson:"key"`
		Value   any    `json:"value" bson:"value"`
	} `json:"id" bson:"_id"`
	Count int64 `json:"count" bson:"count"`
}
