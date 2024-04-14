FROM golang:1.21.1

WORKDIR /src

COPY follower_notify.go ./
COPY main.go ./
COPY notifier.go ./
COPY post_storage_america.go ./
COPY metrics.go ./
COPY weaver.toml ./

RUN go install github.com/TiagoMalhadas/xcweaver/cmd/xcweaver@latest
RUN export PATH="$PATH:$HOME/go/bin"
RUN go mod init us_deployment
RUN go mod tidy
RUN xcweaver generate
RUN go build

CMD ["xcweaver", "multi", "deploy", "weaver.toml"]

#CMD ["xcweaver", "single", "deploy", "weaver.toml"]
