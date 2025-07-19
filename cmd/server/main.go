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
	"github.com/ksj/car-auction/internal/ws"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// healthHandlerはヘルスチェック用のエンドポイント
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Logger.Error("health write failed", zap.Error(err))
	}
}

func main() {
	// 1) ロギングの初期化
	log.Init()
	defer func() {
		_ = log.Logger.Sync()
	}()

	// 2) 設定のロード（.env + 環境変数）
	config.Load()
	log.Logger.Info("Config loaded", zap.Any("cfg", config.Cfg))

	// 3) DB 接続の設定
	db, err := gorm.Open(mysql.Open(config.Cfg.DSN), &gorm.Config{})
	if err != nil {
		stdlog.Fatal(err)
	}

	// 4) AutoMigrate: スキーマの自動生成／更新
	if err := db.AutoMigrate(&model.Auction{}, &model.Bid{}, &model.User{}); err != nil {
		stdlog.Fatal(err)
	}

	// 5) WebSocket ハブの生成および実行
	hub := ws.NewHub()
	go hub.Run()

	// 6) リポジトリおよびサービスの初期化
	auctionRepo := repo.NewAuctionRepo(db)
	bidRepo := repo.NewBidRepo(db)
	userRepo := repo.NewUserRepo(db)

	auctionSvc := service.NewAuctionService(auctionRepo)
	bidSvc := service.NewBidService(bidRepo, hub)
	userSvc := service.NewUserService(userRepo)

	// 7) トレーシングの初期化
	shutdown := tracing.Init()
	defer func() {
		_ = shutdown(context.Background())
	}()

	// 8) ルーター設定
	r := mux.NewRouter()
	r.Use(api.CORSMiddleware)

	// メトリクスとヘルスチェックのエンドポイント
	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/healthz", metrics.InstrumentHandler("/healthz", http.HandlerFunc(healthHandler)))

	// ビジネスドメインルートの登録
	api.RegisterUserRoutes(r, userSvc)
	api.RegisterAuctionRoutes(r, auctionSvc)
	api.RegisterWSRoutes(r, hub)
	api.RegisterBidRoutes(r, bidSvc)

	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	stdlog.Printf("Swagger UI: http://localhost:%s/swagger/index.html\n", config.Cfg.Port)
	api.RegisterHealthRoute(r)

	// トレーシングミドルウェア: 各リクエストに対してスパンを生成
	r.Use(func(next http.Handler) http.Handler {
		tracer := otel.Tracer("car-auction")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := tracer.Start(r.Context(), r.URL.Path)
			defer span.End()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// ファイルアップロードおよび静的ファイルのサーブ
	r.HandleFunc("/upload", api.UploadHandler).Methods("POST")
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./uploads"))))

	// 9) サーバ起動
	addr := ":" + config.Cfg.Port
	log.Logger.Info("Server starting", zap.String("addr", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		stdlog.Fatal(err)
	}
}
