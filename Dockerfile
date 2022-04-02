FROM golang:1.18 as builder

# first (build) stage
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o dionysos-server

# final (target) stage
FROM alpine
COPY --from=builder /app/dionysos-server /
CMD ["/dionysos-server"]
