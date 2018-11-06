package context

import (
	"github.com/gin-gonic/gin"
)

type GenerateResult struct {
	Work string `json:"work"`
}

type ValidateResult struct {
	Valid string `json:"valid"`
}

type CancelResult struct {
}

func (result *GenerateResult) ToResponse() gin.H {
	return gin.H{
		"work": result.Work,
	}
}

func (result *ValidateResult) ToResponse() gin.H {
	return gin.H{
		"valid": result.Valid,
	}
}

func (result *CancelResult) ToResponse() gin.H {
	return gin.H{}
}
