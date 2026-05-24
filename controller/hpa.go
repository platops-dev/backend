package controller

import (
	"context"
	"net/http"
	"time"

	"dashboard-api/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var K8sHpa k8sHpa

type k8sHpa struct{}

// GetK8sHpas 获取 HPA 列表
func (h *k8sHpa) GetK8sHpas(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
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

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	data, err := service.K8sHpa.GetK8sHpas(c, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取HPA列表成功",
		"data": data,
	})
}

// GetK8sHpaDetail 获取 HPA 详情
func (h *k8sHpa) GetK8sHpaDetail(ctx *gin.Context) {
	params := new(struct {
		HpaName   string `form:"hpa_name"`
		Namespace string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	data, err := service.K8sHpa.GetK8sHpaDetail(c, params.HpaName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取HPA详情成功",
		"data": data,
	})
}

// CreateK8sHpa 创建 HPA
func (h *k8sHpa) CreateK8sHpa(ctx *gin.Context) {
	var hpaCreate service.HpaCreate
	if err := ctx.ShouldBindJSON(&hpaCreate); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 15*time.Second)
	defer cancel()

	if err := service.K8sHpa.CreateK8sHpa(c, &hpaCreate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建HPA成功",
		"data": nil,
	})
}

// DeleteK8sHpa 删除 HPA
func (h *k8sHpa) DeleteK8sHpa(ctx *gin.Context) {
	params := new(struct {
		HpaName   string `json:"hpa_name"`
		Namespace string `json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	if err := service.K8sHpa.DeleteK8sHpa(c, params.HpaName, params.Namespace); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除HPA成功",
		"data": nil,
	})
}

// UpdateK8sHpa 更新 HPA
func (h *k8sHpa) UpdateK8sHpa(ctx *gin.Context) {
	var hpaUpdate service.HpaCreate
	if err := ctx.ShouldBindJSON(&hpaUpdate); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c, cancel := context.WithTimeout(ctx.Request.Context(), 15*time.Second)
	defer cancel()

	if err := service.K8sHpa.UpdateK8sHpa(c, hpaUpdate.Name, hpaUpdate.Namespace, &hpaUpdate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新HPA成功",
		"data": nil,
	})
}