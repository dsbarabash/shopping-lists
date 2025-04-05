package app

import (
	"context"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/dsbarabash/shopping-lists/internal/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg *config.Config
	ctx context.Context
}

func NewService(ctx context.Context) (*App, error) {
	// Инит баз клиентов
	return &App{
		ctx: ctx,
		cfg: config.NewConfig(),
	}, nil
}

func (a *App) Start() error {
	ctx, stop := signal.NotifyContext(a.ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		handler.Base(a.ctx, w, r)
	})
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		handler.Login(a.ctx, w, r)
	})
	mux.HandleFunc("POST /api/item/", func(w http.ResponseWriter, r *http.Request) {
		handler.AddItem(a.ctx, w, r)
	})
	mux.HandleFunc("POST /api/shopping_list/", func(w http.ResponseWriter, r *http.Request) {
		handler.AddShoppingList(a.ctx, w, r)
	})
	mux.HandleFunc("GET /api/items", func(w http.ResponseWriter, r *http.Request) {
		handler.GetItems(a.ctx, w, r)
	})
	mux.HandleFunc("GET /api/shopping_lists", func(w http.ResponseWriter, r *http.Request) {
		handler.GetSls(a.ctx, w, r)
	})
	mux.HandleFunc("GET /api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetItemById(a.ctx, w, r)
	})
	mux.HandleFunc("GET /api/shopping_list/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetSlById(a.ctx, w, r)
	})
	mux.HandleFunc("DELETE /api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.DeleteItemById(a.ctx, w, r)
	})
	mux.HandleFunc("DELETE /api/shopping_list/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.DeleteSlById(a.ctx, w, r)
	})
	mux.HandleFunc("PUT /api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.UpdateItemById(a.ctx, w, r)
	})
	mux.HandleFunc("PUT /api/shopping_list/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.UpdateSlById(a.ctx, w, r)
	})

	serverHTTP := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port),
		Handler: mux,
	}

	go func() {
		if err := serverHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("got interruption signal")
	ctxT, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := serverHTTP.Shutdown(ctxT); err != nil {
		return fmt.Errorf("shutdown server: %s\n", err)
	}
	log.Println("FINAL server shutdown")
	return nil
}
