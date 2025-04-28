package auth

type Request struct {
	Phone     string `json:"phone" binding:"required"`
	AuthViaTg bool   `json:"auth_via_tg"`
}
