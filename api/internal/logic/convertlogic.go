package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"net/url"
	"path"
	"shortner/api/internal/svc"
	"shortner/api/internal/types"
	"shortner/api/pkg/base62"
	"shortner/api/pkg/connect"
	"shortner/api/pkg/md5"
	"shortner/model"
	"time"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 转链业务逻辑
	//1.校验url是否有效
	//1.1 req不能为空
	//在函数入口前使用validate校验

	//1.2 传入的url必须应该能ping通，而不是个无效url
	ok := connect.Get(req.LongUrl)
	if !ok {
		logx.Error("Convert:Get:invalid url")
		return nil, fmt.Errorf("invaild url")
	}

	//1.3 传入url是否转链过
	//1.3.1 将longUrl转为md5
	urlMd5 := md5.Sum([]byte(req.LongUrl))

	//1.3.2 将md5值拿去查数据库中是否存在
	One, err := l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: urlMd5, Valid: true})
	if err != sqlx.ErrNotFound {
		if err == nil {

			return nil, fmt.Errorf("该LongUrl已经转过，结果为：%s", l.svcCtx.ShortDoamin+"/"+One.Surl.String)
		}
		logx.Errorf("Convert:find one by Md5 failed:err:%v", err.Error())
		return nil, errors.New("undefine err")
	}

	//1.4 传入url不能是短链，不能循环转链
	myUrl, err := url.Parse(req.LongUrl)
	if err != nil {
		logx.Error("url parse failed,err:", err.Error())
		return nil, err
	}

	// 拿到url的最后的路径string
	basePath := path.Base(myUrl.Path)
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{
		String: basePath,
		Valid:  true,
	})
	if err != sqlx.ErrNotFound {
		if err == nil {
			logx.Errorf("该链接是短链")
			return nil, errors.New("该链接是短链")
		}
		return nil, err
	}
	var seq uint64
	var short string
	for {
		//2.取号 基于mysql实现的发号器
		seq, err = l.svcCtx.Sequence.Next() //实现了mysql自增id取号器、redis取号器
		if err != nil {
			logx.Errorf("取号失败，err:", err)
			return nil, errors.New("取号失败")
		}
		//3.转链
		//3.1安全性：为了避免被人恶意请求查看我的发号器使用情况，将base62的str打乱顺序
		//3.2 黑名单 避免输出短链出现敏感词或者health、version这些关键词
		short = base62.GetBase62(seq)
		if _, ok := l.svcCtx.ShortBlackMap[short]; !ok {
			break //不存在黑名单就不用重新生成短链，退出循环
		}
		logx.Infof("生成了一次敏感词短链:", short)
	}

	//4.存入数据库
	if _, err := l.svcCtx.ShortUrlModel.Insert(
		l.ctx,
		&model.ShortUrlMap{
			Surl:     sql.NullString{String: short, Valid: true},
			Lurl:     sql.NullString{String: req.LongUrl, Valid: true},
			Md5:      sql.NullString{String: urlMd5, Valid: true},
			CreateBy: "Jayb",
			CreateAt: time.Now(),
		},
	); err != nil {
		logx.Errorf("insert into model.ShortUrlMap failed:", err)
		return nil, err
	}

	//将存入数据库中的短链放入布隆过滤器中
	err = l.svcCtx.Filter.Add([]byte(short))
	if err != nil {
		logx.Errorf("insert into bloom filter failed:", err)
		return nil, err
	}

	//5.返回响应
	ShortUrl := l.svcCtx.ShortDoamin + "/" + short
	return &types.ConvertResponse{ShortUrl: ShortUrl}, nil
}
