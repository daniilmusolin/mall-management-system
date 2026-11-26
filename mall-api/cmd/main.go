package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"mall-api/internal/adapter/storage/mysql"
	"mall-api/internal/db"
	"mall-api/internal/handler"
	"mall-api/internal/service"
)

func main() {
	_ = godotenv.Load()

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	dbConn, err := db.InitDB(dsn, os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal("Database initialization error:", err)
	}
	defer dbConn.Close()

	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(25)
	dbConn.SetConnMaxLifetime(5 * time.Minute)

	userRepo := mysql.NewUserRepo(dbConn)
	premiseRepo := mysql.NewPremiseRepo(dbConn)
	tenantRepo := mysql.NewTenantRepo(dbConn)
	contractRepo := mysql.NewContractRepo(dbConn)
	paymentRepo := mysql.NewPaymentRepo(dbConn)
	requestRepo := mysql.NewRequestRepo(dbConn)

	txManager := mysql.NewTransactionManager(dbConn)

	authSvc := service.NewAuthService(userRepo, os.Getenv("JWT_SECRET"))
	premiseSvc := service.NewPremiseService(premiseRepo)
	tenantSvc := service.NewTenantService(tenantRepo)
	contractSvc := service.NewContractService(contractRepo, premiseRepo, tenantRepo, txManager)
	paymentSvc := service.NewPaymentService(paymentRepo, contractRepo)
	requestSvc := service.NewRequestService(requestRepo, premiseRepo)

	authHandler := handler.NewAuthHandler(authSvc)
	premiseHandler := handler.NewPremiseHandler(premiseSvc)
	tenantHandler := handler.NewTenantHandler(tenantSvc)
	contractHandler := handler.NewContractHandler(contractSvc)
	paymentHandler := handler.NewPaymentHandler(paymentSvc)
	requestHandler := handler.NewRequestHandler(requestSvc)

	authMid := handler.AuthMiddleware(authSvc)
	mux := http.NewServeMux()

	authHandler.RegisterRoutes(mux, authMid)
	premiseHandler.RegisterRoutes(mux, authMid)
	tenantHandler.RegisterRoutes(mux, authMid)
	contractHandler.RegisterRoutes(mux, authMid)
	paymentHandler.RegisterRoutes(mux, authMid)
	requestHandler.RegisterRoutes(mux, authMid)

	finalHandler := handler.CORSMiddleware(handler.RecoveryMiddleware(mux))

	log.Println("Server started on :" + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), finalHandler))
}
