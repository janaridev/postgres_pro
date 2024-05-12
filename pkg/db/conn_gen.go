package db

import "fmt"

func GenerateConnString(user, pass, db, host string, port int, useSSL string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", user, pass, host, port, db, useSSL)
}
