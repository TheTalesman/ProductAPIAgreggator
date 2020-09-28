# ProductAPIAgreggator

## REQUIREMENTS
go 1.14 +, gem(ruby), docker, Linux OS in hosting machine

## DESCRIPTION
This project is made up of 2 parts. 
The first part features are:
- Upsert large batches via json in a statless application. 
- Filter duplicate requests with a cache approach (redis) between containers (docker)
- Persists data on mongo (also in docker)
- Traefik routes the clusters as a load balancer using round robin.

The second part features are:
- Sanitizer: Reads a dump file with data, make requests to an external images api, validates and transforms it in structure to be sent to an upgrader.
- Upgrader: receives requests via POST and upserts into the db.


## HOW TO RUN?

1 - Clone repo:
`git clone https://github.com/TheTalesman/ProductAPIAgreggator.git` or unzip the folder in your $GOPATH/src

2 - Vendorize dependencies:
```
cd $GOPATH/go/src/ProductAPIAggregator
go mod init
go mod vendor
```

This project includes dockerfiles and deploy.sh scripts, some scripts may ask for sudo permissions, you can rest assured that there's no catchup on the scripts, it just may be necessary.

This project should run in any linux machine, but if you intend to run it on a RHEL family OS (Fedora, CentOS, RedHat...) you should set SELinux to permissive because of permission problems with docker.
```
### IF RHEL OS ONLY ###
sudo setenforce 0
```


## Part 1
Please free up ports 80, 8080, 8081, 8082, 4567, 27017. 

1 - .env file
```
cd part-1/src/
touch .env
```

edit .env file and insert the environment variables:
```
API_GIN_DEBUG_MODE=true
DB_USER={dbuser}
DB_PASS={password}
DB_HOST=mongo
DB_NAME=linx
```


2 - deploy with sudo (if you have docker installed for user with permissions setup correctly you can just ./deploy)
```
cd part-1/src/
sudo ./deploy.sh
```

3- Fixture File for Part-1
If you need a file to post in the api, you can fixture one. Go to docker-compose.yaml and change the line 26
`command: ./product-api false`
to 
`command: ./product-api true`
do a ./deploy.sh and it will fixture a rbf.json you can use as request body
change it back again to false to run the api. (parameter true will only fixture)

4- Post src/rbf.json as body request to localhost/products, the application will persist

5- GET localhost/products will return persisted objects. As this is a async operation, the response will return 200 ok but the objects will still be in processing, so you should look in the logs to the end of processing before request this GET.


## Part 2

1 - .env file
From the root folder:
```
cd part-2/src/
touch .env
```

edit .env file and insert the environment variables:
```
API_GIN_DEBUG_MODE=true
DB_USER={dbuser}
DB_PASS={password}
DB_HOST=mongo
DB_NAME=linx
```

Start the updater, mongo and traefik lb
`sudo /deploy.sh`


Drop a file named input-dump in ProductApiAggregator/part-2/src/dump some json like:
```
{"productId":"pid482","image":"http://localhost:4567/images/167410.png"}
{"productId":"pid1613","image":"http://localhost:4567/images/122577.png"}
{"productId":"pid7471","image":"http://localhost:4567/images/177204.png"}

```

Start the sanitizer
`sanitizer/san`

Watch the logs!

### SANITIZER 
#### Description
Sanitizer reads a dump file (/part-2/external/input-dump) parses each line into a request in the external images api (/part-2/external/url-aggregator-api.rb).
Sanitizer will read lines, order by "pid" and send all lines in a channel.

The workers read product lines, do a filter, inputing the first 3 valid image url's in an array atribute of the product. Batches of products will be made and be sent to the updater for persistance.

#### Fine Tuning
Using nWorkers =5 and batchSize 10 was the best case cenario found until now. Running in a notebook with a processor  i5-4200M CPU @ 2.50GHz and 8GB RAM. 
Speed ~ 1200 products/min.
Improvements in database connections, clustering mongo and so should be next step in speeding up.
If your speed is below 1200 products/min or if you have network rates/limits, try playing around with nWorkers and batchSize.

If sanitizer drows updater or images api you may want to addup sleepTime, which is the time between jobs are sent. Limit your WIP! ;)
##### CheatSheet
"+ batchSize = - network calls"
"+ batchSize = + request body size"
"+ nWorkers = +network calls"

### UPDATER
#### Description
Updater will receive parsed products as requests from sanitizer. It will them aggregate the lines in bulkWrites (to lessen number of connections to mongo) and wiull upsert the products in db.

Updater use the daos of part-1 for code reuse.




## TODO

- [ ] Dockerize url-aggregator-api.rb
- [ ] Implement Logger with debug level
- [x] Check async  part-1
- [ ] TRANSLATE ALL COMMENTS
- [x] ORGANIZE PROJECT

## PART 1


- [X] PRODUCT API
- [ ] HANDLE JSON SIZE ~5 GB (TO BE TESTED)
- [X] HANDLE REDUDANT REQUESTS BETWEEN 10 MIN WITH 403
- [X] STATELESS API

## PART 2
- [ ] Install CRON JOB
- [X] ENDPOINT FOR IMAGES DUMP
- [x] FILTER IMAGES FOR CODE 200
- [x] SANITIZE AND AGGREGATE IN IMAGES ARRAY (MAX 3) TO USE WITH API



## REQS
- [X] Deve funcionar em um ambiente Linux
- [ ] Deve ter testes automatizados
- [x] Deve ter um README explicando como instalar as dependências, executar as soluções e os testes.


# A Word about tests
Started the project trying with TDD, but I really noticed that if I was going to test everything the right way it would take long enough to lose the project deadline. So I opted leaving the tests behind to workup as next steps. 
To make tests work part-1/src:$ go test ./...
It will not find .env file because go needs a fix in local imports. Theres only a small number of tests in part-1 that hasn`t been used and probably are not useful enough for now. 

Besides that, it has been a wonderful project and lesson. This application relies heavily on logging, as a antifragile way to react fast. With a better logging tool/exposing logging files it can get more robust and even anti fragile with some kind of setup via SNS and Cloudwath, for example.
Traefik is a good option too, but it needs some more setup to enable tracing, metris and access log.

# Final Words
If you need any assistance feel free to reach me at alisonsvargas@gmail.com
