package model

//Product struct
type Product struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}
