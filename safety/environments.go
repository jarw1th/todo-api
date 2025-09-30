package safety

import (
	"fmt"
	"os"
)

func GetJwt() []byte {
	secret := os.Getenv("jwt_secret_key")
	if secret == "" {
		fmt.Println("WARNING: jwt_secret_key is not set")
		secret = "none"
	}
	return []byte(secret)
}

func GetDBConnStr() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if password == "" {
		password = ""
	}
	if dbname == "" {
		dbname = "todo"
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
