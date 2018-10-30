package service

import (
	"github.com/gin-gonic/gin"
	"github.com/vitelabs/go-vite/common/types"
	"math/big"
	"powClient/service/context"
	"powClient/util"
	"strconv"
)

const (
	FullThreshold = 0xffffffc000000000
)

var DefaultDifficulty = new(big.Int).SetUint64(FullThreshold)

func WorkDetail(c *gin.Context) {
	var generateContext context.GenerateContext

	if err := c.Bind(&generateContext); err != nil {
		util.RespondError(c, 400, err)
		return
	}
	dataHash := []byte(generateContext.DataHash)
	if len(dataHash) != types.HashSize {
		util.RespondFailed(c, util.ErrLengthOfHashIncorrect.Code, util.ErrLengthOfHashIncorrect.ErrMsg, "")
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
	if len(work) < 8 {
		util.RespondFailed(c, util.ErrLengthNotEnough.Code, util.ErrLengthNotEnough.ErrMsg, "")
		return
	}
	generateResult := &context.GenerateResult{
		Work: work[:8],
	}
	util.RespondSuccess(c, generateResult, "")
	return
}

func VaildDetail(c *gin.Context) {
	var validateContext context.ValidateContext
	if err := c.Bind(&validateContext); err != nil {
		util.RespondError(c, 400, err)
		return
	}
	work := []byte(validateContext.Work)
	if len(work) < 8 {
		util.RespondFailed(c, util.ErrLengthNotEnough.Code, util.ErrLengthNotEnough.ErrMsg, "")
		return
	}
	dataHash := []byte(validateContext.DataHash)
	if len(dataHash) != types.HashSize {
		util.RespondFailed(c, util.ErrLengthOfHashIncorrect.Code, util.ErrLengthOfHashIncorrect.ErrMsg, "")
		return
	}

	var difficulty string
	if validateContext.Threshold == nil {
		difficulty = strconv.FormatUint(FullThreshold, 16)
	} else {
		difficulty = *validateContext.Threshold
	}
	vaild, err := VaildateWork(validateContext.DataHash, difficulty, work[0:8])
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
	var cancelContext context.CancelContext
	if err := c.Bind(&cancelContext); err != nil {
		util.RespondError(c, 400, err)
		return
	}
	dataHash := []byte(cancelContext.DataHash)
	if len(dataHash) != types.HashSize {
		util.RespondFailed(c, util.ErrLengthOfHashIncorrect.Code, util.ErrLengthOfHashIncorrect.ErrMsg, "")
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
