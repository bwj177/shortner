package connect

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"time"
)

var client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
	Timeout: 2 * time.Second,
}

// 通过get请求看code是否为200
func Get(url string) bool {
	resp, err := client.Get(url)
	if err != nil {
		logx.Errorw("connect client get failed", logx.Field("err:", err.Error()))
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK //忽略重定向url，不允许通过
}
