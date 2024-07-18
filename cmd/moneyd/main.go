package main

import (
	"context"
	"log/slog"
	nhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sneakynet/moneyprinter/pkg/bill"
	"github.com/sneakynet/moneyprinter/pkg/db"
	"github.com/sneakynet/moneyprinter/pkg/http"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	db, err := db.New()
	if err != nil {
		slog.Error("Error creating database", "error", err)
		return
	}

	dbPath := os.Getenv("MONEYPRINTER_DB")
	if dbPath == "" {
		dbPath = "moneyprinter.db"
	}

	if err := db.Connect(dbPath); err != nil {
		slog.Error("Error connecting to database", "error", err)
		return
	}

	if err := db.Migrate(); err != nil {
		slog.Error("Error migrating database", "error", err)
		return
	}

	bp := bill.NewProcessor(bill.WithDatabase(db))

	s, err := http.New(http.WithDB(db), http.WithBillProcessor(bp))
	if err != nil {
		slog.Error("Error creating webserver", "error", err)
		return
	}

	go func() {
		if err := s.Serve(":8080"); err != nil && err != nhttp.ErrServerClosed {
			slog.Error("Error initializing", "error", err)
			quit <- syscall.SIGINT
		}
	}()

	slog.Info("MoneyPrinter is ready to go BRRRRRRRR")
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		slog.Error("Error during shutdown", "error", err)
		return
	}
}
