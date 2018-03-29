# Telegram bot lambda
1. Build: `GOOS=linux go build -o main`
2. Archive: `rm -rf deployment.zip && zip deployment.zip main && rm -rf main`
3. Upload...
3.1. Via aws cli: `aws lambda update-function-code --profile self-mike --function-name telegram-bot --zip-file fileb://deployment.zip`

## Config
- rename .env.default to .env 

## Run tests
- `go test ./...`

## gocov Installation
- `go get github.com/axw/gocov/gocov`
- `go get -u gopkg.in/matm/v1/gocov-html`

## Calculate coverage
- `gocov test ./... | gocov-html > coverage.html`


## Godep
- clear vendor folder (except .gitignore)
- `go get github.com/tools/godep`
- Save deps: `godep save -t ./...`
- Restore deps: `godep restore`
