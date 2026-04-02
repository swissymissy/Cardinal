package handler

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
}

type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type LoginUser struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

type UserProfile struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	CreatedAt      time.Time `json:"created_at"`
	FollowerCount  int64     `json:"followers_count"`
	FollowingCount int64     `json:"followings_count"`
}

type Chirp struct {
	Body string `json:"body"`
}

type CreatedChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
}

type ResponseAccessToken struct {
	Token string `json:"token"`
}

type NewFollow struct {
	FolloweeID uuid.UUID `json:"followee_id"`
}

type Follower struct {
	FollowerID uuid.UUID `json:"follower_id"`
	FolloweeID uuid.UUID `json:"followee_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FollowList struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"followed_at"`
}

type Notification struct {
	ID        uuid.UUID `json:"notif_id"`
	CreatedAt time.Time `json:"created_at"`
	Body      string    `json:"body"`
	Receiver  uuid.UUID `json:"receiver"`
	Username  string    `json:"author"`
	ChirpID   uuid.UUID `json:"chirp_id"`
	IsRead    bool      `json:"is_read"`
}
