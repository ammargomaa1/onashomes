package utils

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	// contextMap stores the mapping of goroutine ID to gin.Context
	contextMap = make(map[uint64]*gin.Context)
	// contextMutex protects concurrent access to contextMap
	contextMutex sync.RWMutex
)

// getGoroutineID returns the current goroutine ID
func getGoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	
	var id uint64
	_, err := fmt.Sscanf(string(b), "goroutine %d", &id)
	if err != nil {
		// Fallback: try a larger buffer or return 0
		b = make([]byte, 512)
		runtime.Stack(b, false)
		_, err = fmt.Sscanf(string(b), "goroutine %d", &id)
	}
	
	return id
}

// SetContextForGoroutine registers the gin.Context for the current goroutine
func SetContextForGoroutine(c *gin.Context) {
	goroutineID := getGoroutineID()
	contextMutex.Lock()
	defer contextMutex.Unlock()
	contextMap[goroutineID] = c
}

// GetContextForGoroutine retrieves the gin.Context for the current goroutine
func GetContextForGoroutine() *gin.Context {
	goroutineID := getGoroutineID()
	contextMutex.RLock()
	defer contextMutex.RUnlock()
	return contextMap[goroutineID]
}

// CleanupContextForGoroutine removes the context entry for the current goroutine
func CleanupContextForGoroutine() {
	goroutineID := getGoroutineID()
	contextMutex.Lock()
	defer contextMutex.Unlock()
	delete(contextMap, goroutineID)
}
