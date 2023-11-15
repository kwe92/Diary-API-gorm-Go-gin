package helper

import (
	"diary_api/model"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// convert private key string to slice of bytes
var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

// GenerateJWT: generate JWT containing user id, issue date, expiry date as claims
func GenerateJWT(user model.User) (string, error) {

	// convert TTL from string to integer
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))

	// create new token with HMAC signing method and JWT claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,

		// token issued time
		"iat": time.Now().Unix(),

		// token expiry date
		"eat": time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})

	// sign token with shared key
	return token.SignedString(privateKey)
}

//Helper functions to validate an extract JWT's--------------------------------------------------
// TODO: Add Comments
// TODO: should the CurrentUser function be moved?

// CurrentUser: retrieve authorized user and all associated entries
func CurrentUser(context *gin.Context) (model.User, error) {

	err := ValidateJWT(context)

	if err != nil {
		return model.User{}, err
	}

	// TODO: why not get token from ValidateJWT?
	token, _ := getToken(context)

	// retrieve the token claims initially setup when the JWT was created
	claims, _ := token.Claims.(jwt.MapClaims)

	// TODO: why float to unsigned int?
	userID := uint(claims["id"].(float64))

	user, err := model.FindUserByID(userID)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// ValidateJWT: ensure valid token in request Authorization Header
// TODO: try returning (*jwt.Token, error)
func ValidateJWT(context *gin.Context) error {
	// ensure valid token in request header
	token, err := getToken(context)
	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return nil
	}

	return errors.New("invalid token provided")

}

// getToken: parse retrieved JWT string with private key.
func getToken(context *gin.Context) (*jwt.Token, error) {

	// retrieve JWT token as string
	tokenString := getTokenFromRequest(context)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// verify the correct signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// return privateKey if the signing method matches the expected
		return privateKey, nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)

	return token, err
}

// getTokenFromRequest: retrieve bearer token from Authentication Header of request
// TODO: add an error if there was no token retrieved
func getTokenFromRequest(context *gin.Context) string {

	// JWT token with bearer prefix
	bearerToken := context.Request.Header.Get("Authorization")

	// split bearer prefix and JWT
	splitToken := strings.Split(bearerToken, " ")

	// if only the bearer prefix and JWT string are elements then return JWT string
	if len(splitToken) == 2 {
		return splitToken[1]
	}

	return ""
}

//?----------------------------------------------------------------------------------------------------

// JSON Web Tokens | Token Based Authentication

//   - a way of performing user authentication over the internet
//   - client sends login details to the HTTP server
//   - server generates JWT
//   - JWT are created with a private key on the server (typically an environment variable)
//   - the JWT is then sent to the client to be held in local storage
//   - removes the need to save cookie ID's in your database
//   - the JWT is added to the Authorization Header of future requests with a Bearer prefix in the value
//   - the server then validates the signature of the request without the need to do database lookups
//   - tokens are managed on the client

// JWT Signing Algorithm

//   - the signing method is an algorithm
//   - there are many signing algorithms, HMAC being the most common

// HMAC - Hashed Based Authentication Codes

//   - a way of signing a JWT
//   - signs messages by using a shared key

// JWT Claims

//   - metadata that asserts information about the token

// JSON Web Token Issues

//   - can be hijacked by attackers
//   - difficult to invalidate

// Access Token Time-To-Live: TTL

//   - the duration of a tokens lifetime (validation period)
//   - assigned to a token during creation
//   - once a tokens TTL has expired that token is no longer valid within the application it is trying to access

// Sessions | Session Based Authentication

//   - authentication state is handled on the server
