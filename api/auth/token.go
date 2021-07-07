package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(userId uint32) (string, error) {
	fmt.Println("userId", userId)
	claims := jwt.MapClaims{}
	fmt.Println("claims", claims)
	claims["authorized"] = true
	claims["id"] = userId
	// claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("cliams", claims)
	tk, tk2 := token.SignedString([]byte(os.Getenv("API_SECRET")))
	fmt.Println("token return ", tk)
	fmt.Println("token return tk2", tk2)
	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func TokenValid(req *http.Request) error {
	tokenString := ExtractToken(req)

	token, err := jwt.Parse(tokenString, func(tkn *jwt.Token) (interface{}, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexcepted Error signing method: %v", tkn.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok && token.Valid {
		Pretty(claims)
	}
	return nil

}

func Pretty(data interface{}) {
	b, err := json.MarshalIndent(data, "", "")
	if err != nil {
		log.Println(err, "errr on marshal indent")
		return
	}

	fmt.Println(string(b))
}

func ExtractToken(req *http.Request) string {
	keys := req.URL.Query()
	fmt.Println(keys, "keys")
	token := keys.Get("token")
	if token != "" {
		return token
	}

	bearerToken := req.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenID(req *http.Request) (uint32, error) {
	tokenString := ExtractToken(req)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexcepted signing methos %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["id"]), 10, 32)

		if err != nil {
			return 0, err
		}
		return uint32(uid), nil
	}
	return 0, nil
}
