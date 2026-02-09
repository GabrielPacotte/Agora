package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GabrielPacotte/Agora/internal/auth"
	httpserver "github.com/GabrielPacotte/Agora/internal/http"
	authhttp "github.com/GabrielPacotte/Agora/internal/http/auth"
	commenthttp "github.com/GabrielPacotte/Agora/internal/http/comment"
	"github.com/GabrielPacotte/Agora/internal/http/middleware"
	posthandler "github.com/GabrielPacotte/Agora/internal/http/post"
	userhttp "github.com/GabrielPacotte/Agora/internal/http/user"
	"github.com/GabrielPacotte/Agora/internal/repositories/postgres"
	"github.com/GabrielPacotte/Agora/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := os.Getenv("KAINE_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:Gabriel2@127.0.0.1:5432/agora-local?sslmode=disable"
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("pgxpool.New: %v", err)
	}
	defer pool.Close()

	// Repositories
	sessionRepo := postgres.NewSessionRepositoryPG(pool)
	userRepo := postgres.NewUserRepositoryPG(pool)
	prefRepo := postgres.NewUserPreferencesRepositoryPG(pool)
	savedContentRepo := postgres.NewSavedContentRepositoryPG(pool)
	postRepo := postgres.NewPostRepositoryPG(pool)
	commentRepo := postgres.NewCommentRepositoryPG(pool)

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtIssuer := auth.NewJWTIssuer(
		jwtSecret,
		"agora",
	)
	tokenVerifier := auth.NewAccessTokenVerifier([]byte(jwtSecret), "agora")

	authMiddleware := middleware.RequireAuth(tokenVerifier)

	opaqueGen := auth.NewOpaqueTokenGenerator()

	// Services
	authSvc := services.NewAuthService(userRepo, sessionRepo, jwtIssuer, opaqueGen, time.Minute*15, time.Hour*720)
	postSvc := services.NewPostService(postRepo)
	commentSvc := services.NewCommentService(commentRepo, postRepo)
	userSvc := services.NewUserService(userRepo, prefRepo, savedContentRepo)

	// Handlers
	authHandler := authhttp.NewHandler(authSvc, false)
	postHandler := posthandler.NewHandler(postSvc, commentSvc)
	commentHandler := commenthttp.NewHandler(postSvc, commentSvc)
	userHandler := userhttp.NewHandler(userSvc)

	router := httpserver.NewRouter(httpserver.Handlers{
		AuthMw:   authMiddleware,
		Auth:     authHandler,
		Posts:    postHandler,
		Comments: commentHandler,
		Users:    userHandler,
	})

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("HTTP server listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
