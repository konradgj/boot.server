# boot.chirpy
boot.dev http server project


## Goals of This Course

-  Understand what web servers are and how they power real-world web applications
-  Build a production-style HTTP server in Go, without the use of a framework
-  Use JSON, headers, and status codes to communicate with clients via a RESTful API
-  Learn what makes Go a great language for building fast web servers
-  Use type safe SQL to store and retrieve data from a Postgres database
-  Implement a secure authentication/authorization system with well-tested cryptography libraries
-  Build and understand webhooks and API keys
-  Document the REST API with markdown


## ⚙️ Environment Variables

Before running the server, ensure you have a `.env` file or environment variables set:

| Variable     | Description                                 | Required |
|--------------|---------------------------------------------|----------|
| `DB_URL`     | PostgreSQL connection string                | ✅       |
| `PLATFORM`   | Platform environment (e.g., `dev`, `prod`) | ✅       |
| `JWT_SECRET` | Secret key for signing JWT tokens           | ✅       |
| `POLKA_KEY`  | API key for Polka webhook verification      | ✅       |

Example `.env`:
```env
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=supersecretjwtkey
POLKA_KEY=somepolkasecretkey
```

## 🧩 Key Endpoints

| Method | Endpoint                  | Description |
|--------|---------------------------|-------------|
| `GET`    | `/api/healthz`              | Health check |
| `GET`    | `/api/chirps`               | List chirps (optionally filter by author) |
| `GET`    | `/api/chirps/{id}`          | Get single chirp by ID |
| `POST`   | `/api/chirps`               | Create a new chirp |
| `DELETE` | `/api/chirps/{id}`          | Delete a chirp |
| `POST`   | `/api/users`                | Create a new user |
| `PUT`    | `/api/users`                | Update user info |
| `POST`   | `/api/login`                | Log in and issue JWT |
| `POST`   | `/api/refresh`              | Refresh JWT |
| `POST`   | `/api/revoke`               | Revoke JWT |
| `GET`    | `/admin/metrics`            | View admin metrics |
| `POST`   | `/admin/reset`              | Reset metrics |
| `POST`   | `/api/polka/webhooks`       | Handle Polka webhooks