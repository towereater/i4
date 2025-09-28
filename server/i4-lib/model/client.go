package model

type Client struct {
	Code     string    `json:"code" bson:"code"`
	Name     string    `json:"name" bson:"name"`
	ApiKey   string    `json:"apiKey" bson:"apiKey"`
	Machines []Machine `json:"machines,omitempty" bson:"machines"`
}

type InsertClientInput struct {
	Code     string    `json:"code" bson:"code"`
	Name     string    `json:"name" bson:"name"`
	ApiKey   string    `json:"apiKey"`
	Machines []Machine `json:"machines" bson:"machines"`
}
