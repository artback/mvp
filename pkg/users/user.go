package users

import (
	"github.com/artback/mvp/pkg/api/middleware/security"
	"github.com/artback/mvp/pkg/change"
)

type User struct {
	Username string        `json:"username" binding:"required"`
	Password string        `json:"password" binding:"required"`
	Role     security.Role `json:"role"  binding:"required"`
	Deposit  int           `json:"deposit" binding:"required"`
}

type Response struct {
	Username string         `json:"username" binding:"required"`
	Role     security.Role  `json:"role"  binding:"required"`
	Deposit  change.Deposit `json:"deposit" binding:"required"`
}
