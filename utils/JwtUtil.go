package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"time"
)

var secret = []byte("spacej")

func NewTokenWithCellphone(cellphone string, role string, uid uint) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"cellphone": cellphone,
		"role":      role,
		"uid":       uid,
		"exp":       time.Now().AddDate(0, 0, 7).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)

	if err != nil {
		return ""
	} else {
		return tokenString
	}

}

func CheckJwtToken(tokenString string) bool {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secret, nil
	})

	if err != nil {

		return false
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		return true
	} else {
		return false
	}

}

func GetClaim(tokenString string, clmKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secret, nil
	})

	if err != nil {

		return nil, fmt.Errorf("pa err")
	}

	if clm, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		return clm, nil
	} else {
		return nil, fmt.Errorf("token err")
	}
}
func GetToken(ctx iris.Context) string {
	token := ctx.GetHeader("Authorization")
	if token != "" && len(token) > 7 {
		token = token[7:]
	}
	return token
}
