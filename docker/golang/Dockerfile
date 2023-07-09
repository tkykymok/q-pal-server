FROM golang:1.19.3-alpine
COPY ./app /app
COPY ./app/go.mod /app

RUN apk update &&  \
    apk add alpine-sdk &&  \
    go install github.com/cosmtrek/air@latest && \
    go install github.com/volatiletech/sqlboiler/v4@latest && \
    go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app
RUN go mod tidy

WORKDIR /app/api
RUN go get
CMD ["air"]
