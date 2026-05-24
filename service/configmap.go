package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ConfigMap configMap

type configMap struct{}

type ConfigMapResp struct {
	Items []corev1.ConfigMap	`json:"items"`
	Total	int					`json:"total"`
}

func (cm *configMap) toCells(std []corev1.ConfigMap) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = ConfigMapCell(std[i])
	}
	return cells
}
func (cm *configMap) fromCells(cells []DataCell) []corev1.ConfigMap {
	ConfigMap := make([]corev1.ConfigMap, len(cells))
	for i := range cells {
		ConfigMap[i] = corev1.ConfigMap(cells[i].(ConfigMapCell))
	}
	return ConfigMap
}

func (cm *configMap) GetConfigMaps(filterName, namespace string, limit, page int) (configMapResp *ConfigMapResp, err error) {
	ConfigMapList, err := K8s.ClientSet.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Namespace: %s 下的ConfigMapList列表失败. " + err.Error()), namespace)
		return nil, errors.New("获取Namespace下的ConfigMapList列表失败. " + err.Error())
	}

	selectableData := &DataSelector{
		GenericDataList: cm.toCells(ConfigMapList.Items),
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

	configMapResps := cm.fromCells(data.GenericDataList)

	return &ConfigMapResp{
		Items: configMapResps,
		Total: total,
	}, nil
}

func (cm *configMap) GetConfigMapDetail(configMapName, namespace string) (configMap *corev1.ConfigMap, err error) {
	ConfigMap, err := K8s.ClientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Namespace: %s 下的ConfigMap %s 详情失败. " +err.Error()), namespace, configMapName)
		return nil, errors.New("获取Namespace下的ConfigMap 详情失败. " + err.Error())
	}
	return ConfigMap, nil
}

func (cm *configMap) DeleteConfigMap(configMapName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), configMapName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Namespace: %s 下的ConfigMap %s 失败. "+ err.Error()), namespace, configMapName)
		return errors.New("删除Namespace下的ConfigMap 失败. "+ err.Error())
	}
	return nil
}

func (cm *configMap) UpdateConfigMap(namespace, content string) (err error) {
	var configMap = &corev1.ConfigMap{}

	err = json.Unmarshal([]byte(content), configMap)
	if err != nil {
		logger.Error(errors.New("JONS反序列化失败." + err.Error()))
		return errors.New("JONS反序列化失败." + err.Error())
	}
	_, err = K8s.ClientSet.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Namespace: %s 下的ConfigMap %s 失败. " + err.Error()), namespace, configMap.Name)
		return errors.New("更新Namespace下的ConfigMap 失败. " + err.Error())
	}
	return nil
}