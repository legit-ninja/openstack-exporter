package exporters

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/amphorae"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/quotas"
	"github.com/prometheus/client_golang/prometheus"
)

// Octavia API v2 entities have two status codes present in the response body.
// The provisioning_status describes the lifecycle status of the entity while the operating_status provides the observed status of the entity.
// Here we put operating_status in metrics value and provisioning_status in metrics label
var knownLoadbalancerStatuses = map[string]int{
	"ONLINE":     0, // Entity is operating normally. All pool members are healthy
	"DRAINING":   1, // The member is not accepting new connections
	"OFFLINE":    2, // Entity is administratively disabled
	"ERROR":      3, // The entity has failed. The member is failing it's health monitoring checks. All of the pool members are in ERROR
	"NO_MONITOR": 4, // No health monitor is configured for this entity and it's status is unknown
}

// The status of the amphora. One of: BOOTING, ALLOCATED, READY, PENDING_CREATE, PENDING_DELETE, DELETED, ERROR.
var knownAmphoraStatuses = map[string]int{
	"BOOTING":        0,
	"ALLOCATED":      1,
	"READY":          2,
	"PENDING_CREATE": 3,
	"PENDING_DELETE": 4,
	"DELETED":        5,
	"ERROR":          6,
}

func mapLoadbalancerStatus(current string) int {
	return mapStatus(knownLoadbalancerStatuses, current)
}

func mapAmphoraStatus(current string) int {
	return mapStatus(knownAmphoraStatuses, current)
}

func mapPoolStatus(current string) int {
	return mapStatus(knownLoadbalancerStatuses, current)
}

type LoadbalancerExporter struct {
	BaseOpenStackExporter
}

// Labels of the per-project quota metrics, identical to the Nova quota metrics.
var defaultLoadbalancerQuotaLabels = []string{"type", "tenant", "tenant_id"}

var defaultLoadbalancerMetrics = []Metric{
	{Name: "total_loadbalancers", Fn: ListAllLoadbalancers},
	{Name: "loadbalancer_status", Labels: []string{"id", "name", "project_id", "operating_status", "provisioning_status", "provider", "vip_address"}},
	{Name: "total_amphorae", Fn: ListAllAmphorae},
	{Name: "amphora_status", Labels: []string{"id", "loadbalancer_id", "compute_id", "status", "role", "lb_network_ip", "ha_ip", "cert_expiration"}},
	{Name: "total_pools", Fn: ListAllPools},
	{Name: "pool_status", Labels: []string{"id", "provisioning_status", "name", "loadbalancers", "protocol", "lb_algorithm", "operating_status", "project_id"}},
	{Name: "quota_loadbalancer", Labels: defaultLoadbalancerQuotaLabels, Fn: ListLoadbalancerQuotas, Slow: true},
	{Name: "quota_listener", Labels: defaultLoadbalancerQuotaLabels, Slow: true},
	{Name: "quota_pool", Labels: defaultLoadbalancerQuotaLabels, Slow: true},
	{Name: "quota_member", Labels: defaultLoadbalancerQuotaLabels, Slow: true},
	{Name: "quota_healthmonitor", Labels: defaultLoadbalancerQuotaLabels, Slow: true},
	{Name: "quota_l7policy", Labels: defaultLoadbalancerQuotaLabels, Slow: true},
	{Name: "quota_l7rule", Labels: defaultLoadbalancerQuotaLabels, Slow: true},
}

func NewLoadbalancerExporter(config *ExporterConfig, logger *slog.Logger) (*LoadbalancerExporter, error) {
	exporter := LoadbalancerExporter{
		BaseOpenStackExporter{
			Name:           "loadbalancer",
			ExporterConfig: *config,
			logger:         logger,
		},
	}

	for _, metric := range defaultLoadbalancerMetrics {
		if exporter.isDeprecatedMetric(&metric) {
			continue
		}
		if !exporter.isSlowMetric(&metric) {
			exporter.AddMetric(metric.Name, metric.Fn, metric.Labels, metric.DeprecatedVersion, nil)
		}
	}

	return &exporter, nil
}

func ListAllLoadbalancers(ctx context.Context, exporter *BaseOpenStackExporter, ch chan<- prometheus.Metric) error {
	var allLoadbalancers []loadbalancers.LoadBalancer
	allPagesLoadbalancers, err := loadbalancers.List(exporter.ClientV2, loadbalancers.ListOpts{}).AllPages(ctx)
	if err != nil {
		return err
	}

	allLoadbalancers, err = loadbalancers.ExtractLoadBalancers(allPagesLoadbalancers)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(exporter.Metrics["total_loadbalancers"].Metric,
		prometheus.GaugeValue, float64(len(allLoadbalancers)))

	// Loadbalancer status metrics
	for _, loadbalancer := range allLoadbalancers {
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["loadbalancer_status"].Metric,
			prometheus.GaugeValue, float64(mapLoadbalancerStatus(loadbalancer.OperatingStatus)), loadbalancer.ID, loadbalancer.Name, loadbalancer.ProjectID,
			loadbalancer.OperatingStatus, loadbalancer.ProvisioningStatus, loadbalancer.Provider, loadbalancer.VipAddress)
	}

	return nil
}

