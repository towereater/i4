package model

type UploadContent struct {
	Hash    string `json:"hash" bson:"hash"`
	Content []byte `json:"content"`
}
