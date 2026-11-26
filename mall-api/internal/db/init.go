package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

//go:embed schema.sql
var schemaSQL string

func InitDB(dsn string, dbName string) (*sql.DB, error) {
	rootDSN := removeDatabaseFromDSN(dsn, dbName)

	if strings.Contains(rootDSN, "?") {
		rootDSN += "&multiStatements=true"
	} else {
		rootDSN += "?multiStatements=true"
	}

	rootDB, err := sql.Open("mysql", rootDSN)
	if err != nil {
		return nil, fmt.Errorf("open root db: %w", err)
	}
	defer rootDB.Close()

	if err := rootDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping root db: %w", err)
	}

	sqlToRun := "USE `" + dbName + "`;\n" + schemaSQL
	if _, err := rootDB.Exec(sqlToRun); err != nil {
		return nil, fmt.Errorf("exec schema: %w", err)
	}

	appDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open app db: %w", err)
	}
	if err := appDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping app db: %w", err)
	}

	if err := seedAdmin(appDB); err != nil {
		return nil, fmt.Errorf("seed admin: %w", err)
	}

	return appDB, nil
}

func removeDatabaseFromDSN(dsn, dbName string) string {
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' {
			return dsn[:i+1] + dsn[i+1+len(dbName):]
		}
	}
	return dsn
}

func seedAdmin(db *sql.DB) error {
	ctx := context.Background()

	email := os.Getenv("ADMIN_EMAIL")
	if email == "" {
		email = "root.admin@mall-system.local"
	}

	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		if os.Getenv("APP_ENV") == "production" {
			return fmt.Errorf("ADMIN_PASSWORD must be set in production")
		}
		password = "X7$kL9#mQ2!vN4@p"
		log.Println("WARNING: using default admin password (dev only)")
	}

	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM Users WHERE Email = ?)", email).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check admin exists: %w", err)
	}
	if exists {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	_, err = db.ExecContext(ctx,
		"INSERT INTO Users (Email, PasswordHash, FullName, Role) VALUES (?, ?, ?, ?)",
		email, string(hash), "System Administrator", "admin",
	)
	if err != nil {
		return fmt.Errorf("insert admin: %w", err)
	}

	log.Printf("=====================================================")
	log.Printf("SUCCESS: ADMIN USER CREATED")
	log.Printf("Email:    %s", email)
	log.Printf("=====================================================")

	return nil
}
