FROM golang:latest
WORKDIR /tix-worker
COPY . .
RUN go mod download
RUN env GOOS=linux GOARCH=amd64 go build -o worker ./cmd/main.go
EXPOSE 8080
CMD ["./worker"]
