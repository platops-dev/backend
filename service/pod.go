package service

import (
	"bytes"
	"context"
	"dashboard-api/config"
	"encoding/json"
	"errors"
	"io"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//定义pod类型和Pod对象，用于包外的调用(包是指service目录)，例如Controller
var Pod pod
type pod struct {}


//定义返回内容
type PodsResp struct {
	Items []corev1.Pod	`json:"items"`
	Total int			`json:"total"`
}


//定义PodsNp类型，用于返回namespace中pod的数量
type PodsNp struct {
	Namespace string `json:"namespace"`
	PodNum int `json:"pod_num"`
}


//定义DataCell到Pod类型转换的方法
func (p *pod) toCells(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}

func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}



func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取pod列表失败, " + err.Error()))
		return nil, errors.New("获取pod列表失败, " + err.Error())
	}
	
	selectableData := &DataSelector{
		GenericDataList: p.toCells(podList.Items),
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
	podsResp = &PodsResp{
		Items: p.fromCells(data.GenericDataList),
		Total: total,
	}
	return podsResp, nil
}


func (p *pod) GetPodDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Pod详情失败, " + err.Error()))
		return nil, errors.New("获取Pod详情失败, " + err.Error())
	}
	return pod, nil
}

func (p *pod) DeletePod(podName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Pod失败, " + err.Error()))
		return errors.New("删除Pod失败, " + err.Error())
	}
	return nil
}

func (p *pod) UpdatePod(namespace, content string) (err error) {
	var pod = &corev1.Pod{}
	//反序列化为pod对象
	err = json.Unmarshal( []byte(content), pod)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Pod失败, " + err.Error()))
		return errors.New("更新Pod失败, " + err.Error())
	}

	return nil
}

func (p *pod) GetPodContainer(podName, namespace string) (containers []string, err error) {
	//获取pod详情
	pod, err := p.GetPodDetail(podName, namespace)
	if err != nil {
		return nil, err
	}
	//从pod对象中拿到容器名
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil
}


func (p *pod) GetPodLog(containerName, podName, namespace string) (log string, err error) {
	//设置日志的配置, 容器名、tail的行数
	LineLimit := int64(config.PodLogTailLine)
	option :=  &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &LineLimit,
	}

	//获取request实例
	req := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, option)

	//发起request请求
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		logger.Error(errors.New("获取podLogs失败, " + err.Error()))
		return "", errors.New("获取podLogs失败, " + err.Error())
	}

	defer podLogs.Close()

	//将io.ReadCloser写入缓冲区, 目的是为了转成string返回
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		logger.Error(errors.New("复制podLogs失败, " + err.Error()))
		return "", errors.New("复制podLogs失败, " + err.Error())
	}
	return buf.String(), nil
}

func (p *pod) GetPodNumPerNp() (podsNps []*PodsNp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaceList.Items {
		podList, err := K8s.ClientSet.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		//组装数据
		podsNp := &PodsNp{
			Namespace: namespace.Name,
			PodNum: len(podList.Items),
		}
		//添加到podsNp切片中
		podsNps = append(podsNps, podsNp)
	}
	return podsNps, nil
}

