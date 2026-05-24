package controller

import (
	"fmt"
	"net/http"
	"dashboard-api/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var ConfigMap configMap

type configMap struct{}

func (cm *configMap) GetConfigMaps(ctx *gin.Context) {
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
	if params.Limit <0 || params.Page <0 {
		logger.Error("Limit/Page 参数不合法...")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Limit/Page参数不合法...",
			"data": nil,
		})
		return
	}
	data, err := service.ConfigMap.GetConfigMaps(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取Namespace: %s 下的ConfigMaps 列表成功", params.Namespace),
		"data": data,
	})
}

func (cm *configMap) GetConfigMapDetail(ctx *gin.Context)  {
	params := new(struct{
		ConfigMapName	string	`form:"configmap_name"`
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
	data, err := service.ConfigMap.GetConfigMapDetail(params.ConfigMapName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取Namespace: %s 下的ConfigMap %s 详情成功", params.Namespace, params.ConfigMapName),
		"data": data,
	})
}

func (cm *configMap) DeleteConfigMap(ctx *gin.Context)  {
	params := new(struct{
		ConfigMapName	string	`json:"configmap_name"`
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
	err := service.ConfigMap.DeleteConfigMap(params.ConfigMapName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("删除Namespace: %s 下的ConfigMap %s 成功", params.Namespace, params.ConfigMapName),
		"data": nil,
	})
}

func (cm *configMap) UpdateConfigMap(ctx *gin.Context)  {
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
	err := service.ConfigMap.UpdateConfigMap(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("更新Namespace: %s 下的ConfigMap 成功", params.Namespace),
		"data": nil,
	})
}