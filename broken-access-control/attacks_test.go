package broken_access_control

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	ssgo "github.com/ra9dev/ss-go"
)

func TestURLHacker_Attack(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wait := ssgo.RunTestServer(t, ctx, ssgo.ServerWithRoute(NewURLAttackTarget()))
	defer wait()

	baseURL := fmt.Sprintf("http://localhost:%d", ssgo.DefaultServerPort)
	hackPath := baseURL + appPath
	hacker := NewURLHacker(
		hackPath+publicAppInfoPath,
		hackPath+privateAppInfoPath,
	)

	if err := hacker.Attack(); err != nil {
		t.Errorf("failed to attack: %+v", err)

		return
	}

	cancel()
}

func TestQueryHacker_Attack(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wait := ssgo.RunTestServer(t, ctx, ssgo.ServerWithRoute(NewQueryAttackTarget()))
	defer wait()

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
		if err := hacker.Attack(); err != nil {
			t.Errorf("failed to attack: %+v", err)

			return
		}
	}

	cancel()
}
