package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/AnangM/livestok-erp/internal/repository"
	"github.com/AnangM/livestok-erp/internal/service"
	"github.com/AnangM/livestok-erp/internal/transport/http/handler"
	"github.com/AnangM/livestok-erp/internal/transport/http/middleware"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	dbUrl := os.Getenv("DB_URL")
	supabaseUrl := os.Getenv("SUPABASE_URL")
	supabaseAnonKey := os.Getenv("SUPABASE_ANON_KEY")

	if dbUrl == "" || supabaseAnonKey == "" || supabaseUrl == "" {
		log.Fatal("ENV Variables (DB_URL, SUPABASE_URL, SUPABASE_ANON_KEY) must be set")
		return
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(middleware.ReqeustId())

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	authService := service.NewAuthService(supabaseUrl, supabaseAnonKey)
	authHandler := handler.NewAuthHandler(authService)
	r.Post("/api/v1/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.SupabaseAuth(supabaseUrl, supabaseAnonKey))
		r.Mount("/api/v1/animals", ApiRoutes(db))
	})

	log.Println("Starting server on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func ApiRoutes(db *sql.DB) chi.Router {
	r := chi.NewRouter()
	animalRepo := repository.NewPostgresAnimalRepository(db)
	animalService := service.NewAnimalService(animalRepo)
	animalHandler := handler.NewAnimalHandler(animalService)
	r.Post("/", animalHandler.Create)
	r.Get("/", animalHandler.List)
	r.Get("/{id}", animalHandler.Get)
	r.Put("/{id}", animalHandler.Update)
	r.Delete("/{id}", animalHandler.Delete)
	return r
}
