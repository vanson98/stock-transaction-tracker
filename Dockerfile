FROM golang:1.22-alpine3.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o stk-service-app ./cmd

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/stk-service-app .
COPY --from=builder /app/app.production.env .
# COPY --from=builder /app/.env .
ENV STK_SERVICE_RUN_MODE=PRODUCTION
EXPOSE 5050
CMD ["./stk-service-app" ]

