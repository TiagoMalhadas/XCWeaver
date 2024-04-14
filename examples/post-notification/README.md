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
go generate
go build
weaver multi deploy weaver.toml
```

Deploy us_deployment:

``` zsh
./manager.py init-social-graph --local
```

Run benchmark:

``` zsh
./manager.py wrk2 --local
```

Gather metrics:
``` zsh
./manager.py metrics --local
```

Clean datastores:

``` zsh
./manager.py storage-clean --local
```

### Additional

#### Manual Testing of HTTP Requests

**Publish Post**: {post}

``` zsh
curl "localhost:12345/post_notification?post=POST"

# e.g.
curl "localhost:12345/post_notification?post=my_first_post"
```

