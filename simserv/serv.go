package simserv

import (
	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
	sgc7pbutils "github.com/zhs007/slotsgamecore7/pbutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.uber.org/zap"
)

// BasicURL - basic url
const BasicURL = "/game"

// Serv -
type Serv struct {
	*sgc7http.Serv
	Service IService
	Cfg     *Config
}

// NewServ - new a serv
func NewServ(service IService, cfg *Config) *Serv {
	s := &Serv{
		sgc7http.NewServ(cfg.BindAddr, cfg.IsDebugMode),
		service,
		cfg,
	}

	s.RegHandle(goutils.AppendString(BasicURL, "/config"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsGet() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ret := s.Service.GetGame().GetConfig()
			if ret == nil {
				s.SetStringResponse(ctx, "{}")
			} else {
				s.SetResponse(ctx, ret)
			}
		})

	s.RegHandle(goutils.AppendString(BasicURL, "/initialize"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsGet() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ps := s.Service.GetGame().Initialize()
			if ps == nil {
				s.SetStringResponse(ctx, "{\"playerStatePublic\":{},\"playerStatePrivate\":{}}")

				return
			}

			pbps, err := s.Service.BuildPBPlayerState(ps)
			if err != nil {
				goutils.Warn("gatiserv.Serv.initialize:BuildPBPlayerState",
					zap.Error(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)
			}

			s.SetResponse(ctx, pbps)
		})

	s.RegHandle(goutils.AppendString(BasicURL, "/play"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsPost() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			params := &sgc7pb.RequestPlay{}
			err := s.ParseBody(ctx, params)
			if err != nil {
				goutils.Warn("gatiserv.Serv.play:ParseBody",
					zap.Error(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ret, err := s.onPlay(params)
			if err != nil {
				goutils.Warn("gatiserv.Serv.play:Play",
					zap.Error(err))

				if err == sgc7game.ErrInvalidStake {
					s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

					return
				}

				s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

				return
			}

			ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
			ctx.Response.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie")

			s.SetPBResponse(ctx, ret)
		})

	return s
}

// ProcCheat - process cheat
func (serv *Serv) ProcCheat(plugin sgc7plugin.IPlugin, cheat string) error {
	if cheat != "" {
		str := goutils.AppendString("[", cheat, "]")

		rngs := []int{}
		err := sonic.Unmarshal([]byte(str), &rngs)
		if err != nil {
			return err
		}

		plugin.SetCache(rngs)
	}

	return nil
}

// Play - play game
func (serv *Serv) onPlay(req *sgc7pb.RequestPlay) (*sgc7pb.ReplyPlay, error) {
	ips := serv.Service.GetGame().NewPlayerState()
	if req.PlayerState != nil {
		err := serv.Service.BuildPlayerStateFromPB(ips, req.PlayerState)
		if err != nil {
			goutils.Error("BasicService.onPlay:BuildPlayerStateFromPB",
				zap.Error(err))

			return nil, err
		}
	}

	plugin := serv.Service.GetGame().NewPlugin()
	defer serv.Service.GetGame().FreePlugin(plugin)

	serv.ProcCheat(plugin, req.Cheat)

	stake := sgc7pbutils.BuildStake(req.Stake)
	err := serv.Service.GetGame().CheckStake(stake)
	if err != nil {
		goutils.Error("BasicService.onPlay:CheckStake",
			goutils.JSON("stake", stake),
			zap.Error(err))

		return nil, err
	}

	results := []*sgc7game.PlayResult{}
	gameData := serv.Service.GetGame().NewGameData(stake)

	cmd := req.Command

	for {
		if cmd == "" {
			cmd = "SPIN"
		}

		pr, err := serv.Service.GetGame().Play(plugin, cmd, req.ClientParams, ips, stake, results, gameData)
		if err != nil {
			goutils.Error("BasicService.onPlay:Play",
				zap.Int("results", len(results)),
				zap.Error(err))

			return nil, err
		}

		if pr == nil {
			break
		}

		results = append(results, pr)
		if pr.IsFinish {
			break
		}

		if pr.IsWait {
			break
		}

		if len(pr.NextCmds) > 0 {
			cmd = pr.NextCmds[0]
		} else {
			cmd = ""
		}
	}

	pr := &sgc7pb.ReplyPlay{
		RandomNumbers: sgc7pbutils.BuildPBRngs(plugin.GetUsedRngs()),
	}

	ps, err := serv.Service.BuildPBPlayerState(ips)
	if err != nil {
		goutils.Error("BasicService.onPlay:BuildPlayerState",
			zap.Error(err))

		return nil, err
	}

	pr.PlayerState = ps

	if len(results) > 0 {
		AddPlayResult(serv.Service, pr, results)

		lastr := results[len(results)-1]

		pr.Finished = lastr.IsFinish
		pr.NextCommands = lastr.NextCmds
	}

	return pr, nil
}
