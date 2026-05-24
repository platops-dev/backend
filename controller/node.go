package controller

import (
	"errors"
	"fmt"
	"net/http"
	"dashboard-api/service"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var K8sNode k8sNode

type k8sNode struct{}

func (kn *k8sNode) GetK8sNodes(ctx *gin.Context) {
	params := new(struct{
		FilterName	string	`form:"filter_name"`
		Limit		int		`form:"limit"`
		Page		int		`form:"page"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error(errors.New("Bind请求参数绑定失败. " + err.Error()))
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
	data, err := service.K8sNode.GetK8sNodes(params.FilterName, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取node列表成功",
		"data": data,
	})
}

func (kn *k8sNode) GetK8sNodeDetail(ctx *gin.Context)  {
	params := new(struct{
		K8sNodeName	string	`form:"k8s_node_name"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error(errors.New("Bind请求参数绑定失败. " + err.Error()))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.K8sNode.GetK8sNodeDetail(params.K8sNodeName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取node: %s 详情成功", params.K8sNodeName),
		"data": data,
	})
}