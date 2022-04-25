FROM golang:1.18.1-buster

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o connector .

ENTRYPOINT ["/app/connector"]