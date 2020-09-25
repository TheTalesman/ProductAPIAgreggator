#just while dev
#rm test.json

CGO_ENABLED=0 go build -o product-api


#sudo docker run -p 6379:6379  --name=myRedis -d redis redis-server --requirepass Redis2019!

#sudo docker build -t product_api:001 .
#sudo docker run -p 8080:8080 product_api:001 ./product-api

sudo docker-compose up -d

