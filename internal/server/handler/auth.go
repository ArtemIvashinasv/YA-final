package handler

import (
	"YandexPracticum-go-final-TODO/internal/config"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Password struct {
	Password string `json:"password"`
}

const secretKey = "my_secret_key"

var getPass = config.Password()

func SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received reqest POST SignIn")
		var pass Password

		err := json.NewDecoder(r.Body).Decode(&pass)
		if err != nil {
			log.Println("JSON deserialization error", err)
			http.Error(w, `{"error":"JSON deserialization error"}`, http.StatusBadRequest)
			return
		}

		if pass.Password == "" {
			json.NewEncoder(w).Encode(map[string]string{"error": "need to enter a password"})
			return
		}

		storedPasword := getPass

		if pass.Password == storedPasword {
			token, err := GetToken(pass.Password)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": string(err.Error())})
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			json.NewEncoder(w).Encode(map[string]string{"token": token})
		}

		if pass.Password != storedPasword {
			json.NewEncoder(w).Encode(map[string]string{"error": "incorrect password"})
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass := getPass

		fmt.Println("Password:", pass)

		if len(pass) > 0 {
			var cookieToken string

			cookie, err := r.Cookie("token")
			if err == nil {
				cookieToken = cookie.Value
			}

			fmt.Println("JWT:", cookieToken)

			var isValid bool

			isValid, err = ValidateToken(cookieToken)
			if err != nil {
				json.NewEncoder(w).Encode(map[string]string{"error": "can't validate token"})
				http.Error(w, "invalid token", http.StatusInternalServerError)
				return
			}

			fmt.Println("validation successful")

			if !isValid {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "authentification required"})
				return
			}
		}

		next(w, r)
	}
}

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.StandardClaims
}

func GetToken(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))
	hashString := hex.EncodeToString(hash[:])

	claims := &Claims{
		PasswordHash: hashString,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("failed to sign jwt:", err)
		return "", err
	}

	log.Println("Result token:", signedToken)

	return signedToken, nil
}

func ValidateToken(tokenString string) (bool, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unknown signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return false, err
	}

	if ok := token.Valid; ok {
		return true, nil
	}

	return false, nil
}
