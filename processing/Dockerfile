FROM registry.emcdtech.com/docker-images/golang:1.23-alpine as modules
WORKDIR /modules
ADD go.mod go.sum ./
RUN go mod download

FROM registry.emcdtech.com/docker-images/golang:1.23-alpine as builder
WORKDIR /app
COPY --from=modules /go/pkg /go/pkg
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o /app/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o /app/coinwatch_consumer ./cmd/coinwatch_consumer
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o /app/workers ./cmd/workers

FROM registry.emcdtech.com/docker-images/alpine:latest
WORKDIR /app
COPY --from=builder /app /app

CMD ["/app/api"]
