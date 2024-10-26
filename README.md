# PhotoGallery Web Application

This is a **PhotoGallery Web Application** built with a **Go REST API backend** and a **Vite/Vue frontend**. The application uses **SQLite** to store image metadata and features image thumbnail generation, EXIF data retrieval, and album organization. It is designed to be served on a local or hosted server with a reverse proxy for routing.

## Features

- Go REST API backend
- Vue.js frontend built with Vite
- SQLite database for image metadata storage
- Automatic thumbnail generation
- EXIF metadata retrieval
- Basic single-user authentication with cookies
- SSL/TLS certificate generation for local development

## Project Structure

- **api/** - Go backend REST API
- **frontend/** - Vite/Vue3 frontend application
- **Caddyfile** - Configuration for the Caddy server (local reverse proxy)
- **.env** - Environment variables for application settings (admin credentials, session details, etc.)

## Prerequisites

- Go
- Node.js & npm
- Caddy (for local development reverse proxy)
- SQLite

## Getting Started

1. **Clone the repository:**

```bash
git clone https://github.com/uvmain/gallery
cd gallery
```

2. **Install dependencies:**
```bash
npm install
```

3. **Set up environment variables:**
Update the following environment variables in the /package.json (for local dev) and in the docker_compose.yml (or sibling .env file):
```plaintext
IMAGE_PATH=path/to/your/images
ADMIN_NAME=your_admin_name
ADMIN_PASSWORD=your_admin_password
```

4. **Generate SSL/TLS certificates (optional for local development):**
```bash
npm run create-cert
```

5. **Start development environment::**
```bash
npm run dev
```
This command:
- Starts the Caddy server for reverse proxy
- Runs the frontend in development mode (MHR)
- Runs the backend in development mode

Open [photogallery.localhost](https://[photogallery.localhost) in your browser

## Scripts

- `npm run backend:dev` - Runs the Go backend in development mode.
- `npm run frontend:dev` - Runs the frontend in development mode.
- `npm run dev` - Starts Caddy and concurrently runs the frontend and backend.
- `npm run build` - Builds both the backend and frontend for production.
- `npm run create-cert` - Generates and installs certificates for local SSL.

## Deployment

```
docker build -t photogallery:latest . && docker compose up -d
```

## Contributing

Feel free to submit issues or pull requests for improvements.

## License

This project is licensed under the MIT License.