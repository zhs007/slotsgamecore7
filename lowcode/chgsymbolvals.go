package lowcode

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const ChgSymbolValsTypeName = "chgSymbolVals"

type ChgSymbolValsType int

const (
	CSVTypeInc ChgSymbolValsType = 0 // ++
	CSVTypeDec ChgSymbolValsType = 1 // --
	CSVTypeMul ChgSymbolValsType = 2 // *= multi
)

func parseChgSymbolValsType(strType string) ChgSymbolValsType {
	if strType == "dec" {
		return CSVTypeDec
	} else if strType == "mul" {
		return CSVTypeMul
	}

	return CSVTypeInc
}

type ChgSymbolValsSourceType int

const (
	CSVSTypePositionCollection ChgSymbolValsSourceType = 0 // positionCollection
	CSVSTypeWinResult          ChgSymbolValsSourceType = 1 // winResult
	CSVSTypeRow                ChgSymbolValsSourceType = 2 // row
	CSVSTypeColumn             ChgSymbolValsSourceType = 3 // column
)

func parseChgSymbolValsSourceType(strType string) ChgSymbolValsSourceType {
	if strType == "positioncollection" {
		return CSVSTypePositionCollection
	} else if strType == "row" {
		return CSVSTypeRow
	} else if strType == "column" {
		return CSVSTypeColumn
	}

	return CSVSTypeWinResult
}

type ChgSymbolValsData struct {
	BasicComponentData
	PosComponentData
}

// OnNewGame -
func (chgSymbolValsData *ChgSymbolValsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	chgSymbolValsData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (chgSymbolValsData *ChgSymbolValsData) onNewStep() {
	if !gIsReleaseMode {
		chgSymbolValsData.PosComponentData.Clear()
	}
}

// Clone
func (chgSymbolValsData *ChgSymbolValsData) Clone() IComponentData {
	if !gIsReleaseMode {
		target := &ChgSymbolValsData{
			BasicComponentData: chgSymbolValsData.CloneBasicComponentData(),
			PosComponentData:   chgSymbolValsData.PosComponentData.Clone(),
		}

		return target
	}

	target := &ChgSymbolValsData{
		BasicComponentData: chgSymbolValsData.CloneBasicComponentData(),
	}

	return target
}

// BuildPBComponentData
func (chgSymbolValsData *ChgSymbolValsData) BuildPBComponentData() proto.Message {
	return &sgc7pb.BasicComponentData{
		BasicComponentData: chgSymbolValsData.BuildPBBasicComponentData(),
	}
}

// GetPos -
func (chgSymbolValsData *ChgSymbolValsData) GetPos() []int {
	return chgSymbolValsData.Pos
}

// AddPos -
func (chgSymbolValsData *ChgSymbolValsData) AddPos(x, y int) {
	chgSymbolValsData.PosComponentData.Add(x, y)
}

// ChgSymbolValsConfig - configuration for ChgSymbolVals
type ChgSymbolValsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                  `yaml:"type" json:"type"`
	Type                 ChgSymbolValsType       `yaml:"-" json:"-"`
	StrSourceType        string                  `yaml:"sourceType" json:"sourceType"`
	SourceType           ChgSymbolValsSourceType `yaml:"-" json:"-"`
	PositionCollection   string                  `yaml:"positionCollection" json:"positionCollection"`
	WinResultComponents  []string                `yaml:"winResultComponents" json:"winResultComponents"`
	MaxNumber            int                     `yaml:"maxNumber" json:"maxNumber"`
	MaxVal               int                     `yaml:"maxVal" json:"maxVal"`
	MinVal               int                     `yaml:"minVal" json:"minVal"`
	Row                  int                     `yaml:"minVal" json:"row"`
	Column               int                     `yaml:"minVal" json:"column"`
	Multi                int                     `yaml:"multi" json:"multi"`
}

// SetLinkComponent
func (cfg *ChgSymbolValsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type ChgSymbolVals struct {
	*BasicComponent `json:"-"`
	Config          *ChgSymbolValsConfig `json:"config"`
}

// Init -
func (chgSymbolVals *ChgSymbolVals) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ChgSymbolVals.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ChgSymbolValsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ChgSymbolVals.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return chgSymbolVals.InitEx(cfg, pool)
}

