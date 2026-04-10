# Cardinal

Cardinal is a simple social media platform that lets users share their thoughts, react to posts, comment, follow others, and stay connected.

**Motivation:** Inspired by the Chirpy project from Boot.dev, Cardinal expands on that foundation with more features to get closer to a real social media platform like Twitter.

---

## Prerequisites

- [Go](https://go.dev/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Docker](https://www.docker.com/get-started/) (for RabbitMQ)
- [Goose](https://github.com/pressly/goose) — database migration tool

---

## Setup

### Option A — Run with Docker (recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/swissymissy/Cardinal.git
   cd Cardinal
   ```

2. **Create a Docker environment file**
   ```bash
   cp .env.example .env.docker
   ```
   Fill in the values. For Docker, the database host must use the service name, not localhost:
   ```
   DB_URL=postgres://postgres:postgres@postgres:5432/cardinal?sslmode=disable
   RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
   DB_HOST=postgres
   DB_PORT=5432
   ```

3. **Generate a JWT secret**
   ```bash
   openssl rand -base64 64
   ```

4. **Start all services**
   ```bash
   docker compose up --build
   ```
   This spins up PostgreSQL, RabbitMQ, the web server, and the workers in one command. Migrations run automatically on startup.

5. **Visit the app**
   ```
   http://localhost:8080
   ```

To stop:
```bash
docker compose down
```

To rebuild after code changes:
```bash
docker compose up --build
```

---

### Option B — Run locally (without Docker)

1. **Clone the repository**
   ```bash
   git clone https://github.com/swissymissy/Cardinal.git
   cd Cardinal
   ```

2. **Download dependencies**
   ```bash
   go mod tidy
   ```

3. **Create environment file**
   ```bash
   cp .env.example .env
   ```
   Fill in the values. For local development, use `localhost` for database and RabbitMQ:
   ```
   DB_URL=postgres://user:password@localhost:5432/cardinal?sslmode=disable
   RABBITMQ_URL=amqp://guest:guest@localhost:5672/
   ```

4. **Generate a JWT secret**
   ```bash
   openssl rand -base64 64
   ```

5. **Start RabbitMQ**
   ```bash
   docker run -d --name rabbitmq -p 5672:5672 rabbitmq:3
   ```

6. **Run database migrations**
   ```bash
   goose -dir sql/schema postgres "$DB_URL" up
   ```

7. **Run the server**
   ```bash
   go run .
   ```

8. **Run the workers** (in a separate terminal)
   ```bash
   go run ./cmd/workers
   ```

---

## Features

| Feature | Description |
|---|---|
| Users | Register and log in to your account |
| Security | Passwords hashed with Argon2id, sessions managed with JWT access + refresh tokens |
| Chirps | Post, view, and delete short messages (up to 500 characters) |
| Comments | Comment on any chirp |
| Reactions | React to chirps with emoji |
| Follow | Follow other users to build your feed |
| Notifications | In-app and email notifications for new chirps from followed users |
| Email Verification | Verify your email address on sign-up |

---

## API Reference

### Auth

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/newuser` | Register a new user |
| `POST` | `/api/userlogin` | Log in and receive tokens |
| `POST` | `/api/refresh` | Refresh an access token |
| `POST` | `/api/revoke` | Revoke a refresh token (logout) |

### Users

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/users/{identifier}` | Get a user by username or ID |
| `GET` | `/api/users/{userID}/followers` | List a user's followers |
| `GET` | `/api/users/{userID}/followings` | List who a user follows |

### Follow

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/newfollow` | Follow a user |
| `DELETE` | `/api/unfollow` | Unfollow a user |

### Chirps

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/newchirp` | Create a new chirp |
| `GET` | `/api/getallchirps` | Get all chirps |
| `POST` | `/api/feed` | Get chirps from followed users |
| `GET` | `/api/chirps/{chirpID}` | Get a single chirp |
| `DELETE` | `/api/chirps/{chirpID}` | Delete a chirp (owner only) |

### Reactions

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/chirps/{chirpID}/react` | Add a reaction to a chirp |
| `DELETE` | `/api/chirps/{chirpID}/react` | Remove your reaction from a chirp |
| `GET` | `/api/chirps/{chirpID}/react` | Get all reactions for a chirp |

### Comments

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/chirps/{chirpID}/comments` | Add a comment to a chirp |
| `GET` | `/api/chirps/{chirpID}/comments` | Get all comments on a chirp |
| `PUT` | `/api/comments/{commentID}` | Edit a comment (owner only) |
| `DELETE` | `/api/comments/{commentID}` | Delete a comment (owner only) |

### Notifications

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/notifications` | Get your notifications |
| `PUT` | `/api/notifications` | Mark all notifications as read |
| `PUT` | `/api/notifications/{notifID}` | Mark a single notification as read |

### Email Verification

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/verify/request` | Request a verification email |
| `GET` | `/api/verify` | Verify email via token link |

### Admin

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/admin/reset` | Reset all users (dev only) |


---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go |
| Database | PostgreSQL |
| Queries | sqlc |
| Migrations | Goose |
| Messaging | RabbitMQ (AMQP) |
| Auth | JWT + Argon2id |
| Frontend | Vanilla JS / HTML / CSS |

---

## Future Features

- **Profile picture** — upload and display a profile avatar
- **Edit profile** — change password
- **Password strength** — enforce strong passwords on registration
- **Edit chirp** — allow users to edit their posted chirps
- **Image attachments** — attach images to chirps
