package model

//Product struct
type Product struct {
	_id  string `bson:"_id"`
	ID   string `bson:"id"`
	Name string `bson:"name"`
}