func ListAllAmphorae(ctx context.Context, exporter *BaseOpenStackExporter, ch chan<- prometheus.Metric) error {
	var allAmphorae []amphorae.Amphora
	allPagesAmphorae, err := amphorae.List(exporter.ClientV2, amphorae.ListOpts{}).AllPages(ctx)
	if err != nil {
		return err
	}

	allAmphorae, err = amphorae.ExtractAmphorae(allPagesAmphorae)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(exporter.Metrics["total_amphorae"].Metric,
		prometheus.GaugeValue, float64(len(allAmphorae)))

	// Loadbalancer status metrics
	for _, amphora := range allAmphorae {
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["amphora_status"].Metric,
			prometheus.GaugeValue, float64(mapAmphoraStatus(amphora.Status)), amphora.ID, amphora.LoadbalancerID, amphora.ComputeID, amphora.Status,
			amphora.Role, amphora.LBNetworkIP, amphora.HAIP, amphora.CertExpiration.Format(time.RFC3339))
	}

	return nil
}

func ListAllPools(ctx context.Context, exporter *BaseOpenStackExporter, ch chan<- prometheus.Metric) error {
	var allPools []pools.Pool
	allPagesPools, err := pools.List(exporter.ClientV2, pools.ListOpts{}).AllPages(ctx)
	if err != nil {
		return err
	}

	allPools, err = pools.ExtractPools(allPagesPools)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(exporter.Metrics["total_pools"].Metric,
		prometheus.GaugeValue, float64(len(allPools)))

	for _, pool := range allPools {
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["pool_status"].Metric,
			prometheus.GaugeValue, float64(mapPoolStatus(pool.OperatingStatus)), pool.ID, pool.ProvisioningStatus, pool.Name,
			lbsLabels(pool.Loadbalancers), pool.Protocol, pool.LBMethod, pool.OperatingStatus, pool.ProjectID)
	}

	return nil
}

// loadbalancerQuotaUsage holds the number of Octavia objects owned by a project.
type loadbalancerQuotaUsage struct {
	loadbalancers int
	listeners     int
	pools         int
}

// ListLoadbalancerQuotas emits the Octavia quotas of every project returned
// by GetProjects. Octavia only reports the limits, so the in_use values of
// the load balancer, listener and pool quotas are counted from a single
// load balancer listing (each load balancer lists its listeners and pools).
// Members, health monitors, L7 policies and L7 rules cannot be counted
// without additional per-object calls, so only their limit is emitted.
func ListLoadbalancerQuotas(ctx context.Context, exporter *BaseOpenStackExporter, ch chan<- prometheus.Metric) error {
	allProjects, err := GetProjects(ctx, exporter)
	if err != nil {
		return err
	}

	allPagesLoadbalancers, err := loadbalancers.List(exporter.ClientV2, loadbalancers.ListOpts{ProjectID: exporter.TenantID}).AllPages(ctx)
	if err != nil {
		return err
	}

	allLoadbalancers, err := loadbalancers.ExtractLoadBalancers(allPagesLoadbalancers)
	if err != nil {
		return err
	}

	usage := make(map[string]*loadbalancerQuotaUsage, len(allProjects))
	for _, p := range allProjects {
		usage[p.ID] = &loadbalancerQuotaUsage{}
	}
	for _, lb := range allLoadbalancers {
		u, ok := usage[lb.ProjectID]
		if !ok {
			continue
		}
		u.loadbalancers++
		u.listeners += len(lb.Listeners)
		u.pools += len(lb.Pools)
	}

	for _, p := range allProjects {
		quota, err := quotas.Get(ctx, exporter.ClientV2, p.ID).Extract()
		if err != nil {
			return err
		}
		if quota == nil {
			return fmt.Errorf("empty load balancer quota response for project %s", p.ID)
		}

		u := usage[p.ID]
		collectLoadbalancerQuota(ch, exporter, "quota_loadbalancer", quota.Loadbalancer, &u.loadbalancers, p.Name, p.ID)
		collectLoadbalancerQuota(ch, exporter, "quota_listener", quota.Listener, &u.listeners, p.Name, p.ID)
		collectLoadbalancerQuota(ch, exporter, "quota_pool", quota.Pool, &u.pools, p.Name, p.ID)
		collectLoadbalancerQuota(ch, exporter, "quota_member", quota.Member, nil, p.Name, p.ID)
		collectLoadbalancerQuota(ch, exporter, "quota_healthmonitor", quota.Healthmonitor, nil, p.Name, p.ID)
		collectLoadbalancerQuota(ch, exporter, "quota_l7policy", quota.L7Policy, nil, p.Name, p.ID)
		collectLoadbalancerQuota(ch, exporter, "quota_l7rule", quota.L7Rule, nil, p.Name, p.ID)
	}

	return nil
}

// collectLoadbalancerQuota emits the limit of an Octavia quota and, when the
// usage is known, its in_use value. Octavia has no reservations, so no
// reserved value is emitted. The metric is skipped when it has been disabled
// with --disable-metric.
func collectLoadbalancerQuota(ch chan<- prometheus.Metric, exporter *BaseOpenStackExporter, name string, limit int, inUse *int, projectName, projectID string) {
	metric, ok := exporter.Metrics[name]
	if !ok {
		return
	}

	if inUse != nil {
		ch <- prometheus.MustNewConstMetric(metric.Metric, prometheus.GaugeValue, float64(*inUse), "in_use", projectName, projectID)
	}
	ch <- prometheus.MustNewConstMetric(metric.Metric, prometheus.GaugeValue, float64(limit), "limit", projectName, projectID)
}

func lbsLabels(lbs []pools.LoadBalancerID) string {
	label := ""
	for i, l := range lbs {
		if i != 0 {
			label += ","
		}
		label += l.ID
	}

	return label
}
