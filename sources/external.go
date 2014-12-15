package sources

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"

    "github.com/golang/glog"
)

// While updating this, also update heapster/deploy/Dockerfile.
const HostsFile = "etc/hosts"

type ExternalSource struct {
    cadvisor *cadvisorSource
}

func (self *ExternalSource) getCadvisorHosts() (*CadvisorHosts, error) {
    fi, err := os.Stat(HostsFile)
    if err != nil {
        return nil, err
    }
    if fi.Size() == 0 {
        return &CadvisorHosts{}, nil
    }
    contents, err := ioutil.ReadFile(HostsFile)
    if err != nil {
        return nil, err
    }
    var cadvisorHosts CadvisorHosts
    err = json.Unmarshal(contents, &cadvisorHosts)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal contents of file %s. Error: %s", HostsFile, err)
    }
    return &cadvisorHosts, nil
}

func (self *ExternalSource) GetInfo() (ContainerData, error) {
    hosts, err := self.getCadvisorHosts()
    if err != nil {
        return ContainerData{}, err
    }

    containers, nodes, err := self.cadvisor.fetchData(hosts)
    if err != nil {
        glog.Error(err)
        return ContainerData{}, nil
    }

    return ContainerData{
        Containers: containers,
        Machine:    nodes,
    }, nil
}

func newExternalSource() (Source, error) {
    if _, err := os.Stat(HostsFile); err != nil {
        return nil, fmt.Errorf("Cannot stat hosts_file %s. Error: %s", HostsFile, err)
    }
    cadvisorSource := newCadvisorSource()
    return &ExternalSource{
        cadvisor: cadvisorSource,
    }, nil
}
