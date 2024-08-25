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

func (s *jwtService) UserFromClaims(ctx interfaces.IRequestContext) (*interfaces.User, error) {
	userLocal := ctx.Locals("user")
	if userLocal == nil {
		return nil, nil
	}
	userToken := ctx.Locals("user").(*jwt.Token)
	if userToken == nil {
		return nil, nil
	}
	claims := userToken.Claims.(jwt.MapClaims)
	if claims == nil {
		return nil, nil
	}
	retUser := &interfaces.User{}
	if username, ok := claims["username"].(string); ok {
		retUser.Username = username
	}
	if email, ok := claims["email"].(string); ok {
		retUser.Email = email
	}
	if role, ok := claims["role"].(string); ok {
		retUser.Role = role
	}
	if id, ok := claims["sub"].(uint); ok {
		retUser.ID = id
	}
	return retUser, nil
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
