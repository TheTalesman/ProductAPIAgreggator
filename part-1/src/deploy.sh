#just while dev
#rm test.json

CGO_ENABLED=0 go build -o product-api

#shutdown any container still running
docker-compose down
docker-compose up  --scale product_api=3

