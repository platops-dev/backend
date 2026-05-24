package controller

import (
	"dashboard-api/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var K8sService k8sService

type k8sService struct{}

func (svc *k8sService) GetK8sServices(ctx *gin.Context) {
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
	data, err := service.K8sService.GetK8sServices(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf( "获取namespace: %s 下的service列表成功", params.Namespace),
		"data": data,
	})
}

func (svc *k8sService) GetK8sServiceDetail(ctx *gin.Context)  {
	params := new(struct{
		K8sServiceName	string	`form:"k8s_service_name"`
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
	data, err := service.K8sService.GetK8sServiceDetail(params.K8sServiceName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取Namespace: %s 下 Service: %s 详情成功. ", params.Namespace, params.K8sServiceName),
		"data": data,
	})
}

func (svc *k8sService) DeleteK8sService(ctx *gin.Context)  {
	params := new(struct{
		K8sServiceName	string	`json:"k8s_service_name"`
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
	err := service.K8sService.DeleteK8sService(params.K8sServiceName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("删除Namespace: %s 下 Service: %s 成功. ", params.Namespace, params.K8sServiceName),
		"data": nil,
	})
}

func (svc *k8sService) CreateService(ctx *gin.Context)  {
	var (
		ServiceCreate = new(service.ServiceCreate)
		err error
	)
	if err := ctx.ShouldBindJSON(ServiceCreate); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err = service.K8sService.CreateService(ServiceCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("创建Namespace: %s 下 Service: %s 成功. ", ServiceCreate.Namespace, ServiceCreate.Name),
		"data": nil,
	})
}

func (svc *k8sService) UpdateK8sService(ctx *gin.Context)  {
	params := new(struct{
		Namespace	string	`json:"namespace"`
		Content		string	`json:"content"`
	})
	fmt.Println(1111)
	// PUT请求, 绑定参数方法为ctx.SouldBindJSON
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	err := service.K8sService.Updatek8sService(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("创建Namespace: %s 下 Service 成功. ", params.Namespace),
		"data": nil,
	})
}