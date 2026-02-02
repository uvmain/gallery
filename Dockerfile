FROM node:24-alpine AS frontend-build

WORKDIR /frontend

COPY ./frontend .

RUN npm install

RUN npm run build

FROM golang:1.25.6 AS backend-build

WORKDIR /app

COPY . .

COPY --from=frontend-build /frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 go build -o gallery .

FROM gcr.io/distroless/static-debian12

COPY --from=backend-build /app/gallery .

EXPOSE 8080

CMD ["./gallery"]