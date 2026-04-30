package models

import (
	"time"

	admtypes "github.com/Authula/authula/plugins/admin/types"
)

// User mirrors authula's models.User with Metadata typed as an arbitrary JSON object.
type User struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Email         string         `json:"email"`
	EmailVerified bool           `json:"email_verified"`
	Image         *string        `json:"image"`
	Metadata      map[string]any `json:"metadata"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// Session mirrors authula's models.Session.
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	IPAddress *string   `json:"ip_address"`
	UserAgent *string   `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetMeResponse mirrors authula's internal GetMeResponse.
type GetMeResponse struct {
	User    *User    `json:"user"`
	Session *Session `json:"session"`
}

// SignOutRequest mirrors authula's internal SignOutRequest.
type SignOutRequest struct {
	SessionID  *string `json:"session_id,omitempty"`
	SignOutAll bool    `json:"sign_out_all,omitempty"`
}

// SignOutResponse mirrors authula's internal SignOutResponse.
type SignOutResponse struct {
	Message string `json:"message"`
}

// SignUpRequest mirrors authula's email-password SignUpRequest with object Metadata.
type SignUpRequest struct {
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	Password    string         `json:"password"`
	Image       *string        `json:"image,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	CallbackURL *string        `json:"callback_url,omitempty"`
}

// SignUpResponse mirrors authula's email-password SignUpResponse using local User/Session.
type SignUpResponse struct {
	User    *User    `json:"user"`
	Session *Session `json:"session"`
}

// SignInRequest mirrors authula's email-password SignInRequest.
type SignInRequest struct {
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	CallbackURL *string `json:"callback_url,omitempty"`
}

// SignInResponse mirrors authula's email-password SignInResponse using local User/Session.
type SignInResponse struct {
	User    *User    `json:"user"`
	Session *Session `json:"session"`
}

// CreateUserRequest mirrors authula's admin CreateUserRequest with object Metadata.
type CreateUserRequest struct {
	Name          string         `json:"name"`
	Email         string         `json:"email"`
	EmailVerified *bool          `json:"email_verified,omitempty"`
	Image         *string        `json:"image,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

// UpdateUserRequest mirrors authula's admin UpdateUserRequest with object Metadata.
type UpdateUserRequest struct {
	Name          *string        `json:"name,omitempty"`
	Email         *string        `json:"email,omitempty"`
	EmailVerified *bool          `json:"email_verified,omitempty"`
	Image         *string        `json:"image,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

// CreateUserResponse mirrors authula's admin CreateUserResponse using local User.
type CreateUserResponse struct {
	User *User `json:"user"`
}

// GetUserByIDResponse mirrors authula's admin GetUserByIDResponse using local User.
type GetUserByIDResponse struct {
	User *User `json:"user"`
}

// UpdateUserResponse mirrors authula's admin UpdateUserResponse using local User.
type UpdateUserResponse struct {
	User *User `json:"user"`
}

// UsersPage mirrors authula's admin UsersPage using local User.
type UsersPage struct {
	Users      []User  `json:"users"`
	NextCursor *string `json:"next_cursor,omitempty"`
}

// AdminUserSession mirrors authula's admin AdminUserSession using local Session.
type AdminUserSession struct {
	Session Session                     `json:"session"`
	State   *admtypes.AdminSessionState `json:"state,omitempty"`
}
