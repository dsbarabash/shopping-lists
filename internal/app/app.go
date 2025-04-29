package app

import (
	"context"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/dsbarabash/shopping-lists/internal/handler"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg        *config.Config
	ctx        context.Context
	controller *handler.Controller
}

func NewService(ctx context.Context, mongoDB *repository.MongoDb) (*App, error) {
	// Инит баз клиентов
	c := &handler.Controller{
		MongoDb: mongoDB,
	}
	return &App{
		ctx:        ctx,
		cfg:        config.NewConfig(),
		controller: c,
	}, nil
}

func (a *App) Start() error {
	ctx, stop := signal.NotifyContext(a.ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /registration", func(w http.ResponseWriter, r *http.Request) {
		a.controller.Registration(w, r)
	})
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		a.controller.Login(w, r)
	})
	mux.HandleFunc("POST /api/item", func(w http.ResponseWriter, r *http.Request) {
		handler.UserIdentity(w, r, a.controller.AddItem)
	})
	mux.HandleFunc("POST /api/shopping_list", func(w http.ResponseWriter, r *http.Request) {
		handler.UserIdentity(w, r, a.controller.AddShoppingList)
	})
	mux.HandleFunc("GET /api/items", func(w http.ResponseWriter, r *http.Request) {
		handler.UserIdentity(w, r, a.controller.GetItems)
	})
	mux.HandleFunc("GET /api/shopping_lists", func(w http.ResponseWriter, r *http.Request) {
		handler.UserIdentity(w, r, a.controller.GetShoppingLists)
	})
	mux.HandleFunc("GET /api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.UserIdentity(w, r, a.controller.GetItemById)
	})
	mux.HandleFunc("GET /api/shopping_list/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.UserIdentity(w, r, a.controller.GetShoppingListById)
	})
	mux.HandleFunc("DELETE /api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.UserIdentity(w, r, a.controller.DeleteItemById)
	})
	mux.HandleFunc("DELETE /api/shopping_list/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.UserIdentity(w, r, a.controller.DeleteShoppingListById)
	})
	mux.HandleFunc("PUT /api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.UserIdentity(w, r, a.controller.UpdateItemById)
	})
	mux.HandleFunc("PUT /api/shopping_list/{id}", func(w http.ResponseWriter, r *http.Request) {
		handler.UserIdentity(w, r, a.controller.UpdateShoppingListById)
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
