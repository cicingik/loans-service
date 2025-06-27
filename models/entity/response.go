// Package entity ...
package entity

type (
	// LoginResponse ...
	LoginResponse struct {
		AccessToken string `json:"access_token"`
		AccessRoles string `json:"role"`
		UserID      uint64 `json:"user_id"`
		// RefreshToken string `json:"refresh_token"`
	}
)
