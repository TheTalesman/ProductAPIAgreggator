package dao

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Client is the DB global client
var Client *mongo.Client

//RClient Redis global client
var RClient *redis.Client

func rClient() (err error) {
	RClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "Redis2019!",
	})
	_, err = RClient.Ping(context.Background()).Result()
	return
}

//ConnectRedis to Redis
func ConnectRedis() (ok bool) {
	err := rClient()
	if err != nil {
		log.Fatal(err)
		ok = false
	}
	return
}

//Connect to database
func Connect() (ok bool) {
	//Assert once an deny on any error.
	ok = true

	url, ok := getEnvVars()

	//TODO FIND OR DO BETTER LOGGER FOR DEBUGING.
	// Use for debug only, security issues on production log.Println(dbString.String())
	log.Println("DB STRING:", url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))

	if err != nil {
		log.Fatal(err)
		ok = false
	}
	Client = client
	return
}

//Upsert an map to db
func UpsertMany(entity []map[string]interface{}, col string) (ok bool, err error) {
	ok = true
	collection := Client.Database("linx").Collection(col)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()

	var operations []mongo.WriteModel

	log.Println("Initing Ops")
	Total := len(entity)
	i := 0
	factor := Total / 10
	h := factor
	var wg sync.WaitGroup
	for i < Total {

		wg.Add(1)

		go func(i int, h int, wg *sync.WaitGroup) {
			for i < h {
				nEntity := entity[i]
				op := mongo.NewUpdateOneModel()
				op.SetFilter(bson.M{"ID": nEntity["ID"]})

				op.SetUpdate(bson.M{"$set": nEntity})

				op.SetUpsert(true)
				operations = append(operations, op)

				i++

			}
			log.Println("Finished appending Ops:", i)
			log.Println("of a Total: ", Total)
			bulkOption := options.BulkWriteOptions{}
			bulkOption.SetOrdered(true)

			_, err := collection.BulkWrite(ctx, operations, &bulkOption)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(i, h, &wg)
		wg.Wait()
		h += factor
		i += factor
		if h > Total {
			h = Total
		}
	}

	return

}

//Upsert an map to db
func Upsert(entity map[string]interface{}, col string) (ok bool, id interface{}) {
	ok = true
	collection := Client.Database("linx").Collection(col)
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

func getEnvVars() (url string, ok bool) {
	ok = true

	err := godotenv.Load("daos/.env")
	if err != nil {

		log.Fatal(err)
		ok = false

	}
	u := os.Getenv("DB_USER")
	p := os.Getenv("DB_PASS")
	//dbn := os.Getenv("DB_NAME")
	dbh := os.Getenv("DB_HOST")

	var dbString strings.Builder

	//for the sake of memory
	dbString.WriteString("mongodb://")
	dbString.WriteString(u)
	dbString.WriteString(":")
	dbString.WriteString(p)
	dbString.WriteString("@")
	dbString.WriteString(dbh)
	dbString.WriteString("/")
	url = dbString.String()
	return
}
