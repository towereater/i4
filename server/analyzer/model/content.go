package model

type UploadContent struct {
	Hash    uint32 `json:"hash" bson:"hash"`
	Content []byte `json:"content"`
}

type DataContent struct {
	Type    string `json:"type"`
	Content any    `json:"content"`
}

type DataGauge struct {
	Machine   string   `json:"machine,omitempty" bson:"machine"`
	Key       string   `json:"key" bson:"key"`
	Value     any      `json:"value" bson:"value"`
	Timestamp string   `json:"timestamp" bson:"timestamp"`
	Params    []string `json:"params,omitempty" bson:"params,omitempty"`
}

type DataInterval struct {
	Machine string   `json:"machine,omitempty" bson:"machine"`
	Key     string   `json:"key" bson:"key"`
	Start   string   `json:"start" bson:"start"`
	End     string   `json:"end" bson:"end"`
	Params  []string `json:"params,omitempty" bson:"params,omitempty"`
}
