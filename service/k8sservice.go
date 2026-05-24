package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var K8sService k8sService

type k8sService struct{}

type K8sServiceResp struct {
	Items	[]corev1.Service	`json:"items"`
	Total	int					`json:"total"`
}


type ServiceCreate struct {
	Name			string				`json:"name"`
	Namespace		string				`json:"namespace"`
	Type			string				`json:"type"`
	ContainerPort	int32				`json:"container_port"`
	Port			int32				`json:"port"`
	NodePort		int32				`json:"node_port"`
	Label			map[string]string	`json:"label"`
	Protocol		string				`json:"protocol"`
	Selector		map[string]string	`json:"selector"`	
}

func (svc *k8sService) toCells(std []corev1.Service) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = k8sServiceCell(std[i])
	}
	return cells
}

func (svc *k8sService) fromCells(cells []DataCell) []corev1.Service {
	service := make([]corev1.Service, len(cells))
	for i := range cells {
		service[i] = corev1.Service(cells[i].(k8sServiceCell))
	}
	return service
}

func (svc *k8sService) GetK8sServices(filterName, namespace string, limit, page int) (k8sServiceResp *K8sServiceResp, err error) {
	ServiceList, err := K8s.ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取ServiceList列表失败, " + err.Error()))
		return nil, errors.New("获取ServiceList列表失败, " + err.Error())
	}

	selectableData := &DataSelector{
		GenericDataList: svc.toCells(ServiceList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page: page,
			},
		},
	}

	filtered  := selectableData.Filter()
	total	  := len(filtered.GenericDataList)
	data	  := filtered.Sort().Paginate()

	return &K8sServiceResp{
		Items: svc.fromCells(data.GenericDataList),
		Total: total,
	}, nil
}


func (svc *k8sService) GetK8sServiceDetail(k8sServiceName, namespace string) (service *corev1.Service, err error) {
	Service, err := K8s.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), k8sServiceName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Namespace: %s 下 Service: %s 详情失败, " + err.Error()), namespace, k8sServiceName)
		return nil, errors.New("获取Service列表失败, " + err.Error())
	}
	return Service, nil
}

func (svc *k8sService) CreateService(data *ServiceCreate) (err error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: data.Name,
			Namespace: data.Namespace,
			Labels: data.Label,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(data.Type),
			Ports: []corev1.ServicePort{
				{
					Name: "http",
					Port: data.Port,
					Protocol: corev1.Protocol(data.Protocol),
					TargetPort: intstr.IntOrString{
						Type: 0,
						IntVal: data.ContainerPort,
					},
				},
			},
			Selector: data.Selector,
		},
	}

	if data.NodePort != 0 && data.Type == "NodePort" {
		service.Spec.Ports[0].NodePort = data.NodePort
	}

	_, err = K8s.ClientSet.CoreV1().Services(data.Namespace).Create(context.TODO(),service, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Namespace: %s 下 Service: %s 失败, " + err.Error()), data.Namespace, service.Name)
		return errors.New("创建Service失败, " + err.Error())
	}
	return nil
}

func (svc *k8sService) DeleteK8sService(k8sServiceName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), k8sServiceName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Namespace: %s 下 Service: %s 失败, " + err.Error()), namespace, k8sServiceName)
		return errors.New("删除Service失败, " + err.Error())
	}
	return nil
}

func (svc *k8sService) Updatek8sService(namespace, content string) (err error) {
	var service = &corev1.Service{}

	err = json.Unmarshal([]byte(content), service)
	if err != nil {
		logger.Error(errors.New("JSON反序列化失败." + err.Error()))
		return nil
	}

	_, err = K8s.ClientSet.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Namespace: %s 下 Service: %s 失败, " + err.Error()), namespace, service.Name)
		return errors.New("更新Service失败, " + err.Error())
	}
	return nil
}