package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	thisDocs "github.com/mibrgmv/document-service/docs"
	"github.com/mibrgmv/document-service/internal/config"
	"github.com/mibrgmv/document-service/internal/handlers"
	"github.com/mibrgmv/document-service/internal/repository/postgres"
	"github.com/mibrgmv/document-service/internal/repository/redis"
	"github.com/mibrgmv/document-service/internal/service"
	"github.com/mibrgmv/document-service/pkg/database"
	"github.com/mibrgmv/document-service/pkg/jwt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	cfg    *config.Config
	server *http.Server
}

func New(cfg *config.Config) *Server {
	if err := database.RunMigrations(cfg); err != nil {
		log.Fatal("failed to run migrations: ", err)
	}

	pg, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	rdb, err := database.NewRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}

	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.Expiration)

	userRepo := postgres.NewUserRepository(pg)
	docRepo := postgres.NewDocumentRepository(pg)
	cacheRepo := redis.NewCacheRepository(rdb)

	authService := service.NewAuthService(userRepo, cacheRepo, jwtManager, cfg.AdminToken)
	docService := service.NewDocumentService(docRepo, cacheRepo, userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	docHandler := handlers.NewDocumentHandler(docService)

	router := gin.Default()
	router.Use(handlers.CORSMiddleware())

	thisDocs.SwaggerInfo.Host = cfg.Server.Port
	if cfg.Server.Port != "" {
		thisDocs.SwaggerInfo.Host = "localhost:" + cfg.Server.Port
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/auth", authHandler.Auth)
		api.DELETE("/auth/:token", authHandler.Logout)

		docs := api.Group("/docs")
		docs.Use(handlers.AuthMiddleware(jwtManager))
		{
			docs.GET("", docHandler.GetDocuments)
			docs.HEAD("", docHandler.GetDocumentsHead)
			docs.POST("", docHandler.UploadDocument)
			docs.GET("/:id", docHandler.GetDocument)
			docs.HEAD("/:id", docHandler.GetDocumentHead)
			docs.DELETE("/:id", docHandler.DeleteDocument)
		}
	}

	httpServer := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		cfg:    cfg,
		server: httpServer,
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	log.Println("Gracefully shutting down server...")
	return s.server.Shutdown(ctx)
}
