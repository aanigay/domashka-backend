package jwt

import (
	"domashka-backend/config"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"strconv"
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

func (u *UseCase) ValidateJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return u.cfg.Secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

func (u *UseCase) GenerateJWT(userID int64, chefID *int64, role string) (string, error) {
	tokenUUID := uuid.New()
	var claims jwt.MapClaims
	if chefID == nil {
		claims = jwt.MapClaims{
			"uuid":    tokenUUID.String(),
			"exp":     time.Now().Add(u.cfg.Exp).Unix(),
			"iat":     time.Now().Unix(),
			"role":    role,
			"user_id": strconv.FormatInt(userID, 10),
		}
	} else {
		claims = jwt.MapClaims{
			"uuid":    tokenUUID.String(),
			"exp":     time.Now().Add(u.cfg.Exp).Unix(),
			"iat":     time.Now().Unix(),
			"role":    role,
			"user_id": strconv.FormatInt(userID, 10),
			"chef_id": *chefID,
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(u.cfg.Secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
