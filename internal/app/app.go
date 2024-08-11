package app

import (
	"github.com/atopos31/code-sandbox/internal/coder"
	"github.com/atopos31/code-sandbox/internal/sandbox"
	"github.com/atopos31/code-sandbox/pkg/model"
	"github.com/gin-gonic/gin"
)

type App struct {
	sandboxPool *sandbox.SandboxPool
	coders      map[string]func() coder.Coder
	server      *gin.Engine
}

func New(*sandbox.SandboxPool) *App {
	server := gin.Default()
	coders := map[string]func() coder.Coder{
		"go":     coder.NewGOCoder,
		"c":      coder.NewCCoder,
		"cpp":    coder.NewCPPCoder,
		"java":   coder.NewJavaCoder,
		"python": coder.NewPythonCoder,
	}
	return &App{
		sandboxPool: sandbox.NewSandboxPool(10),
		coders:      coders,
		server:      server,
	}
}

func (a *App) Run() {
	a.server.POST("/run", func(c *gin.Context) {
		req := new(model.CodeRequest)
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(400, &model.CodeResponse{
				Code: 400,
				Msg:  err.Error(),
			})
			return
		}
		coder := a.coders[req.Language]()
		defer coder.Clean()
		sandbox, err := a.sandboxPool.GetSandbox()
		defer a.sandboxPool.ReleaseSandbox(sandbox)
		if err != nil {
			c.JSON(500, &model.CodeResponse{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}
		coder.SetSandbox(sandbox)
		meta, err := coder.Build(req.Code)
		if err != nil {
			c.JSON(500, &model.CodeResponse{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}

		if meta.Status != "" {
			c.JSON(400, &model.CodeResponse{
				Code: 400,
				Msg:  meta.Stderr,
			})
			return
		}
		var metas []*model.CodeMETA
		for _, stdin := range req.Stdin {
			meta, _ = coder.Run(req.MaxTime, req.MaxMem, stdin)
			metas = append(metas, meta)
		}

		c.JSON(200, &model.CodeResponse{
			Code: 200,
			Meta: metas,
		})

	})
	a.server.Run(":6758")
}
