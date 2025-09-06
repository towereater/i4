FROM golang:1.24.5

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

ENTRYPOINT ["/main"]
CMD ["./config.json"]