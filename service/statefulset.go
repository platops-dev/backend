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

var StatefulSet statefulSet

type statefulSet struct{}

type StatefulSetResp struct {
	Items	[]appsv1.StatefulSet	`json:"items"`
	Total	int						`json:"total"`
}

type StatefulSetCreate struct {
	Metadata		metav1.ObjectMeta			`json:"metadata"`
	Spec			appsv1.StatefulSetSpec		`json:"spec"`
	Status			appsv1.StatefulSetStatus	`json:"status"`		
}



type StatefulSetNp struct {
	Namespace		string	`json:"namespace"`
	StatefulSetNum 	int		`json:"statefulset_num"`
}


func (sts *statefulSet) toCells(std []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = statefulSetCell(std[i])
	}
	return cells
}

func (sts *statefulSet) fromCells(cells []DataCell) []appsv1.StatefulSet {
	statefulSet := make([]appsv1.StatefulSet, len(cells))
	for i := range cells {
		statefulSet[i] = appsv1.StatefulSet(cells[i].(statefulSetCell))
	}
	return statefulSet
}


func (sts *statefulSet) GetStatefulSets(filterName, namespace string, limit, page int) (statefulSetResp *StatefulSetResp, err error) {
	StatefulSetList, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取StatefulSet列表失败. " + err.Error()))
		return nil, errors.New("获取StatefulSet列表失败. " + err.Error())
	}

	selectableData := &DataSelector{
		GenericDataList: sts.toCells(StatefulSetList.Items),
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

	return &StatefulSetResp{
		Items: sts.fromCells(data.GenericDataList),
		Total: total,
	}, nil
}




func (sts *statefulSet) GetStatefulSetDetail(statefulSetName, namespace string) (statefulSet *appsv1.StatefulSet, err error) {
	StatefulSet, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Statefulset详情失败. " + err.Error()))
		return nil, errors.New("获取Statefulset详情失败. " + err.Error())
	}
	return StatefulSet, nil
}

func (sts *statefulSet) ScaleStatefulSet(statefulSetName, namespace string, scalenum int) (replicas int32, err error) {
	scale, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).GetScale(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取副本数失败. " + err.Error()))
		return 0, errors.New("获取副本数失败. " + err.Error())
	}

	scale.Spec.Replicas = int32(scalenum)

	newScale, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).UpdateScale(context.TODO(), statefulSetName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新副本数失败. " + err.Error()))
		return 0, errors.New("更新副本数失败. " + err.Error())
	}
	return newScale.Spec.Replicas, nil
}

// func (sts *statefulSet) CreateStatefulSet(data *StatefulSetCreate) (err error) {
// 	statefulSet := &appsv1.StatefulSet{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name: data.StatefulSetName,
// 			Namespace: data.Namespace,
// 			Labels: data.Labels,
// 		},
// 		Spec: appsv1.StatefulSetSpec{
// 			ServiceName: data.ServiceName,
// 			Replicas: &data.Resplicas,
// 			Selector: &metav1.LabelSelector{
// 				MatchLabels: data.Labels,
// 			},
// 			Template: corev1.PodTemplateSpec{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Labels: data.Labels,
// 				},
// 				Spec: corev1.PodSpec{
// 					Containers: data.Containers,
// 				},
// 			},
// 			VolumeClaimTemplates: data.VolumeClaimTemplates,
// 		},
// 		Status: appsv1.StatefulSetStatus{},
// 	}

