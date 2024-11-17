package server

import (
	"context"
	"hot-coffee/internal/handler"
	"hot-coffee/logging"
	"hot-coffee/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start(Port string) {
	defer utils.CatchCriticalPoint()

	http.HandleFunc("/menu/", handler.MenuHandler)
	http.HandleFunc("/menu", handler.MenuHandler)
	http.HandleFunc("/inventory/", handler.InventoryHandler)
	http.HandleFunc("/inventory", handler.InventoryHandler)
	http.HandleFunc("/order/", handler.OrderHandler)
	http.HandleFunc("/order", handler.OrderHandler)
	http.HandleFunc("/reports/", handler.ReportHandler)

	srv := &http.Server{
		Addr:         ":" + Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go gracefulShutdown(srv)

	logging.Info("Starting server on port", "port", Port)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logging.Error("Server failed to start", err, "port", Port)
		return
	}

	logging.Info("Server stopped gracefully", "port", Port)
}

func gracefulShutdown(srv *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logging.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logging.Error("Server forced to shutdown", err)
	}
}
