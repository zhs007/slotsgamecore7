package sgc7game

import (
	"errors"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/zhs007/slotsgamecore7/gati"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// testGame - game
type testGame struct {
	*BasicGame
}

func newtestGame() (*testGame, error) {
	game := &testGame{
		NewBasicGame(func() sgc7plugin.IPlugin {
			return gati.NewPluginGATI(&gati.Config{
				GameID:  "936207324",
				RNGURL:  "http://127.0.0.1:50000/numbers",
				RngNums: 100,
			})
		}),
	}

	game.Cfg.Width = 5
	game.Cfg.Height = 3
	game.Cfg.Reels = make(map[string]*ReelsData)

	// game.Plugin = gati.NewPluginGATI(&gati.Config{
	// 	GameID:  "936207324",
	// 	RNGURL:  "http://127.0.0.1:50000/numbers",
	// 	RngNums: 100,
	// })

	r, err := LoadReels5JSON("../unittestdata/reels.json")
	if err != nil {
		return nil, err
	}

	game.Cfg.Reels["bg"] = r

	r, err = LoadReels5JSON("../unittestdata/reels2.json")
	if err != nil {
		return nil, err
	}

	game.Cfg.Reels["fg1"] = r

	return game, nil
}

func Test_NewGameScene(t *testing.T) {
	game, err := newtestGame()
	if err != nil {
		t.Fatalf("Test_NewGameScene newtestGame err %v",
			err)
	}

	gs, err := NewGameScene(game.Cfg.Width, game.Cfg.Height)
	assert.NoError(t, err)

	if len(gs.Arr) != 5 {
		t.Fatalf("Test_NewGameScene NewGameScene width err %d",
			len(gs.Arr))
	}

	for x := 0; x < 5; x++ {
		if len(gs.Arr[x]) != 3 {
			t.Fatalf("Test_NewGameScene NewGameScene height err %d",
				len(gs.Arr[x]))
		}

		for y := 0; y < 3; y++ {
			if gs.Arr[x][y] != -1 {
				t.Fatalf("Test_NewGameScene NewGameScene value err [%d][%d] %d",
					x, y, gs.Arr[x][y])
			}
		}
	}

	t.Logf("Test_NewGameScene OK")
}

func Test_RandGameScene(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET",
		"http://127.0.0.1:50000/numbers?size=100",
		httpmock.NewStringResponder(200, "[0, 1, 2, 3, 4, 0, 1, 2, 3, 32]"))

	game, err := newtestGame()
	if err != nil {
		t.Fatalf("Test_RandGameScene newtestGame err %v",
			err)
	}

	plugin := game.NewPlugin()
	defer game.FreePlugin(plugin)

	gs, err := NewGameScene(game.Cfg.Width, game.Cfg.Height)
	assert.NoError(t, err)

	err = gs.RandReels(game, plugin, "bg")
	if err != nil {
		t.Fatalf("Test_RandGameScene RandReels err %v",
			err)
	}

	assert.Equal(t, gs.Arr[0][0], 1, "they should be equal")
	assert.Equal(t, gs.Arr[0][1], 8, "they should be equal")
	assert.Equal(t, gs.Arr[0][2], 10, "they should be equal")

	assert.Equal(t, gs.Arr[1][0], 5, "they should be equal")
	assert.Equal(t, gs.Arr[1][1], 7, "they should be equal")
	assert.Equal(t, gs.Arr[1][2], 11, "they should be equal")

	assert.Equal(t, gs.Arr[2][0], 8, "they should be equal")
	assert.Equal(t, gs.Arr[2][1], 2, "they should be equal")
	assert.Equal(t, gs.Arr[2][2], 7, "they should be equal")

	assert.Equal(t, gs.Arr[3][0], 11, "they should be equal")
	assert.Equal(t, gs.Arr[3][1], 6, "they should be equal")
	assert.Equal(t, gs.Arr[3][2], 1, "they should be equal")

	assert.Equal(t, gs.Arr[4][0], 10, "they should be equal")
	assert.Equal(t, gs.Arr[4][1], 6, "they should be equal")
	assert.Equal(t, gs.Arr[4][2], 0, "they should be equal")

	err = gs.RandReels(game, plugin, "fg1")
	if err != nil {
		t.Fatalf("Test_RandGameScene RandReels err %v",
			err)
	}

	assert.Equal(t, gs.Arr[0][0], 1, "they should be equal")
	assert.Equal(t, gs.Arr[0][1], 8, "they should be equal")
	assert.Equal(t, gs.Arr[0][2], 10, "they should be equal")

	assert.Equal(t, gs.Arr[1][0], 5, "they should be equal")
	assert.Equal(t, gs.Arr[1][1], 7, "they should be equal")
	assert.Equal(t, gs.Arr[1][2], 11, "they should be equal")

	assert.Equal(t, gs.Arr[2][0], 8, "they should be equal")
	assert.Equal(t, gs.Arr[2][1], 2, "they should be equal")
	assert.Equal(t, gs.Arr[2][2], 7, "they should be equal")

	assert.Equal(t, gs.Arr[3][0], 11, "they should be equal")
	assert.Equal(t, gs.Arr[3][1], 6, "they should be equal")
	assert.Equal(t, gs.Arr[3][2], 1, "they should be equal")

	assert.Equal(t, gs.Arr[4][0], 5, "they should be equal")
	assert.Equal(t, gs.Arr[4][1], 4, "they should be equal")
	assert.Equal(t, gs.Arr[4][2], 9, "they should be equal")

	t.Logf("Test_RandGameScene OK")
}

