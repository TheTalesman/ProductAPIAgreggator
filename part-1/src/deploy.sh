#just while dev
#rm test.json
cd ../../docker/ 
sudo docker-compose up -d
cd ../part-1/src
go build 
./src