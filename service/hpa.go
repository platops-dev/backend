package service

import (
	"context"
	"fmt"
	"time"

	"github.com/wonderivan/logger"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// var K8sHpa = k8sHpa
var K8sHpa k8sHpa

type k8sHpa struct{}

// 响应结构体
type K8sHpaResp struct {
	Items []autoscalingv2.HorizontalPodAutoscaler `json:"items"`
	Total int                                     `json:"total"`
}

// 创建请求结构体
type HpaCreate struct {
	Name           string            `json:"name" binding:"required"`
	Namespace      string            `json:"namespace" binding:"required"`
	Label          map[string]string `json:"label"`
	MinReplicas    int32             `json:"min_replicas" binding:"min=1"`
	MaxReplicas    int32             `json:"max_replicas" binding:"required,min=1"`
	CpuPercent     int32             `json:"cpu_percent" binding:"required,min=1,max=100"`
	MemoryPercent  int32             `json:"memory_percent"`
	ScaleTargetRef ScaleTargetRef    `json:"scale_target_ref" binding:"required"`
}

type ScaleTargetRef struct {
	ApiVersion string `json:"api_version" binding:"required"`
	Kind       string `json:"kind" binding:"required"`
	Name       string `json:"name" binding:"required"`
}

// 分页查询参数
type HpaListParams struct {
	FilterName string `form:"filter_name"`
	Namespace  string `form:"namespace" binding:"required"`
	Limit      int    `form:"limit" binding:"min=1,max=100"`
	Page       int    `form:"page" binding:"min=1"`
}

// toCells 转换为 DataCell 切片
func (h *k8sHpa) toCells(std []autoscalingv2.HorizontalPodAutoscaler) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = HorizontalPodAutoscalerCell(std[i])
		// cells[i] = HorizontalPodAutoscalerCell{
		//     HorizontalPodAutoscaler: &std[i],
		// }
	}
	return cells
}

// fromCells 从 DataCell 切片转换回 HPA 切片
func (h *k8sHpa) fromCells(cells []DataCell) []autoscalingv2.HorizontalPodAutoscaler {
	hpas := make([]autoscalingv2.HorizontalPodAutoscaler, len(cells))
	for i := range cells {
		hpas[i] = autoscalingv2.HorizontalPodAutoscaler(cells[i].(HorizontalPodAutoscalerCell))
	}
	return hpas
}

// GetK8sHpas 获取 HPA 列表（带分页和过滤）
func (h *k8sHpa) GetK8sHpas(ctx context.Context, filterName, namespace string, limit, page int) (*K8sHpaResp, error) {
	// 添加超时控制
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	hpaList, err := K8s.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		logger.Error("获取HPA列表失败", "namespace", namespace, "error", err)
		return nil, fmt.Errorf("获取HPA列表失败: %w", err)
	}

	// 数据选择和过滤
	selectableData := &DataSelector{
		GenericDataList: h.toCells(hpaList.Items),
		DataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{Name: filterName},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()

	return &K8sHpaResp{
		Items: h.fromCells(data.GenericDataList),
		Total: total,
	}, nil
}

// GetK8sHpaDetail 获取 HPA 详情
func (h *k8sHpa) GetK8sHpaDetail(ctx context.Context, hpaName, namespace string) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	hpaDetail, err := K8s.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(namespace).Get(ctx, hpaName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取HPA详情失败",
			"namespace", namespace,
			"hpaName", hpaName,
			"error", err,
		)
		return nil, fmt.Errorf("获取HPA %s/%s 详情失败: %w", namespace, hpaName, err)
	}
	return hpaDetail, nil
}

// buildMetrics 构建 HPA 指标配置
func (h *k8sHpa) buildMetrics(cpuPercent, memoryPercent int32) ([]autoscalingv2.MetricSpec, error) {
	var metrics []autoscalingv2.MetricSpec

	// 添加 CPU 指标
	if cpuPercent > 0 {
		metrics = append(metrics, autoscalingv2.MetricSpec{
			Type: autoscalingv2.ResourceMetricSourceType,
			Resource: &autoscalingv2.ResourceMetricSource{
				Name: corev1.ResourceCPU,
				Target: autoscalingv2.MetricTarget{
					Type:               autoscalingv2.UtilizationMetricType,
					AverageUtilization: &cpuPercent,
				},
			},
		})
	}

	// 添加内存指标
	if memoryPercent > 0 {
		metrics = append(metrics, autoscalingv2.MetricSpec{
			Type: autoscalingv2.ResourceMetricSourceType,
			Resource: &autoscalingv2.ResourceMetricSource{
				Name: corev1.ResourceMemory,
				Target: autoscalingv2.MetricTarget{
					Type:               autoscalingv2.UtilizationMetricType,
					AverageUtilization: &memoryPercent,
				},
			},
		})
	}

	if len(metrics) == 0 {
		return nil, fmt.Errorf("至少需要配置一个指标（CPU或内存）")
	}

	return metrics, nil
}