// InitEx -
func (chgSymbolVals *ChgSymbolVals) InitEx(cfg any, pool *GamePropertyPool) error {
	chgSymbolVals.Config = cfg.(*ChgSymbolValsConfig)
	chgSymbolVals.Config.ComponentType = ChgSymbolValsTypeName

	chgSymbolVals.Config.Type = parseChgSymbolValsType(chgSymbolVals.Config.StrType)
	chgSymbolVals.Config.SourceType = parseChgSymbolValsSourceType(chgSymbolVals.Config.StrSourceType)

	chgSymbolVals.onInit(&chgSymbolVals.Config.BasicComponentConfig)

	return nil
}

func (chgSymbolVals *ChgSymbolVals) rebuildPos(pos []int, plugin sgc7plugin.IPlugin) ([]int, error) {
	if chgSymbolVals.Config.MaxNumber <= 0 {
		return pos, nil
	}

	if len(pos)/2 <= chgSymbolVals.Config.MaxNumber {
		return pos, nil
	}

	npos := []int{}

	for i := 0; i < chgSymbolVals.Config.MaxNumber; i++ {
		cr, err := plugin.Random(context.Background(), len(pos)/2)
		if err != nil {
			goutils.Error("ChgSymbolVals.rebuildPos:Random",
				goutils.Err(err))

			return nil, err
		}

		npos = append(npos, pos[cr*2], pos[cr*2+1])

		pos = append(pos[:cr*2], pos[(cr+1)*2:]...)
	}

	return npos, nil
}

func (chgSymbolVals *ChgSymbolVals) GetMulti(basicCD *BasicComponentData) int {
	multi, isok := basicCD.GetConfigIntVal(CCVMulti)
	if isok {
		return multi
	}

	return chgSymbolVals.Config.Multi
}

