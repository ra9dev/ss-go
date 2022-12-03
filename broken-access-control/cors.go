package broken_access_control

import (
	"log"
	"net/http"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	ssgo "github.com/ra9dev/ss-go"
)

const (
	locationPath    = "/location"
	originHeaderKey = "Origin"
)

// NewCORSAttackTarget to imitate vulnerability target for CORSHacker
func NewCORSAttackTarget() ssgo.ServerRoute {
	router := chi.NewRouter()

	router.Use(
		cors.AllowAll().Handler,
	)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get(locationPath)

		log.Printf("%s request processed for origin %s", locationPath, origin)

		render.JSON(w, r, gofakeit.Address())
	})

	return ssgo.ServerRoute{
		Pattern: locationPath,
		Handler: router,
	}
}
