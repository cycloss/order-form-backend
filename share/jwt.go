package share

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthHeader struct {
	Token string `header:"Authorize" binding:"required"`
}

const forgotAuthHeaderMessage = "Authorize header missing from request"

func GetAuthHeader(c *gin.Context) (*AuthHeader, error) {
	var h AuthHeader
	// not mandatory bind so have control over response
	var err = c.ShouldBindHeader(&h)
	// if error then wasn't present
	if err != nil {
		return nil, NewApiErr(http.StatusBadRequest, forgotAuthHeaderMessage, "")
	}
	return &h, err

}

// Extract the jwt claims from the Authorize header.
//
// EXTREME DANGER WARNING: This method parses the jwt but doesn't validate the signature.
// It's only ever useful in cases where you know the signature is valid (because it has
// been checked previously in the stack) and you want to extract values from it.
func GetUnverifiedJwtClaimsFromHeader(c *gin.Context) (*UInvestClaims, error) {
	h, err := GetAuthHeader(c)

	if err != nil {
		return nil, NewApiErr(http.StatusBadRequest, "Authorize header missing from request", "")
	}
	split := strings.Split(h.Token, "Bearer ")
	if len(split) < 2 {
		return nil, NewApiErr(http.StatusBadRequest, "Authorize header present, but Bearer token prefix not found", "")
	}

	rawJwt := split[1]
	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(rawJwt, jwt.MapClaims{})
	if err != nil {
		msg := fmt.Sprintf("failed to parse raw jwt: %v", err)
		return nil, NewApiErr(http.StatusUnauthorized, "bad jwt", msg)
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, NewApiErr(http.StatusUnauthorized, "bad jwt", "")
	}

	claimsBundle, err := UnmarshalClaims(claims)

	if err != nil {
		log.Printf("could not unmarshal jwt. claims: %v", claims)
		return nil, err
	}

	return claimsBundle, nil

}

type UInvestClaims struct {
	ClientId int
	Iat      int64
	Exp      int64
}

func UnmarshalClaims(claims jwt.MapClaims) (*UInvestClaims, error) {
	id, ok := claims["clientId"].(int)
	if !ok || id == 0 {
		return nil, NewApiErr(http.StatusUnauthorized, "clientId not found in jwt claims", "")
	}

	iat, ok := claims["iat"].(float64)
	if !ok || iat == 0 {
		return nil, NewApiErr(http.StatusUnauthorized, "iat not found in jwt claims", "")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || exp == 0 {
		return nil, NewApiErr(http.StatusUnauthorized, "exp not found in jwt claims", "")

	}

	return &UInvestClaims{ClientId: id, Iat: int64(iat), Exp: int64(exp)}, nil
}
