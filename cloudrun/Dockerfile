FROM golang:1.23.2-alpine3.20 AS build

WORKDIR /app

COPY . /app/src

RUN cd ./src && go build -o /app/main . && rm -rf /app/src

CMD ["/app/main"]

EXPOSE 8080
