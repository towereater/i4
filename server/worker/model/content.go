package model

type FileContent struct {
	Hash    uint32 `json:"hash" bson:"hash"`
	Content []byte `json:"content"`
}
