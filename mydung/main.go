package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"mydung/internal/api"
	"mydung/internal/models"
	"mydung/internal/queue"
	"mydung/internal/storage"

	"github.com/spf13/viper"
)

// loadConfig: Purana logic sahi hai, bas ise error handling ke liye thoda saaf rakha hai
func loadConfig() models.Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Config read error: %v", err)
	}

	var cfg models.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Config decode error: %v", err)
	}
	return cfg
}

func main() {
	config := loadConfig()
	fmt.Printf("🚀 Starting My Dung with %s storage...\n", config.Storage.Type)

	// 1. Factory se JobStore initialize karein (Memory ya Redis)
	db, err := storage.NewJobStore(config)
	if err != nil {
		log.Fatalf("Store setup failed: %v", err)
	}

	// 2. Global context banayein workers ko signal bhejne ke liye
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// 3. Nested loop: Har queue ke liye naye workers start karein
	for _, qName := range config.Worker.Queues {
		for i := 1; i <= config.Worker.Count; i++ {
			wg.Add(1)
			// Worker function ab 'ctx' aur 'qName' dono lega
			go queue.Worker(ctx, i, db, &wg, qName)
		}
	}

	// 4. API Routes setup (Interface 'db' pass kar rahe hain)
	handlerWithValidation := api.ValidateQueueMiddleware(
    http.HandlerFunc(api.MakeJobHandler(db)), 
    config.Worker.Queues,
)
	http.Handle("/submit-job", handlerWithValidation)
	http.HandleFunc("/status", api.StatusHandler(db))

	serverAddress := fmt.Sprintf(":%d", config.App.Port)
	server := &http.Server{Addr: serverAddress}

	go func() {
		fmt.Printf("🌐 REST API Server running on http://localhost%s\n", serverAddress)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server fail ho gaya: %v", err)
		}
	}()

	// 🛑 Graceful Shutdown Logic
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan // Ctrl+C ka intezar
	fmt.Println("\n🛑 Shutdown signal received! Cleaning up...")

	// 5. Workers ko 'cancel' signal bhejein
	cancel() 

	// 6. Server ko band karein
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	server.Shutdown(shutdownCtx)
	
	// 7. Wait karein jab tak saare workers 'Done' na bol dein
	wg.Wait() 

	fmt.Println("✅ Sab workers safe log out ho gaye. System Shutdown. Bye!")
}