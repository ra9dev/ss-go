package broken_access_control

import (
	"context"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"

	ssgo "github.com/ra9dev/ss-go"
)

func TestQueryHacker_Attack(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	baseURL, wait := ssgo.RunTestServer(t, ctx, ssgo.ServerWithRoute(NewQueryAttackTarget()))
	defer wait()

	defer cancel()

	hackPath := baseURL + accountPath
	hackers := []QueryHacker{
		NewQueryHacker(
			hackPath,
			map[string]string{
				nicknameParam: "admin",
			},
		),
		NewQueryHacker(
			hackPath,
			map[string]string{
				nicknameParam: gofakeit.Username(),
				"otherParam":  gofakeit.HackerNoun(),
			},
		),
	}

	for _, hacker := range hackers {
		err := hacker.Attack()
		assert.NoError(t, err)
	}
}
