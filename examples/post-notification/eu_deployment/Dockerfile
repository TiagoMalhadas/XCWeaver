FROM golang:1.21.1

WORKDIR /src


COPY main.go ./
COPY post_storage_europe.go ./
COPY post_upload.go ./
COPY notifier.go ./
COPY weaver.toml ./

RUN go install github.com/TiagoMalhadas/xcweaver/cmd/xcweaver@latest
RUN export PATH="$PATH:$HOME/go/bin"
RUN go mod init eu_deployment
RUN go mod tidy
RUN xcweaver generate
RUN go build

EXPOSE 12345

CMD ["xcweaver", "multi", "deploy", "weaver.toml"]

#CMD ["xcweaver", "single", "deploy", "weaver.toml"]
