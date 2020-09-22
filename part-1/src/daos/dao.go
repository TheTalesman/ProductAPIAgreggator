package dao

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Client is the DB global client
var Client *mongo.Client

//Connect to database
func Connect() (ok bool) {
	//Assert once an deny on any error.
	ok = true

	u, p, dbn, dbh, ok := getEnvVars()
	var dbString strings.Builder

	dbString.WriteString("mongodb+srv://")
	dbString.WriteString(u)
	dbString.WriteString(":")
	dbString.WriteString(p)
	dbString.WriteString("@")
	dbString.WriteString(dbh)
	dbString.WriteString("/")
	dbString.WriteString(dbn)
	dbString.WriteString("retryWrites=true&w=majority")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbString.String()))
	Client = client
	if err != nil {
		log.Fatal(err)
		ok = false

	}
	return
}

func stringBuilder() {

}
func getEnvVars() (u string, pass string, dbName string, dbHost string, ok bool) {
	ok = true
	err := godotenv.Load("daos/.env")
	if err != nil {
		//	log.Fatal("Error loading .env file")
		log.Fatal(err)
		ok = false

	}
	u = os.Getenv("DB_USER")
	pass = os.Getenv("DB_PASS")
	dbName = os.Getenv("DB_NAME")
	dbHost = os.Getenv("DB_HOST")
	return
}

//Upsert an map to db
func Upsert(entity map[string]interface{}) bool {
	return false

}
