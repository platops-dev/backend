package controller

import (
	"fmt"
	"net/http"
	"dashboard-api/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var Secret secret

type secret struct{}

func (st *secret) GetSecrets(ctx *gin.Context)  {
	params := new(struct{
		FilterName	string	`form:"filter_name"`
		Namespace	string	`form:"namespace"`
		Page		int		`form:"page"`
		Limit		int		`form:"limit"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数绑定失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	if params.Limit < 0 || params.Page < 0 {
		logger.Error("Limit/Page 参数不合法...")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Limit/Page 参数不合法...",
			"data": nil,
		})
		return
	}
	data, err := service.Secret.GetSecrets(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取Namespace: %s 下的Secrets 列表成功", params.Namespace),
		"data": data,
	})
}


func (st *secret) GetSecretDetail(ctx *gin.Context)  {
	params := new(struct{
		SecretName		string	`form:"secret_name"`
		Namespace		string	`form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数绑定失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Secret.GetSecretDetail(params.SecretName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取Namespace: %s 下的Secret %s 详情成功", params.Namespace, params.SecretName),
		"data": data,
	})
}

func (st *secret) DeleteSecret(ctx *gin.Context)  {
	params := new(struct{
		SecretName		string	`json:"secret_name"`
		Namespace		string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数绑定失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Secret.DeleteSecret(params.SecretName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("删除Namespace: %s 下的Secret %s 成功", params.Namespace, params.SecretName),
		"data": nil,
	})
}

func (st *secret) UpdateSecret(ctx *gin.Context)  {
	params := new(struct{
		Content			string	`json:"content"`
		Namespace		string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数绑定失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Secret.UpdateSecret(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("更新Namespace: %s 下的Secret 成功", params.Namespace),
		"data": nil,
	})
}