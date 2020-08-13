package b_controller

import (
	"context"
	"github.com/djghostghost/go-beebasic/b_globals"
	"github.com/siddontang/go/bson"
	"time"
)

func (t *BasicController) BPrepare() {
	t.startTime = time.Now().UnixNano()

	t.Ctx.Input.SetData("_____t", t)
	// 设置request生命周期的参数
	ctx := context.TODO()
	ctx = b_globals.WithRequestID(ctx, bson.NewObjectId().Hex())
	t.ReqCtx = ctx
}

func (t *BasicController) BClearCtx()  {
	// 删掉reqId相关资源
	b_globals.RemoveOrmer(t.ReqCtx)
	t.ReqCtx.Done()
}
