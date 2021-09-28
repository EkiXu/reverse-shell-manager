package typing

import (
	"github.com/dgrijalva/jwt-go"
)

// Custom claims structure
type CustomClaims struct {
	Name string
	jwt.StandardClaims
}
