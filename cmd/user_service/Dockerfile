# build stage
FROM golang:1.24-alpine AS build-env

RUN apk --no-cache add git
ENV GOPRIVATE=github.com/Jozzo6

WORKDIR /go/src/github.com/Jozzo6/casino_loyalty_reward_system/

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o user_service cmd/user_service/main.go

# final image
FROM alpine:3.19
RUN apk update && apk add ca-certificates
WORKDIR /app
COPY --from=build-env /go/src/github.com/Jozzo6/casino_loyalty_reward_system/user_service /app/

ENTRYPOINT /app/user_service