FROM golang:1.21.5

#RUN mkdir /socialnetwork

WORKDIR /scr


#COPY main.go ./
#COPY pkg ./pkg
#COPY weaver-eu.toml ./

COPY . .

RUN go install github.com/ServiceWeaver/weaver/cmd/weaver@v0.22.1-0.20231019162801-c2294d1ae0e8
RUN export PATH="$PATH:$HOME/go/bin"
RUN go generate
RUN go build

EXPOSE 9000

CMD ["weaver", "multi", "deploy", "weaver-eu.toml"]

#CMD ["weaver", "single", "deploy", "weaver-eu.toml"]