// playgame
func (chgSymbolVals *ChgSymbolVals) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// symbolVal2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*ChgSymbolValsData)

	cd.onNewStep()

	os := chgSymbolVals.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	if os != nil {
		nos := os

		if chgSymbolVals.Config.SourceType == CSVSTypePositionCollection {
			pc, isok := gameProp.Components.MapComponents[chgSymbolVals.Config.PositionCollection]
			if isok {
				pccd := gameProp.GetComponentData(pc)
				pos := pccd.GetPos()
				if len(pos) > 0 {
					npos, err := chgSymbolVals.rebuildPos(pos, plugin)
					if err != nil {
						goutils.Error("ChgSymbolVals.OnPlayGame:rebuildPos",
							goutils.Err(err))

						return "", nil
					}

					if chgSymbolVals.Config.Type == CSVTypeInc {
						for i := 0; i < len(npos)/2; i++ {
							if nos == os {
								nos = os.CloneEx(gameProp.PoolScene)
							}

							if nos.Arr[npos[i*2]][npos[i*2+1]] < chgSymbolVals.Config.MaxVal {
								nos.Arr[npos[i*2]][npos[i*2+1]]++

								if !gIsReleaseMode {
									cd.AddPos(npos[i*2], npos[i*2+1])
								}
							}
						}
					} else if chgSymbolVals.Config.Type == CSVTypeDec {
						for i := 0; i < len(npos)/2; i++ {
							if nos.Arr[npos[i*2]][npos[i*2+1]] > chgSymbolVals.Config.MinVal {
								if nos == os {
									nos = os.CloneEx(gameProp.PoolScene)
								}

								nos.Arr[npos[i*2]][npos[i*2+1]]--

								if !gIsReleaseMode {
									cd.AddPos(npos[i*2], npos[i*2+1])
								}
							}
						}
					} else if chgSymbolVals.Config.Type == CSVTypeMul {
						multi := chgSymbolVals.GetMulti(&cd.BasicComponentData)

						for i := 0; i < len(npos)/2; i++ {
							if nos == os {
								nos = os.CloneEx(gameProp.PoolScene)
							}

							nos.Arr[npos[i*2]][npos[i*2+1]] *= multi

							if nos.Arr[npos[i*2]][npos[i*2+1]] > chgSymbolVals.Config.MaxVal {
								nos.Arr[npos[i*2]][npos[i*2+1]] = chgSymbolVals.Config.MaxVal
							}

							if !gIsReleaseMode {
								cd.AddPos(npos[i*2], npos[i*2+1])
							}
						}
					}
				}
			}
		} else if chgSymbolVals.Config.SourceType == CSVSTypeWinResult {
			for _, cn := range chgSymbolVals.Config.WinResultComponents {
				ccd := gameProp.GetComponentDataWithName(cn)
				// ccd := gameProp.MapComponentData[cn]
				lst := ccd.GetResults()
				for _, ri := range lst {
					pos := curpr.Results[ri].Pos
					if len(pos) > 0 {
						npos, err := chgSymbolVals.rebuildPos(pos, plugin)
						if err != nil {
							goutils.Error("ChgSymbolVals.OnPlayGame:rebuildPos",
								goutils.Err(err))

							return "", nil
						}

						if chgSymbolVals.Config.Type == CSVTypeInc {
							for i := 0; i < len(npos)/2; i++ {
								if nos.Arr[npos[i*2]][npos[i*2+1]] < chgSymbolVals.Config.MaxVal {
									if nos == os {
										nos = os.CloneEx(gameProp.PoolScene)
									}

									nos.Arr[npos[i*2]][npos[i*2+1]]++

									if !gIsReleaseMode {
										cd.AddPos(npos[i*2], npos[i*2+1])
									}
								}
							}
						} else if chgSymbolVals.Config.Type == CSVTypeDec {
							for i := 0; i < len(npos)/2; i++ {
								if nos.Arr[npos[i*2]][npos[i*2+1]] > chgSymbolVals.Config.MinVal {
									if nos == os {
										nos = os.CloneEx(gameProp.PoolScene)
									}

									nos.Arr[npos[i*2]][npos[i*2+1]]--

									if !gIsReleaseMode {
										cd.AddPos(npos[i*2], npos[i*2+1])
									}
								}
							}
						} else if chgSymbolVals.Config.Type == CSVTypeMul {
							multi := chgSymbolVals.GetMulti(&cd.BasicComponentData)

							for i := 0; i < len(npos)/2; i++ {
								if nos == os {
									nos = os.CloneEx(gameProp.PoolScene)
								}

								nos.Arr[npos[i*2]][npos[i*2+1]] *= multi

								if nos.Arr[npos[i*2]][npos[i*2+1]] > chgSymbolVals.Config.MaxVal {
									nos.Arr[npos[i*2]][npos[i*2+1]] = chgSymbolVals.Config.MaxVal
								}

								if !gIsReleaseMode {
									cd.AddPos(npos[i*2], npos[i*2+1])
								}
							}
						}
					}
				}
			}
		} else if chgSymbolVals.Config.SourceType == CSVSTypeRow {
			if chgSymbolVals.Config.Type == CSVTypeInc {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				y := chgSymbolVals.Config.Row

				for x := 0; x < os.Width; x++ {
					if nos.Arr[x][y] < chgSymbolVals.Config.MaxVal {
						nos.Arr[x][y]++

						if !gIsReleaseMode {
							cd.AddPos(x, y)
						}
					}
				}
			} else if chgSymbolVals.Config.Type == CSVTypeDec {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				y := chgSymbolVals.Config.Row

				for x := 0; x < os.Width; x++ {
					if nos.Arr[x][y] > chgSymbolVals.Config.MinVal {
						nos.Arr[x][y]--

						if !gIsReleaseMode {
							cd.AddPos(x, y)
						}
					}
				}
			} else if chgSymbolVals.Config.Type == CSVTypeMul {
				multi := chgSymbolVals.GetMulti(&cd.BasicComponentData)

				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				y := chgSymbolVals.Config.Row

				for x := 0; x < os.Width; x++ {
					nos.Arr[x][y] *= multi
					if nos.Arr[x][y] > chgSymbolVals.Config.MaxVal {
						nos.Arr[x][y] = chgSymbolVals.Config.MaxVal
					}

					if !gIsReleaseMode {
						cd.AddPos(x, y)
					}
				}

			}
		} else if chgSymbolVals.Config.SourceType == CSVSTypeColumn {
			if chgSymbolVals.Config.Type == CSVTypeInc {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				x := chgSymbolVals.Config.Column

				for y := 0; y < os.Height; y++ {
					if nos.Arr[x][y] < chgSymbolVals.Config.MaxVal {
						nos.Arr[x][y]++

						if !gIsReleaseMode {
							cd.AddPos(x, y)
						}
					}
				}
			} else if chgSymbolVals.Config.Type == CSVTypeDec {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				x := chgSymbolVals.Config.Column

				for y := 0; y < os.Height; y++ {
					if nos.Arr[x][y] > chgSymbolVals.Config.MinVal {
						nos.Arr[x][y]--

						if !gIsReleaseMode {
							cd.AddPos(x, y)
						}
					}
				}
			} else if chgSymbolVals.Config.Type == CSVTypeMul {
				multi := chgSymbolVals.GetMulti(&cd.BasicComponentData)

				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				x := chgSymbolVals.Config.Column

				for y := 0; y < os.Height; y++ {
					nos.Arr[x][y] *= multi
					if nos.Arr[x][y] > chgSymbolVals.Config.MaxVal {
						nos.Arr[x][y] = chgSymbolVals.Config.MaxVal
					}

					if !gIsReleaseMode {
						cd.AddPos(x, y)
					}
				}
			}
		}

		if nos == os {
			nc := chgSymbolVals.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		chgSymbolVals.AddOtherScene(gameProp, curpr, nos, &cd.BasicComponentData)

		nc := chgSymbolVals.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := chgSymbolVals.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (chgSymbolVals *ChgSymbolVals) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*ChgSymbolValsData)

	if len(cd.UsedOtherScenes) > 0 {
		asciigame.OutputOtherScene("after ChgSymbolVals", pr.OtherScenes[cd.UsedOtherScenes[0]])
	}

	return nil
}

