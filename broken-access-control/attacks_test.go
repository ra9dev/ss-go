package broken_access_control

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	ssgo "github.com/ra9dev/ss-go"
	"github.com/stretchr/testify/assert"
)

func TestURLHacker_Attack(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	baseURL, wait := ssgo.RunTestServer(t, ctx, ssgo.ServerWithRoute(NewURLAttackTarget()))
	defer wait()
	defer cancel()

	hackPath := baseURL + appPath
	hacker := NewURLHacker(
		hackPath+publicAppInfoPath,
		hackPath+privateAppInfoPath,
	)

	err := hacker.Attack()
	assert.NoError(t, err)
}

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
