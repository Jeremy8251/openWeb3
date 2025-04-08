package validate

import (
	"io"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models/request"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PoolBaseInfo struct{}

func NewPoolBaseInfo() *PoolBaseInfo {
	return &PoolBaseInfo{}
}

func (v *PoolBaseInfo) PoolBaseInfo(c *gin.Context, req *request.PoolBaseInfo) int {
	// 将 JSON 请求体解析为结构体PoolBaseInfo
	err := c.ShouldBind(req)
	if err == io.EOF {
		return statecode.ParameterEmptyErr
	} else if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			if e.Field() == "ChainId" && e.Tag() == "required" {
				return statecode.ChainIdEmpty
			}
		}
		return statecode.CommonErrServerErr
	}

	if req.ChainId != 97 && req.ChainId != 56 {
		return statecode.ChainIdErr
	}

	return statecode.CommonSuccess
}
