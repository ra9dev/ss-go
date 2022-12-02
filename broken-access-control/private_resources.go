package broken_access_control

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	ssgo "github.com/ra9dev/ss-go"
	"net/http"
)

const (
	appPath            = "/app"
	publicAppInfoPath  = "/info"
	privateAppInfoPath = "/adminInfo"
)

const (
	accountPath   = "/account"
	nicknameParam = "nickname"
)

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

func NewQueryAttackTarget() ssgo.ServerRoute {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		nickname := r.URL.Query().Get(nicknameParam)

		render.JSON(w, r, map[string]any{
			nicknameParam: nickname,
			"name":        gofakeit.Name(),
			"phone":       gofakeit.Phone(),
			"email":       gofakeit.Email(),
			"address":     gofakeit.Address(),
		})
	})

	return ssgo.ServerRoute{
		Pattern: accountPath,
		Handler: router,
	}
}
