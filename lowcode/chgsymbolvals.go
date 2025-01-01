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
	CSVTypeMul ChgSymbolValsType = 2 // *=
	CSVTypeAdd ChgSymbolValsType = 3 // +=
	CSVTypeSet ChgSymbolValsType = 4 // =
)

func parseChgSymbolValsType(strType string) ChgSymbolValsType {
	if strType == "dec" {
		return CSVTypeDec
	} else if strType == "mul" {
		return CSVTypeMul
	} else if strType == "add" {
		return CSVTypeAdd
	} else if strType == "set" {
		return CSVTypeSet
	}

	return CSVTypeInc
}

type ChgSymbolValsSourceType int

const (
	CSVSTypePositionCollection ChgSymbolValsSourceType = 0 // positionCollection
	CSVSTypeWinResult          ChgSymbolValsSourceType = 1 // winResult
	CSVSTypeRow                ChgSymbolValsSourceType = 2 // row
	CSVSTypeColumn             ChgSymbolValsSourceType = 3 // column
	CSVSTypeAll                ChgSymbolValsSourceType = 4 // all
)

func parseChgSymbolValsSourceType(strType string) ChgSymbolValsSourceType {
	if strType == "positioncollection" {
		return CSVSTypePositionCollection
	} else if strType == "row" {
		return CSVSTypeRow
	} else if strType == "column" {
		return CSVSTypeColumn
	} else if strType == "winResult" {
		return CSVSTypeWinResult
	}

	return CSVSTypeAll
}

type ChgSymbolValsTargetType int

const (
	CSVTTypeNumber     ChgSymbolValsTargetType = 0 // number
	CSVTTypeWeight     ChgSymbolValsTargetType = 1 // weight
	CSVTTypeEachWeight ChgSymbolValsTargetType = 2 // each weight
)

func parseChgSymbolValsTargetType(strType string) ChgSymbolValsTargetType {
	if strType == "weight" {
		return CSVTTypeWeight
	} else if strType == "eachWeight" {
		return CSVTTypeEachWeight
	}

	return CSVTTypeNumber
}

