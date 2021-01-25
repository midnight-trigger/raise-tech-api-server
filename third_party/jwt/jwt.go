package jwt

import (
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/midnight-trigger/todo/configs"
)

// todo complete Claims
type Claims struct {
	Sub    string `json:"sub"`
	UserId string `json:"user_id"`
}

type Config struct {
	SecretToken   []byte
	SigningMethod *jwt.SigningMethodHMAC
}

type Meta struct {
	Code         int    `json:"code"`
	ErrorType    string `json:"error_type"`
	ErrorMessage string `json:"error_message"`
}

// get middleware JWTConfig
func GetMiddlewareJWTConfig() middleware.JWTConfig {
	jwtConfig := configs.GetJwtConfig()
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(jwtConfig.JWTSecretToken))
	return middleware.JWTConfig{
		SigningKey:    publicKey,
		Skipper:       middleware.DefaultSkipper,
		SigningMethod: jwt.SigningMethodRS256.Name,
		ContextKey:    "jwt",
		TokenLookup:   "header:" + echo.HeaderAuthorization,
		AuthScheme:    "Bearer",
		Claims:        jwt.MapClaims{},
	}
}

// get jwt token from request
func GetJWTClaims(ctx echo.Context) (*Claims, *Meta) {

	claims := new(Claims)
	t := ctx.Request().Header["Authorization"]

	tokenString := ""
	if len(t) > 0 {
		tokenString = t[0]
	}

	if len(tokenString) > 7 && strings.ToUpper(tokenString[0:6]) == "BEARER" {
		tokenString = tokenString[7:]
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return jwtError()
	}

	if c, ok := token.Claims.(jwt.MapClaims); ok {
		claims.UserId = c["cognito:username"].(string)
	} else {
		return jwtError()
	}

	return claims, nil
}

func jwtError() (*Claims, *Meta) {
	meta := &Meta{
		Code:         http.StatusBadRequest,
		ErrorType:    http.StatusText(http.StatusBadRequest),
		ErrorMessage: "jwt解析エラー",
	}
	return nil, meta
}
