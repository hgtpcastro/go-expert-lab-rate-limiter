# go-expert-lab-rate-limiter
Rate limiter em Go que pode ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

O rate limiter é configurável para realizar a checagem de limites por IP ou token `API_KEY`, e utiliza o Redis como _storage_.

A configuração é realizada através de variáveis de ambiente declaradas no arquivo `.env`

## Configuração
Ajuste-o o arquivo `.env` presente na pasta _./cmd/app_, conforme necessidade. Por padrão, os seguintes valores são utilizados:

```sh
APP_PORT=8080 # Porta do servidor Web

# Configurações do Redis
REDIS_HOST="localhost"
REDIS_PORT=6379
REDIS_PASSWORD=""
REDIS_DB=0

RATE_MAX_REQUESTS_BY_IP=10 # Número máximo de requisições por IP
RATE_MAX_REQUESTS_BY_TOKEN=100 # Número máximo de requisições por token
RATE_PERIOD_WINDOW_SECONDS=60 # Período de tempo em segundos
```

### Buildar a imagem docker e inicar a aplicação
```bash
    make start
```

### Parar a aplicação
```bash
    make stop
```

### Remover os containers
```bash
    make clean
```

### Rodar os testes de carga
```bash
    make test-smoke # Teste de carga do tipo smoke (duração de 1 minuto);
    make test-stress # Teste de carga do tipo stress(duração de 4 minuto);
```

### Exemplos de requisições via curl

- **Requisição por IP com sucesso:**

```sh
$ curl -vvv http://localhost:8080
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< X-Ratelimit-Limit: 10
< X-Ratelimit-Remaining: 8
< X-Ratelimit-Reset: 1738976031
< Date: Sat, 08 Feb 2025 00:53:08 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```
- **Requisição por IP bloqueada:**

```sh
$ curl -vvv http://localhost:8080
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 429 Too Many Requests
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< X-Ratelimit-Limit: 10
< X-Ratelimit-Remaining: 0
< X-Ratelimit-Reset: 1738976134
< Date: Sat, 08 Feb 2025 00:54:45 GMT
< Content-Length: 95
< 
you have reached the maximum number of requests or actions allowed within a certain time frame
* Connection #0 to host localhost left intact
```

- **Requisição por Token com sucesso:**

```sh
$ curl -H 'API_KEY: abc123' -vvv http://localhost:8080
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> API_KEY: abc123
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< X-Ratelimit-Limit: 100
< X-Ratelimit-Remaining: 98
< X-Ratelimit-Reset: 1738976316
< Date: Sat, 08 Feb 2025 00:58:14 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```

- **Requisição com checagem via token bloqueada:**

```sh
$ curl -H 'API_KEY: abc123' -vvv http://localhost:8080
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> API_KEY: abc123
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 429 Too Many Requests
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< X-Ratelimit-Limit: 100
< X-Ratelimit-Remaining: 0
< X-Ratelimit-Reset: 1738976421
< Date: Sat, 08 Feb 2025 00:59:43 GMT
< Content-Length: 95
< 
you have reached the maximum number of requests or actions allowed within a certain time frame
* Connection #0 to host localhost left intact
```

### Exemplos de requisições via arquivo .http (VSCode: REST Client Plugin)

Navegue até a pasta api no diretório raiz do projeto


```sh
request_with_token.http
request_without_token.http
```

## <a name="license"></a> License

Copyright (c) 2025 [Hugo Castro Costa]

[Hugo Castro Costa]: https://github.com/hgtpcastro