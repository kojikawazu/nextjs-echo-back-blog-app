# ビルドステージ
FROM golang:1.19 as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

# 実行ステージ
FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /app/main /app/main
CMD ["/app/main"]