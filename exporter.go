package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	blockstorageLimits "github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/limits"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/snapshots"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/volumes"
	computeLimits "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/limits"
	"github.com/gophercloud/gophercloud/v2/openstack/config"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
	loadbalancers "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
	loadbalancerQuotas "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/quotas"
	networkQuotas "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/quotas"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// OpenStackExporter collects OpenStack metrics
type OpenStackExporter struct {
	ctx                context.Context
	computeClient      *gophercloud.ServiceClient
	storageClient      *gophercloud.ServiceClient
	loadbalancerClient *gophercloud.ServiceClient
	networkClient      *gophercloud.ServiceClient
	tenantID           string
	maxTotalInstances  *prometheus.GaugeVec
	totalInstancesUsed *prometheus.GaugeVec
	maxTotalCores      *prometheus.GaugeVec
	totalCoresUsed     *prometheus.GaugeVec
	maxTotalRAM        *prometheus.GaugeVec
	totalRAMUsed       *prometheus.GaugeVec
	maxTotalVolumes    *prometheus.GaugeVec
	totalVolumesUsed   *prometheus.GaugeVec
	maxTotalSnapshots  *prometheus.GaugeVec
	totalSnapshotsUsed *prometheus.GaugeVec
	maxTotalVolumeGigabytes *prometheus.GaugeVec
	totalVolumeGigabytesUsed *prometheus.GaugeVec
	maxTotalLoadBalancers   *prometheus.GaugeVec
	totalLoadBalancersUsed *prometheus.GaugeVec
	maxTotalFloatingIPs     *prometheus.GaugeVec
	totalFloatingIPsUsed   *prometheus.GaugeVec
	volumeSize   *prometheus.GaugeVec
	volumeStatus *prometheus.GaugeVec
	snapshotSize *prometheus.GaugeVec
	snapshotStatus *prometheus.GaugeVec
}

// NewOpenStackExporter creates a new OpenStackExporter
// It takes a context, computeClient and storageClient as parameters
// It returns a pointer to OpenStackExporter
func NewOpenStackExporter(ctx context.Context, computeClient, storageClient, loadbalancerClient, networkClient *gophercloud.ServiceClient, tenantID string) *OpenStackExporter {
	return &OpenStackExporter{
			ctx:                ctx,
			computeClient:      computeClient,
			storageClient:      storageClient,
			loadbalancerClient: loadbalancerClient,
			networkClient:      networkClient,
			tenantID:           tenantID,
			maxTotalInstances: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_max_total_instances",
					Help: "Maximum total instances allowed in OpenStack",
			}, []string{"tenant_id"}),
			totalInstancesUsed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_total_instances_used",
					Help: "Total instances currently in use in OpenStack",
			}, []string{"tenant_id"}),
			maxTotalCores: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_max_total_cores",
					Help: "Maximum total CPU cores allowed in OpenStack",
			}, []string{"tenant_id"}),
			totalCoresUsed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_total_cores_used",
					Help: "Total CPU cores currently in use in OpenStack",
			}, []string{"tenant_id"}),
			maxTotalRAM: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_max_total_ram",
					Help: "Maximum total RAM size allowed in OpenStack",
			}, []string{"tenant_id"}),
			totalRAMUsed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_total_ram_used",
					Help: "Total RAM currently in use in OpenStack",
			}, []string{"tenant_id"}),
			maxTotalVolumes: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_max_total_volumes",
					Help: "Maximum total volumes allowed in OpenStack storage",
			}, []string{"tenant_id"}),
			totalVolumesUsed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_total_volumes_used",
					Help: "Total volumes currently in use in OpenStack storage",
			}, []string{"tenant_id"}),
			maxTotalSnapshots: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_max_total_snapshots",
					Help: "Maximum total snapshots allowed in OpenStack storage",
			}, []string{"tenant_id"}),
			totalSnapshotsUsed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_total_snapshots_used",
					Help: "Total snapshots currently in use in OpenStack storage",
			}, []string{"tenant_id"}),
			maxTotalVolumeGigabytes: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_max_total_volume_gigabytes",
					Help: "Maximum total volume gigabytes allowed in OpenStack storage",
			}, []string{"tenant_id"}),
			totalVolumeGigabytesUsed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_total_volume_gigabytes_used",
					Help: "Total volume gigabytes currently in use in OpenStack storage",
			}, []string{"tenant_id"}),
			maxTotalLoadBalancers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_max_loadbalancers",
					Help: "Maximum load balancers allowed in OpenStack",
			}, []string{"tenant_id"}),
			totalLoadBalancersUsed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_total_loadbalancers_used",
					Help: "Total load balancers currently in use in OpenStack",
			}, []string{"tenant_id"}),
			maxTotalFloatingIPs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_max_floatingips",
					Help: "Maximum floating IPs allowed in OpenStack",
			}, []string{"tenant_id"}),
			totalFloatingIPsUsed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_total_floatingips_used",
					Help: "Total floating IPs currently in use in OpenStack",
			}, []string{"tenant_id"}),
			volumeSize: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Name: "openstack_volume_size",
				Help: "Size of volumes in OpenStack, identified by tenant ID, volume ID, and volume name",
			}, []string{"tenant_id", "volume_id", "volume_name"}),
			volumeStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_volume_status",
					Help: "Status of volumes in OpenStack, identified by tenant ID, volume ID, and volume name",
			}, []string{"tenant_id", "volume_id", "volume_name"}),
			snapshotSize: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Name: "openstack_snapshot_size",
				Help: "Size of snapshots in OpenStack, identified by tenant ID, volume ID, and volume name",
			}, []string{"tenant_id", "snap_id", "snap_name", "volume_id"}),
			snapshotStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: "openstack_snapshot_status",
					Help: "Status of snapshots in OpenStack, identified by tenant ID, volume ID, and volume name",
			}, []string{"tenant_id", "snap_id", "snap_name", "volume_id"}),
	}
}

