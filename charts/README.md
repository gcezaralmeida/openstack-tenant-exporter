# Install

´´´ bash
export OTE_TAG=0.1.2
helm package ./charts
helm upgrade --install openstack-tenant-exporter ./openstack-tenant-exporter-$OTE_TAG.tgz -n kube-system
´´´
