package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	computeLimits "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/limits"
	blockstorageLimits "github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/limits"
	"github.com/gophercloud/gophercloud/v2/openstack/config"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// OpenStackExporter collects OpenStack metrics
type OpenStackExporter struct {
	ctx                context.Context
	computeClient      *gophercloud.ServiceClient
	storageClient      *gophercloud.ServiceClient
	maxTotalInstances  prometheus.Gauge
	totalInstancesUsed prometheus.Gauge
	maxTotalCores      prometheus.Gauge
	totalCoresUsed     prometheus.Gauge
	maxTotalRAM        prometheus.Gauge
	totalRAMUsed       prometheus.Gauge
	maxTotalVolumes    prometheus.Gauge
	totalVolumesUsed   prometheus.Gauge
	maxTotalSnapshots  prometheus.Gauge
	totalSnapshotsUsed prometheus.Gauge
	maxTotalVolumeGigabytes prometheus.Gauge
	totalVolumeGigabytesUsed prometheus.Gauge
}

// NewOpenStackExporter creates a new OpenStackExporter
func NewOpenStackExporter(ctx context.Context, computeClient, storageClient *gophercloud.ServiceClient) *OpenStackExporter {
	return &OpenStackExporter{
		ctx:           ctx,
		computeClient: computeClient,
		storageClient: storageClient,
		maxTotalInstances: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "openstack_max_total_instances",
			Help: "Maximum total instances allowed in OpenStack",
		}),
		totalInstancesUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "openstack_total_instances_used",
			Help: "Total instances currently in use in OpenStack",
		}),
		maxTotalCores: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "openstack_max_total_cores",
			Help: "Maximum total CPU cores allowed in OpenStack",
		}),
		totalCoresUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "openstack_total_cores_used",
			Help: "Total CPU cores currently in use in OpenStack",
		}),
		maxTotalRAM: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "openstack_max_total_ram",
			Help: "Maximum total RAM size allowed in OpenStack",
		}),
		totalRAMUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "openstack_total_ram_used",
			Help: "Total RAM currently in use in OpenStack",
		}),
		maxTotalVolumes: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "openstack_max_total_volumes",
			Help: "Maximum total volumes allowed in OpenStack storage",
		}),
		totalVolumesUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "openstack_total_volumes_used",
			Help: "Total volumes currently in use in OpenStack storage",
		}),
		maxTotalSnapshots: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "openstack_max_total_snapshots",
			Help: "Maximum total snapshots allowed in OpenStack storage",
		}),
		totalSnapshotsUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "openstack_total_snapshots_used",
			Help: "Total snapshots currently in use in OpenStack storage",
		}),
		maxTotalVolumeGigabytes: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "openstack_max_total_volume_gigabytes",
			Help: "Maximum total volume gigabytes allowed in OpenStack storage",
		}),
		totalVolumeGigabytesUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "openstack_total_volume_gigabytes_used",
			Help: "Total volume gigabytes currently in use in OpenStack storage",
		}),
	}
}

// Describe sends the descriptors of each metric over to the provided channel
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
}

// Collect fetches the statistics from OpenStack and delivers them as Prometheus metrics
func (e *OpenStackExporter) Collect(ch chan<- prometheus.Metric) {
	computeLimitsInfo, err := computeLimits.Get(e.ctx, e.computeClient, computeLimits.GetOpts{}).Extract()
	if err != nil {
		log.Println("Failed to get compute limits:", err)
		return
	}

	e.maxTotalInstances.Set(float64(computeLimitsInfo.Absolute.MaxTotalInstances))
	e.totalInstancesUsed.Set(float64(computeLimitsInfo.Absolute.TotalInstancesUsed))
	e.maxTotalCores.Set(float64(computeLimitsInfo.Absolute.MaxTotalCores))
	e.totalCoresUsed.Set(float64(computeLimitsInfo.Absolute.TotalCoresUsed))
	e.maxTotalRAM.Set(float64(computeLimitsInfo.Absolute.MaxTotalRAMSize))
	e.totalRAMUsed.Set(float64(computeLimitsInfo.Absolute.TotalRAMUsed))

	storageLimitsInfo, err := blockstorageLimits.Get(e.ctx, e.storageClient).Extract()
	if err != nil {
		log.Println("Failed to get storage limits:", err)
		return
	}

	e.maxTotalVolumes.Set(float64(storageLimitsInfo.Absolute.MaxTotalVolumes))
	e.totalVolumesUsed.Set(float64(storageLimitsInfo.Absolute.TotalVolumesUsed))
	e.maxTotalSnapshots.Set(float64(storageLimitsInfo.Absolute.MaxTotalSnapshots))
	e.totalSnapshotsUsed.Set(float64(storageLimitsInfo.Absolute.TotalSnapshotsUsed))
	e.maxTotalVolumeGigabytes.Set(float64(storageLimitsInfo.Absolute.MaxTotalVolumeGigabytes))
	e.totalVolumeGigabytesUsed.Set(float64(storageLimitsInfo.Absolute.TotalGigabytesUsed))

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
}

// CustomRoundTripper handles token refreshing
type CustomRoundTripper struct {
	originalTransport http.RoundTripper
	providerClient    *gophercloud.ProviderClient
	ctx               context.Context
}

func (rt *CustomRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := rt.originalTransport.RoundTrip(req)
	if err != nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}

	// Token expired, refresh it
	err = rt.providerClient.ReauthFunc(rt.ctx)
	if err != nil {
		return nil, err
	}

	// Retry the request
	req.Header.Set("X-Auth-Token", rt.providerClient.TokenID)
	return rt.originalTransport.RoundTrip(req)
}

func main() {
	ctx := context.Background()

	authOptions, _, tlsConfig, err := clouds.Parse()
	if err != nil {
		log.Fatalf("Failed to parse cloud config: %v", err)
	}

	providerClient, err := config.NewProviderClient(ctx, authOptions, config.WithTLSConfig(tlsConfig))
	if err != nil {
		log.Fatalf("Failed to create provider client: %v", err)
	}

	providerClient.HTTPClient.Transport = &CustomRoundTripper{
		originalTransport: http.DefaultTransport,
		providerClient:    providerClient,
		ctx:               ctx,
	}

	computeClient, err := openstack.NewComputeV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("Failed to create compute client: %v", err)
	}

	storageClient, err := openstack.NewBlockStorageV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}

	exporter := NewOpenStackExporter(ctx, computeClient, storageClient)
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Prometheus metrics server started at :9183")
	log.Fatal(http.ListenAndServe(":9183", nil))
}
