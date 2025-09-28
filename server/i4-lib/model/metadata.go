package model

type UploadMetadata struct {
	Timestamp string `json:"timestamp" bson:"ts"`
	Size      uint32 `json:"size" bson:"size"`
	Extension string `json:"extension" bson:"ext"`
	Hash      string `json:"hash" bson:"hash"`
}

type InsertMetadataInput struct {
	Timestamp string `json:"ts"`
	Size      uint32 `json:"size"`
	Extension string `json:"ext"`
	Hash      string `json:"hash"`
}
