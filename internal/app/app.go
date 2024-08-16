package app

import (
	"fmt"

	"github.com/atopos31/code-sandbox/internal/newcoder"
	"github.com/atopos31/code-sandbox/internal/sandbox"
	"github.com/atopos31/code-sandbox/pkg/model"
	"github.com/gin-gonic/gin"
)

type App struct {
	sandboxPool *sandbox.SandboxPool
	coders      map[string]newcoder.NewCoderFunc
	server      *gin.Engine
}

func New(sandboxPool *sandbox.SandboxPool) *App {
	// gin.SetMode(gin.ReleaseMode)
	server := gin.New()
	coders := map[string]newcoder.NewCoderFunc{
		"go":     newcoder.NewGoCoder,
		"c":      newcoder.NewCCoder,
		"cpp":    newcoder.NewCppCoder,
		"java":   newcoder.NewJavaCoder,
		"python": newcoder.NewPyCoder,
	}
	return &App{
		sandboxPool: sandboxPool,
		coders:      coders,
		server:      server,
	}
}

func (a *App) Run(port string) {
	a.server.POST("/run", func(c *gin.Context) {
		req := new(model.CodeRequest)
		if err := c.ShouldBindJSON(req); err != nil {
			c.JSON(200, &model.CodeResponse[any]{
				Code: 400,
				Msg:  err.Error(),
			})
			return
		}
		coder := a.coders[req.Language](req.Code)
		defer coder.Clean()
		sandboxBuild, err := a.sandboxPool.GetSandbox()

		if err != nil {
			c.JSON(200, &model.CodeResponse[any]{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}
		meta, err := coder.Build(sandboxBuild)
		a.sandboxPool.ReleaseSandbox(sandboxBuild)
		if err != nil {
			c.JSON(200, &model.CodeResponse[any]{
				Code: 500,
				Msg:  err.Error(),
			})
			return
		}

		if meta.Status != "" {
			c.JSON(200, &model.CodeResponse[model.BuildMeta]{
				Code: 400,
				Msg:  "Build Error",
				Meta: []model.BuildMeta{*meta},
			})
			return
		}

		var metas = []model.RunMeta{}

		var metac = make(chan model.RunMeta)
		for _, stdin := range req.Stdin {
			go func() {
				sandboxRun, _ := a.sandboxPool.GetSandbox()
				defer a.sandboxPool.ReleaseSandbox(sandboxRun)
				err := coder.Run(sandboxRun, stdin, req.MaxTime, req.MaxMem, metac)
				if err != nil {
					fmt.Println(err)
				}
			}()
		}

		for i := 0; i < len(req.Stdin); i++ {
			meta := <-metac
			metas = append(metas, meta)
		}

		c.JSON(200, &model.CodeResponse[model.RunMeta]{
			Code: 200,
			Meta: metas,
		})

	})
	a.server.Run(":" + port)
}