type ChgSymbolValsData struct {
	BasicComponentData
	PosComponentData
	targetVal int
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
	StrTargetType        string                  `yaml:"targetType" json:"targetType"`
	TargetType           ChgSymbolValsTargetType `yaml:"-" json:"-"`
	PositionCollection   string                  `yaml:"positionCollection" json:"positionCollection"`
	WinResultComponents  []string                `yaml:"winResultComponents" json:"winResultComponents"`
	MaxNumber            int                     `yaml:"maxNumber" json:"maxNumber"`
	MaxVal               int                     `yaml:"maxVal" json:"maxVal"`
	MinVal               int                     `yaml:"minVal" json:"minVal"`
	Row                  int                     `yaml:"minVal" json:"row"`
	Column               int                     `yaml:"minVal" json:"column"`
	Multi                int                     `yaml:"multi" json:"multi"`
	Number               int                     `yaml:"number" json:"number"`
	TargetWeight         string                  `yaml:"weight" json:"-"`
	TargetWeightVW2      *sgc7game.ValWeights2   `yaml:"-" json:"-"`
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
	chgSymbolVals.Config.TargetType = parseChgSymbolValsTargetType(chgSymbolVals.Config.StrTargetType)

	if chgSymbolVals.Config.TargetWeight != "" {
		vw2, err := pool.LoadIntWeights(chgSymbolVals.Config.TargetWeight, chgSymbolVals.Config.UseFileMapping)
		if err != nil {
			goutils.Error("ChgSymbolVals.Init:LoadStrWeights",
				slog.String("Weight", chgSymbolVals.Config.TargetWeight),
				goutils.Err(err))

			return err
		}

		chgSymbolVals.Config.TargetWeightVW2 = vw2
	}

	// 兼容性配置
	if chgSymbolVals.Config.TargetType == CSVTTypeNumber {
		if chgSymbolVals.Config.Number == 0 && chgSymbolVals.Config.Multi != 0 {
			chgSymbolVals.Config.Number = chgSymbolVals.Config.Multi
		}
	}

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

func (chgSymbolVals *ChgSymbolVals) GetNumber(cd *ChgSymbolValsData) int {
	multi, isok := cd.GetConfigIntVal(CCVMulti)
	if isok {
		return multi
	}

	number, isok := cd.GetConfigIntVal(CCVNumber)
	if isok {
		return number
	}

	return chgSymbolVals.Config.Number
}

func (chgSymbolVals *ChgSymbolVals) GetTarget(cd *ChgSymbolValsData, plugin sgc7plugin.IPlugin) (int, error) {
	if chgSymbolVals.Config.TargetType == CSVTTypeEachWeight {
		ival, err := chgSymbolVals.Config.TargetWeightVW2.RandVal(plugin)
		if err != nil {
			goutils.Error("ChgSymbolVals.GetTarget:RandVal",
				goutils.Err(err))

			return 0, err
		}

		return ival.Int(), nil
	}

	return cd.targetVal, nil
}

// playgame
func (chgSymbolVals *ChgSymbolVals) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*ChgSymbolValsData)

	cd.onNewStep()

	if chgSymbolVals.Config.TargetType == CSVTTypeNumber {
		cd.targetVal = chgSymbolVals.GetNumber(cd)
	} else if chgSymbolVals.Config.TargetType == CSVTTypeWeight && (chgSymbolVals.Config.Type == CSVTypeMul || chgSymbolVals.Config.Type == CSVTypeAdd || chgSymbolVals.Config.Type == CSVTypeSet) {
		ival, err := chgSymbolVals.Config.TargetWeightVW2.RandVal(plugin)
		if err != nil {
			goutils.Error("ChgSymbolVals.OnPlayGame:RandVal",
				goutils.Err(err))

			return "", err
		}

		cd.targetVal = ival.Int()
	}

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
						for i := 0; i < len(npos)/2; i++ {
							multi, err := chgSymbolVals.GetTarget(cd, plugin)
							if err != nil {
								goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
									goutils.Err(err))

								return "", nil
							}

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
					} else if chgSymbolVals.Config.Type == CSVTypeAdd {
						for i := 0; i < len(npos)/2; i++ {
							off, err := chgSymbolVals.GetTarget(cd, plugin)
							if err != nil {
								goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
									goutils.Err(err))

								return "", nil
							}

							if nos == os {
								nos = os.CloneEx(gameProp.PoolScene)
							}

							nos.Arr[npos[i*2]][npos[i*2+1]] += off

							if nos.Arr[npos[i*2]][npos[i*2+1]] > chgSymbolVals.Config.MaxVal {
								nos.Arr[npos[i*2]][npos[i*2+1]] = chgSymbolVals.Config.MaxVal
							}

							if !gIsReleaseMode {
								cd.AddPos(npos[i*2], npos[i*2+1])
							}
						}
					} else if chgSymbolVals.Config.Type == CSVTypeSet {
						for i := 0; i < len(npos)/2; i++ {
							val, err := chgSymbolVals.GetTarget(cd, plugin)
							if err != nil {
								goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
									goutils.Err(err))

								return "", nil
							}

							if nos == os {
								nos = os.CloneEx(gameProp.PoolScene)
							}

							nos.Arr[npos[i*2]][npos[i*2+1]] = val

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
							for i := 0; i < len(npos)/2; i++ {
								if nos == os {
									nos = os.CloneEx(gameProp.PoolScene)
								}

								multi, err := chgSymbolVals.GetTarget(cd, plugin)
								if err != nil {
									goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
										goutils.Err(err))

									return "", nil
								}

								nos.Arr[npos[i*2]][npos[i*2+1]] *= multi

								if nos.Arr[npos[i*2]][npos[i*2+1]] > chgSymbolVals.Config.MaxVal {
									nos.Arr[npos[i*2]][npos[i*2+1]] = chgSymbolVals.Config.MaxVal
								}

								if !gIsReleaseMode {
									cd.AddPos(npos[i*2], npos[i*2+1])
								}
							}
						} else if chgSymbolVals.Config.Type == CSVTypeAdd {
							for i := 0; i < len(npos)/2; i++ {
								if nos == os {
									nos = os.CloneEx(gameProp.PoolScene)
								}

								off, err := chgSymbolVals.GetTarget(cd, plugin)
								if err != nil {
									goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
										goutils.Err(err))

									return "", nil
								}

								nos.Arr[npos[i*2]][npos[i*2+1]] += off

								if nos.Arr[npos[i*2]][npos[i*2+1]] > chgSymbolVals.Config.MaxVal {
									nos.Arr[npos[i*2]][npos[i*2+1]] = chgSymbolVals.Config.MaxVal
								}

								if !gIsReleaseMode {
									cd.AddPos(npos[i*2], npos[i*2+1])
								}
							}
						} else if chgSymbolVals.Config.Type == CSVTypeSet {
							for i := 0; i < len(npos)/2; i++ {
								if nos == os {
									nos = os.CloneEx(gameProp.PoolScene)
								}

								val, err := chgSymbolVals.GetTarget(cd, plugin)
								if err != nil {
									goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
										goutils.Err(err))

									return "", nil
								}

								nos.Arr[npos[i*2]][npos[i*2+1]] = val

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
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				y := chgSymbolVals.Config.Row

				for x := 0; x < os.Width; x++ {
					multi, err := chgSymbolVals.GetTarget(cd, plugin)
					if err != nil {
						goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
							goutils.Err(err))

						return "", nil
					}

					nos.Arr[x][y] *= multi
					if nos.Arr[x][y] > chgSymbolVals.Config.MaxVal {
						nos.Arr[x][y] = chgSymbolVals.Config.MaxVal
					}

					if !gIsReleaseMode {
						cd.AddPos(x, y)
					}
				}

			} else if chgSymbolVals.Config.Type == CSVTypeAdd {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				y := chgSymbolVals.Config.Row

				for x := 0; x < os.Width; x++ {
					off, err := chgSymbolVals.GetTarget(cd, plugin)
					if err != nil {
						goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
							goutils.Err(err))

						return "", nil
					}

					nos.Arr[x][y] += off
					if nos.Arr[x][y] > chgSymbolVals.Config.MaxVal {
						nos.Arr[x][y] = chgSymbolVals.Config.MaxVal
					}

					if !gIsReleaseMode {
						cd.AddPos(x, y)
					}
				}

			} else if chgSymbolVals.Config.Type == CSVTypeSet {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				y := chgSymbolVals.Config.Row

				for x := 0; x < os.Width; x++ {
					val, err := chgSymbolVals.GetTarget(cd, plugin)
					if err != nil {
						goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
							goutils.Err(err))

						return "", nil
					}

					nos.Arr[x][y] = val

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
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				x := chgSymbolVals.Config.Column

				for y := 0; y < os.Height; y++ {
					multi, err := chgSymbolVals.GetTarget(cd, plugin)
					if err != nil {
						goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
							goutils.Err(err))

						return "", nil
					}

					nos.Arr[x][y] *= multi
					if nos.Arr[x][y] > chgSymbolVals.Config.MaxVal {
						nos.Arr[x][y] = chgSymbolVals.Config.MaxVal
					}

					if !gIsReleaseMode {
						cd.AddPos(x, y)
					}
				}
			} else if chgSymbolVals.Config.Type == CSVTypeAdd {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				x := chgSymbolVals.Config.Column

				for y := 0; y < os.Height; y++ {
					off, err := chgSymbolVals.GetTarget(cd, plugin)
					if err != nil {
						goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
							goutils.Err(err))

						return "", nil
					}

					nos.Arr[x][y] += off
					if nos.Arr[x][y] > chgSymbolVals.Config.MaxVal {
						nos.Arr[x][y] = chgSymbolVals.Config.MaxVal
					}

					if !gIsReleaseMode {
						cd.AddPos(x, y)
					}
				}
			} else if chgSymbolVals.Config.Type == CSVTypeSet {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				x := chgSymbolVals.Config.Column

				for y := 0; y < os.Height; y++ {
					val, err := chgSymbolVals.GetTarget(cd, plugin)
					if err != nil {
						goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
							goutils.Err(err))

						return "", nil
					}

					nos.Arr[x][y] = val

					if !gIsReleaseMode {
						cd.AddPos(x, y)
					}
				}
			}
		} else if chgSymbolVals.Config.SourceType == CSVSTypeAll {
			if chgSymbolVals.Config.Type == CSVTypeInc {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				for x := range os.Arr {
					for y := 0; y < os.Height; y++ {
						if nos.Arr[x][y] < chgSymbolVals.Config.MaxVal {
							nos.Arr[x][y]++

							if !gIsReleaseMode {
								cd.AddPos(x, y)
							}
						}
					}
				}
			} else if chgSymbolVals.Config.Type == CSVTypeDec {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				for x := range os.Arr {
					for y := 0; y < os.Height; y++ {
						if nos.Arr[x][y] > chgSymbolVals.Config.MinVal {
							nos.Arr[x][y]--

							if !gIsReleaseMode {
								cd.AddPos(x, y)
							}
						}
					}
				}
			} else if chgSymbolVals.Config.Type == CSVTypeMul {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				for x := range os.Arr {
					for y := 0; y < os.Height; y++ {
						multi, err := chgSymbolVals.GetTarget(cd, plugin)
						if err != nil {
							goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
								goutils.Err(err))

							return "", nil
						}

						nos.Arr[x][y] *= multi
						if nos.Arr[x][y] > chgSymbolVals.Config.MaxVal {
							nos.Arr[x][y] = chgSymbolVals.Config.MaxVal
						}

						if !gIsReleaseMode {
							cd.AddPos(x, y)
						}
					}
				}
			} else if chgSymbolVals.Config.Type == CSVTypeAdd {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				for x := range os.Arr {
					for y := 0; y < os.Height; y++ {
						off, err := chgSymbolVals.GetTarget(cd, plugin)
						if err != nil {
							goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
								goutils.Err(err))

							return "", nil
						}

						nos.Arr[x][y] += off
						if nos.Arr[x][y] > chgSymbolVals.Config.MaxVal {
							nos.Arr[x][y] = chgSymbolVals.Config.MaxVal
						}

						if !gIsReleaseMode {
							cd.AddPos(x, y)
						}
					}
				}
			} else if chgSymbolVals.Config.Type == CSVTypeSet {
				if nos == os {
					nos = os.CloneEx(gameProp.PoolScene)
				}

				for x := range os.Arr {
					for y := 0; y < os.Height; y++ {
						val, err := chgSymbolVals.GetTarget(cd, plugin)
						if err != nil {
							goutils.Error("ChgSymbolVals.OnPlayGame:GetTarget",
								goutils.Err(err))

							return "", nil
						}

						nos.Arr[x][y] = val

						if !gIsReleaseMode {
							cd.AddPos(x, y)
						}
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
// "targetType": "eachWeight",
// "number": 0,
// "targetWeight": "hotzone"
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
	TargetType          string   `json:"targetType"`
	Number              int      `json:"number"`
	TargetWeight        string   `json:"targetWeight"`
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
		Number:              jcfg.Number,
		StrTargetType:       jcfg.TargetType,
		TargetWeight:        jcfg.TargetWeight,
	}

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
