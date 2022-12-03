package broken_access_control

import (
	"net/http"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	ssgo "github.com/ra9dev/ss-go"
)

const (
	accountPath   = "/account"
	nicknameParam = "nickname"
)

// NewQueryAttackTarget to imitate vulnerability target for QueryHacker
func NewQueryAttackTarget() ssgo.ServerRoute {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		nickname := r.URL.Query().Get(nicknameParam)

		render.JSON(w, r, map[string]any{
			nicknameParam: nickname,
			"name":        gofakeit.Name(),
			"phone":       gofakeit.Phone(),
			"email":       gofakeit.Email(),
		})
	})

	return ssgo.ServerRoute{
		Pattern: accountPath,
		Handler: router,
	}
}