// Describe sends the descriptors of each metric over to the provided channel.
//
// It describes the metrics that this exporter provides. This function is called by
// the Prometheus client when it wants to scrape the metrics provided by this
// exporter.
func (e *OpenStackExporter) Describe(ch chan<- *prometheus.Desc) {
	e.maxTotalInstances.Describe(ch)
	e.totalInstancesUsed.Describe(ch)
	e.maxTotalCores.Describe(ch)
	e.totalCoresUsed.Describe(ch)
	e.maxTotalRAM.Describe(ch)
	e.totalRAMUsed.Describe(ch)
	e.maxTotalVolumes.Describe(ch)
	e.totalVolumesUsed.Describe(ch)
	e.maxTotalSnapshots.Describe(ch)
	e.totalSnapshotsUsed.Describe(ch)
	e.maxTotalVolumeGigabytes.Describe(ch)
	e.totalVolumeGigabytesUsed.Describe(ch)
	e.maxTotalLoadBalancers.Describe(ch)
	e.totalLoadBalancersUsed.Describe(ch)
	e.maxTotalFloatingIPs.Describe(ch)
	e.totalFloatingIPsUsed.Describe(ch)
	e.volumeSize.Describe(ch)
	e.volumeStatus.Describe(ch)
	e.snapshotSize.Describe(ch)
	e.snapshotStatus.Describe(ch)
}

func(e *OpenStackExporter) GetTenantID() string {
	return e.tenantID
}

func(e *OpenStackExporter) collectComputeLimits() {
	// Fetch compute limits from OpenStack
	computeLimitsInfo, err := computeLimits.Get(e.ctx, e.computeClient, computeLimits.GetOpts{}).Extract()
	if err != nil {
			log.Println("Failed to get compute limits:", err)
			return
	}

	// Set compute limits Prometheus metrics
	e.maxTotalInstances.WithLabelValues(e.tenantID).Set(float64(computeLimitsInfo.Absolute.MaxTotalInstances))
	e.totalInstancesUsed.WithLabelValues(e.tenantID).Set(float64(computeLimitsInfo.Absolute.TotalInstancesUsed))
	e.maxTotalCores.WithLabelValues(e.tenantID).Set(float64(computeLimitsInfo.Absolute.MaxTotalCores))
	e.totalCoresUsed.WithLabelValues(e.tenantID).Set(float64(computeLimitsInfo.Absolute.TotalCoresUsed))
	e.maxTotalRAM.WithLabelValues(e.tenantID).Set(float64(computeLimitsInfo.Absolute.MaxTotalRAMSize))
	e.totalRAMUsed.WithLabelValues(e.tenantID).Set(float64(computeLimitsInfo.Absolute.TotalRAMUsed))
}

