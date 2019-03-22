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
## API

* POST `http://localhost:8080/compress-url`

Body

```
{
  "long-url": "https://lh3.googleusercontent.com/krcdFE-XBYxynEbovjtyN6jkMpg5QfZeS1ohq1eAOdyXk3-T6Iu1bnFyTO5U3jS7KIyK_mNiHEIYUKO5eSWHgnSqyz85dtAezNplFwDpxVUzZkJIDc7nLLHu--SnPuDO6EyO3jbj1E9RnEj70UQuUvC116fqoKGunahjNUATuR0OJ_emvKKdiMW8732vmdL_S6otXghWhBccXKOUfFs-jD7yUTIYU-k2xgYWEEs53pgR8wKwi1sJQI8PWa9620wAHJw9ppK9xSANzcA31LIjBQp-AIQ-GEu3yOu2NQ1JPZwh3KEkSMSFk2MpiCeHPgnFLY0wlAfqvNsf33n9VwvtOdfrcwmLia27kzOyCISgTpLx7m72hZa9HMbxVlfa9z5P7lEGBOTXtwaxSKjC08CCLkYka-mXEvX_AtJnMix4krMnWfrDXfMaYPMFwVFJU7yqppMk7G-SjukFm2aBzyaIkWcmh0haKXnz93PxAYx37Fx5F6QwncQpuadCcWu3vf3RWcc5k2WUnntyNvmMmWyoiIa7HguHpStykJCOcHi8CtL9ceLCQQe93OfSMTczKXMXqlRZpc3UC9Pb9jtRqi5G7-pfpME9FTaZ5FHZ6IqP7jv8OcNRd0tvW7SwBRZ-4deKU5_yDqbI33exldHEQX56b5Xx7VrHxFR_tkC828uA2gDMlSVQA35HftwJOs6Qqx8fNmfXhdSnPSlShJTgv7323QfX=w680-h640-no"
}
```

* POST `http://localhost:8080/uncompress-url`

Body

```
{
  "short-url": "http://localhost:8080/t"
}
```


* PUT `http://localhost:8080/short-url/{short-url}`

Body

```
{
  "long-url": "https://lh3.googleusercontent.com/krcdFE-XBYxynEbovjtyN6jkMpg5QfZeS1ohq1eAOdyXk3-T6Iu1bnFyTO5U3jS7KIyK_mNiHEIYUKO5eSWHgnSqyz85dtAezNplFwDpxVUzZkJIDc7nLLHu--SnPuDO6EyO3jbj1E9RnEj70UQuUvC116fqoKGunahjNUATuR0OJ_emvKKdiMW8732vmdL_S6otXghWhBccXKOUfFs-jD7yUTIYU-k2xgYWEEs53pgR8wKwi1sJQI8PWa9620wAHJw9ppK9xSANzcA31LIjBQp-AIQ-GEu3yOu2NQ1JPZwh3KEkSMSFk2MpiCeHPgnFLY0wlAfqvNsf33n9VwvtOdfrcwmLia27kzOyCISgTpLx7m72hZa9HMbxVlfa9z5P7lEGBOTXtwaxSKjC08CCLkYka-mXEvX_AtJnMix4krMnWfrDXfMaYPMFwVFJU7yqppMk7G-SjukFm2aBzyaIkWcmh0haKXnz93PxAYx37Fx5F6QwncQpuadCcWu3vf3RWcc5k2WUnntyNvmMmWyoiIa7HguHpStykJCOcHi8CtL9ceLCQQe93OfSMTczKXMXqlRZpc3UC9Pb9jtRqi5G7-pfpME9FTaZ5FHZ6IqP7jv8OcNRd0tvW7SwBRZ-4deKU5_yDqbI33exldHEQX56b5Xx7VrHxFR_tkC828uA2gDMlSVQA35HftwJOs6Qqx8fNmfXhdSnPSlShJTgv7323QfX=w680-h640-no"
}
```

* GET `http://localhost:8080/short-url`

* GET `http://localhost:8080/short-url/{short-url}`

* DELETE `http://localhost:8080/short-url/{short-url}`

* GET `http://localhost:8080/health`

## Deploy

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