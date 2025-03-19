package jwt

import (
	"domashka-backend/config"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type UseCase struct {
	cfg *config.JWTConfig
}

func New(cfg *config.JWTConfig) *UseCase {
	return &UseCase{
		cfg: cfg,
	}
}

func (u *UseCase) ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return u.cfg.Secret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenUUID := claims["uuid"].(string)
		return tokenUUID, nil
	}

	return "", jwt.ErrSignatureInvalid
}

func (u *UseCase) GenerateJWT() (string, error) {
	tokenUUID := uuid.New()

	claims := jwt.MapClaims{
		"uuid": tokenUUID.String(),
		"exp":  time.Now().Add(u.cfg.Exp).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(u.cfg.Secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
