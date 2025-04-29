package model

type UploadMetadata struct {
	Client    string `json:"client" bson:"client"`
	Machine   string `json:"machine" bson:"machine"`
	Timestamp string `json:"timestamp" bson:"ts"`
	Size      int64  `json:"size" bson:"size"`
	Extension string `json:"extension" bson:"ext"`
	Hash      uint32 `json:"hash" bson:"hash"`
}

type InsertMetadataInput struct {
	Client    string `json:"client"`
	Machine   string `json:"machine"`
	Timestamp string `json:"ts"`
	Size      int64  `json:"size"`
	Extension string `json:"ext"`
	FileHash  uint32 `json:"fileHash"`
}

type InsertMetadataOutput struct {
	Id   uint32                   `json:"id"`
	Urls InsertMetadataOutputUrls `json:"urls"`
}

type InsertMetadataOutputUrls struct {
	UploadContent string `json:"uploadContent"`
}
