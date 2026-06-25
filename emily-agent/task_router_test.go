package main

import "testing"

func TestRouteTask_LocalWhenNoClusters(t *testing.T) {
	result := routeTask(nil, "local", 0.3, "gpu")
	if result != "" {
		t.Errorf("expected local execution, got route to %q", result)
	}
}

func TestRouteTask_LocalWhenOnlySelf(t *testing.T) {
	clusters := []clusterHeartbeatView{
		{ClusterID: "local", LoadScore: 0.0, Capabilities: []string{"gpu"}},
	}
	result := routeTask(clusters, "local", 0.3, "gpu")
	if result != "" {
		t.Errorf("expected local, got %q", result)
	}
}

func TestRouteTask_RoutesToLowerLoadRemote(t *testing.T) {
	clusters := []clusterHeartbeatView{
		{ClusterID: "local", LoadScore: 0.8, Capabilities: []string{"emily_prime"}},
		{ClusterID: "colab-gpu", LoadScore: 0.1, Capabilities: []string{"gpu", "emily_prime"}},
		{ClusterID: "prod", LoadScore: 0.5, Capabilities: []string{"emily_prime"}},
	}
	// GPU task should go to colab-gpu (only one with gpu + lower load).
	result := routeTask(clusters, "local", 0.8, "gpu")
	if result != "colab-gpu" {
		t.Errorf("expected colab-gpu, got %q", result)
	}
}

func TestRouteTask_NoRouteWhenRemoteHigherLoad(t *testing.T) {
	clusters := []clusterHeartbeatView{
		{ClusterID: "local", LoadScore: 0.1, Capabilities: []string{"emily_prime"}},
		{ClusterID: "remote", LoadScore: 0.9, Capabilities: []string{"emily_prime", "gpu"}},
	}
	result := routeTask(clusters, "local", 0.1, "")
	if result != "" {
		t.Errorf("expected local (lower load), got %q", result)
	}
}

func TestRouteTask_NoRouteWhenCapMissing(t *testing.T) {
	clusters := []clusterHeartbeatView{
		{ClusterID: "remote", LoadScore: 0.0, Capabilities: []string{"emily_prime"}},
	}
	result := routeTask(clusters, "local", 0.9, "gpu")
	if result != "" {
		t.Errorf("expected local (remote lacks gpu), got %q", result)
	}
}
