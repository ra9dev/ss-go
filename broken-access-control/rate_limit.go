package broken_access_control

import (
	"net/http"
	"sync"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	ssgo "github.com/ra9dev/ss-go"
)

const (
	stocksPath = "/stocks"

	stocksPathRPSLimit = 3

	ipHeaderKey = "X-Real-IP"
)

// NewRateLimitAttackTarget to imitate vulnerability target for RateLimitHacker
func NewRateLimitAttackTarget() ssgo.ServerRoute {
	router := chi.NewRouter()

	reqPerAddr := make(map[string]uint64)
	mu := new(sync.Mutex)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		resp := make(map[string]any)

		remoteAddr := r.Header.Get(ipHeaderKey)

		mu.Lock()
		reqPerAddr[remoteAddr]++
		ssgo.ServerLogger.Printf(
			"%s made %d hits, limit is %d RPS",
			remoteAddr,
			reqPerAddr[remoteAddr],
			stocksPathRPSLimit,
		)
		mu.Unlock()

		for i := 0; i < 10; i++ {
			resp[gofakeit.Company()] = gofakeit.Price(0, float64(gofakeit.Uint32()))
		}

		render.JSON(w, r, resp)
	})

	return ssgo.ServerRoute{
		Pattern: stocksPath,
		Handler: router,
	}
}
