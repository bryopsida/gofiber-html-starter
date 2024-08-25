package jwt

import (
	"os"
	"time"

	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService(settings interfaces.ISettingsService) interfaces.IJWTService {
	key, err := settings.GetString("jwt_signing_key")
	if err != nil {
		panic(err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return &jwtService{
		secretKey: key,
		issuer:    hostname,
	}
}

func (s *jwtService) Generate(user *interfaces.User) (string, error) {
	claims := jwt.MapClaims{
		"iss":      s.issuer,
		"sub":      user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(s.secretKey), nil
}
func (s *jwtService) Validate(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, s.keyFunc, jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return token, nil
}
