package main

//import (
//	"awesomeProject1/config"
//	"awesomeProject1/internal/repository"
//	_ "awesomeProject1/internal/repository"
//	"awesomeProject1/internal/server"
//	_ "awesomeProject1/internal/server"
//	"awesomeProject1/internal/service"
//	_ "awesomeProject1/internal/service"
//	_ "context"
//	"database/sql"
//	"errors"
//	_ "errors"
//	"github.com/gorilla/mux"
//	"github.com/prometheus/client_golang/prometheus/promhttp"
//	"log"
//	"net/http"
//	_ "net/http"
//	"os"
//	"os/signal"
//	"strings"
//	"syscall"
//
//	"github.com/caarlos0/env"
//	_ "github.com/gorilla/mux"
//	"github.com/sirupsen/logrus"
//
//	"github.com/go-sql-driver/mysql"
//	_ "github.com/prometheus/client_golang/prometheus/promhttp"
//)
//
//const dbName = "orders"
//
//func main() {
//	// Config setup
//	var cfg config.Config
//	if err := env.Parse(&cfg); err != nil {
//		panic(err)
//	}
//
//	// Logs setup
//	logger := logrus.New()
//	level, err := logrus.ParseLevel(strings.ToLower(cfg.LogLevel))
//	if err != nil {
//		panic(err)
//	}
//	logger.SetLevel(level)
//
//	sigCh := make(chan os.Signal, 1)
//	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
//
//	pgCfg := mysql.Config{
//		User:   cfg.User,
//		Passwd: cfg.Password,
//		Net:    "tcp",
//		Addr:   cfg.Host,
//		DBName: dbName,
//	}
//
//	db, err := sql.Open("mysql", pgCfg.FormatDSN())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	repo := repository.NewDB(db, logger)
//
//	orderService := service.NewOrderService(repo, logger)
//
//	orderServer := server.NewServer(orderService, logger)
//
//	err = orderServer.Run(cfg.Port)
//	if err != nil {
//		log.Fatalf("cant run parser service controller")
//	}
//
//	<-sigCh
//
//}
//
//func startMetricsServer(ctx context.Context, port string, logger *logrus.Logger) {
//	router := mux.NewRouter()
//
//	router.Path("/metrics").Handler(promhttp.Handler())
//
//	s := &http.Server{
//		Addr:    ":" + port,
//		Handler: router,
//	}
//
//	go func() {
//		logger.Infof("Starting metrics server at :%s/metrics", port)
//		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
//			logger.Errorf("Metrics server error: %v", err)
//		}
//	}()
//
//	go func() {
//		<-ctx.Done()
//		if err := s.Shutdown(ctx); err != nil {
//			logger.Errorf("Error shutting down metrics server: %v", err)
//		} else {
//			logger.Info("Metrics server stopped.")
//		}
//	}()
//}
