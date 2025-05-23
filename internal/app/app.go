package app

import (
	"context"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/config"
	"github.com/dsbarabash/shopping-lists/internal/frontend/rest"
	"github.com/dsbarabash/shopping-lists/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	cfg    *config.Config
	ctx    context.Context
	server *rest.RestServer
}

func NewService(ctx context.Context, service service.Service, userService service.UserService) (*App, error) {
	// Инит баз клиентов
	c := &rest.RestServer{
		Service:     service,
		UserService: userService,
	}
	return &App{
		ctx:    ctx,
		cfg:    config.NewConfig(),
		server: c,
	}, nil
}

func (a *App) Start() error {
	ctx, stop := signal.NotifyContext(a.ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /registration", func(w http.ResponseWriter, r *http.Request) {
		a.server.Registration(w, r)
	})
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		a.server.Login(w, r)
	})
	mux.HandleFunc("POST /api/item", func(w http.ResponseWriter, r *http.Request) {
		rest.UserIdentity(w, r, a.server.AddItem)
	})
	mux.HandleFunc("GET /api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		rest.UserIdentity(w, r, a.server.GetItemById)
	})
	mux.HandleFunc("GET /api/items", func(w http.ResponseWriter, r *http.Request) {
		rest.UserIdentity(w, r, a.server.GetItems)
	})
	mux.HandleFunc("PUT /api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		rest.UserIdentity(w, r, a.server.UpdateItemById)
	})
	mux.HandleFunc("DELETE /api/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		rest.UserIdentity(w, r, a.server.DeleteItemById)
	})
	mux.HandleFunc("POST /api/shopping_list", func(w http.ResponseWriter, r *http.Request) {
		rest.UserIdentity(w, r, a.server.AddShoppingList)
	})
	mux.HandleFunc("GET /api/shopping_list/{id}", func(w http.ResponseWriter, r *http.Request) {
		rest.UserIdentity(w, r, a.server.GetShoppingListById)
	})
	mux.HandleFunc("GET /api/shopping_lists", func(w http.ResponseWriter, r *http.Request) {
		rest.UserIdentity(w, r, a.server.GetShoppingLists)
	})
	mux.HandleFunc("PUT /api/shopping_list/{id}", func(w http.ResponseWriter, r *http.Request) {
		rest.UserIdentity(w, r, a.server.UpdateShoppingListById)
	})
	mux.HandleFunc("DELETE /api/shopping_list/{id}", func(w http.ResponseWriter, r *http.Request) {
		rest.UserIdentity(w, r, a.server.DeleteShoppingListById)
	})
	mux.HandleFunc("GET /api/shopping_list_items/{id}", func(w http.ResponseWriter, r *http.Request) {
		rest.UserIdentity(w, r, a.server.GetItemsByShoppingListId)
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
