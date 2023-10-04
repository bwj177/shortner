package handler

import (
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shortner/api/internal/logic"
	"shortner/api/internal/svc"
	"shortner/api/internal/types"
)

func ConvertHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConvertRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//参数校验规则
		if err := validator.New().StructCtx(r.Context(), &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewConvertLogic(r.Context(), svcCtx)
		resp, err := l.Convert(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
