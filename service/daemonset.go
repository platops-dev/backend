package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	_"k8s.io/api/core/v1"
	_"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_"k8s.io/apimachinery/pkg/util/intstr"
)

var DaemonSet daemonSet

type daemonSet struct{}

type DaemonSetResp struct {
	Items	[]appsv1.DaemonSet	`json:"items"`
	Total	int					`json:"total"`
}

type DaemonSetCreate struct {
	Metadata		metav1.ObjectMeta			`json:"metadata"`
	Spec			appsv1.DaemonSetSpec		`json:"spec"`
	Status			appsv1.DaemonSetStatus		`json:"status"`		
}



type DaemonSetNp struct {
	Namespace		string	`json:"namespace"`
	DaemonSetNum	int		`json:"daemonset_num"` 
}


func (ds *daemonSet) toCells(std []appsv1.DaemonSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = daemonSetCell(std[i])
	}
	return cells
}

func (ds *daemonSet) formcells(cells []DataCell) []appsv1.DaemonSet {
	daemonSet := make([]appsv1.DaemonSet, len(cells))
	for i := range cells {
		daemonSet[i] = appsv1.DaemonSet(cells[i].(daemonSetCell))
	}
	return daemonSet
}



func (ds *daemonSet) GetDaemonSets(filterName, namespace string, limit, page int) (daemonSetResp *DaemonSetResp, err error) {
	DaemonSetList, err := K8s.ClientSet.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Daemonset 列表失败. " + err.Error()))
		return nil, errors.New("获取Daemonset 列表失败. " + err.Error())
	}

	selectableData := &DataSelector{
		GenericDataList: ds.toCells(DaemonSetList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page: page,
			},
		},
	}

	filtered := selectableData.Filter()
	total	:= len(filtered.GenericDataList)
	data    := filtered.Sort().Paginate()

	return &DaemonSetResp{
		Items: ds.formcells(data.GenericDataList),
		Total: total,
	}, nil
}


func (ds *daemonSet) GetDaemonSetDetail(daemonSetName, namespace string) (daemonSet *appsv1.DaemonSet, err error)  {
	Daemonset, err := K8s.ClientSet.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取daemonset 详情失败. " + err.Error()))
		return nil, errors.New("获取daemonset 详情失败. " + err.Error())
	}
	return Daemonset, nil
}


func (ds *daemonSet) CreateDaemonSet(data *DaemonSetCreate) (err error) {
	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: data.Metadata,
		Spec: data.Spec,
		Status: data.Status,
	}
	_, err = K8s.ClientSet.AppsV1().DaemonSets(data.Metadata.Namespace).Create(context.TODO(), daemonSet, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建DaemonSet失败. " + err.Error()))
		return errors.New("创建DaemonSet失败. " + err.Error())
	}
	return nil
}



func (ds *daemonSet) DeleteDaemonSet(daemonSetName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Delete(context.TODO(), daemonSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除DaemonSet失败. " + err.Error()))
		return errors.New("删除DaemonSet失败. " + err.Error())
	}
	return nil
}


func (ds *daemonSet) RestartDaemonSet(daemonSetName, namespace string) (err error)  {
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]string{
						"kubectl.kubenetes.io/restarteAt": metav1.Now().Format(time.RFC3339),
					},
				},
			},
		},
	}
	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error(errors.New("序列化json失败. " + err.Error()))
		return errors.New("序列化json失败. " + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Patch(context.TODO(), daemonSetName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error(errors.New("重启Daemonset失败. " + err.Error()))
		return errors.New("重启Daemonset失败. " + err.Error())
	}
	return nil
}




func (ds *daemonSet) UpdateDaemonSet(namespace, content string) (err error) {
	var daemonSet = &appsv1.DaemonSet{}

	err = json.Unmarshal([]byte(content), daemonSet)
	if err != nil {
		logger.Error(errors.New("反序列化json失败. " + err.Error()))
		return errors.New("反序列化json失败. " + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Update(context.TODO(), daemonSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Daemonset失败. " + err.Error()))
		return errors.New("更新Daemonset失败. " + err.Error())
	}
	return nil
}

func (ds *daemonSet) GetDaemonSetsNumPerNp() (daemonSetsNps []*DaemonSetNp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaceList.Items {
		DaemonsetList, err := K8s.ClientSet.AppsV1().DaemonSets(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		daemonSetsNp := &DaemonSetNp{
			Namespace: namespace.Name,
			DaemonSetNum: len(DaemonsetList.Items),
		}

		daemonSetsNps = append(daemonSetsNps, daemonSetsNp)
	}
	return daemonSetsNps, nil
}