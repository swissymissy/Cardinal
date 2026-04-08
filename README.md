# Cardinal

Cardinal is a simple social media platform that lets you express your thoughts, follow people you like, react, comment, and stay in the loop.

---

## Features

### Users
- Create a new account and log in
- Email verification on sign-up
- Search for users by username or user ID

### Posts (Chirps)
- Create and delete posts
- React to posts with emojis (❤️ 😂 😮 😢 👍)
- Comment on posts
- Pagination — chirps are loaded in batches, with a "load more" option

### Social
- Follow and unfollow users
- In-app notifications when someone you follow posts a new chirp
- Email notifications delivered via SMTP

### Security
- Passwords hashed with **Argon2id**
- Authentication via **JWT** — short-lived access tokens and long-lived refresh tokens
- Token revocation on logout

---

## Prerequisites

Make sure the following are installed before setting up Cardinal:

- **Go** (1.21+) — [https://go.dev/doc/install](https://go.dev/doc/install)
- **PostgreSQL** — [https://www.postgresql.org/download/](https://www.postgresql.org/download/)
- **Docker** — [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/) (used to run RabbitMQ locally)
- **Goose** — database migration tool

  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```

---

## Setup (Draft version before Docker-compose)

### 1. Clone the repository

```bash
git clone https://github.com/swissymissy/Cardinal
cd Cardinal
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Start RabbitMQ with Docker

```bash
docker run -d --name rabbitmq -p 5672:5672 rabbitmq:3
```

To verify it's running:

```bash
docker ps
```

### 4. Create the PostgreSQL database

```bash
createdb cardinal
```

If `createdb` is not available, you can use `psql` instead:

```sql
psql -U postgres -c "CREATE DATABASE cardinal;"
```

### 5. Configure environment variables

Copy the example file and fill in your values:

```bash
cp .env.example .env
```

| Variable | Description |
|---|---|
| `PORT` | Port the server listens on (e.g. `8080`) |
| `PLATFORM` | Set to `dev` for local development |
| `DB_URL` | PostgreSQL connection string (see format below) |
| `RABBITMQ_URL` | RabbitMQ connection string (default: `amqp://guest:guest@localhost:5672/`) |
| `JWT_SECRET` | Secret key for signing JWT tokens (see generation below) |
| `BASE_URL` | Base URL of the server — used in email verification links (e.g. `http://localhost:8080`) |
| `SMTP_HOST` | SMTP server host (e.g. `smtp.gmail.com`) |
| `SMTP_PORT` | SMTP server port (e.g. `587`) |
| `SMTP_USERNAME` | Your sender email address |
| `SMTP_PASSWORD` | SMTP app password (see Gmail note below) |

**DB_URL format:**
```
postgres://username:password@localhost:5432/cardinal?sslmode=disable
```

**Generate a JWT secret:**
```bash
openssl rand -base64 64
```
Paste the output as the value of `JWT_SECRET` in your `.env`.

**Gmail SMTP setup:**
Gmail requires an App Password rather than your account password. To generate one:
1. Enable 2-Step Verification on your Google account
2. Go to [myaccount.google.com/apppasswords](https://myaccount.google.com/apppasswords)
3. Create a new app password and paste it as `SMTP_PASSWORD`

### 6. Run database migrations

```bash
goose -dir sql/schema postgres "your-db-url-here" up
```

Or if you have `DB_URL` set in your `.env`:

```bash
export $(cat .env | xargs) && goose -dir sql/schema postgres "$DB_URL" up
```

### 7. Start the server

```bash
go run .
```

The app will be available at [http://localhost:8080](http://localhost:8080).

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go |
| Database | PostgreSQL |
| Query generation | sqlc |
| Migrations | Goose |
| Auth | JWT (golang-jwt) |
| Password hashing | Argon2id (alexedwards/argon2id) |
| Message queue | RabbitMQ (amqp091-go) |
| Email | go-mail |
| Frontend | Vanilla HTML, CSS, JavaScript |
