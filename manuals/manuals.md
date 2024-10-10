# マニュアル

## プロジェクトの立ち上げ

```bash
mkdir -p backend
cd backend
go mod init backend
```

## モジュールのインストール

```bash
go get github.com/labstack/echo/v4
go get github.com/joho/godotenv
go get github.com/jackc/pgx/v4/pgxpool
go get github.com/stretchr/testify/assert
```

## モジュールの整理

```bash
go mod tidy
```

## サーバーの起動

```bash
go run main.go
```

## テスト

```bash
JWT_SECRET_KEY=xxxxxx go test -count=1 ./...
```