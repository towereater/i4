package model

type InsertMetadataInput struct {
	Client    string `json:"client"`
	Machine   string `json:"machine"`
	Timestamp string `json:"ts"`
	Size      int32  `json:"size"`
	Extension string `json:"ext"`
}
