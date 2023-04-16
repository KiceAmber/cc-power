package model

import "powerTrading/web/groupSignature"

type User struct {
	Role            string               `json:"role"`
	UserId          string               `json:"user_id"`
	Username        string               `json:"username"`
	Password        string               `json:"password"`
	CreatedAt       string               `json:"created_at"`
	GroupPrivateKey *groupSignature.Cert `json:"group_private_key"`
}
