package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// パスワードをハッシュ化
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// JSONエラーレスポンスを送信する関数
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// ユーザー登録ハンドラー
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.Username == "" || requestBody.Password == "" {
		sendErrorResponse(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// パスワードをハッシュ化
	hashedPassword, err := hashPassword(requestBody.Password)
	if err != nil {
		sendErrorResponse(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// ユーザーをデータベースに挿入
	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", requestBody.Username, hashedPassword)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			sendErrorResponse(w, "Username already exists", http.StatusConflict)
			return
		}
		sendErrorResponse(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// /login ハンドラー
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// ユーザー情報をデータベースから取得
	var hashedPassword string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username = ?", requestBody.Username).Scan(&hashedPassword)
	if err != nil {
		sendErrorResponse(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// パスワードを検証
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(requestBody.Password)); err != nil {
		sendErrorResponse(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// JWTトークンを生成
	token, err := generateJWT(requestBody.Username)
	if err != nil {
		sendErrorResponse(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// トークンをJSONで返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// /protected ハンドラー
func protectedHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		sendErrorResponse(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Bearer トークンの取り出し
	tokenString := authHeader[len("Bearer "):]
	claims, err := validateJWT(tokenString)
	if err != nil {
		sendErrorResponse(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// 保護されたデータを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Welcome to the protected route!",
		"claims":  claims,
	})
}
