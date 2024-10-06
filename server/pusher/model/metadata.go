package model

type FileMetadata struct {
	Client    string `json:"client" bson:"client"`
	Machine   string `json:"machine" bson:"machine"`
	Timestamp string `json:"timestamp" bson:"ts"`
	Size      int32  `json:"size" bson:"size"`
	Extension string `json:"extension" bson:"ext"`
	Hash      string `json:"hash" bson:"hash"`
}
