FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/amlak-api ./cmd/api

FROM alpine:3.22
WORKDIR /app
COPY --from=build /bin/amlak-api /usr/local/bin/amlak-api
EXPOSE 8080
CMD ["amlak-api"]

