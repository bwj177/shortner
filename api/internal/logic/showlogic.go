package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shortner/api/internal/svc"
	"shortner/api/internal/types"
)

var Err404 = errors.New("404，该短链不存在")

type ShowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowLogic {
	return &ShowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShowLogic) Show(req *types.ShowRequest) (resp *types.ShowResponse, err error) {
	// 查看短链接

	//1.根据短链查询到对应长链
	//1.0 查询前提前使用布隆过滤器来判断该短链是否存在于数据库、避免缓存穿透

	ok, err := l.svcCtx.Filter.Exists([]byte(req.ShortUrl))
	if err != nil {
		logx.Error("check exist by bloom filter failed,err:", err)
		return nil, err
	}
	if !ok { //不存在的短链
		return nil, Err404
	}
	//1.1查询数据库前增加了缓存层，go-zero封装
	// go-zero缓存自带的go-singlefilght 用于合并请求，避免缓存击穿
	u, err := l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: req.ShortUrl, Valid: true})
	if err != nil {
		if err == sqlx.ErrNotFound {
			logx.Error("收到未知短链请求,err:", err.Error())
			return nil, err
		}
		logx.Errorf("surl get lurl failed by mysql,err:", err)
		return nil, err
	}

	//2.跳转到长链的url
	//将长链返回给上一层通过302重定向跳转
	l.svcCtx.UserVisitResourceTotal.Inc(u.Lurl.String, u.Surl.String)
	return &types.ShowResponse{LongUrl: u.Lurl.String}, nil
}
