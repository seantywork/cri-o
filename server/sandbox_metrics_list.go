package server

import (
	"context"

	types "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// ListPodSandboxMetrics lists all pod sandbox metrics.
func (s *Server) ListPodSandboxMetrics(ctx context.Context, req *types.ListPodSandboxMetricsRequest) (*types.ListPodSandboxMetricsResponse, error) {
	sboxList := s.ListSandboxes()
	metricsList := s.MetricsForPodSandboxList(sboxList)
	responseMetricsList := make([]*types.PodSandboxMetrics, 0, len(metricsList))

	for _, metrics := range metricsList {
		if current := metrics.GetMetric(); current != nil {
			responseMetricsList = append(responseMetricsList, current)
		} else {
			// Iterate over container metrics within each PodSandboxMetrics.
			containerMetricsList := metrics.GetMetric().GetContainerMetrics()
			for _, containerMetrics := range containerMetricsList {
				containerPodSandboxMetrics := &types.PodSandboxMetrics{
					PodSandboxId:     metrics.GetMetric().GetPodSandboxId(),
					ContainerMetrics: []*types.ContainerMetrics{containerMetrics},
				}
				responseMetricsList = append(responseMetricsList, containerPodSandboxMetrics)
			}
		}
	}

	return &types.ListPodSandboxMetricsResponse{
		PodMetrics: responseMetricsList,
	}, nil
}
