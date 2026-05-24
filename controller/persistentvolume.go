package controller

import (
	"errors"
	"fmt"
	"net/http"
	"dashboard-api/service/cluster"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var PersistentVolume persistentVolume

type persistentVolume struct{}

func (pv *persistentVolume) GetPersistentVolumes(ctx *gin.Context) {
	params := new(struct{
		FilterName	string	`form:"filter_name"`
		Limit		int		`form:"limit"`
		Page		int		`form:"page"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error(errors.New("Bind请求参数绑定失败. " + err.Error()))
		ctx.JSON(http .StatusInternalServerError, gin.H{
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
	data, err := cluster.PersistentVolume.GetPersistentVolumes(params.FilterName, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取PersistentVolume列表成功",
		"data": data,
	})
}

func (pv *persistentVolume) GetPersistentVolumeDetail(ctx *gin.Context)  {
	params := new(struct{
		PersistentVolumeName	string	`form:"persistent_volume_name"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error(errors.New("Bind请求参数绑定失败. " + err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := cluster.PersistentVolume.GetPersistentVolumeDetail(params.PersistentVolumeName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取PersistentVolume: %s 详情成功", params.PersistentVolumeName),
		"data": data,
	})
}


func (pv *persistentVolume) DeletePersistentVolume(ctx *gin.Context)  {
	params := new(struct{
		PersistentVolumeName	string	`json:"persistent_volume_name"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error(errors.New("ShouldBind请求参数绑定失败. " + err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := cluster.PersistentVolume.DeletePersistentVolume(params.PersistentVolumeName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("删除PersistentVolume: %s 成功", params.PersistentVolumeName),
		"data": nil,
	})
}