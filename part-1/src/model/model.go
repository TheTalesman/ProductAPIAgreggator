package model

import "reflect"

//Product struct
type Product struct {
	_id  string `bson:"_id"`
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

//StructIdField Retorna o primeiro campo, por padrão ID
func (q *Product) StructIdField() string {
	return reflect.TypeOf(q).Elem().Field(0).Tag.Get("json")

}
