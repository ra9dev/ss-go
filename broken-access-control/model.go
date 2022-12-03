package broken_access_control

import (
	"encoding/json"
	"net/http"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	ssgo "github.com/ra9dev/ss-go"
	"github.com/ra9dev/ss-go/pkg/slice"
)

const (
	cardPath = "/card"

	cardAccessJWTSecret = "card-secret"
	permissionsCardKey  = "card"
)

// Card is an example structure for a credit/debit card
type Card struct {
	UserID uuid.UUID `json:"user_id"`
	*gofakeit.CreditCardInfo
}

// NewModelAccessControlAttackTarget to imitate vulnerability target for ModelAccessControlHacker
func NewModelAccessControlAttackTarget() ssgo.ServerRoute {
	router := chi.NewRouter()

	router.Put("/", func(w http.ResponseWriter, r *http.Request) {
		card := new(Card)

		if err := json.NewDecoder(r.Body).Decode(card); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		jwtToken := r.Header.Get(ssgo.AuthorizationHeaderKey)

		claims, err := ssgo.ClaimsFromJWT(jwtToken, []byte(cardAccessJWTSecret))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		cardPermissions, hasCardPermissions := claims.Permissions[permissionsCardKey]
		if !hasCardPermissions || !slice.Has(cardPermissions, ssgo.PermissionWrite) {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		ssgo.ServerLogger.Printf("updated card %s for user %s", card.Number, card.UserID)

		render.JSON(w, r, card)
	})

	return ssgo.ServerRoute{
		Pattern: cardPath,
		Handler: router,
	}
}
