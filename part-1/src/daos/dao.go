package dao

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	bson "go.mongodb.org/mongo-driver/bson"
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

	//for the sake of memory
	dbString.WriteString("mongodb+srv://")
	dbString.WriteString(u)
	dbString.WriteString(":")
	dbString.WriteString(p)
	dbString.WriteString("@")
	dbString.WriteString(dbh)
	dbString.WriteString("/")
	dbString.WriteString(dbn)
	dbString.WriteString("?retryWrites=true&w=majority")

	//TODO FIND OR DO BETTER LOGGER FOR DEBUGING.
	// Use for debug only, security issues on production log.Println(dbString.String())

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

//Upsert an map to db
func Upsert(entity map[string]interface{}) (ok bool, id interface{}) {
	ok = true
	collection := Client.Database("linx").Collection("product")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eBson, err := bson.Marshal(entity)
	if err != nil {
		log.Fatal(err)
		ok = false
	}
	var rOptions options.ReplaceOptions
	rOptions.SetUpsert(true)
	res, err := collection.ReplaceOne(ctx, bson.M{"ID": entity["ID"]}, eBson, &rOptions)
	if err != nil {
		log.Fatal(err)
		ok = false
	}

	id = res.UpsertedID
	return

}

//FindAll gets all elements from collection
func FindAll(col string) (ok bool, entity []map[string]interface{}) {
	log.Println(col)
	ok = true
	collection := Client.Database("linx").Collection(col)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.M{})
	err = cur.All(ctx, &entity)

	if err != nil {
		log.Println(" Erro no FindAll")
		log.Fatal(err)

		ok = false
		return
	}
	return
}

//FindByID receives id and returns a ok true if found. Product will contain data from the entity db.
func FindByID(id string, col string) (ok bool, entity map[string]interface{}) {
	ok = true
	collection := Client.Database("linx").Collection(col)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"ID": id}).Decode(&entity)
	if err != nil {
		log.Println(" Erro no Find")
		log.Fatal(err)

		ok = false
		return
	}

	return

}

func getEnvVars() (u string, pass string, dbName string, dbHost string, ok bool) {
	ok = true

	err := godotenv.Load("daos/.env")
	if err != nil {

		log.Fatal(err)
		ok = false

	}
	u = os.Getenv("DB_USER")
	pass = os.Getenv("DB_PASS")
	dbName = os.Getenv("DB_NAME")
	dbHost = os.Getenv("DB_HOST")
	return
}
