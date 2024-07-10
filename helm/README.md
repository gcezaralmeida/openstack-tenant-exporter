# Install

helm package .
helm upgrade --install openstack-tenant-exporter ./openstack-tenant-exporter-0.1.1.tgz -n kube-system