// CreateK8sHpa 创建 HPA（使用 autoscaling/v2 API）
func (h *k8sHpa) CreateK8sHpa(ctx context.Context, data *HpaCreate) error {
	// 参数校验
	if err := validateHpaCreate(data); err != nil {
		return fmt.Errorf("参数校验失败: %w", err)
	}

	// 设置默认值
	minReplicas := data.MinReplicas
	if minReplicas == 0 {
		minReplicas = 1
	}

	// 构建指标配置
	metrics, err := h.buildMetrics(data.CpuPercent, data.MemoryPercent)
	if err != nil {
		return err
	}

	// 构建 v2 版本的 HPA 对象
	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: data.ScaleTargetRef.ApiVersion,
				Kind:       data.ScaleTargetRef.Kind,
				Name:       data.ScaleTargetRef.Name,
			},
			MinReplicas: &minReplicas,
			MaxReplicas: data.MaxReplicas,
			Metrics:     metrics,
		},
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// 使用 v2 API 创建 HPA
	_, err = K8s.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(data.Namespace).Create(ctx, hpa, metav1.CreateOptions{})
	if err != nil {
		logger.Error("创建HPA失败",
			"namespace", data.Namespace,
			"name", data.Name,
			"error", err,
		)
		return fmt.Errorf("创建HPA %s/%s 失败: %w", data.Namespace, data.Name, err)
	}

	logger.Info("创建HPA成功",
		"namespace", data.Namespace,
		"name", data.Name,
		"cpuPercent", data.CpuPercent,
		"memoryPercent", data.MemoryPercent,
	)
	return nil
}

// UpdateK8sHpa 更新 HPA（使用 autoscaling/v2 API）
func (h *k8sHpa) UpdateK8sHpa(ctx context.Context, hpaName, namespace string, data *HpaCreate) error {
	// 参数校验
	if err := validateHpaCreate(data); err != nil {
		return fmt.Errorf("参数校验失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// 先获取现有的 HPA（v2 版本）
	existingHpa, err := K8s.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(namespace).Get(ctx, hpaName, metav1.GetOptions{})
	if err != nil {
		logger.Error("获取HPA失败，无法更新",
			"namespace", namespace,
			"name", hpaName,
			"error", err,
		)
		return fmt.Errorf("获取HPA %s/%s 失败，无法更新: %w", namespace, hpaName, err)
	}

	// 设置默认值
	minReplicas := data.MinReplicas
	if minReplicas == 0 {
		minReplicas = 1
	}

	// 构建指标配置
	metrics, err := h.buildMetrics(data.CpuPercent, data.MemoryPercent)
	if err != nil {
		return err
	}

	// 更新 HPA 的字段
	existingHpa.ObjectMeta.Labels = data.Label
	existingHpa.Spec.ScaleTargetRef = autoscalingv2.CrossVersionObjectReference{
		APIVersion: data.ScaleTargetRef.ApiVersion,
		Kind:       data.ScaleTargetRef.Kind,
		Name:       data.ScaleTargetRef.Name,
	}
	existingHpa.Spec.MinReplicas = &minReplicas
	existingHpa.Spec.MaxReplicas = data.MaxReplicas
	existingHpa.Spec.Metrics = metrics

	// 执行更新
	_, err = K8s.ClientSet.AutoscalingV2().HorizontalPodAutoscalers(namespace).Update(ctx, existingHpa, metav1.UpdateOptions{})
	if err != nil {
		logger.Error("更新HPA失败",
			"namespace", namespace,
			"name", hpaName,
			"error", err,
		)
		return fmt.Errorf("更新HPA %s/%s 失败: %w", namespace, hpaName, err)
	}

	logger.Info("更新HPA成功",
		"namespace", namespace,
		"name", hpaName,
		"cpuPercent", data.CpuPercent,
		"memoryPercent", data.MemoryPercent,
	)
	return nil
}

// DeleteK8sHpa 删除 HPA
func (h *k8sHpa) DeleteK8sHpa(ctx context.Context, hpaName, namespace string) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := K8s.ClientSet.AutoscalingV1().HorizontalPodAutoscalers(namespace).Delete(ctx, hpaName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error("删除HPA失败",
			"namespace", namespace,
			"name", hpaName,
			"error", err,
		)
		return fmt.Errorf("删除HPA %s/%s 失败: %w", namespace, hpaName, err)
	}

	logger.Info("删除HPA成功",
		"namespace", namespace,
		"name", hpaName,
	)
	return nil
}

// validateHpaCreate 校验 HPA 创建/更新参数
func validateHpaCreate(data *HpaCreate) error {
	if data.Name == "" {
		return fmt.Errorf("HPA名称不能为空")
	}
	if data.Namespace == "" {
		return fmt.Errorf("命名空间不能为空")
	}
	if data.MaxReplicas < data.MinReplicas {
		return fmt.Errorf("最大副本数不能小于最小副本数")
	}
	if data.CpuPercent < 0 || data.CpuPercent > 100 {
		return fmt.Errorf("CPU百分比必须在0-100之间")
	}
	if data.MemoryPercent < 0 || data.MemoryPercent > 100 {
		return fmt.Errorf("内存百分比必须在0-100之间")
	}
	if data.CpuPercent == 0 && data.MemoryPercent == 0 {
		return fmt.Errorf("至少需要配置CPU或内存指标之一")
	}
	if data.ScaleTargetRef.Name == "" {
		return fmt.Errorf("目标资源名称不能为空")
	}
	return nil
}
