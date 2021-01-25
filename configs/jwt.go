package configs

import "os"

// JwtConfig is base configuration of jwt
type JwtConfig struct {
	JWTSecretToken string
	JWTIssuer      string
}

func GetJwtConfig() *JwtConfig {
	jwtConfig := new(JwtConfig)
	if v := os.Getenv("JWT_SECRET_TOKEN"); v != "" {
		jwtConfig.JWTSecretToken = v
	}

	if v := os.Getenv("JWT_ISSUER"); v != "" {
		jwtConfig.JWTIssuer = v
	}

	return jwtConfig
}
