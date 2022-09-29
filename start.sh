echo "start Twitta"

go mod tidy && go build -o ./build/start main.go && ./build/start
