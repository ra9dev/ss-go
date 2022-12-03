package broken_access_control

import (
	"context"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"

	ssgo "github.com/ra9dev/ss-go"
)

func TestCORSHacker_Attack(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	baseURL, wait := ssgo.RunTestServer(t, ctx, ssgo.ServerWithRoute(NewCORSAttackTarget()))
	defer wait()

	defer cancel()

	hackPath := baseURL + locationPath
	hacker := NewCORSHacker(
		hackPath,
		gofakeit.URL(),
		gofakeit.URL(),
		gofakeit.URL(),
	)

	err := hacker.Attack()
	assert.NoError(t, err)
}
