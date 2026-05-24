package controller

import (
	"fmt"
	"net/http"
	"dashboard-api/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var PersistentVolumeClaim persistentvolumeClaim

type persistentvolumeClaim struct{}

func (pvc *persistentvolumeClaim) PersistentVolumeClaims(ctx *gin.Context) {
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
	if params.Limit <=0 || params.Page <=0 {
		logger.Error("Limit/Page 参数不合法或为空...")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("Limit(%d)/Page(%d) 参数不合法或为空...", params.Limit, params.Page),
			"data": nil,
		})
		return
	}
	data, err := service.PersistentVolumeClaim.GetPersistentVolumeClaims(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取Namespace: %s 下的PersistentVolumeClaim 列表成功", params.Namespace),
		"data": data,
	})
}

func (pvc *persistentvolumeClaim) GetPersistentVolumeClaimDetail(ctx *gin.Context)  {
	params := new(struct{
		PersistentvolumeClaimName	string	`form:"persistent_volume_claim_name"`
		Namespace					string	`form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数绑定失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.PersistentVolumeClaim.GetPersistentVolumeClaimDetail(params.PersistentvolumeClaimName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取Namespace: %s 下的PersistentVolumeClaim: %s 详情成功", params.Namespace, params.PersistentvolumeClaimName),
		"data": data,
	})
}


func (pvc *persistentvolumeClaim) DeletePersistentVolumeClaim(ctx *gin.Context)  {
	params := new(struct{
		PersistentvolumeClaimName	string	`json:"persistent_volume_claim_name"`
		Namespace					string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数绑定失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.PersistentVolumeClaim.DeletePersistentVolumeClaim(params.PersistentvolumeClaimName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("删除Namespace: %s 下的PersistentVolumeClaim: %s 成功", params.Namespace, params.PersistentvolumeClaimName),
		"data": nil,
	})
}

func (pvc *persistentvolumeClaim) UpdatePersistentVolumeClaim(ctx *gin.Context)  {
	params := new(struct{
		Content		string	`json:"content"`
		Namespace	string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数绑定失败," + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.PersistentVolumeClaim.UpdatePersistentVolumeClaim(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("更新Namespace: %s 下的PersistentVolumeClaim 成功", params.Namespace),
		"data": nil,
	})
}