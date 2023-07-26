FROM golang:1.19-alpine

RUN apk update && apk upgrade --no-cache

WORKDIR /app

RUN apk --no-cache add ca-certificates \
    && apk add --no-cache make bash git

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

COPY . .
COPY .env .

RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN go install -mod=mod golang.org/x/tools/gopls
RUN go install github.com/swaggo/swag/cmd/swag@latest

# ENTRYPOINT CompileDaemon --build="go build -a -o ./build/app/webapi ./cmd/daemon/main.go" --command="./build/app/main.go"
ENTRYPOINT CompileDaemon --build="make go-build" --command="make run"
