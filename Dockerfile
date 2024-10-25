FROM node:20-alpine AS frontend-build

WORKDIR /frontend

COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

FROM golang:1.23 AS backend-build

WORKDIR /app

COPY api/ .
RUN CGO_ENABLED=0 go build -o server .

FROM gcr.io/distroless/static-debian12

COPY --from=backend-build /app/server .
COPY --from=frontend-build /frontend/dist ./dist

EXPOSE 8080

CMD ["./server"]