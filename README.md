# OpenStack Tenant Exporter

## Build openstack-tenant-exporter

´´´ bash
export OTE_TAG=0.1.2
docker build -t openstack-tenant-exporter:$OTE_TAG .
docker tag openstack-tenant-exporter:$OTE_TAG brocolis/openstack-tenant-exporter:$OTE_TAG
docker tag openstack-tenant-exporter:$OTE_TAG brocolis/openstack-tenant-exporter:latest
docker push brocolis/openstack-tenant-exporter:$OTE_TAG
docker push brocolis/openstack-tenant-exporter:latest
´´´
