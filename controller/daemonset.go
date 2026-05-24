package controller

import (
	"dashboard-api/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)


var DaemonSet daemonSet

type daemonSet struct {}

func (ds *daemonSet) GetDaemonSets(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		namespace  string `form:"namespace"`
		Page       int    `form:"page"`
		Limit      int    `form:"limit"`
	})

	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败. " + err.Error())
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

	data, err := service.DaemonSet.GetDaemonSets(params.FilterName, params.namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取DaemonSet 列表成功",
		"data": data,
	})
}

func (ds *daemonSet) GetDaemonSetDetail(ctx *gin.Context) {
	params := new(struct {
		DaemonSetName string `form:"daemonset_name"`
		Namespace       string `form:"namespace"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败. " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.DaemonSet.GetDaemonSetDetail(params.DaemonSetName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  fmt.Sprintf("daemonset: %s 详情获取成功", params.DaemonSetName),
		"data": data,
	})
}


func (ds *daemonSet) CreateDaemonSet(ctx *gin.Context) {
	var (
		DaemonSetCreate	= new(service.DaemonSetCreate)
		err error
	)
	if err := ctx.ShouldBindJSON(DaemonSetCreate); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.DaemonSet.CreateDaemonSet(DaemonSetCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  fmt.Sprintf("daemonset: %s 创建成功", DaemonSetCreate.Metadata.Name),
		"data": nil,
	})

}

func (ds *daemonSet) DeleteDaemonSet(ctx *gin.Context) {
	params := new(struct {
		DaemonSetName string `json:"daemonset_name"`
		Namespace       string `json:"namespace"`
	})

	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.DaemonSet.DeleteDaemonSet(params.DaemonSetName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  fmt.Sprintf("daemonset: %s 删除成功", params.DaemonSetName),
		"data": nil,
	})
}


func (ds *daemonSet) RestartDaemonSet(ctx *gin.Context)  {
	params := new(struct {
		DaemonSetName 	string `json:"daemonset_name"`
		Namespace       string `json:"namespace"`
	})

	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	err := service.DaemonSet.RestartDaemonSet(params.DaemonSetName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  fmt.Sprintf("daemonset: %s 重启成功", params.DaemonSetName),
		"data": nil,
	})
}


func (ds *daemonSet) UpdateDaemonSet(ctx *gin.Context)  {
	params := new(struct{
		Content		string	`json:"content"`
		Namespace	string	`json:"namespace"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err := service.DaemonSet.UpdateDaemonSet(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "daemonset更新成功",
		"data": nil,
	})
}

func (ds *daemonSet) GetDaemonSetsNumPerNp(ctx *gin.Context)  {
	data, err := service.DaemonSet.GetDaemonSetsNumPerNp()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取每个namespace下的daemonset成功",
		"data": data,
	})
}