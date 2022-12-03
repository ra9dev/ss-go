package broken_access_control

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	ssgo "github.com/ra9dev/ss-go"
	"net/http"
)

func NewCORSAttackTarget() ssgo.ServerRoute {
	router := chi.NewRouter()

	router.Use(
		cors.AllowAll().Handler,
	)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]any{
			"name":  gofakeit.Name(),
			"phone": gofakeit.Phone(),
			"email": gofakeit.Email(),
		})
	})

	return ssgo.ServerRoute{
		Pattern: accountPath,
		Handler: router,
	}
}
