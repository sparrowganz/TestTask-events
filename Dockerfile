FROM golang:latest as builder

WORKDIR /app

COPY . .
COPY cmd/web/main.go .
RUN go mod vendor

RUN CGO_ENABLED=0 GO111MODULE=on go build -mod=vendor -v -o app .


FROM alpine:latest

WORKDIR /home

COPY --from=builder /app/config/config.yml .
COPY --from=builder /app/app .

ENTRYPOINT ["./app"]
CMD ["--config=config.yml"]