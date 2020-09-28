
docker-compose down
cd updater/
go build -o updater
cd ../sanitizer/
go build -o san
cd ../../external/
nohup ruby url-aggregator-api.rb &
cd ../src
docker-compose up --force-recreate --scale updater=5