func(e *OpenStackExporter) collectStorageLimits() {
	// Fetch storage limits from OpenStack
	storageLimitsInfo, err := blockstorageLimits.Get(e.ctx, e.storageClient).Extract()
	if err != nil {
			log.Println("Failed to get storage limits:", err)
			return
	}

	// Set storage limits Prometheus metrics
	e.maxTotalVolumes.WithLabelValues(e.tenantID).Set(float64(storageLimitsInfo.Absolute.MaxTotalVolumes))
	e.totalVolumesUsed.WithLabelValues(e.tenantID).Set(float64(storageLimitsInfo.Absolute.TotalVolumesUsed))
	e.maxTotalSnapshots.WithLabelValues(e.tenantID).Set(float64(storageLimitsInfo.Absolute.MaxTotalSnapshots))
	e.totalSnapshotsUsed.WithLabelValues(e.tenantID).Set(float64(storageLimitsInfo.Absolute.TotalSnapshotsUsed))
	e.maxTotalVolumeGigabytes.WithLabelValues(e.tenantID).Set(float64(storageLimitsInfo.Absolute.MaxTotalVolumeGigabytes))
	e.totalVolumeGigabytesUsed.WithLabelValues(e.tenantID).Set(float64(storageLimitsInfo.Absolute.TotalGigabytesUsed))
}

func(e *OpenStackExporter) collectLoadBalancerLimits() {
	// Fetch load balancer quotas
	quotasLBInfo, err := loadbalancerQuotas.Get(e.ctx, e.loadbalancerClient, e.tenantID).Extract()
	if err != nil {
			log.Println("Failed to get load balancer quotas:", err)
			return
	}

	// Set load balancer quotas Prometheus metrics
	e.maxTotalLoadBalancers.WithLabelValues(e.tenantID).Set(float64(quotasLBInfo.Loadbalancer))
	listOpts := loadbalancers.ListOpts{}
	allPages, err := loadbalancers.List(e.loadbalancerClient, listOpts).AllPages(e.ctx)
	if err != nil {
			log.Println("Failed to list load balancers:", err)
			return
	}
	allLoadbalancers, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil {
			log.Println("Failed to extract load balancers:", err)
			return
	}
	e.totalLoadBalancersUsed.WithLabelValues(e.tenantID).Set(float64(len(allLoadbalancers)))
}

func(e *OpenStackExporter) collectNetworkLimits() {
	// Fetch network quotas
	quotasNetInfo, err := networkQuotas.GetDetail(e.ctx, e.networkClient, e.tenantID).Extract()
	if err != nil {
			log.Println("Failed to get network quotas:", err)
			return
	}
	e.maxTotalFloatingIPs.WithLabelValues(e.tenantID).Set(float64(quotasNetInfo.FloatingIP.Limit))
	e.totalFloatingIPsUsed.WithLabelValues(e.tenantID).Set(float64(quotasNetInfo.FloatingIP.Used))
}

func (e *OpenStackExporter) collectStorageVolumes() {
	// Fetch storage volumes
	allPages, err := volumes.List(e.storageClient, volumes.ListOpts{}).AllPages(e.ctx)
	if err != nil {
		log.Println("Failed to list load balancers:", err)
		return
	}

	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
			log.Println("Failed to extract load balancers:", err)
			return
	}

	for _, volume := range allVolumes {
		e.volumeSize.WithLabelValues(e.tenantID, volume.ID, volume.Name).Set(float64(volume.Size))
		e.volumeStatus.WithLabelValues(e.tenantID, volume.ID, volume.Name).Set(float64(VolumeStatusToNumber(volume.Status)))
	}
}

