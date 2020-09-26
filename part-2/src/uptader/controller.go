package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
	go dao.UpsertMany(ps, "product", "linx-2")
	ok = true
	return
}
