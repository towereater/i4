FROM golang:1.24.5

WORKDIR /app

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

ENTRYPOINT ["/main"]
CMD ["./config.json"]