// 	for i := range statefulSet.Spec.Template.Spec.Containers {
// 		if data.Health_Check {
// 			Container_Port := statefulSet.Spec.Template.Spec.Containers[i].Ports[0].ContainerPort
// 			statefulSet.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
// 				ProbeHandler: corev1.ProbeHandler{
// 					HTTPGet: &corev1.HTTPGetAction{
// 						Path: data.Health_Path,
// 						//instr.InOrString 的作用时端口可以定义为整型, 也可以定义为字符串
// 						//Type=0 则表示该结构体实例内的数据为整型, 转json时只使用IntVal的数据
// 						//Type=1 则表示该结构体实例内的数据为字符串, 转json时只使用StrVal的数据
// 						Port: intstr.IntOrString{
// 							Type: 0,
// 							IntVal: Container_Port,
// 						},
// 					},
// 				},
// 			}
// 			statefulSet.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
// 				ProbeHandler: corev1.ProbeHandler{
// 					HTTPGet: &corev1.HTTPGetAction{
// 						Path: data.Health_Path,
// 						//instr.InOrString 的作用时端口可以定义为整型, 也可以定义为字符串
// 						//Type=0 则表示该结构体实例内的数据为整型, 转json时只使用IntVal的数据
// 						//Type=1 则表示该结构体实例内的数据为字符串, 转json时只使用StrVal的数据
// 						Port: intstr.IntOrString{
// 							Type: 0,
// 							IntVal: Container_Port,
// 						},
// 					},
// 				},
// 			}
// 		}
// 		if data.Resource_Check {
// 			statefulSet.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
// 				corev1.ResourceCPU : resource.MustParse(data.Cpu),
// 				corev1.ResourceMemory : resource.MustParse(data.Memory),
// 			}
// 			statefulSet.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
// 				corev1.ResourceCPU : resource.MustParse(data.Cpu),
// 				corev1.ResourceMemory : resource.MustParse(data.Memory),
// 			}
// 		}
		
// 	}
// 	_, err = K8s.ClientSet.AppsV1().StatefulSets(data.Namespace).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
// 	if err != nil {
// 		logger.Error(errors.New("创建statefulset失败. " + err.Error()))
// 		return errors.New("创建statefulset失败. " + err.Error())
// 	}
// 	return nil
// }


func (sts *statefulSet) CreateStatefulSet(data *StatefulSetCreate) (err error) {
		statefulSet := &appsv1.StatefulSet{
		ObjectMeta: data.Metadata,
		Spec: data.Spec,
		Status: data.Status,
	}
	// fmt.Println(*&statefulSet.Spec.VolumeClaimTemplates)
	_, err = K8s.ClientSet.AppsV1().StatefulSets(data.Metadata.Namespace).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建statefulset失败. " + err.Error()))
		return errors.New("创建statefulset失败. " + err.Error())
	}
	 return nil
}



func (sts *statefulSet) DeleteStatefulSet(statefulSetName, namespace  string) (err error) {
	err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除statefulset失败. " + err.Error()))
		return errors.New("删除statefulset失败. " + err.Error())
	}
	return nil
}

func (sts *statefulSet) RestartStatefulSet(statefulSetName, namespace string) (err error) {
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]string{
						"kubectl.kubenetes.io/restarteAt" : metav1.Now().Format(time.RFC3339),
					},
				},
			},
		},
	}

	patchByte, err := json.Marshal(patchData)
	if err != nil {
		logger.Error(errors.New("json序列化失败. " + err.Error()))
		return errors.New("json序列化失败. " + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Patch(context.TODO(), statefulSetName, "application/strategic-merge-patch+json", patchByte, metav1.PatchOptions{})
	if err != nil {
		logger.Error(errors.New("重启statefulset失败. " + err.Error()))
		return errors.New("重启statefulset失败. " + err.Error())
	}
	return nil
}

func (sts *statefulSet) UpdateStatefulSet(namespace, content string) (err error) {
	var statefulSet = &appsv1.StatefulSet{}

	err = json.Unmarshal([]byte(content), statefulSet)
	if err != nil {
		logger.Error(errors.New("json反序列化失败. " + err.Error()))
		return errors.New("json反序列化失败. " + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新statefulset失败. " + err.Error()))
		return errors.New("更新statefulset失败. " + err.Error())
	}
	return nil
}

func (sts *statefulSet) GetStatefulSetsNumPerNp() (statefulSetNps []*StatefulSetNp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	
	for _, namespace := range namespaceList.Items {
		statefulSetList, err := K8s.ClientSet.AppsV1().StatefulSets(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		statefulSetNp := &StatefulSetNp{
			Namespace: namespace.Name,
			StatefulSetNum: len(statefulSetList.Items),
		}

		statefulSetNps = append(statefulSetNps, statefulSetNp)
	}
	return statefulSetNps, nil
}