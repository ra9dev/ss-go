package broken_access_control

import (
	"context"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ssgo "github.com/ra9dev/ss-go"
)

func newCardJWT(t *testing.T, userID uuid.UUID) string {
	t.Helper()

	token, err := ssgo.NewSignedJWT(
		[]byte(cardAccessJWTSecret),
		userID,
		map[string][]ssgo.Permission{
			permissionsCardKey: {ssgo.PermissionRead, ssgo.PermissionWrite},
		},
	)
	require.NoError(t, err)

	return token
}

func TestModelAccessControlHacker_Attack(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	baseURL, wait := ssgo.RunTestServer(t, ctx, ssgo.ServerWithRoute(NewModelAccessControlAttackTarget()))
	defer wait()

	defer cancel()

	hackPath := baseURL + cardPath

	hackerUserID := uuid.New()
	hackerJWTToken := newCardJWT(t, hackerUserID)
	ssgo.HackerLogger.Printf("Obtained jwt token for hacker account %s", hackerUserID)

	userIDToHack := uuid.New()
	cardToHack := Card{
		UserID:         userIDToHack,
		CreditCardInfo: gofakeit.CreditCard(),
	}

	hacker := NewModelAccessControlHacker(
		hackPath,
		hackerJWTToken,
		cardToHack,
	)

	err := hacker.Attack()
	assert.NoError(t, err)
}
