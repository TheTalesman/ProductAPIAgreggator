# ProductAPIAgreggator

## REQS
- [ ] Deve funcionar em um ambiente Linux
- [ ] Deve ter testes automatizados
- [ ] Deve ter um README explicando como instalar as dependências, executar as soluções e os testes.

# TODO

## PART 1
- [ ] PRODUCT CRUD
- [ ] HANDLE JSON SIZE ~5 GB
- [ ] HANDLE REDUDANT REQUESTS BETWEEN 10 MIN WITH 403
- [ ] STATELESS API

## PART 2

- [ ] ENDPOINT FOR IMAGES DUMP
- [ ] FILTER IMAGES FOR CODE 200
- [ ] SANITIZE AND AGGREGATE IN IMAGES ARRAY (MAX 3) TO USE WITH API


# DESCRIPTION
## Parte 1 - API de produtos

Precisamos de uma API para receber a atualização de dados cadastrais de produtos. Ela deve receber um corpo no formato JSON, onde o tamanho varia desde alguns poucos Kb até alguns Gb.




Experiências anteriores mostram que alguns clientes costumam enviar o mesmo corpo repetidas vezes ao longo de um curto espaço de tempo.
Isso nos causou alguns problemas, como o fato de ter que escalar nossos bancos de dados muito além do necessário afim de aguentar a carga extra desnecessária.
Para evitar que isto ocorra, precisamos que esta API negue requisições que tem o mesmo corpo num intervalo de 10 minutos.


Aqui está um exemplo do comportamento esperado:
```bash
# 2018-03-01T13:00:00 - primeira requisição, durante 10 minutos requests com o mesmo corpo serão negadas
curl -XPOST http://your-api.chaordic.com.br/v1/products -d '[{"id": "123", "name": "mesa"}]' #=> 200 OK

# 2018-03-01T13:09:59 - mesmo corpo que a request anterior.
curl -XPOST http://your-api.chaordic.com.br/v1/products -d '[{"id": "123", "name": "mesa"}]' #=> 403 Forbidden

# 2018-03-01T13:10:00 - agora a API deve voltar a aceitar o corpo
curl -XPOST http://your-api.chaordic.com.br/v1/products -d '[{"id": "123", "name": "mesa"}]' #=> 200 OK
```
Como esta API atenderá milhares de requisições simultâneas, ela precisa funcionar em um cluster.
É esperado que o comportamento descrito acima se mantenha, independente do nó que receber a requisição.


## Parte 2 - Agregador de URLs


Recebemos um dump com lista de URLs de imagens de produtos que vamos utilizar para manter nossa base de dados atualizada.
Este dump contém imagens de milhões de produtos e URLs, e é atualizado a cada 10 minutos:

```json
{"productId": "pid2", "image": "http://www.linx.com.br/platform-test/6.png"}
{"productId": "pid1", "image": "http://www.linx.com.br/platform-test/1.png"}
{"productId": "pid1", "image": "http://www.linx.com.br/platform-test/2.png"}
{"productId": "pid1", "image": "http://www.linx.com.br/platform-test/7.png"}
{"productId": "pid1", "image": "http://www.linx.com.br/platform-test/3.png"}
{"productId": "pid1", "image": "http://www.linx.com.br/platform-test/1.png"}
{"productId": "pid2", "image": "http://www.linx.com.br/platform-test/5.png"}
{"productId": "pid2", "image": "http://www.linx.com.br/platform-test/4.png"}
```

As URLs pertencem a uma empresa terceirizada que hospeda a maioria destas imagens, e ela nos cobra um valor fixo por cada request.
Já sabemos que o dump de origem não tem uma boa confiabilidade, pois encontramos várias imagens repetidas e boa parte delas também retornam status 404.
Como não é interessante atualizar nossa base com dados ruins, filtramos apenas as URLs que retornam status 200.


O processo de atualização deve receber como input um dump sanitizado, onde o formato é ligeiramente diferente da entrada:

```json
{"productId": "pid1", "images": ["http://www.linx.com.br/platform-test/1.png", "http://www.linx.com.br/platform-test/2.png", "http://www.linx.com.br/platform-test/7.png"]}
{"productId": "pid2", "images": ["http://www.linx.com.br/platform-test/3.png", "http://www.linx.com.br/platform-test/5.png", "http://www.linx.com.br/platform-test/6.png"]}
```

Para diminuir a quantidade de requests necessárias para validar as URLs, decidimos limitar a quantidade de imagens por produto em até 3.
O seu objetivo é criar um programa que gera o dump final no menor tempo possível e com o mínimo de requests desnecessárias (já que existe um custo fixo por requisição).

O arquivo [input-dump.gz](./input-dump.gz) é um exemplo do dump de entrada. E você pode usá-lo para testar sua implementação.
Também criamos uma api que responde as URLs do `input-dump.gz`. Ela é apenas um mock, mas vai te ajudar a implementar a solução do desafio. Para executá-la, basta:

```shell
gem install sinatra
ruby url-aggregator-api.rb
```