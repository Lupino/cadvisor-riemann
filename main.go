package main

import (
    "fmt"
    "time"
    "flag"
    "github.com/golang/glog"
    "github.com/bigdatadev/goryman"
    "github.com/Lupino/cadvisor-riemann/sources"
)

var riemannAddress = flag.String("riemann_address", "localhost:5555", "specify the riemann server location")
var sampleInterval = flag.Duration("interval", 5*time.Second, "Interval between sampling (default: 5s)")

func pushToRiemann(r *goryman.GorymanClient, service string, metric int, host string, tags []string) {
    err := r.SendEvent(&goryman.Event{
        Service: service,
        Metric:  metric,
        Tags:    tags,
        Host:    host,
    })
    if err != nil {
        glog.Fatalf("unable to write to riemann: %s", err)
    }
}

func main() {
    defer glog.Flush()
    flag.Parse()

    source, err := sources.NewSource()
    if err != nil {
        glog.Fatalf("unable to setup source: %s", err)
    }
    // Setting up the Riemann client
    r := goryman.NewGorymanClient(*riemannAddress)
    err = r.Connect()
    if err != nil {
        glog.Fatalf("unable to connect to riemann: %s", err)
    }
    //defer r.Close()
    // Setting up the ticker
    ticker := time.NewTicker(*sampleInterval).C
    for {
        select {
        case <-ticker:
            // Make the call to get all the possible data points
            info, err := source.GetInfo()
            if err != nil {
                glog.Fatalf("unable to retrieve machine data: %s", err)
            }
            // Start dumping data into riemann
            // Loop into each ContainerInfo
            // Get stats
            // Push into riemann
            for _, container := range info.Containers {
                pushToRiemann(r, fmt.Sprintf("Cpu.Load %s", container.Name), int(container.Stats[0].Cpu.Load), container.Hostname, []string{})
                pushToRiemann(r, fmt.Sprintf("Cpu.Usage.Total %s", container.Name), int(container.Stats[0].Cpu.Usage.Total), container.Hostname, []string{})
                pushToRiemann(r, fmt.Sprintf("Memory.Usage %s", container.Name), int(container.Stats[0].Memory.Usage), container.Hostname, []string{})
            }
        }
    }
}
