package model

type UploadContent struct {
	Hash    uint32 `json:"hash" bson:"hash"`
	Content []byte `json:"content"`
}

type DataContent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type DataInterval struct {
	Start  string   `json:"start"`
	End    string   `json:"end"`
	Key    string   `json:"key"`
	Params []string `json:"params,omitempty"`
}

type DataGauge struct {
	Timestamp string   `json:"timestamp"`
	Key       string   `json:"key"`
	Value     any      `json:"value"`
	Params    []string `json:"params,omitempty"`
}
