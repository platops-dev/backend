package service

import (
	"context"
	"encoding/json"
	"errors"
	_"fmt"
	"time"

	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	_ "k8s.io/api/core/v1"
	_ "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/util/intstr"
)

var Deployment deployment

type deployment struct{}

// 定义列表返回内容
type DeploymentsResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

// 定义DeploymentCreate 结构体, 用于创建deployment需要的参数属性的定义
type DeployCreate struct {
	Metadata metav1.ObjectMeta       `json:"metadata"`
	Spec     appsv1.DeploymentSpec   `json:"spec"`
	Status   appsv1.DeploymentStatus `json:"status"`
}

type DeployNp struct {
	Namespace string `json:"namespace"`
	DeployNum int    `json:"deployment_num"`
}

// 定义DataCell 到 Deployment 类型转换的方法
func (d *deployment) toCells(std []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = deploymentCell(std[i])
	}
	return cells
}

func (d *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return deployments
}

// 获取deployment 列表
func (d *deployment) GetDeployments(filterName, namespace string, limit, page int) (deploymentsResp *DeploymentsResp, err error) {
	deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Deployment列表失败, " + err.Error()))
		return nil, errors.New("获取Deployment列表失败, " + err.Error())
	}

	//将deploymentList中的deploment列表(Items), 放进dataselector对象中，进行排序
	selectableData := &DataSelector{
		GenericDataList: d.toCells(deploymentList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	filtered := selectableData.Filter()

	total := len(filtered.GenericDataList)

	data := filtered.Sort().Paginate()

	deploymentsResp = &DeploymentsResp{
		Items: d.fromCells(data.GenericDataList),
		Total: total,
	}
	return deploymentsResp, nil
}

// 获取deployment详情
func (d *deployment) GetDeploymentDetail(deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = K8s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Deployment详情失败. " + err.Error()))
		return nil, errors.New("获取Deployment详情失败. " + err.Error())
	}
	return deployment, nil
}

func (d *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int) (replica int32, err error) {
	scale, err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Deployment副本数信息失败. " + err.Error()))
		return 0, errors.New("获取Deployment副本数信息失败. " + err.Error())
	}

	scale.Spec.Replicas = int32(scaleNum)

	newScale, err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Deployment副本数信息失败. " + err.Error()))
		return 0, errors.New("更新Deployment副本数信息失败. " + err.Error())
	}

	return newScale.Spec.Replicas, nil
}

func (d *deployment) CreateDeployment(data *DeployCreate) (err error) {
	deployment := &appsv1.Deployment{
		ObjectMeta: data.Metadata,
		Spec:       data.Spec,
		Status:     data.Status,
	}
	//logger.Info(fmt.Sprintf("创建的 Deployment 内容: %+v", deployment))
	_, err = K8s.ClientSet.AppsV1().Deployments(data.Metadata.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Deployment失败. " + err.Error()))
		return errors.New("创建Deployment失败. " + err.Error())
	}
	return nil
}

func (d *deployment) RestartDeployment(deploymentName, namespace string) (err error) {
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]string{
						"kubectl.kubenetes.io/restartAt": metav1.Now().Format(time.RFC3339),
					},
				},
			},
		},
	}

	patchByte, err := json.Marshal(patchData)

	if err != nil {
		logger.Error(errors.New("JSON序列化失败. " + err.Error()))
		return errors.New("JSON序列化失败. " + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error(errors.New("重启Deployment失败. " + err.Error()))
		return errors.New("重启Deployment失败. " + err.Error())
	}

	return nil
}

func (d *deployment) DeleteDeployment(deploymentName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Deployment失败. " + err.Error()))
		return errors.New("删除Deployment失败. " + err.Error())
	}
	return
}

func (d *deployment) UpdateDeployment(namespace, content string) (err error) {
	var deploy = &appsv1.Deployment{}
	err = json.Unmarshal([]byte(content), deploy)
	if err != nil {
		logger.Error(errors.New("反序列化失败. " + err.Error()))
		return errors.New("反序列化失败. " + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Deployment失败. " + err.Error()))
		return errors.New("更新Deployment失败. " + err.Error())
	}
	return nil
}

func (d *deployment) GetDeployNumPerNP() (deploysNps []*DeployNp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaceList.Items {
		deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		deploysNp := &DeployNp{
			Namespace: namespace.Name,
			DeployNum: len(deploymentList.Items),
		}

		deploysNps = append(deploysNps, deploysNp)
	}
	return deploysNps, nil
}
