FROM registry.emcdtech.com/docker-images/golang:1.22-alpine as modules
WORKDIR /modules
ADD go.mod go.sum ./
RUN go mod download

FROM registry.emcdtech.com/docker-images/golang:1.22-alpine as builder
ARG SWAGGER=false
WORKDIR /app
COPY --from=modules /go/pkg /go/pkg
COPY . .
RUN GOOS=linux go build -ldflags '-w -s' -a -o /app/application main.go

FROM registry.emcdtech.com/docker-images/alpine:latest
WORKDIR /app
COPY --from=builder /app/application /app/app

CMD ["/app/app"]
