package service

import (
	"context"
	"errors"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Namepsace namespace

type namespace struct{}

type NamespaceResp struct {
	Items []corev1.Namespace	`json:"items"`
	Total	int					`json:"total"`
}

type NamepsaceCreate struct {
	Metadata	metav1.ObjectMeta		`json:"metadata"`
	Spec		corev1.NamespaceSpec	`json:"spec"`
	Status		corev1.NamespaceStatus	`json:"status"`
}


//类型转换
func (ns *namespace) toCells(std []corev1.Namespace) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = namespaceCell(std[i])
	}
	return cells
}

func (ns *namespace) fromCells(cells []DataCell) []corev1.Namespace {
	namespace := make([]corev1.Namespace, len(cells))
	for i := range cells {
		namespace[i] = corev1.Namespace(cells[i].(namespaceCell))
	}
	return namespace
}


func (ns *namespace) GetNamespaces(filterName string, limit, page int) (namespaceResp *NamespaceResp, err error) {
	NamespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取NamespaceList列表失败." + err.Error()))
		return nil, errors.New("获取NamespaceList列表失败." + err.Error())
	}

	selectableData := &DataSelector{
		GenericDataList: ns.toCells(NamespaceList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page: page,
			},
		},
	}

	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()

	namespaceResps := ns.fromCells(data.GenericDataList)

	return &NamespaceResp{
		Items: namespaceResps,
		Total: total,
	}, nil
}

func (ns *namespace) GetNamespaceDetail(namespaceName string) (namespace *corev1.Namespace, err error) {
	Namespace, err := K8s.ClientSet.CoreV1().Namespaces().Get(context.TODO(), namespaceName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Namespace: %s 详情失败." + err.Error()), namespaceName)
		return nil, errors.New("获取Namespace详情失败." + err.Error())
	}
	return Namespace, nil	
}

func (ns *namespace) CreateNamespace(data *NamepsaceCreate) (err error) {
	namespace := &corev1.Namespace{
		ObjectMeta: data.Metadata,
		Spec: data.Spec,
		Status: data.Status,
	}
	_, err = K8s.ClientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Namepsace失败. " + err.Error()))
		return errors.New("创建Namepsace失败. " + err.Error())
	}
	return nil
}


func (ns *namespace) DeleteNamespace(namespaceName string) (err error) {
	err = K8s.ClientSet.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Namespace: %s 失败." + err.Error()), namespaceName)
		return errors.New("删除Namespace失败." + err.Error())
	}
	return nil
}