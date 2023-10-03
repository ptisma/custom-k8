FROM golang:1.19 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/sample-api ./cmd/sample-api


FROM alpine:3.14

WORKDIR /app

COPY --from=builder /app/sample-api .

COPY configs/ .   

EXPOSE 8080

CMD ["./sample-api"]
