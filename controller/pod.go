package controller

import (
	"dashboard-api/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var Pod pod

type pod struct{}

func (p *pod) GetPods(ctx *gin.Context) {
	params := new(struct{
		FilterName	string	`form:"filter_name"`
		Namespace	string	`form:"namespace"`
		Page		int		`form:"page"`
		Limit		int		`form:"limit"`
	})

	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数绑定失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}

	if 	params.Limit < 0 || params.Page < 0 {
		logger.Error("Limit/Page 参数不合法...")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Limit/Page 参数不合法...",
			"data": nil,
		})
		return 
	}

	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Pod列表成功",
		"data": data,
	})
}

func (p *pod) GetPodDetail(ctx *gin.Context) {
	params := new(struct{
		PodName		string	`form:"pod_name"`
		Namespace	string	`form:"namespace"`
	})
	
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数绑定失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPodDetail(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取Pod详情成功",
		"data": data,
	})
}


func (p *pod) DeletePod(ctx *gin.Context)  {
	params := new(struct{
		PodName		string	`json:"pod_name"`
		Namespace	string	`json:"namespace"`
	})

	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind参数绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("删除Pod %s 成功", params.PodName),
		"data": nil,
	})
}


func (p *pod) UpdatePod(ctx *gin.Context)  {
	params := new(struct{
		PodName		string	`json:"pod_name"`
		Namespace	string	`json:"namespace"`
		Content		string	`json:"content"`
	})

	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("ShouldBind参数绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Pod.UpdatePod(params.Namespace, params.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("Pod %s 更新成功", params.PodName),
		"data": nil,
	})
}

func (p *pod) GetPodContainer(ctx *gin.Context)  {
	params := new(struct{
		PodName		string	`form:"pod_name"`
		Namespace	string	`form:"namespace"`
	})

	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind参数绑定失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPodContainer(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取 Pod %s 中的容器名成功", params.PodName),
		"data": data,
	})
}


func (p *pod) GetPodLog(ctx *gin.Context)  {
	params := new(struct{
		PodName			string	`form:"pod_name"`
		Namespace		string	`form:"namespace"`
		ContainerName	string	`form:"container_name"`
	})

	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数绑定失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Pod.GetPodLog(params.ContainerName, params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("获取容器 %s 中的日志成功", params.ContainerName),
		"data": data,
	})
}


func (p *pod) GetPodNumPerNp(ctx *gin.Context)  {
	data, err := service.Pod.GetPodNumPerNp()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
			"data": nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "获取每个namespace的pod数量成功",
		"data": data,
	})
}