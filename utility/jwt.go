package utility

import (
	"errors"
	"fmt"
	"journal_api/model"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// convert private key string to slice of bytes
var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

// GenerateJWT: generate JWT containing user id, issue date, expiry date as claims
func GenerateJWT(user model.User) (string, error) {

	// convert TTL from string to integer
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))

	// create new token with HMAC signing method and JWT claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// unique user id
		"id": user.ID,

		// token issued time
		"iat": time.Now().Unix(),

		// token expiry date
		"eat": time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})

	// sign token with shared key
	return token.SignedString(privateKey)
}

//--------------------Helper Functions: Validate and Extract JWT--------------------//

// CurrentUser: retrieve authorized user record instance as a struct with all associated entries.
func CurrentUser(ctx *gin.Context, db *gorm.DB) (model.User, error) {

	// validate JWT of accessing user
	token, err := ValidateJWT(ctx)

	if err != nil {
		return model.User{}, err
	}

	// retrieve token claims initially setup when the JWT was created
	claims, _ := token.Claims.(jwt.MapClaims)

	// type assertion as float64
	userID := uint(claims["id"].(float64))

	// find user by id from the JWT map of claims
	user, err := model.FindUserByID(userID, db)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// ValidateJWT: ensure valid token in request Authorization Header.
func ValidateJWT(ctx *gin.Context) (*jwt.Token, error) {

	// ensure valid token in request header
	token, err := getToken(ctx)

	// log.Println("JWT token:", token)

	if err != nil {
		return &jwt.Token{}, err
	}

	// type assert claims type
	_, ok := token.Claims.(jwt.MapClaims)

	// if correct claim type and valid token return token
	if ok && token.Valid {
		return token, nil
	}

	return &jwt.Token{}, errors.New("invalid token provided")

}

// getToken: parse retrieved JWT string with private key.
func getToken(ctx *gin.Context) (*jwt.Token, error) {

	// retrieve JWT from request
	tokenString, err := getTokenFromRequest(ctx)

	if err != nil {
		return &jwt.Token{}, err
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// verify correct signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// return privateKey if signing method matches expected
		return privateKey, nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)

	return token, err
}

// getTokenFromRequest: retrieve bearer token from Authorization Header of request.
func getTokenFromRequest(ctx *gin.Context) (string, error) {

	// retrieve Bearer token from request Authorization Header
	bearerToken := ctx.Request.Header.Get("Authorization")

	// split Bearer string from JWT string
	splitToken := strings.Split(bearerToken, " ")

	// log.Println("Token string:", splitToken)

	// if Bearer string and JWT string are the only elements return JWT string
	if len(splitToken) == 2 {
		return splitToken[1], nil
	}

	return "", errors.New("there was an issue retrieving JWT from request Authorization Header.")
}

// JSON Web Tokens | Token Based Authentication

//   - perform user authentication over the internet
//   - client sends login details to HTTP server
//   - server generates JWT
//   - JWT's are created with a private key located on the HTTP server (typically an environment variable)
//   - JWT's are sent to the client to be held in local storage
//   - removes the need to save cookie ID's in your database
//   - JWT's are added to the Authorization Header of future requests with a Bearer prefix
//   - server validates signature of request without requiring database lookups
//   - tokens are managed on the client

// JWT Signing Algorithm

//   - signing methods are algorithms
//   - there are many signing algorithms, HMAC being the most common

// HMAC - Hashed Based Authentication Codes

//   - a way of signing a JWT
//   - signs messages by using a shared key (private key)

// JWT Claims

//   - metadata that asserts information about the token
//   - key / value pair part of a tokens attributes (typically a hashmap defined by the JWT package you use)

// JSON Web Token Issues

//   - can be hijacked by attackers
//   - difficult to invalidate

// Access Token Time-To-Live: TTL

//   - the duration of a tokens lifetime (validation period)
//   - assigned to a token during creation in jwt claims
//   - TTL expiration invalidates a token

// Sessions | Session Based Authentication

//   - authentication state is handled on the server
//   - authentication information must be saved and retreived from the database
