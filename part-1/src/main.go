package main

import (
	dao "ProductAPIAgreggator/part-1/src/daos"
	. "ProductAPIAgreggator/part-1/src/model"
	"log"
	"net/http"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	_ "gopkg.in/mgo.v2/bson"
)

func main() {
	dao.Connect()
	r := gin.Default()
	r.GET("/products", FindProducts)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

//FindProducts find all products
func FindProducts(c *gin.Context) {

	ok, res := dao.FindAll("product")
	if !ok {
		c.JSON(http.StatusNotFound, res)

	}
	c.JSON(http.StatusOK, res)
}

//Upsert a p Product as parameter.
//Returns OK flags success.
func Upsert(p Product) bool {
	pMap := structs.Map(p)
	ok, _ := dao.Upsert(pMap)
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
