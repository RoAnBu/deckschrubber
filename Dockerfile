FROM golang:latest

COPY . /go/src/github.com/RoAnBu/deckschrubber

WORKDIR /go/src/github.com/RoAnBu/deckschrubber
# Repo file path necessary for testing if file is valid
ENV DS_REPO_FILE_PATH=/go/src/github.com/RoAnBu/deckschrubber/repositories.yml
RUN go get
RUN go test ./... && CGO_ENABLED=0 GOOS=linux go build -a .

FROM alpine:latest

RUN apk --update upgrade && apk add ca-certificates && mkdir -p /deckschrubber/config && rm -rf /var/cache/apk/*

WORKDIR /deckschrubber
COPY --from=0 /go/src/github.com/RoAnBu/deckschrubber/deckschrubber deckschrubber
COPY repositories.yml repositories.yml

ENV DS_REPO_FILE_PATH=/deckschrubber/repositories.yml

ENTRYPOINT ["./deckschrubber"]