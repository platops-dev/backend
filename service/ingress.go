package service

import (
	"context"
	"encoding/json"
	"errors"
	_"fmt"

	"github.com/wonderivan/logger"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Ingress ingress

type ingress struct{}

type IngressResp struct {
	Items	[]nwv1.Ingress	`json:"items"`
	Total	int				`json:"total"`
}

type IngressCreate struct {
	Name		string	`json:"name"`
	Namespace	string	`json:"namespace"`
	Label		map[string]string	`json:"label"`
	Hosts		map[string][]*HttpPath	`json:"hosts"`
}

type HttpPath struct {
	Path	string				`json:"path"`
	PathType	nwv1.PathType	`json:"path_type"`
	ServiceName	string			`json:"service_name"`
	ServicePort	int32			`json:"service_port"`	
}

func (i *ingress) toCells(std []nwv1.Ingress) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = ingressCell(std[i])
	}
	return cells
}

func (i *ingress) fromCells(cells []DataCell) []nwv1.Ingress {
	ingress := make([]nwv1.Ingress, len(cells))
	for i := range cells {
		ingress[i] = nwv1.Ingress(cells[i].(ingressCell))
	}
	return ingress
}



func (i *ingress) GetIngress(filterName, namespace string, limit, page int) (ingress *IngressResp, err error) {
	IngressList, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取IngressList列表失败, " + err.Error()))
		return nil, errors.New("获取IngressList列表失败, " + err.Error())
	}

	selectableData := &DataSelector{
		GenericDataList: i.toCells(IngressList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page: page,
			},
		},
	}

	filtered := selectableData.Filter()
	total    := len(filtered.GenericDataList)
	data     := filtered.Sort().Paginate()

	return &IngressResp{
		Items: i.fromCells(data.GenericDataList),
		Total: total,
	}, nil
}

func (i *ingress) GetIngressDetail(ingressName, namespace string) (ingress *nwv1.Ingress, err error) {
	Ingress, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Ingress 详情失败, " + err.Error()))
		return nil, errors.New("获取Ingress 详情失败, " + err.Error())
	}
	return Ingress, nil
}

func (i *ingress) DeleteIngress(ingressName, namespace string) (err error) {
	err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Ingress: %s 失败, " + err.Error()), ingressName)
		return errors.New("删除 Ingress 失败, " + err.Error())
	}
	return nil
}


func (i *ingress) CreateIngress(data *IngressCreate) (err error) {
	var ingressRules []nwv1.IngressRule
	var httpIngressPaths []nwv1.HTTPIngressPath

	ingress := &nwv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: data.Name,
			Namespace: data.Namespace,
			Labels: data.Label,
		},
		Status: nwv1.IngressStatus{},
	}

	for key, value := range data.Hosts {
		ir := nwv1.IngressRule{
			Host: key,
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{Paths: nil},
			},
		}
		for _, httpPath := range value {
			hip := nwv1.HTTPIngressPath{
				Path: httpPath.Path,
				PathType: &httpPath.PathType,
				Backend: nwv1.IngressBackend{
					Service: &nwv1.IngressServiceBackend{
						Name: httpPath.ServiceName,
						Port: nwv1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
				},
			}

			httpIngressPaths = append(httpIngressPaths, hip)

		}

		ir.IngressRuleValue.HTTP.Paths = httpIngressPaths

		ingressRules = append(ingressRules, ir)
	}

	ingress.Spec.Rules = ingressRules

	_, err = K8s.ClientSet.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Namespace: %s 下的Ingress: %s 失败, " + err.Error()), data.Namespace, data.Name)
		return errors.New("创建 Ingress 失败, " + err.Error())
	}
	return nil
}

func (i *ingress) UpdateIngress(namespace, content string) (err error) {
	var ingress = &nwv1.Ingress{}

	err = json.Unmarshal([]byte(content), ingress)
	if err != nil {
		logger.Error(errors.New("JONS反序列化失败." + err.Error()))
		return errors.New("JONS反序列化失败." + err.Error())
	}

	_, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Namespace: %s 下的Ingress: %s 失败, " + err.Error()), namespace, ingress.Name)
		return errors.New("更新 Ingress 失败, " + err.Error())
	}
	return nil
}