package broken_access_control

import (
	"context"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"

	ssgo "github.com/ra9dev/ss-go"
)

func TestRateLimitHacker_Attack(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	baseURL, wait := ssgo.RunTestServer(t, ctx, ssgo.ServerWithRoute(NewRateLimitAttackTarget()))
	defer wait()

	defer cancel()

	hackPath := baseURL + stocksPath
	hacker := NewRateLimitHacker(
		hackPath,
		gofakeit.IPv4Address()+":8080",
	)

	err := hacker.Attack()
	assert.NoError(t, err)
}
