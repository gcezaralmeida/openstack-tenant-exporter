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

// Define individual Prometheus Gauge metrics for compute limits
var (
	maxTotalInstances = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openstack_max_total_instances",
		Help: "Maximum total instances allowed in OpenStack",
	})
	totalInstancesUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openstack_total_instances_used",
		Help: "Total instances currently in use in OpenStack",
	})
	maxTotalCores = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openstack_max_total_cores",
		Help: "Maximum total CPU cores allowed in OpenStack",
	})
	totalCoresUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openstack_total_cores_used",
		Help: "Total CPU cores currently in use in OpenStack",
	})
	maxTotalRAM = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openstack_max_total_ram",
		Help: "Maximum total RAM size allowed in OpenStack",
	})
	totalRAMUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openstack_total_ram_used",
		Help: "Total RAM currently in use in OpenStack",
	})
)

// Define individual Prometheus Gauge metrics for storage limits
var (
	maxTotalVolumes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openstack_max_total_volumes",
		Help: "Maximum total volumes allowed in OpenStack storage",
	})
	totalVolumesUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openstack_total_volumes_used",
		Help: "Total volumes currently in use in OpenStack storage",
	})
	maxTotalSnapshots = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openstack_max_total_snapshots",
		Help: "Maximum total snapshots allowed in OpenStack storage",
	})
	totalSnapshotsUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openstack_total_snapshots_used",
		Help: "Total snapshots currently in use in OpenStack storage",
	})
	maxTotalVolumeGigabytes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openstack_max_total_volume_gigabytes",
		Help: "Maximum total volume gigabytes allowed in OpenStack storage",
	})
	totalVolumeGigabytesUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "openstack_total_volume_gigabytes_used",
		Help: "Total volume gigabytes currently in use in OpenStack storage",
	})
)

func init() {
	// Register compute metrics with Prometheus
	prometheus.MustRegister(maxTotalInstances)
	prometheus.MustRegister(totalInstancesUsed)
	prometheus.MustRegister(maxTotalCores)
	prometheus.MustRegister(totalCoresUsed)
	prometheus.MustRegister(maxTotalRAM)
	prometheus.MustRegister(totalRAMUsed)

	// Register storage metrics with Prometheus
	prometheus.MustRegister(maxTotalVolumes)
	prometheus.MustRegister(totalVolumesUsed)
	prometheus.MustRegister(maxTotalSnapshots)
	prometheus.MustRegister(totalSnapshotsUsed)
	prometheus.MustRegister(maxTotalVolumeGigabytes)
	prometheus.MustRegister(totalVolumeGigabytesUsed)
}

func main() {
	// Create a context
	ctx := context.Background()

	// Parse the cloud configuration
	authOptions, _, tlsConfig, err := clouds.Parse()
	if err != nil {
		log.Fatalf("failed to parse cloud config: %v", err)
	}

	// Create a provider client
	providerClient, err := config.NewProviderClient(ctx, authOptions, config.WithTLSConfig(tlsConfig))
	if err != nil {
		log.Fatalf("failed to create provider client: %v", err)
	}

	// Create a compute client
	computeClient, err := openstack.NewComputeV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("failed to create compute client: %v", err)
	}

	// Retrieve compute limits
	computeLimitsOpts := computeLimits.GetOpts{}
	computeLimitsInfo, err := computeLimits.Get(ctx, computeClient, computeLimitsOpts).Extract()
	if err != nil {
		log.Fatalf("failed to get compute limits: %v", err)
	}

	// Set Prometheus metrics for compute limits
	maxTotalInstances.Set(float64(computeLimitsInfo.Absolute.MaxTotalInstances))
	totalInstancesUsed.Set(float64(computeLimitsInfo.Absolute.TotalInstancesUsed))
	maxTotalCores.Set(float64(computeLimitsInfo.Absolute.MaxTotalCores))
	totalCoresUsed.Set(float64(computeLimitsInfo.Absolute.TotalCoresUsed))
	maxTotalRAM.Set(float64(computeLimitsInfo.Absolute.MaxTotalRAMSize))
	totalRAMUsed.Set(float64(computeLimitsInfo.Absolute.TotalRAMUsed))

	// Create a storage client
	storageClient, err := openstack.NewBlockStorageV2(providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		log.Fatalf("failed to create storage client: %v", err)
	}

	// Retrieve storage limits
	storageLimitsInfo, err := blockstorageLimits.Get(ctx, storageClient).Extract()
	if err != nil {
		log.Fatalf("failed to get storage limits: %v", err)
	}

	// Set Prometheus metrics for storage limits
	maxTotalVolumes.Set(float64(storageLimitsInfo.Absolute.MaxTotalVolumes))
	totalVolumesUsed.Set(float64(storageLimitsInfo.Absolute.TotalVolumesUsed))
	maxTotalSnapshots.Set(float64(storageLimitsInfo.Absolute.MaxTotalSnapshots))
	totalSnapshotsUsed.Set(float64(storageLimitsInfo.Absolute.TotalSnapshotsUsed))
	maxTotalVolumeGigabytes.Set(float64(storageLimitsInfo.Absolute.MaxTotalVolumeGigabytes))
	totalVolumeGigabytesUsed.Set(float64(storageLimitsInfo.Absolute.TotalGigabytesUsed))
	//openstackStorageLimits.WithLabelValues("max_total_volumes").Set(float64(limits.Absolute.MaxTotalVolumes))

	// Register HTTP handler for Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	// Start HTTP server in a goroutine
	go func() {
		log.Fatal(http.ListenAndServe(":9183", nil))
	}()

	// Log server start
	log.Println("Prometheus metrics server started at :9183")

	// Keep the program running
	select {}
}
