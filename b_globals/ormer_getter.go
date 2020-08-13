package b_globals

import (
	"context"
	"errors"
	"github.com/astaxie/beego/orm"
	"sync"
)

var ormerMapper sync.Map

func GetOrmer(ctx context.Context) (orm.Ormer, error) {
	reqId, exists := GetRequestID(ctx)
	if !exists {
		return nil, errors.New("get ormer error")
	}
	ormer, exists := ormerMapper.Load(reqId)
	if !exists || ormer == nil {
		ormer = orm.NewOrm()
	}
	ormerMapper.Store(reqId, ormer)
	return ormer.(orm.Ormer), nil
}

func RemoveOrmer(ctx context.Context) {
	reqId, exists := GetRequestID(ctx)
	if !exists {
		return
	}
	ormerMapper.Delete(reqId)
}
