package handler

import (
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"shortner/api/internal/logic"
	"shortner/api/internal/svc"
	"shortner/api/internal/types"
)

func ShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShowRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//参数校验规则
		if err := validator.New().StructCtx(r.Context(), &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewShowLogic(r.Context(), svcCtx)
		resp, err := l.Show(&req)

		if err != nil {
			if err == logic.Err404 {
				httpx.ErrorCtx(r.Context(), w, err)
			} else {
				httpx.ErrorCtx(r.Context(), w, err)

			}
		} else {
			http.Redirect(w, r, resp.LongUrl, http.StatusFound)

			//httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
