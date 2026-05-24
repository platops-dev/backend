package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var PersistentVolumeClaim persistentVolumeClaim

type persistentVolumeClaim struct {}

type PersistentVolumeClaimResp struct {
	Items []corev1.PersistentVolumeClaim	`json:"items"`
	Total	int								`json:"total"`
}

func (pvc *persistentVolumeClaim) toCells(std []corev1.PersistentVolumeClaim) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = PersistentVolumeClaimCell(std[i])
	}
	return cells
}
func (pvc *persistentVolumeClaim) fromCells(cells []DataCell) []corev1.PersistentVolumeClaim {
	PersistentVolumeClaim := make([]corev1.PersistentVolumeClaim, len(cells))
	for i := range cells {
		PersistentVolumeClaim[i] = corev1.PersistentVolumeClaim(cells[i].(PersistentVolumeClaimCell))
	}
	return PersistentVolumeClaim
}

func (pvc *persistentVolumeClaim) GetPersistentVolumeClaims(filterName, namespace string, limit, page int) (persistentVolumeClaimResp *PersistentVolumeClaimResp, err error) {
	PersistentVolumeClaimList, err := K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Namespace: %s 下的PersistentVolumeClaimList列表失败. " + err.Error()), namespace)
		return nil, errors.New("获取Namespace下的PersistentVolumeClaimList列表失败. " + err.Error())
	}

	selectableData := &DataSelector{
		GenericDataList: pvc.toCells(PersistentVolumeClaimList.Items),
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

	persistentVolumeClaimResps := pvc.fromCells(data.GenericDataList)

	return &PersistentVolumeClaimResp{
		Items: persistentVolumeClaimResps,
		Total: total,
	}, nil
}


func (pvc *persistentVolumeClaim) GetPersistentVolumeClaimDetail(persistentVolumeClaimName, namespace string) (persistentVolumeClaim *corev1.PersistentVolumeClaim, err error) {
	PersistentVolumeClaim, err := K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), persistentVolumeClaimName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Namespace: %s 下的PersistentVolumeClaim: %s 详情失败. " + err.Error()), namespace, persistentVolumeClaimName)
		return nil, errors.New("获取Namespace下的PersistentVolumeClaim 详情失败. " + err.Error())
	}
	return PersistentVolumeClaim, nil
}

func (pvc *persistentVolumeClaim) DeletePersistentVolumeClaim(persistentVolumeClaimName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), persistentVolumeClaimName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Namespace: %s 下的PersistentVolumeClaim: %s 失败. " + err.Error()), namespace, persistentVolumeClaimName)
		return errors.New("删除Namespace下的PersistentVolumeClaim 失败. " + err.Error())
	}
	return nil
}

func (pvc *persistentVolumeClaim) UpdatePersistentVolumeClaim(namespace, content string) (err error) {
	var persistentVolumeClaim = &corev1.PersistentVolumeClaim{}

	err = json.Unmarshal([]byte(content), persistentVolumeClaim)
	if err != nil {
		logger.Error(errors.New("JSON反序化失败." + err.Error()))
		return errors.New("JSON反序化失败." + err.Error())
	}

	_, err = K8s.ClientSet.CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), persistentVolumeClaim, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Namespace: %s 下的PersistentVolumeClaim: %s 失败. " + err.Error()), namespace, persistentVolumeClaim.Name)
		return errors.New("更新Namespace下的PersistentVolumeClaim 失败. " + err.Error())
	}
	return nil
}
