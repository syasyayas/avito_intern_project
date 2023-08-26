# Step 1: Modules caching
FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.21-alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd
RUN chmod +x /bin/app


# Step 3: Final
FROM alpine
COPY --from=builder /app/config.yml /
COPY --from=builder /bin/app /app
CMD ["/app", "-config=config.yml"]
