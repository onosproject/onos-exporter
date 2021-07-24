# onos-exporter
The exporter for ONOS SD-RAN (µONOS Architecture) to scrape, format, and export KPIs to TSDB databases (e.g., Prometheus).

## Overview
The onos-exporter realizes the collection of KPIs from multiple ONOS SD-RAN components via gRPC interfaces, properly label them according to their namespace and subsystem, and turn them available to be pulled (or pushed to) TSDBs. Currently the implementation supports Prometheus.


## Enable 

To enable logging/monitoring in RiaB, in the sdran-in-a-box-values.yaml (available for the latest version) file remove the comments of the lines below:

```yaml
  fluent-bit:
    enabled: true
  opendistro-es:
    enabled: true
  prometheus-stack:
    enabled: true
```

Associated with the monitoring of sdran components is the onos-exporter, the exporter for ONOS SD-RAN (µONOS Architecture) to scrape, format, and export onos KPIs to TSDB databases (e.g., Prometheus). Currently the implementation supports Prometheus. In order to enable onos-exporter, as shown below, make sure the prometheus-stack is enabled too.

```yaml
  prometheus-stack:
    enabled: true
  onos-exporter:
    enabled: true
```

Be sure to enable onos-kpimon and onos-pci to also look at their metrics in the Grafana dashboard.

## Visualize Grafana

After modified the values file, then run the make command to instantiate RiaB. After deployed, the services and pods related to logging and monitoring will be shown as:

```text
$ kubectl -n riab get svc
... 
alertmanager-operated                     ClusterIP   None              <none>        9093/TCP,9094/TCP,9094/UDP            90s
prometheus-operated                       ClusterIP   None              <none>        9090/TCP                              90s
sd-ran-fluent-bit                         ClusterIP   192.168.205.134   <none>        2020/TCP                              90s
sd-ran-grafana                            ClusterIP   192.168.209.213   <none>        80/TCP                                90s
sd-ran-kube-prometheus-sta-alertmanager   ClusterIP   192.168.166.174   <none>        9093/TCP                              90s
sd-ran-kube-prometheus-sta-operator       ClusterIP   192.168.152.79    <none>        443/TCP                               90s
sd-ran-kube-prometheus-sta-prometheus     ClusterIP   192.168.199.115   <none>        9090/TCP                              90s
sd-ran-kube-state-metrics                 ClusterIP   192.168.155.231   <none>        8080/TCP                              90s
sd-ran-opendistro-es-client-service       ClusterIP   192.168.183.47    <none>        9200/TCP,9300/TCP,9600/TCP,9650/TCP   90s
sd-ran-opendistro-es-data-svc             ClusterIP   None              <none>        9300/TCP,9200/TCP,9600/TCP,9650/TCP   90s
sd-ran-opendistro-es-discovery            ClusterIP   None              <none>        9300/TCP                              90s
sd-ran-opendistro-es-kibana-svc           ClusterIP   192.168.129.238   <none>        5601/TCP                              90s
sd-ran-prometheus-node-exporter           ClusterIP   192.168.137.224   <none>        9100/TCP                              90s
```

Make a port-forward rule to the grafana service on port 3000.

```bash
kubectl -n riab port-forward svc/sd-ran-grafana 3000:80
```

Open a browser and access `localhost:3000`. The credentials to access grafana are: username: admin and password: prom-operator.

To look at the grafana dashboard for the sdran component logs and KPIs, check in the left menu of grafana the option dashboards and select the submenu Manage (or just access in the browser the address http://localhost:3000/dashboards).

In the menu that shows, look for the dashboard named `Kubernetes / Logs / Pod` to check the logs of the sd-ran Kubernetes pods.

In the menu that shows, look for the dashboard named `Kubernetes / SD-RAN KPIs` to check the KPIs of the sd-ran components (e.g., kpimon, pci, topo, uenib and e2t).

In the top menu, the dropdown menus allow the selection of the Namespace riab and one of its Pods. It is also possible to type a string to be found in the logs of a particular pod using the field String.

Similarly, other dashboards can be found in the left menu of grafana, showing for instance each pod workload in the dashboad `Kubernetes / Compute Resources / Workload`.



## Visualize onos-exporter

To look at the onos-exporter metrics, it's possible to access the onos-exporter directly or visualize the metrics in grafana.

To access the metrics directly have a port-forward kubectl command for onos-exporter service:

```bash
kubectl -n riab port-forward svc/onos-exporter 9861
```

Then access the address `localhost:9861/metrics` in the browser. The exporter shows golang related metrics too.

To access the metrics using grafana, proceed with the access to grafana. After accessing grafana go to the Explore item on the left menu, on the openned window select the Prometheus data source, and type the name of the metrics to see its visualization and click on the Run query button.