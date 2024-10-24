# Stage 1: Build Go application
FROM golang:1.23 AS go-builder
WORKDIR /app

COPY api/go.mod api/go.sum ./
RUN go mod download

COPY api/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o backend .

# Stage 2: Build Vite app and serve with Node.js
FROM node:20-alpine
WORKDIR /app

# Copy frontend dependencies and install
COPY frontend/package*.json ./
RUN npm ci

# Copy frontend source code and build
COPY frontend/ .
RUN npm run build

# Copy server.js to serve the frontend
COPY frontend/server.js /app/server.js

# Copy Go binary from go-builder stage
COPY --from=go-builder /app/backend /app/backend

# Install any additional dependencies, including concurrently
RUN npm install concurrently

# Expose port for the frontend app
EXPOSE 5173

# Run both the Go backend and Node.js frontend together
CMD ["npx", "concurrently", "./backend", "node server.js"]