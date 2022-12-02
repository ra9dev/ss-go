package broken_access_control

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	ssgo "github.com/ra9dev/ss-go"
)

func TestURLHacker_Attack(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wait := ssgo.RunTestServer(t, ctx, ssgo.ServerWithRoute(NewURLAttackTarget()))
	defer wait()
	defer cancel()

	baseURL := fmt.Sprintf("http://localhost:%d", ssgo.DefaultServerPort)
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
	wait := ssgo.RunTestServer(t, ctx, ssgo.ServerWithRoute(NewQueryAttackTarget()))
	defer wait()
	defer cancel()

	baseURL := fmt.Sprintf("http://localhost:%d", ssgo.DefaultServerPort)
	hackPath := baseURL + accountPath
	hackers := []ssgo.Hacker{
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
			},
		),
	}

	for _, hacker := range hackers {
		err := hacker.Attack()
		assert.NoError(t, err)
	}
}
