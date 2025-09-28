package model

type Machine struct {
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}

type InsertMachineInput struct {
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name"`
}
