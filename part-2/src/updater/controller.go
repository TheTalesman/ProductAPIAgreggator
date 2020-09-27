package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

//ProductInput objeto transformados para ser enviado ao updater
type ProductInput struct {
	ProductID string   `json:"productId"`
	Images    []string `json:"images"`
}

//StructIdField Retorna o primeiro campo, por padrão ID
func (q *ProductInput) StructIdField() string {
	return reflect.TypeOf(q).Elem().Field(0).Tag.Get("json")

}

//UpsertProducts upserts many products
func UpsertProducts(c *gin.Context) {
	var input []map[string]interface{}
	var body ProductInput
	mapstructure.Decode(c.Request.Body, &body)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err, "body": body})
		return
	}
	//log.Println(input)
	ok, err := UpsertMany(input)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

}

//Upsert a p Product as parameter.
//Returns OK flags success.
func UpsertMany(ps []map[string]interface{}) (ok bool, err error) {
	idField := (&ProductInput{}).StructIdField()
	go dao.UpsertMany(ps, "product", idField, "linx-2")
	ok = true
	return
}

//FindProducts find all products
func FindProducts(c *gin.Context) {

	ok, res := dao.FindAll("product", "linx-2")
	if !ok {
		c.JSON(http.StatusNotFound, res)

	}
	c.JSON(http.StatusOK, res)
}
