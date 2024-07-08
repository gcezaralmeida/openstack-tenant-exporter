# Install

helm package .
helm install release-name ./openstack-tenant-exporter-0.1.0.tgz -n kube-system
