package model

type InsertMetadataInput struct {
	Client    string `json:"client"`
	Machine   string `json:"machine"`
	Timestamp string `json:"ts"`
	Size      int64  `json:"size"`
	Extension string `json:"ext"`
	FileHash  uint32 `json:"fileHash"`
}
