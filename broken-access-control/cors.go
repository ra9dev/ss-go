package broken_access_control

import (
	"net/http"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"

	ssgo "github.com/ra9dev/ss-go"
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
