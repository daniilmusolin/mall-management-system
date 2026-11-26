//go:build integration

package mysql_test

import (
	"fmt"
	"os"
	"testing"

	"mall-api/internal/db"
)

func TestMain(m *testing.M) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASS", "root"),
		getEnv("DB_HOST", "127.0.0.1:3306"),
		"shopping_center_test",
	)

	conn, err := db.InitDB(dsn, "shopping_center_test")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	os.Exit(m.Run())
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
