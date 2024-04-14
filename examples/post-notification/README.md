# Post-Notification / XCWeaver

Implementation of a geo-replicated Post-Notification application that shows the occurence of cross-service inconsistencies.

## Requirements

- [Golang >= 1.21](https://go.dev/doc/install)
- [XCWeaver](https://github.com/TiagoMalhadas/xcweaver)


## LOCAL Deployment

### Running Locally with XCWeaver in Multi Process

Deploy datastores:

``` zsh
chmod +x redis.sh
chmod +x rabbitMQ.sh
./redis.sh
./rabbitMQ.sh
```

Deploy eu_deployment:

``` zsh
cd eu_deployment
docker build -t eu_deployment .
docker run -it --rm --net rabbits -p 12345:12345 eu_deployment
```

Deploy us_deployment:

``` zsh
cd ../eu_deployment
docker build -t us_deployment .
docker run -it --rm --net rabbits us_deployment
```

Run benchmark:

``` zsh
cd ..
chmod +x test.sh
./test.sh
```

Gather metrics:
``` zsh
TO-DO
```

### Additional

#### Manual Testing of HTTP Requests

**Publish Post**: {post}

``` zsh
curl "localhost:12345/post_notification?post=POST"

# e.g.
curl "localhost:12345/post_notification?post=my_first_post"
```

