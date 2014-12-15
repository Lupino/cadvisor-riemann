package sources

import (
    cadvisor "github.com/google/cadvisor/info")

type Container struct {
    Name  string                     `json:"name,omitempty"`
    Spec  cadvisor.ContainerSpec     `json:"spec,omitempty"`
    Stats []*cadvisor.ContainerStats `json:"stats,omitempty"`
}

func newContainer() *Container {
    return &Container{Stats: make([]*cadvisor.ContainerStats, 0)}
}

type RawContainer struct {
    Hostname string `json:"hostname,omitempty"`
    Container
}

type ContainerData struct {
    Containers []RawContainer
    Machine    []RawContainer
}

type CadvisorHosts struct {
    Port  int               `json:"port"`
    Hosts map[string]string `json:"hosts"`
}

type Source interface {
    // Fetches containers or pod information from all the nodes in the cluster.
    // Returns:
    // 1. podsOrContainers: A slice of Pod or a slice of RawContainer
    // 2. nodes: A slice of RawContainer, one for each node in the cluster, that contains
    // root cgroup information.
    GetInfo() (ContainerData, error)
}

func NewSource() (Source, error) {
    return newExternalSource()
}
