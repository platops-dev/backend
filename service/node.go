package service

import (
	"context"
	"errors"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)



var K8sNode k8sNode

type k8sNode struct{}


type K8sNodeResp struct {
	Items	[]corev1.Node	`json:"items"`
	Total	int				`json:"total"`
}


func (kn *k8sNode) toCells(std []corev1.Node) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = k8sNodeCell(std[i])
	}
	return cells
}

func (kn *k8sNode) fromCells(cells []DataCell) []corev1.Node {
	node := make([]corev1.Node, len(cells))
	for i := range cells {
		node[i] = corev1.Node(cells[i].(k8sNodeCell))
	}
	return node
}

func (kn *k8sNode) GetK8sNodes(filterName string, limit, page int) (k8sNodeResp *K8sNodeResp, err error) {
	NodeList, err := K8s.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取NodeList列表失败." + err.Error()))
		return nil, errors.New("获取NodeList列表失败." + err.Error())
	}

	selectableData := &DataSelector{
		GenericDataList: kn.toCells(NodeList.Items),
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

	return &K8sNodeResp{
		Items: kn.fromCells(data.GenericDataList),
		Total: total,
	}, nil
}

func (kn *k8sNode) GetK8sNodeDetail(k8sNodeName string) (node *corev1.Node, err error) {
	Node, err := K8s.ClientSet.CoreV1().Nodes().Get(context.TODO(), k8sNodeName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Node: %s 详情失败." + err.Error()), k8sNodeName)
		return nil, errors.New("获取Node详情失败." + err.Error())
	}
	return Node, nil
}