// NewComponentData -
func (chgSymbolVals *ChgSymbolVals) NewComponentData() IComponentData {
	return &ChgSymbolValsData{}
}

// // OnStats
// func (chgSymbolVals *ChgSymbolVals) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

func NewChgSymbolVals(name string) IComponent {
	return &ChgSymbolVals{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "maxNumber": 0,
// "maxVal": 99,
// "type": "mul",
// "sourceType": "row",
// "multi": 1,
// "row": "row4"
type jsonChgSymbolVals struct {
	Type                string   `json:"type"`
	SourceType          string   `json:"sourceType"`
	PositionCollection  string   `json:"positionCollection"`
	WinResultComponents []string `json:"winResultComponents"`
	MaxNumber           int      `json:"maxNumber"`
	MaxVal              int      `json:"maxVal"`
	MinVal              int      `json:"minVal"`
	Row                 string   `json:"row"`
	Column              string   `json:"column"`
	Multi               int      `json:"multi"`
}

func (jcfg *jsonChgSymbolVals) parseRow() int {
	if jcfg.Row != "" {
		arr := strings.Split(jcfg.Row, "row")
		if len(arr) == 2 {
			i64, err := goutils.String2Int64(arr[1])
			if err != nil {
				goutils.Error("jsonChgSymbolVals.parseRow:String2Int64",
					goutils.Err(err))

				return 0
			}

			return int(i64) - 1
		}
	}

	return 0
}

func (jcfg *jsonChgSymbolVals) parseColumn() int {
	if jcfg.Column != "" {
		arr := strings.Split(jcfg.Column, "column")
		if len(arr) == 2 {
			i64, err := goutils.String2Int64(arr[1])
			if err != nil {
				goutils.Error("jsonChgSymbolVals.parseColumn:String2Int64",
					goutils.Err(err))

				return 0
			}

			return int(i64) - 1
		}
	}

	return 0
}

func (jcfg *jsonChgSymbolVals) build() *ChgSymbolValsConfig {
	cfg := &ChgSymbolValsConfig{
		StrType:             jcfg.Type,
		StrSourceType:       strings.ToLower(jcfg.SourceType),
		PositionCollection:  jcfg.PositionCollection,
		WinResultComponents: jcfg.WinResultComponents,
		MaxNumber:           jcfg.MaxNumber,
		MaxVal:              jcfg.MaxVal,
		MinVal:              jcfg.MinVal,
		Row:                 jcfg.parseRow(),
		Column:              jcfg.parseColumn(),
		Multi:               jcfg.Multi,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseChgSymbolVals(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseChgSymbolVals:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseChgSymbolVals:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonChgSymbolVals{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseChgSymbolVals:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: ChgSymbolValsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
