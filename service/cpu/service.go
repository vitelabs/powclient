package cpu

import (
	"github.com/gin-gonic/gin"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/powclient/metrics"
	"github.com/vitelabs/powclient/service/context"
	"github.com/vitelabs/powclient/util"
	"math/big"
	"time"
)

func WorkDetail(c *gin.Context) {
	var tagsNames []string
	tagsNames = append(tagsNames, "pow", "cpu", "GetPowNonce")
	defer metrics.TimeConsuming(tagsNames, time.Now())

	var generateContext context.GenerateContext
	if err := c.Bind(&generateContext); err != nil {
		util.RespondError(c, 400, err)
		return
	}

	if generateContext.Threshold == nil {
		util.RespondFailed(c, util.ErrThresholdParsingFailed.Code, util.ErrThresholdParsingFailed.ErrMsg, "")
		return
	}
	threshold, ok := new(big.Int).SetString(*generateContext.Threshold, 16)
	if !ok {
		util.RespondFailed(c, util.ErrThresholdParsingFailed.Code, util.ErrThresholdParsingFailed.ErrMsg, "")
		return
	}
	hash, err := types.HexToHash(generateContext.DataHash)
	if err != nil {
		util.RespondFailed(c, util.ErrHashParsingFailed.Code, err, "")
		return
	}

	work, err := GetPowNonce(threshold, hash)
	if err != nil {
		util.RespondFailed(c, util.ErrClientPowFailed.Code, err, "")
		return
	}
	generateResult := &context.GenerateResult{
		Work: *work,
	}
	util.RespondSuccess(c, generateResult, "")
	return
}

func CancelDetail(c *gin.Context) {
}
