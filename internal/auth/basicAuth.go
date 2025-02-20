package auth

import (
	"log"
	"strings"
	"net/http"
    "github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
	"fmt"
	//"log"
)

func HashPassword(password string) (string, error) {
	pBytes, err := bcrypt.GenerateFromPassword([]byte(password), 0)	
	//returns []byte, error
	return string(pBytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//	func CompareHashAndPassword(hashedPassword, password []byte) error
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Issuer" : "chirpy",
		"IssuedAt" : time.Now().UTC(),
		"Subject" : fmt.Sprintf("%s", userID),
		"expiresAt" : fmt.Sprintf("%s", time.Now().UTC().Add(expiresIn)),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var user_id uuid.UUID
	claimStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimStruct,
		func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil },
	)
	userIDString, err := token.Claims.GetSubject()
    if err != nil {
        return user_id, err
    }
	issuer, err := token.Claims.GetIssuer()
    if err != nil {
        return user_id, err
    }
	if issuer != "chirpy" {
		return user_id, errors.New("invalid issuer")
	}
	user_id, err = uuid.Parse(userIDString)
	return user_id, err
}

func GetBearerToken(headers http.Header) (string, error){
	bearer := headers.Get("Authorization") 
	log.Println("[getbearertoken] bearer:", bearer)
	if bearer == "" {
		return "", errors.New("Invalid header")
	}
	split_bearer := strings.Split(bearer, " ")
	if len(split_bearer) != 2 {
		return "", errors.New("Invalid Auth header")
	}
	return split_bearer[1], nil
}

