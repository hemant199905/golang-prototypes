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

	"mydung/internal/api"
	"mydung/internal/models"
	"mydung/internal/queue"
	"mydung/internal/storage"

	"github.com/spf13/viper"
)

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
	fmt.Printf("Starting Dung Beetle with %d workers... 🚀\n", config.Worker.Count)

	jobQueue := make(chan models.Job, config.Worker.QueueSize)
	db, err := storage.NewRedisStore(config.Redis)
	if err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}
	// db := storage.NewMemoryStore()
	// 1. WaitGroup setup karo
	var wg sync.WaitGroup

	for i := 1; i <= config.Worker.Count; i++ {
		wg.Add(1) // Attendance register mein +1
		go queue.Worker(i, jobQueue, &wg, db) // wg pass kiya
	}

	// API Setup
	http.HandleFunc("/submit-job", api.MakeJobHandler(jobQueue, db))
	http.HandleFunc("/status", api.StatusHandler(db))

	serverAddress := fmt.Sprintf(":%d", config.App.Port)
	server := &http.Server{Addr: serverAddress}

	// Server ko background mein start karo
	go func() {
		fmt.Printf("🌐 REST API Server running on http://localhost%s\n", serverAddress)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server fail ho gaya: %v", err)
		}
	}()

	// 🛑 2. OS Signals ko pakadne ka jaal (Trap)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM) // Ctrl+C ka wait karega

	<-sigChan // 🌟 Program yahan ruka rahega jab tak Ctrl+C nahi dabta!
	server.Shutdown(context.Background())
	fmt.Println("\n🛑 Shutdown signal received! Gate lock kar rahe hain...")
	
	// 3. Naye jobs aana band karo (Close channel)
	close(jobQueue) 
	
	// 4. Wait for all chefs to finish (Wait)
	wg.Wait() 

	fmt.Println("✅ Sab workers safe log out ho gaye. System Shutdown. Bye!")
}