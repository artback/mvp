package users

import "github.com/artback/mvp/pkg/change"

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     Role   `json:"role"  binding:"required"`
	Deposit  int    `json:"deposit" binding:"required"`
}

type Response struct {
	Username string         `json:"username" binding:"required"`
	Role     Role           `json:"role"  binding:"required"`
	Deposit  change.Deposit `json:"deposit" binding:"required"`
}

func (u User) IsRole(roles ...Role) bool {
	var isRole = false
	for _, r := range roles {
		if r == u.Role {
			isRole = true
			break
		}
	}
	return isRole
}
