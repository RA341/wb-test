FROM golang:1 AS gobuild

WORKDIR /app

COPY go.* .

RUN go mod download

COPY main.go .

RUN go build -ldflags "-s -w" -o wb "main.go"

FROM scratch

COPY --from=gobuild /app/wb .

ENTRYPOINT ["./wb"]