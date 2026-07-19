package api_middlewares

import (
	"bytes"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type cacheEntry struct {
	body      []byte
	expiresAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]cacheEntry
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: make(map[string]cacheEntry),
	}
}

var GlobalCache *MemoryCache = NewMemoryCache()

func AddToGlobal(r chi.Router) {
	r.Use(GlobalCache.BasicCache())
}

func (c *MemoryCache) BasicCache() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only cache safe GET requests
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			key := r.URL.RequestURI() // e.g., "/api/questions?chapters=1,2,3"
			log.Printf("key %v", key)

			// 1. Check Cache Hit
			c.mu.RLock()
			entry, found := c.store[key]
			c.mu.RUnlock()

			if found {
				w.Header().Set("X-Cache", "HIT")
				byteCount, err := w.Write(entry.body)
				if err != nil {
					log.Printf("Error writhing cache: %v", err)
				}
				log.Printf("Written %v bytes", byteCount)
				return
			}

			// 2. Cache Miss: Intercept the response using a custom recorder
			rec := &responseRecorder{ResponseWriter: w, body: &bytes.Buffer{}}
			w.Header().Set("X-Cache", "MISS")

			next.ServeHTTP(rec, r)

			// 3. Save to memory if it was a successful 200 OK
			if rec.status == http.StatusOK || rec.status == 0 {
				c.mu.Lock()
				c.store[key] = cacheEntry{
					body: rec.body.Bytes(),
				}
				c.mu.Unlock()
			}
		})
	}
}

// responseRecorder captures the status code and body as it gets sent out
type responseRecorder struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b) // Save copy for cache
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
