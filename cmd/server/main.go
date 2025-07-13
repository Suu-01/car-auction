package main

import (
	"context"
	stdlog "log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/ksj/car-auction/docs"
	"github.com/ksj/car-auction/internal/api"
	"github.com/ksj/car-auction/internal/config"
	"github.com/ksj/car-auction/internal/log"
	"github.com/ksj/car-auction/internal/metrics"
	"github.com/ksj/car-auction/internal/model"
	"github.com/ksj/car-auction/internal/repo"
	"github.com/ksj/car-auction/internal/service"
	"github.com/ksj/car-auction/internal/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Logger.Error("health write failed", zap.Error(err))
	}
}

func main() {
	log.Init()
	defer func() {
		_ = log.Logger.Sync()
	}()

	// 1) 설정 로드: .env + 환경변수
	config.Load()
	log.Logger.Info("Config loaded", zap.Any("cfg", config.Cfg))
	db, err := gorm.Open(mysql.Open(config.Cfg.DSN), &gorm.Config{})
	if err != nil {
		stdlog.Fatal(err)
	}

	// 2) AutoMigrate
	if err := db.AutoMigrate(&model.Auction{}, &model.Bid{}, &model.User{}); err != nil {
		stdlog.Fatal(err)
	}

	auctionRepo := repo.NewAuctionRepo(db)
	bidRepo := repo.NewBidRepo(db)
	userRepo := repo.NewUserRepo(db)

	auctionSvc := service.NewAuctionService(auctionRepo)
	bidSvc := service.NewBidService(bidRepo)
	userSvc := service.NewUserService(userRepo)

	shutdown := tracing.Init()
	defer shutdown(context.Background())

	// 5) 라우터 설정
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/healthz", metrics.InstrumentHandler("/healthz", http.HandlerFunc(healthHandler)))
	api.RegisterUserRoutes(r, userSvc)
	api.RegisterAuctionRoutes(r, auctionSvc)
	api.RegisterBidRoutes(r, bidSvc)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	stdlog.Printf("Swagger UI: http://localhost:%s/swagger/index.html\n", config.Cfg.Port)
	api.RegisterHealthRoute(r)
	r.Use(func(next http.Handler) http.Handler {
		tracer := otel.Tracer("car-auction")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracer.Start(r.Context(), r.URL.Path)
			defer span.End()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	addr := ":" + config.Cfg.Port
	log.Logger.Info("Server starting", zap.String("addr", addr), zap.String("dsn", config.Cfg.DSN))
	if err := http.ListenAndServe(addr, r); err != nil {
		stdlog.Fatal(err)
	}
}
