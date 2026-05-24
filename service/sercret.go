package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Secret secret

type secret struct{}

type SecretResp struct {
	Items []corev1.Secret	`json:"items"`
	Total  int				`json:"total"`
}


func (st *secret) toCells(std []corev1.Secret) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = SecretCell(std[i])
	}
	return cells
}
func (st *secret) fromCells(cells []DataCell) []corev1.Secret {
	Secret := make([]corev1.Secret, len(cells))
	for i := range cells {
		Secret[i] = corev1.Secret(cells[i].(SecretCell))
	}
	return Secret
}


func (st *secret) GetSecrets(filterName, namespace string, limit, page int) (secretResp *SecretResp, err error) {
	SecretList, err := K8s.ClientSet.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Namespace: %s 下的SecretList列表失败. " + err.Error()), namespace)
		return nil, errors.New("获取Namespace下的SecretList列表失败. " + err.Error())
	}

	selectableData := &DataSelector{
		GenericDataList: st.toCells(SecretList.Items),
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

	secretResps := st.fromCells(data.GenericDataList)

	return &SecretResp{
		Items: secretResps,
		Total: total,
	}, nil
}

func (st *secret) GetSecretDetail(secretName, namespace string) (secret *corev1.Secret, err error) {
	Secret, err := K8s.ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Namespace: %s Secret %s 详情失败. " + err.Error()), namespace, secretName)
		return nil, errors.New("获取Namespace下的Secret 详情失败. " + err.Error())
	}
	return Secret, nil
}

func (st *secret) DeleteSecret(secretName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Namespace: %s 下的Secret %s 失败. " + err.Error()), namespace, secretName)
		return errors.New("删除Namespace下的Secret 失败. " + err.Error())
	}
	return nil
}


func (st *secret) UpdateSecret(namespace, content string) (err error) {
	var secret = &corev1.Secret{}

	err = json.Unmarshal([]byte(content), secret)
	if err != nil {
		logger.Error(errors.New("JONS反序列化失败." + err.Error()))
		return errors.New("JONS反序列化失败." + err.Error())
	}
	
	_, err = K8s.ClientSet.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Namespace: %s 下的Secret %s 失败. " + err.Error()), namespace, secret.Name)
		return errors.New("更新Namespace下的Secret 失败. " + err.Error())
	}
	return nil
}