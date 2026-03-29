package api


import (	"net/http"
)

// isAllowed helper function hai jo check karta hai ki queue name allowed queues mein hai ya nahi
func ValidateQueueMiddleware(next http.Handler, allowedQueues []string) http.Handler {
    return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
        queueName := r.URL.Query().Get("queue")

        // Agar queue name list mein nahi hai...
        if !isAllowed(queueName, allowedQueues) {
            // 1. User ko batao ki error hai
            http.Error(w, "Invalid queue name", http.StatusBadRequest)
            
            // 2. RETURN kar jao! (next.ServeHTTP call mat karo) ✋
            return 
        }

        // Agar sab sahi hai, toh aage badho
        next.ServeHTTP(w, r)
    })
}

func isAllowed(queueName string, allowedQueues []string) bool {
	for _, allowed := range allowedQueues {
		if queueName == allowed {
			return true
		}
	}
	return false
}