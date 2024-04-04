FROM golang:1.22.2 AS build

WORKDIR /app

COPY go.mod *.go ./
RUN go mod download

COPY . .

RUN go build -o /app/markdown-thing

EXPOSE 8080

CMD ["/app/markdown-thing"]