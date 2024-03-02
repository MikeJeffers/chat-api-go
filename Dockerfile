FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ls -la 
# RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/main.go
# EXPOSE 3000
# CMD ["/server"]