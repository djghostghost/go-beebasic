package b_controller

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/djghostghost/go-beebasic/b_error"
	"runtime"
	"strconv"
	"time"
)

func CustomPanicRecover(ctx *context.Context) {
	if err := recover(); err != nil {
		t := ctx.Input.GetData("_____t").(*BasicController)
		if err == web.ErrAbort {
			t.Finish()
			return
		}
		if !web.BConfig.RecoverPanic {
			t.Finish()
			panic(err)
		}

		switch err.(type) {
		case *b_error.BizError:

			bizE := err.(*b_error.BizError)

			handleBizError(ctx, bizE)
			t.Ctx.Input.SetData("___status", bizE.Code)
			t.Finish()
			return
		case *b_error.SystemError:
			systemError := err.(*b_error.SystemError)

			logs.Error("System Error:", systemError.Error())
			t.Abort(strconv.Itoa(systemError.Code))
			return
		}

		var stack string
		stack = stack + fmt.Sprintln(fmt.Sprintf("the request url is: %s", ctx.Input.URL()))
		stack = stack + fmt.Sprintln(fmt.Sprintf("Handler crashed with error: %s", err))
		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
		}
		logs.Critical(stack)
		t.Finish()
		if ctx.Output.Status != 0 {
			ctx.ResponseWriter.WriteHeader(ctx.Output.Status)
		} else {
			ctx.ResponseWriter.WriteHeader(500)
		}
	}
}

func handleBizError(ctx *context.Context, bizE *b_error.BizError) {
	logs.Warn("[BIZ] biz exception: code", bizE.Code, " message ", bizE.Message)
	rd := &ResponseData{
		Ret:        bizE.Code,
		Message:    bizE.Message,
		ServerTime: time.Now().Unix(),
	}
	hasIndent := web.BConfig.RunMode != web.PROD
	jsonpCallback := ctx.Request.Form.Get("callback")
	if jsonpCallback != "" {
		ctx.Output.JSONP(rd, hasIndent)
	} else {
		ctx.Output.JSON(rd, hasIndent, false)
	}

}
