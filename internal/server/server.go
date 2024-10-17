package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	config "github.com/mrpawan-gupta/depserve/configs"
	"github.com/rs/cors"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// APIServer main api server configuration and settings
type APIServer struct {
	ListenerAddr string
	Host         string
	Router       chi.Router
	httpServer   *http.Server
	config       *config.Config
	db           *gorm.DB
	cors         *cors.Cors
}

func (server *APIServer) Config() *config.Config {
	return server.config
}

func (server *APIServer) Database() *gorm.DB {
	return server.db
}

func (server *APIServer) InitServer() {
	server.setCors()
	server.initDatabase()
}

func (server *APIServer) initDatabase() {
	dataSourceName := fmt.Sprintf("postgres://%s:%d/%s?sslmode=%s&user=%s&password=%s",
		server.Config().Database.Host,
		server.Config().Database.Port,
		server.Config().Database.Name,
		server.Config().Database.SslMode,
		server.Config().Database.User,
		server.Config().Database.Pass,
	)
	log.Printf("dataSourceName %s\n", dataSourceName)
	server.initDataSource(dataSourceName)
}

func (server *APIServer) initDataSource(dataSourceName string) {
	//database, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
	//if err != nil {
	//	log.Fatalf("failed to connect database: %v", err)
	//}
	//server.db = database
}

func (server *APIServer) setCors() {
	server.cors = cors.New(
		cors.Options{
			AllowedOrigins: server.Config().Cors.AllowedOrigins,
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		})
}

func (server *APIServer) InitRouter() {
	server.Router = chi.NewRouter()
}

func (server *APIServer) InitMiddleWare() {
	//server.Router.NotFound(func(writer http.ResponseWriter, request *http.Request) {
	//	utils.HandleError(writer, http.StatusNotFound, "endpoint not found")
	//})
	server.Router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)
}

func NewAPIServer() *APIServer {
	conf := config.New()
	router := chi.NewRouter()
	return &APIServer{
		ListenerAddr: conf.Api.Port,
		Host:         conf.Api.Host,
		Router:       router,
		config:       conf,
	}
}

func (server *APIServer) RunServer() {
	log.Printf("Starting JSON API server on post %s on port %s\n", server.Host, server.ListenerAddr)
	server.httpServer = &http.Server{
		Addr:              server.Host + ":" + server.ListenerAddr,
		Handler:           server.Router,
		ReadHeaderTimeout: server.Config().Api.ReadHeaderTimeout,
	}
	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not listen on host %s and %s: %v\n", server.Host, server.ListenerAddr, err)
		}
	}()
	server.Stop()
}

// Stop stops the API server gracefully.
func (server *APIServer) Stop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), server.Config().Api.GracefulTimeout*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
	}
	log.Println("Server stopped gracefully")
}
