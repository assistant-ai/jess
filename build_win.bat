
set /p VERSION=<VERSION
set "CGO_ENABLED=1"
set "GOOS=windows"
set "GOARCH=amd64"

go build -ldflags="-X main.version=%VERSION%" -o "jess-%GOOS%-%GOARCH%.exe"