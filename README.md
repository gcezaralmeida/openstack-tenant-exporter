# OpenStack Tenant Exporter

This chart installs the OpenStack Tenant Exporter, a Prometheus exporter
that collects metrics about OpenStack tenants.

## Install Using Helm

``` bash
helm upgrade --install openstack-tenant-exporter oci://registry-1.docker.io/brocolis/openstack-tenant-exporter -n kube-system
```

## OpenStack Metrics

- openstack_max_total_instances
- openstack_total_instances_used
- openstack_max_total_cores
- openstack_total_cores_used
- openstack_max_total_ram
- openstack_total_ram_used
- openstack_max_total_volumes
- openstack_total_volumes_used
- openstack_max_total_snapshots
- openstack_total_snapshots_used
- openstack_max_total_volume_gigabytes
- openstack_total_volume_gigabytes_used
- openstack_max_loadbalancers
- openstack_total_loadbalancers_used
- openstack_max_floatingips
- openstack_total_floatingips_used
- openstack_volume_size
- openstack_volume_status
- openstack_snapshot_size
- openstack_snapshot_status


## Dashboards Examples

![OpenStack Overview](dashboards/openstack-overview.png)
![OpenStack Volumes](dashboards/openstack-volumes.png)
![OpenStack Volume Snapshots](dashboards/openstack-volumes-snapshots.png)
