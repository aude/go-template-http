FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o "go-template-http"

FROM scratch
CMD ["/go-template-http"]
COPY --from=builder /app/go-template-http /go-template-http