func(e *OpenStackExporter) collectStorageSnapshots() {
	// Fetch storage snapshots
	allPages, err := snapshots.List(e.storageClient, snapshots.ListOpts{}).AllPages(e.ctx)
	if err != nil {
		log.Println("Failed to list load balancers:", err)
		return
	}	
	allSnapshots, err := snapshots.ExtractSnapshots(allPages)
	if err != nil {
			log.Println("Failed to extract load balancers:", err)
			return
	}

	for _, snapshot := range allSnapshots {
		e.snapshotSize.WithLabelValues(e.tenantID, snapshot.ID, snapshot.Name, snapshot.VolumeID).Set(float64(snapshot.Size))
		e.snapshotStatus.WithLabelValues(e.tenantID, snapshot.ID, snapshot.Name, snapshot.VolumeID).Set(float64(SnapshotStatusToNumber(snapshot.Status)))
	}	

}

func (e *OpenStackExporter) Collect(ch chan<- prometheus.Metric) {
	
	e.collectComputeLimits()

	e.collectStorageLimits()

	e.collectLoadBalancerLimits()

	e.collectNetworkLimits()

	e.collectStorageVolumes()

	e.collectStorageSnapshots()

	// Collect and send Prometheus metrics
	e.maxTotalInstances.Collect(ch)
	e.totalInstancesUsed.Collect(ch)
	e.maxTotalCores.Collect(ch)
	e.totalCoresUsed.Collect(ch)
	e.maxTotalRAM.Collect(ch)
	e.totalRAMUsed.Collect(ch)
	e.maxTotalVolumes.Collect(ch)
	e.totalVolumesUsed.Collect(ch)
	e.maxTotalSnapshots.Collect(ch)
	e.totalSnapshotsUsed.Collect(ch)
	e.maxTotalVolumeGigabytes.Collect(ch)
	e.totalVolumeGigabytesUsed.Collect(ch)
	e.maxTotalLoadBalancers.Collect(ch)
	e.totalLoadBalancersUsed.Collect(ch)
	e.maxTotalFloatingIPs.Collect(ch)
	e.totalFloatingIPsUsed.Collect(ch)
	e.volumeSize.Collect(ch)
	e.volumeStatus.Collect(ch)
	e.snapshotSize.Collect(ch)
	e.snapshotStatus.Collect(ch)
}


// main is the entry point of the application.
//
// It parses the OpenStack cloud configuration, creates a provider client,
// compute client, and storage client, creates an OpenStackExporter,
// registers the exporter with Prometheus, starts a Prometheus metrics
// server, and listens for HTTP requests.
func main() {
	ctx := context.Background()

	log.Println("Starting OpenStack Tenant Exporter")

	authOptions, _, tlsConfig, err := clouds.Parse()
	if err != nil {
			log.Fatalf("Failed to parse cloud config: %v", err)
	}

	authOptions.AllowReauth = true

	providerClient, err := config.NewProviderClient(ctx, authOptions, config.WithTLSConfig(tlsConfig))
	if err != nil {
			log.Fatalf("Failed to create provider client: %v", err)
	}

	log.Println("Provider client created")

	computeClient, err := openstack.NewComputeV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
			log.Fatalf("Failed to create compute client: %v", err)
	}

	log.Println("Compute client created")

	storageClient, err := openstack.NewBlockStorageV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
			log.Fatalf("Failed to create storage client: %v", err)
	}

	log.Println("Storage client created")

	loadbalancerClient, err := openstack.NewLoadBalancerV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
			log.Fatalf("Failed to create loadbalancer client: %v", err)
	}

	log.Println("Loadbalancer client created")

	networkClient, err := openstack.NewNetworkV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
			log.Fatalf("Failed to create network client: %v", err)
	}

	log.Println("Network client created")

	exporter := NewOpenStackExporter(ctx, computeClient, storageClient, loadbalancerClient, networkClient, authOptions.TenantID)

	log.Println("OpenStackExporter created")

	prometheus.MustRegister(exporter)

	log.Println("OpenStackExporter registered with Prometheus")

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Prometheus metrics server started at :9183")

	log.Fatal(http.ListenAndServe(":9183", nil))
}