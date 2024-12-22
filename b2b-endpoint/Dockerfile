FROM registry.emcdtech.com/docker-images/golang:1.23-alpine as modules
WORKDIR /modules
ADD go.mod go.sum ./
RUN go mod download

FROM registry.emcdtech.com/docker-images/golang:1.23-alpine as builder
ARG SWAGGER=false
WORKDIR /app
COPY --from=modules /go/pkg /go/pkg
COPY . .
ENV SWAGGER=$SWAGGER
RUN if [ "$SWAGGER" = "true" ]; \
  then go install github.com/swaggo/swag/cmd/swag@v1.8.4 && \
  swag init --parseDependency; \
  fi;
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o /app/application main.go

FROM registry.emcdtech.com/docker-images/alpine:latest
WORKDIR /app
COPY --from=builder /app/application /app/app

CMD ["/app/app"]
