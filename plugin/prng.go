package sgc7plugin

import (
	"context"
)

const RngA int64 = 165403
const RngC int64 = 51654324
const RngM int64 = 2147483647

// PRNGPlugin - prng plugin
type PRNGPlugin struct {
	PluginBase
	Seed int
}

// NewPRNGPlugin - new a PRNGPlugin
func NewPRNGPlugin() *PRNGPlugin {
	prng := &PRNGPlugin{
		PluginBase: NewPluginBase(),
	}

	return prng
}

// Random - return [0, r)
func (prng *PRNGPlugin) Random(ctx context.Context, r int) (int, error) {
	prng.Seed = int((RngC*int64(prng.Seed) + RngA) % RngM)

	return prng.Seed % r, nil
}

// Init - initial
func (prng *PRNGPlugin) Init() {
}

// SetSeed - set a seed
func (prng *PRNGPlugin) SetSeed(seed int) {
	prng.Seed = seed
}
