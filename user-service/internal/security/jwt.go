package security

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var (
	privateKey   = []byte(os.Getenv("TOKEN_PRIVATE_KEY"))
	BEARER_TOKEN = "Bearer"
)

func InitJWT(path string) {

	if path != "" {
		fmt.Println("Load file in: ", path)
		err := godotenv.Load(path)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		err := godotenv.Load("./internal/security/.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func GenerateJWTToken(userId string) (string, error) {

	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	}).SignedString(privateKey)

	return token, nil
}

func IsTimeExpired(expire int64) error {

	if expire != 0 {
		exp := time.Unix(int64(expire), 0)
		if exp.After(time.Now()) {
			return nil
		}
	}
	return errors.New("token is expired")
}

func ParseClaims(tokenString string) (jwt.MapClaims, error) {

	token, err := ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), err
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(tokenString *jwt.Token) (interface{}, error) {
		if _, ok := tokenString.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tokenString.Header["alg"])
		}
		return privateKey, nil
	})
}
