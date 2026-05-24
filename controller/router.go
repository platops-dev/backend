package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var Router router

type router struct{}

func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/testapi", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "testapi success!",
			"data": nil,
		})
	}).
		//pod操作
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pods/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pods/delete", Pod.DeletePod).
		PUT("/api/k8s/pods/update", Pod.UpdatePod).
		GET("/api/k8s/pods/container", Pod.GetPodContainer).
		GET("/api/k8s/pods/log", Pod.GetPodLog).
		GET("/api/k8s/pods/numnp", Pod.GetPodNumPerNp).
		//deployment操作
		GET("/api/k8s/deployments", Deployment.GetDeployments).
		GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail).
		PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment).
		DELETE("/api/k8s/deployment/delete", Deployment.DeleteDeployment).
		POST("/api/k8s/deployment/create", Deployment.CreateDeployment).
		PUT("/api/k8s/deployment/restart", Deployment.RestartDeployment).
		PUT("/api/k8s/deployment/update", Deployment.UpdateDeployment).
		GET("/api/k8s/deployment/numnp", Deployment.GetDeployNumPerNP).
		//statefulset操作
		GET("/api/k8s/statefulsets", StatefulSet.GetStatefulSets).
		GET("/api/k8s/statefulset/detail", StatefulSet.GetStatefulSetDetail).
		PUT("/api/k8s/statefulset/scale", StatefulSet.ScaleStatefulSet).
		POST("/api/k8s/statefulset/create", StatefulSet.CreateStatefulSet).
		DELETE("/api/k8s/statefulset/delete", StatefulSet.DeleteStatefulSet).
		PUT("/api/k8s/statefulset/restart", StatefulSet.RestartStatefulSet).
		PUT("/api/k8s/statefulset/update", StatefulSet.UpdateStatefulSet).
		GET("/api/k8s/statefulset/numnp", StatefulSet.GetStatefulSetsNumPerNp).
		//daemonset操作
		GET("/api/k8s/daemonsets", DaemonSet.GetDaemonSets).
		GET("/api/k8s/daemonset/detail", DaemonSet.GetDaemonSetDetail).
		POST("/api/k8s/daemonset/create", DaemonSet.CreateDaemonSet).
		DELETE("/api/k8s/daemonset/delete", DaemonSet.DeleteDaemonSet).
		PUT("/api/k8s/daemonset/restart", DaemonSet.RestartDaemonSet).
		PUT("/api/k8s/daemonset/update", DaemonSet.UpdateDaemonSet).
		GET("/api/k8s/daemonset/numnp", DaemonSet.GetDaemonSetsNumPerNp).
		//集群级别-node操作
		GET("/api/k8s/nodes", K8sNode.GetK8sNodes).
		GET("/api/k8s/node/detail", K8sNode.GetK8sNodeDetail).
		//集群级别-namespace操作
		GET("/api/k8s/namespaces", Namepsace.GetNamespaces).
		GET("/api/k8s/namespace/detail", Namepsace.GetNamespaceDetail).
		DELETE("/api/k8s/namespace/delete", Namepsace.DeleteNamespace).
		//集群级别-persistentvolume
		GET("/api/k8s/persistentvolumes", PersistentVolume.GetPersistentVolumes).
		GET("/api/k8s/persistentvolume/detail", PersistentVolume.GetPersistentVolumeDetail).
		DELETE("/api/k8s/persistentvolume/delete", PersistentVolume.DeletePersistentVolume).
		//service操作
		GET("/api/k8s/services", K8sService.GetK8sServices).
		GET("/api/k8s/service/detail", K8sService.GetK8sServiceDetail).
		DELETE("/api/k8s/service/delete", K8sService.DeleteK8sService).
		POST("/api/k8s/service/create", K8sService.CreateService).
		PUT("/api/k8s/service/update", K8sService.UpdateK8sService).
		//Ingress操作
		GET("/api/k8s/ingress", Ingress.GetIngress).
		GET("/api/k8s/ingress/detail", Ingress.GetIngressDetail).
		DELETE("/api/k8s/ingress/delete", Ingress.DeleteIngress).
		POST("/api/k8s/ingress/create", Ingress.CreateIngress).
		PUT("/api/k8s/ingress/update", Ingress.UpdateIngress).
		//ConfigMap操作
		GET("/api/k8s/configmaps", ConfigMap.GetConfigMaps).
		GET("/api/k8s/configmap/detail", ConfigMap.GetConfigMapDetail).
		DELETE("/api/k8s/configmap/delete", ConfigMap.DeleteConfigMap).
		PUT("/api/k8s/configmap/update", ConfigMap.UpdateConfigMap).
		//Secret操作
		GET("/api/k8s/secrets", Secret.GetSecrets).
		GET("/api/k8s/secret/detail", Secret.GetSecretDetail).
		DELETE("/api/k8s/secret/delete", Secret.DeleteSecret).
		PUT("/api/k8s/secret/update", Secret.UpdateSecret).
		//PersistentVolumeClaim操作
		GET("/api/k8s/persistentvolumeclaims", PersistentVolumeClaim.PersistentVolumeClaims).
		GET("/api/k8s/persistentvolumeclaim/detail", PersistentVolumeClaim.GetPersistentVolumeClaimDetail).
		DELETE("/api/k8s/persistentvolumeclaim/delete", PersistentVolumeClaim.DeletePersistentVolumeClaim).
		PUT("/api/k8s/persistentvolumeclaim/update", PersistentVolumeClaim.UpdatePersistentVolumeClaim).
		//Hpa操作
		GET("/api/k8s/hpas", K8sHpa.GetK8sHpas).
		GET("/api/k8s/hpa/detail", K8sHpa.GetK8sHpaDetail).
		DELETE("/api/k8s/hpa/delete", K8sHpa.DeleteK8sHpa).
		POST("/api/k8s/hpa/create", K8sHpa.CreateK8sHpa).
		PUT("/api/k8s/hpa/update", K8sHpa.UpdateK8sHpa)
}
