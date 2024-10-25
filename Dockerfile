FROM node:20-alpine AS frontend-build

WORKDIR /frontend

COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

FROM golang:1.23 AS backend-build

WORKDIR /app

COPY api/ .
RUN go build -o server .

FROM alpine:latest

WORKDIR /app

COPY --from=backend-build /app/server .
COPY --from=frontend-build /frontend/dist ./dist

EXPOSE 8080

CMD ["./server"]
