# OpenStack Tenant Exporter

This chart installs the OpenStack Tenant Exporter, a Prometheus exporter
that collects metrics about OpenStack tenants.

## Install Using Helm

``` bash
helm upgrade --install openstack-tenant-exporter oci://registry-1.docker.io/brocolis/openstack-tenant-exporter -n kube-system
```

