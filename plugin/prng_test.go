package sgc7plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PRNG(t *testing.T) {
	prng := NewPRNGPlugin()

	prng.SetSeed(1627)

	v1, err := prng.Random(context.Background(), 10000)
	assert.NoError(t, err)
	assert.Equal(t, v1, 8318)

	v2, err := prng.Random(context.Background(), 10000)
	assert.NoError(t, err)
	assert.Equal(t, v2, 2600)

	v3, err := prng.Random(context.Background(), 10000)
	assert.NoError(t, err)
	assert.Equal(t, v3, 6065)

	v4, err := prng.Random(context.Background(), 10000)
	assert.NoError(t, err)
	assert.Equal(t, v4, 8439)

	v5, err := prng.Random(context.Background(), 10000)
	assert.NoError(t, err)
	assert.Equal(t, v5, 1345)

	t.Logf("Test_PRNG OK")
}
