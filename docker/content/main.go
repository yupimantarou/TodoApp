package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWTシークレットキー
var jwtSecret = []byte("your_very_secure_secret_key")

// JWTトークンを検証する関数
func validateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// トークン検証後に home.html を返すハンドラー
func homeHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// "Bearer " プレフィックスを削除してトークンを取得
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]

	// トークンの検証
	claims, err := validateJWT(tokenString)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// トークンからユーザー名を取得
	username, ok := claims["username"].(string)
	if !ok || username == "" {
		http.Error(w, "Invalid token: Missing username", http.StatusUnauthorized)
		return
	}

	// `home.html` のパスを取得
	htmlFilePath := filepath.Join(".", "home.html")

	// `home.html` を提供
	tmpl, err := template.ParseFiles(htmlFilePath)
	if err != nil {
		log.Printf("Failed to parse home.html: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// ユーザー名をテンプレートに渡す
	data := struct {
		Username string
	}{
		Username: username,
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Failed to render home.html: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/login/", homeHandler)

	addr := ":8080"
	log.Printf("Server is running on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
