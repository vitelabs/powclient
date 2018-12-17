package gpu

import (
	"github.com/gin-gonic/gin"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/powclient/log15"
	"github.com/vitelabs/powclient/metrics"
	"github.com/vitelabs/powclient/service/context"
	"github.com/vitelabs/powclient/util"
	"strconv"
	"time"
)

const (
	FullThreshold = 0xffffffc000000000
)

var log = log15.New("module", "service_gpu")

func WorkDetail(c *gin.Context) {
	var tagsNames []string
	tagsNames = append(tagsNames, "pow", "gpu", "WorkDetail")
	defer metrics.TimeConsuming(tagsNames, time.Now())

	var generateContext context.GenerateContext
	if err := c.Bind(&generateContext); err != nil {
		util.RespondError(c, 400, err)
		return
	}

	if len([]byte(generateContext.DataHash)) < types.HashSize {
		util.RespondFailed(c, util.ErrLengthNotEnough.Code, util.ErrLengthNotEnough.ErrMsg, "")
		return
	}

	var difficulty string
	if generateContext.Threshold == nil {
		difficulty = strconv.FormatUint(FullThreshold, 16)
	} else {
		difficulty = *generateContext.Threshold
	}
	work, err := GenerateWork(generateContext.DataHash, difficulty)
	if err != nil {
		util.RespondFailed(c, util.ErrServerPostFailed.Code, err, "")
		return
	}
	generateResult := &context.GenerateResult{
		Work: *work,
	}
	util.RespondSuccess(c, generateResult, "")
	return
}

func VaildDetail(c *gin.Context) {
	var tagsNames []string
	tagsNames = append(tagsNames, "pow", "gpu", "VaildDetails")
	defer metrics.TimeConsuming(tagsNames, time.Now())

	var validateContext context.ValidateContext
	if err := c.Bind(&validateContext); err != nil {
		util.RespondError(c, 400, err)
		return
	}

	if len([]byte(validateContext.Work)) < 8 {
		util.RespondFailed(c, util.ErrLengthNotEnough.Code, util.ErrLengthNotEnough.ErrMsg, "")
		return
	}
	if len([]byte(validateContext.DataHash)) < types.HashSize {
		util.RespondFailed(c, util.ErrLengthNotEnough.Code, util.ErrLengthNotEnough.ErrMsg, "")
		return
	}

	var difficulty string
	if validateContext.Threshold == nil {
		difficulty = strconv.FormatUint(FullThreshold, 16)
	} else {
		difficulty = *validateContext.Threshold
	}

	vaild, err := VaildateWork(validateContext.DataHash, difficulty, validateContext.Work)
	if err != nil {
		util.RespondFailed(c, util.ErrServerPostFailed.Code, util.ErrServerPostFailed.ErrMsg, "")
		return
	} else {
		validateResult := &context.ValidateResult{}
		if vaild {
			validateResult.Valid = "1"
		} else {
			validateResult.Valid = "0"
		}
		util.RespondSuccess(c, validateResult, "")
		return
	}
}

func CancelDetail(c *gin.Context) {
	var tagsNames []string
	tagsNames = append(tagsNames, "pow", "gpu", "CancelDetail")
	defer metrics.TimeConsuming(tagsNames, time.Now())

	var cancelContext context.CancelContext
	if err := c.Bind(&cancelContext); err != nil {
		util.RespondError(c, 400, err)
		return
	}

	err := CancelWork(cancelContext.DataHash)
	if err != nil {
		util.RespondFailed(c, util.ErrServerPostFailed.Code, err, "")
		return
	} else {
		util.RespondSuccess(c, &context.CancelResult{}, "")
		return
	}
}
