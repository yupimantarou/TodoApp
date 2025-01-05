module jwt

go 1.23

require (
	github.com/go-sql-driver/mysql v1.6.0 // MySQLドライバ
	github.com/golang-jwt/jwt/v5 v5.0.0 // JWTの生成と検証 (利用可能な最新バージョン)
	golang.org/x/crypto v0.13.0 // パスワードハッシュなど暗号化関連
)
