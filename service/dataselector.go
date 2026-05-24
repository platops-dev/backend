package service

import (
	"sort"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	nwv1 "k8s.io/api/networking/v1"
)

/*
1. 定义数据结构
*/
//dataSelect 用于封装排序、过滤、分页的数据类型
type DataSelector struct {
	GenericDataList []DataCell
	DataSelectQuery *DataSelectQuery
}

// DataCell接口，用于各种资源list的类型转换，转换后可以使用DataSelector的自定义排序方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// DataSelectQuery 定义过滤和分页的属性，过滤：Name， 分页：Limit和Page
// Limit是单页的数据条数
// Page是第几页
type DataSelectQuery struct {
	FilterQuery   *FilterQuery
	PaginateQuery *PaginateQuery
}

type FilterQuery struct {
	Name string
}

type PaginateQuery struct {
	Limit int
	Page  int
}

/*
2. 排序
*/
func (d *DataSelector) Len() int {
	return len(d.GenericDataList)
}

func (d *DataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

func (d *DataSelector) Less(i, j int) bool {
	a := d.GenericDataList[i].GetCreation()
	b := d.GenericDataList[j].GetCreation()
	return b.Before(a)
}

func (d *DataSelector) Sort() *DataSelector {
	sort.Sort(d)
	return d
}

/*
3. 过滤
*/
func (d *DataSelector) Filter() *DataSelector {
	//若Name的传参为空, 则返回所有元素
	if d.DataSelectQuery.FilterQuery.Name == "" {
		return d
	}
	//若Name的传参不为空, 则返回元素名中包含Name的所有元素
	filterList := []DataCell{}
	for _, v := range d.GenericDataList {
		matches := true
		objName := v.GetName()
		if !strings.Contains(objName, d.DataSelectQuery.FilterQuery.Name) {
			matches = false
			continue
		}

		if matches {
			filterList = append(filterList, v)
		}
	}
	d.GenericDataList = filterList
	return d
}

/*
3. 分页
*/
func (d *DataSelector) Paginate() *DataSelector {
	limit := d.DataSelectQuery.PaginateQuery.Limit
	page := d.DataSelectQuery.PaginateQuery.Page

	//验证参数是否合法, 如果不合法则不返回数据
	startIndex := limit * (page - 1)
	endIndex := limit * page

	if limit == 0 || page == 0 || len(d.GenericDataList) == 0 {
		return d
	}

	if startIndex >= len(d.GenericDataList) {
		startIndex = 0
	}

	//处理最后一页
	if startIndex < len(d.GenericDataList) && len(d.GenericDataList) < endIndex && endIndex > startIndex {
		endIndex = len(d.GenericDataList)
	}

	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

// 定义podCell类型，实现DataCell接口，用于类型转换
// 定义podCell类型，实现GetCreation 和GetName方法后，可进行类型转换
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}

type deploymentCell appsv1.Deployment

func (d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d deploymentCell) GetName() string {
	return d.Name
}

type daemonSetCell appsv1.DaemonSet

func (ds daemonSetCell) GetCreation() time.Time {
	return ds.CreationTimestamp.Time
}

func (ds daemonSetCell) GetName() string {
	return ds.Name
}

type statefulSetCell appsv1.StatefulSet

func (sts statefulSetCell) GetCreation() time.Time {
	return sts.CreationTimestamp.Time
}

func (sts statefulSetCell) GetName() string {
	return sts.Name
}

type k8sNodeCell corev1.Node

func (kn k8sNodeCell) GetCreation() time.Time {
	return kn.CreationTimestamp.Time
}

func (kn k8sNodeCell) GetName() string {
	return kn.Name
}

type namespaceCell corev1.Namespace

func (ns namespaceCell) GetCreation() time.Time {
	return ns.CreationTimestamp.Time
}

func (ns namespaceCell) GetName() string {
	return ns.Name
}

type k8sServiceCell corev1.Service

func (svc k8sServiceCell) GetCreation() time.Time {
	return svc.CreationTimestamp.Time
}

func (svc k8sServiceCell) GetName() string {
	return svc.Name
}

type ingressCell nwv1.Ingress

func (i ingressCell) GetCreation() time.Time {
	return i.CreationTimestamp.Time
}

func (i ingressCell) GetName() string {
	return i.Name
}

type ConfigMapCell corev1.ConfigMap

func (cm ConfigMapCell) GetCreation() time.Time {
	return cm.CreationTimestamp.Time
}

func (cm ConfigMapCell) GetName() string {
	return cm.Name
}

type SecretCell corev1.Secret

func (st SecretCell) GetCreation() time.Time {
	return st.CreationTimestamp.Time
}

func (st SecretCell) GetName() string {
	return st.Name
}

type PersistentVolumeClaimCell corev1.PersistentVolumeClaim

func (pvc PersistentVolumeClaimCell) GetCreation() time.Time {
	return pvc.CreationTimestamp.Time
}

func (pvc PersistentVolumeClaimCell) GetName() string {
	return pvc.Name
}

type HorizontalPodAutoscalerCell autoscalingv2.HorizontalPodAutoscaler
// 将类型别名改为结构体嵌入
// type HorizontalPodAutoscalerCell struct {
//     *autoscalingv2.HorizontalPodAutoscaler
// }

func (h HorizontalPodAutoscalerCell) GetCreation() time.Time {
	return h.CreationTimestamp.Time
}

func (h HorizontalPodAutoscalerCell) GetName() string {
	return h.Name
}
