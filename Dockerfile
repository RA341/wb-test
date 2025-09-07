FROM golang:1 AS gobuild

WORKDIR /app

COPY go.* .

RUN go mod download

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o wb "main.go"

FROM scratch

WORKDIR /app

COPY --from=gobuild /app/wb .

COPY index.html index.html

ENTRYPOINT ["./wb"]