# Multi-stage build: Go server + static assets

FROM golang:1.22-alpine AS builder
WORKDIR /src

# Copy server source
COPY board-3d/server/ ./board-3d/server/

# Initialize module (no external deps)
WORKDIR /src/board-3d/server
RUN go mod init board3d && go mod tidy
RUN CGO_ENABLED=0 go build -o /server

FROM alpine:3.19
WORKDIR /app

# Copy binary
COPY --from=builder /server /server

# Copy minimal static tree (AH only)
COPY board-3d/ /app/AH/board-3d/
COPY img/ /app/AH/img/

ENV PORT=8080
EXPOSE 8080

CMD ["/server", "-addr", ":8080", "-root", "/app"]

