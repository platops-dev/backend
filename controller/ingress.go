package controller

import (
	"fmt"
	"net/http"
	"dashboard-api/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var Ingress ingress

type ingress struct{}

func (i *ingress) GetIngress(ctx *gin.Context) {
	params := new(struct{
		FilterName	string	`form:"filter_name"`
		Namespace	string	`form:"namespace"`
		Page		int		`form:"page"`
		Limit		int		`form:"limit"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
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
	data, err := service.Ingress.GetIngress(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf( "获取namespace: %s 下的 Ingress 列表成功", params.Namespace),
		"data": data,
	})
}

func (i *ingress) GetIngressDetail(ctx *gin.Context)  {
	params := new(struct{
		IngressName		string	`form:"ingress_name"`
		Namespace		string	`form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Ingress.GetIngressDetail(params.IngressName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取Namespace: %s 下 Ingress: %s 详情成功. ", params.Namespace, params.IngressName),
		"data": data,
	})
}

func (i *ingress) DeleteIngress(ctx *gin.Context)  {
	params := new(struct{
		IngressName		string	`json:"ingress_name"`
		Namespace		string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.Ingress.DeleteIngress(params.IngressName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("删除Namespace: %s 下 Ingress: %s 成功. ", params.Namespace, params.IngressName),
		"data": nil,
	})
}

func (i *ingress) CreateIngress(ctx *gin.Context)  {
	var (
		IngressCreate = new(service.IngressCreate)
		err error
	)
	if err := ctx.ShouldBindJSON(IngressCreate); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err = service.Ingress.CreateIngress(IngressCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("Namespace: %s 下的 Ingress: %s 创建成功", IngressCreate.Namespace, IngressCreate.Name),
		"data": nil,
	})
}

func (i *ingress) UpdateIngress(ctx *gin.Context)  {
	params := new(struct{
		Content			string	`json:"content"`
		Namespace		string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	if params.Content == "" || params.Namespace == "" {
		logger.Error("content/namespace 未传参, 请传参数")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "content/namespace 未传参, 请传参数",
			"data": nil,
		})
		return
	}
	err := service.Ingress.UpdateIngress(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("Namespace: %s 下的 Ingress更新成功", params.Namespace),
		"data": nil,
	})
}