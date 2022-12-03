package broken_access_control

import (
	"net/http"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	ssgo "github.com/ra9dev/ss-go"
)

const (
	appPath            = "/app"
	publicAppInfoPath  = "/info"
	privateAppInfoPath = "/adminInfo"
)

// NewURLAttackTarget to imitate vulnerability target for URLHacker
func NewURLAttackTarget() ssgo.ServerRoute {
	router := chi.NewRouter()

	router.Get(publicAppInfoPath, func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]any{
			"version":     gofakeit.AppVersion(),
			"description": "Public info",
		})
	})

	router.Get(privateAppInfoPath, func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]any{
			"version":     gofakeit.AppVersion(),
			"author":      gofakeit.AppAuthor(),
			"description": "Admin private info",
		})
	})

	return ssgo.ServerRoute{
		Pattern: appPath,
		Handler: router,
	}
}
