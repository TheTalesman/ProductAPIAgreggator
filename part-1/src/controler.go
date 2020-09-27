package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"
	. "ProductAPIAgreggator/part-1/src/model"
	"log"
	"net/http"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

//FindProducts find all products
func FindProducts(c *gin.Context) {

	ok, res := dao.FindAll("product", "linx")
	if !ok {
		c.JSON(http.StatusNotFound, res)

	}
	c.JSON(http.StatusOK, res)
}

//UpsertProducts upserts many products
func UpsertProducts(c *gin.Context) {
	// Validate input
	var input []map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing body function/ bad formating"})
		return
	}
	log.Println("Json Binded")
	ok, err := UpsertMany(input)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	log.Println(input)

}

//Upsert a p Product as parameter.
//Returns OK flags success.
func UpsertMany(ps []map[string]interface{}) (ok bool, err error) {
	//TODO put dbname in env vars
	idField := (&Product{}).StructIdField()

	go dao.UpsertMany(ps, "product", idField, "linx")
	ok = true
	return
}

//Upsert a p Product as parameter.
//Returns OK flags success.
func Upsert(p Product) bool {
	pMap := structs.Map(p)
	ok, _ := dao.Upsert(pMap, "product", "linx")
	return ok
}

//Find Receive a p Product as parameter
//Returns the relative product in database by ID. OK flags success.
func Find(p Product) (ok bool, pResponse Product) {

	ok, pInterface := dao.FindByID(p.ID, "product")
	err := mapstructure.Decode(pInterface, &pResponse)
	if err != nil {
		log.Println(" Erro no Decode()")
		log.Fatal(err)
		ok = false
	}
	return
}
