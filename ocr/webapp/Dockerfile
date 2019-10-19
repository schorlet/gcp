FROM golang:1.13.1-alpine3.10 as builder

# Copy local code to the container image.
WORKDIR /workspace
COPY . .

# Build the command inside the container.
RUN go get github.com/tdewolff/minify/cmd/minify
RUN go generate
RUN CGO_ENABLED=0 GOOS=linux go build -v -o webapp



# Use a Docker multi-stage build to create a lean production image.
FROM alpine:3.10
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary to the production image from the builder stage.
COPY --from=builder /workspace/webapp /app/webapp
COPY --from=builder /workspace/assets/index.html /app/assets/
COPY --from=builder /workspace/assets/app.min.js /app/assets/
COPY --from=builder /workspace/assets/style.min.css /app/assets/

# Run the web service on container startup.
CMD ["/app/webapp"]
