cd ..
export $(grep -v '^#' .env | xargs -d '\n')
cd chat-api-go
go run .