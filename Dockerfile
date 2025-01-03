FROM golang:1.22

LABEL authors="dimon"

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o TODO-List ./cmd/main.go

CMD ["./TODO-List"]