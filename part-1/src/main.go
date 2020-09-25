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

func main() {
	dao.Connect()
	/* 	var portFind = ""
	   	var port int = 8079

	   	for portFind != "" {
	   		port++
	   		portFind, _ := dao.RClient.Get(context.Background(), strconv.Itoa(port)).Result()

	   		if portFind == "" {
	   			err := dao.RClient.Set(context.Background(), strconv.Itoa(port), "used", 0).Err()
	   			if err != nil {
	   				log.Println(err)
	   				return
	   			}
	   			break
	   		}
	   	} */

	//	b := ":" + strconv.Itoa(port)
	r := gin.Default()
	r.Use(Guard())
	r.GET("/products", FindProducts)
	r.POST("/products", UpsertProducts)
	//a := rand.Intn(8099-8000) + 8000
	//b := ":" + strconv.Itoa(a)
	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

//FindProducts find all products
func FindProducts(c *gin.Context) {

	ok, res := dao.FindAll("product")
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
	ok, err = dao.UpsertMany(ps, "product")
	return
}

//Upsert a p Product as parameter.
//Returns OK flags success.
func Upsert(p Product) bool {
	pMap := structs.Map(p)
	ok, _ := dao.Upsert(pMap, "product")
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
