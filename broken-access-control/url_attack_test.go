package broken_access_control

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	ssgo "github.com/ra9dev/ss-go"
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