func Test_ForEachAround(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET",
		"http://127.0.0.1:50000/numbers?size=100",
		httpmock.NewStringResponder(200, "[0, 1, 2, 3, 4, 0, 1, 2, 3, 32]"))

	game, err := newtestGame()
	assert.NoError(t, err)

	plugin := game.NewPlugin()
	defer game.FreePlugin(plugin)

	gs, err := NewGameScene(game.Cfg.Width, game.Cfg.Height)
	assert.NoError(t, err)

	err = gs.RandReels(game, plugin, "bg")
	assert.NoError(t, err)

	gs.ForEachAround(-1, 0, func(x, y int, val int) {
		assert.NoError(t, errors.New("Test_ForEachAround ForEachAround non-call"))
	})

	gs.ForEachAround(0, 3, func(x, y int, val int) {
		assert.NoError(t, errors.New("Test_ForEachAround ForEachAround non-call"))
	})

	gs.ForEachAround(5, 0, func(x, y int, val int) {
		assert.NoError(t, errors.New("Test_ForEachAround ForEachAround non-call"))
	})

	gs.ForEachAround(0, -1, func(x, y int, val int) {
		assert.NoError(t, errors.New("Test_ForEachAround ForEachAround non-call"))
	})

	lstpos := []int{}
	lstv := []int{}
	gs.ForEachAround(0, 0, func(x, y int, val int) {
		lstpos = append(lstpos, x)
		lstpos = append(lstpos, y)

		lstv = append(lstv, val)
	})

	assert.Equal(t, len(lstv), 3)
	assert.Equal(t, len(lstpos), 6)

	for i, v := range lstv {
		assert.Equal(t, v, gs.Arr[lstpos[i*2]][lstpos[i*2+1]])
	}

	lstpos = []int{}
	lstv = []int{}
	gs.ForEachAround(1, 1, func(x, y int, val int) {
		lstpos = append(lstpos, x)
		lstpos = append(lstpos, y)

		lstv = append(lstv, val)
	})

	assert.Equal(t, len(lstv), 8)
	assert.Equal(t, len(lstpos), 16)

	for i, v := range lstv {
		assert.Equal(t, v, gs.Arr[lstpos[i*2]][lstpos[i*2+1]])
	}

	t.Logf("Test_ForEachAround OK")
}

func Test_ForEach(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET",
		"http://127.0.0.1:50000/numbers?size=100",
		httpmock.NewStringResponder(200, "[0, 1, 2, 3, 4, 0, 1, 2, 3, 32]"))

	game, err := newtestGame()
	assert.NoError(t, err)

	plugin := game.NewPlugin()
	defer game.FreePlugin(plugin)

	gs, err := NewGameScene(game.Cfg.Width, game.Cfg.Height)
	assert.NoError(t, err)

	err = gs.RandReels(game, plugin, "bg")
	assert.NoError(t, err)

	gs.ForEachAround(-1, 0, func(x, y int, val int) {
		assert.NoError(t, errors.New("Test_ForEachAround ForEachAround non-call"))
	})

	gs.ForEachAround(0, 3, func(x, y int, val int) {
		assert.NoError(t, errors.New("Test_ForEachAround ForEachAround non-call"))
	})

	gs.ForEachAround(5, 0, func(x, y int, val int) {
		assert.NoError(t, errors.New("Test_ForEachAround ForEachAround non-call"))
	})

	gs.ForEachAround(0, -1, func(x, y int, val int) {
		assert.NoError(t, errors.New("Test_ForEachAround ForEachAround non-call"))
	})

	lstpos := []int{}
	lstv := []int{}
	gs.ForEachAround(0, 0, func(x, y int, val int) {
		lstpos = append(lstpos, x)
		lstpos = append(lstpos, y)

		lstv = append(lstv, val)
	})

	assert.Equal(t, len(lstv), 3)
	assert.Equal(t, len(lstpos), 6)

	for i, v := range lstv {
		assert.Equal(t, v, gs.Arr[lstpos[i*2]][lstpos[i*2+1]])
	}

	lstpos = []int{}
	lstv = []int{}
	gs.ForEachAround(1, 1, func(x, y int, val int) {
		lstpos = append(lstpos, x)
		lstpos = append(lstpos, y)

		lstv = append(lstv, val)
	})

	assert.Equal(t, len(lstv), 8)
	assert.Equal(t, len(lstpos), 16)

	for i, v := range lstv {
		assert.Equal(t, v, gs.Arr[lstpos[i*2]][lstpos[i*2+1]])
	}

	t.Logf("Test_ForEachAround OK")
}
