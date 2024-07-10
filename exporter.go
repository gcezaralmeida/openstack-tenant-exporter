package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	blockstorageLimits "github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/limits"
	computeLimits "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/limits"
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
	maxTotalCores     prometheus.Gauge
	totalCoresUsed    prometheus.Gauge
	maxTotalRAM       prometheus.Gauge
	totalRAMUsed      prometheus.Gauge
	maxTotalVolumes    prometheus.Gauge
	totalVolumesUsed   prometheus.Gauge
	maxTotalSnapshots  prometheus.Gauge
	totalSnapshotsUsed prometheus.Gauge
	maxTotalVolumeGigabytes prometheus.Gauge
	totalVolumeGigabytesUsed prometheus.Gauge
}

// NewOpenStackExporter creates a new OpenStackExporter
// It takes a context, computeClient and storageClient as parameters
// It returns a pointer to OpenStackExporter
func NewOpenStackExporter(ctx context.Context, computeClient, storageClient *gophercloud.ServiceClient) *OpenStackExporter {
	// Create a new OpenStackExporter
	return &OpenStackExporter{
		ctx:           ctx, // The context
		computeClient: computeClient, // The compute client
		storageClient: storageClient, // The storage client

		// Create the metrics
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

// Describe sends the descriptors of each metric over to the provided channel.
//
// It describes the metrics that this exporter provides. This function is called by
// the Prometheus client when it wants to scrape the metrics provided by this
// exporter.
func (e *OpenStackExporter) Describe(ch chan<- *prometheus.Desc) {
	// Describe the maximum total number of instances allowed in OpenStack
	e.maxTotalInstances.Describe(ch)

	// Describe the total number of instances currently in use in OpenStack
	e.totalInstancesUsed.Describe(ch)

	// Describe the maximum total number of cores allowed in OpenStack
	e.maxTotalCores.Describe(ch)

	// Describe the total number of cores currently in use in OpenStack
	e.totalCoresUsed.Describe(ch)

	// Describe the maximum total RAM size allowed in OpenStack
	e.maxTotalRAM.Describe(ch)

	// Describe the total RAM currently in use in OpenStack
	e.totalRAMUsed.Describe(ch)

	// Describe the maximum total number of volumes allowed in OpenStack storage
	e.maxTotalVolumes.Describe(ch)

	// Describe the total volumes currently in use in OpenStack storage
	e.totalVolumesUsed.Describe(ch)

	// Describe the maximum total number of snapshots allowed in OpenStack storage
	e.maxTotalSnapshots.Describe(ch)

	// Describe the total snapshots currently in use in OpenStack storage
	e.totalSnapshotsUsed.Describe(ch)

	// Describe the maximum total volume gigabytes allowed in OpenStack storage
	e.maxTotalVolumeGigabytes.Describe(ch)

	// Describe the total volume gigabytes currently in use in OpenStack storage
	e.totalVolumeGigabytesUsed.Describe(ch)
}

// Collect fetches the statistics from OpenStack and delivers them as Prometheus metrics
//
// Fetches compute and storage limits from OpenStack and sets the corresponding
// Prometheus metrics accordingly.
func (e *OpenStackExporter) Collect(ch chan<- prometheus.Metric) {
	// Fetch compute limits from OpenStack
	computeLimitsInfo, err := computeLimits.Get(e.ctx, e.computeClient, computeLimits.GetOpts{}).Extract()
	if err != nil {
		log.Println("Failed to get compute limits:", err)
		return
	}

	// Set compute limits Prometheus metrics
	e.maxTotalInstances.Set(float64(computeLimitsInfo.Absolute.MaxTotalInstances))
	e.totalInstancesUsed.Set(float64(computeLimitsInfo.Absolute.TotalInstancesUsed))
	e.maxTotalCores.Set(float64(computeLimitsInfo.Absolute.MaxTotalCores))
	e.totalCoresUsed.Set(float64(computeLimitsInfo.Absolute.TotalCoresUsed))
	e.maxTotalRAM.Set(float64(computeLimitsInfo.Absolute.MaxTotalRAMSize))
	e.totalRAMUsed.Set(float64(computeLimitsInfo.Absolute.TotalRAMUsed))

	// Fetch storage limits from OpenStack
	storageLimitsInfo, err := blockstorageLimits.Get(e.ctx, e.storageClient).Extract()
	if err != nil {
		log.Println("Failed to get storage limits:", err)
		return
	}

	// Set storage limits Prometheus metrics
	e.maxTotalVolumes.Set(float64(storageLimitsInfo.Absolute.MaxTotalVolumes))
	e.totalVolumesUsed.Set(float64(storageLimitsInfo.Absolute.TotalVolumesUsed))
	e.maxTotalSnapshots.Set(float64(storageLimitsInfo.Absolute.MaxTotalSnapshots))
	e.totalSnapshotsUsed.Set(float64(storageLimitsInfo.Absolute.TotalSnapshotsUsed))
	e.maxTotalVolumeGigabytes.Set(float64(storageLimitsInfo.Absolute.MaxTotalVolumeGigabytes))
	e.totalVolumeGigabytesUsed.Set(float64(storageLimitsInfo.Absolute.TotalGigabytesUsed))

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
}

// main is the entry point of the application.
//
// It parses the OpenStack cloud configuration, creates a provider client,
// compute client, and storage client, creates an OpenStackExporter,
// registers the exporter with Prometheus, starts a Prometheus metrics
// server, and listens for HTTP requests.
func main() {
	// Create a context
	ctx := context.Background()

	// Log the start of the application
	log.Println("Starting OpenStack Tenant Exporter")

	// Parse the OpenStack cloud configuration
	authOptions, _, tlsConfig, err := clouds.Parse()
	if err != nil {
		log.Fatalf("Failed to parse cloud config: %v", err)
	}

	// Allow reauthentication
	authOptions.AllowReauth = true

	// Create a provider client
	providerClient, err := config.NewProviderClient(ctx, authOptions, config.WithTLSConfig(tlsConfig))
	if err != nil {
		log.Fatalf("Failed to create provider client: %v", err)
	}

	// Log the creation of the provider client
	log.Println("Provider client created")

	// Create a compute client
	computeClient, err := openstack.NewComputeV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("Failed to create compute client: %v", err)
	}

	// Log the creation of the compute client
	log.Println("Compute client created")

	// Create a storage client
	storageClient, err := openstack.NewBlockStorageV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)
	}

	// Log the creation of the storage client
	log.Println("Storage client created")

	// Create an OpenStackExporter
	exporter := NewOpenStackExporter(ctx, computeClient, storageClient)

	// Log the creation of the exporter
	log.Println("OpenStackExporter created")

	// Register the exporter with Prometheus
	prometheus.MustRegister(exporter)

	// Log the registration of the exporter
	log.Println("OpenStackExporter registered with Prometheus")

	// Start a Prometheus metrics server
	http.Handle("/metrics", promhttp.Handler())

	// Log the start of the Prometheus metrics server
	log.Println("Prometheus metrics server started at :9183")

	// Listen for HTTP requests
	log.Fatal(http.ListenAndServe(":9183", nil))
}