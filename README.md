# short-url

## Easy Project

Short URL is `DEFAULT_SCHEMA + DEFAULT_DOMAIN_NAME + sequence`

I used the [andyxning](https://github.com/andyxning/shortme) code in the project.

>Example:
>
>DEFAULT_SCHEMA = https
>
>DEFAULT_DOMAIN_NAME = short-url.example.com
>
>The system will be automatic generation of sequence
>
>Create Short URL is

```
https://short-url.example.com/tas8
```
 
**Docker RUN**
```
docker run -d -p 8080:8080 --name short-url --restart always \
    -e DEFAULT_DB_NAME=short-url.db \
    -e DEFAULT_SCHEMA=http \
    -e DEFAULT_DOMAIN_NAME=127.0.0.1:8080 \
    xushikuan/short-url:1.0
```

**If you want save database**

```
docker run -d -p 8080:8080 --name short-url --restart always \
    -e DEFAULT_DB_NAME=short-url.db \
    -e DEFAULT_SCHEMA=http \
    -e DEFAULT_DOMAIN_NAME=127.0.0.1:8080 \
    -v `pwd`/data:/go/data \
    xushikuan/short-url:1.0
```

docker-compose.yml

```
version: '3.3'

services:

  short-url:
    build: .
    image: "xushikuan/short-url:1.0"
    environment:
      DEFAULT_DB_NAME: short-url.db
      DEFAULT_SCHEMA: http
      DEFAULT_DOMAIN_NAME: 127.0.0.1:8080
    deploy:
      replicas: 1
    volumes:
      - ${PWD}/data:/go/data
    ports:
      - 8080:8080
    restart: always
```

**docker stack deploy -c docker-compose.yml short**

```
cookie$ docker stack deploy -c docker-compose.yml short
Ignoring unsupported options: build, restart

Creating network short_default
Creating service short_short-url
```