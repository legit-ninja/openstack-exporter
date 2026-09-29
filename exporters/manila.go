package exporters

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
	"github.com/prometheus/client_golang/prometheus"
)

type ManilaExporter struct {
	BaseOpenStackExporter
}

// Microversion required for GET /quota-sets/{project_id}/detail, which
// returns the in_use, reserved and limit value of each quota. Available since
// Ocata (2017).
const manilaQuotaDetailMicroversion = "2.25"

// Labels of the per-project quota metrics, identical to the Nova quota metrics.
var defaultManilaQuotaLabels = []string{"type", "tenant", "tenant_id"}

var defaultManilaMetrics = []Metric{
	{Name: "shares_counter", Fn: CountShares},
	{Name: "share_gb", Labels: []string{"id", "name", "status", "availability_zone", "share_type", "share_proto", "share_type_name", "project_id"}, Fn: nil},
	{Name: "share_status", Labels: []string{"id", "name", "status", "size", "share_type", "share_proto", "share_type_name", "project_id"}, Fn: ListShareStatus},
	{Name: "share_status_counter", Labels: []string{"status"}, Fn: nil},
	{Name: "quota_shares", Labels: defaultManilaQuotaLabels, Fn: ListShareQuotas, Slow: true},
	{Name: "quota_gigabytes", Labels: defaultManilaQuotaLabels, Slow: true},
	{Name: "quota_snapshots", Labels: defaultManilaQuotaLabels, Slow: true},
	{Name: "quota_snapshot_gigabytes", Labels: defaultManilaQuotaLabels, Slow: true},
	{Name: "quota_share_networks", Labels: defaultManilaQuotaLabels, Slow: true},
}

func NewManilaExporter(config *ExporterConfig, logger *slog.Logger) (*ManilaExporter, error) {
	exporter := ManilaExporter{
		BaseOpenStackExporter{
			Name:           "sharev2",
			ExporterConfig: *config,
			logger:         logger,
		},
	}

	for _, metric := range defaultManilaMetrics {
		if exporter.isDeprecatedMetric(&metric) {
			continue
		}
		if !exporter.isSlowMetric(&metric) {
			exporter.AddMetric(metric.Name, metric.Fn, metric.Labels, metric.DeprecatedVersion, nil)
		}
	}

	return &exporter, nil
}

func CountShares(ctx context.Context, exporter *BaseOpenStackExporter, ch chan<- prometheus.Metric) error {
	var allShares []shares.Share

	allPagesShares, err := shares.ListDetail(exporter.ClientV2, shares.ListOpts{AllTenants: true}).AllPages(ctx)
	if err != nil {
		return err
	}

	allShares, err = shares.ExtractShares(allPagesShares)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(exporter.Metrics["shares_counter"].Metric,
		prometheus.GaugeValue, float64(len(allShares)))

	// share_gb metrics
	for _, share := range allShares {
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["share_gb"].Metric,
			prometheus.GaugeValue, float64(share.Size), share.ID, share.Name,
			share.Status, share.AvailabilityZone, share.ShareType, share.ShareProto, share.ShareTypeName, share.ProjectID)
	}

	share_status_counter := map[string]int{
		"creating":              0,
		"available":             0,
		"updating":              0,
		"migrating":             0,
		"migration_error":       0,
		"extending":             0,
		"deleting":              0,
		"shrinking":             0,
		"error":                 0,
		"error_deleting":        0,
		"shrinking_error":       0,
		"reverting_error":       0,
		"restoring":             0,
		"reverting":             0,
		"managing":              0,
		"unmanaging":            0,
		"reverting_to_snapshot": 0,
		"soft_deleting":         0,
		"inactive":              0,
	}

	for _, share := range allShares {
		share_status_counter[share.Status]++
	}

	// Share status counter metrics
	for status, count := range share_status_counter {
		ch <- prometheus.MustNewConstMetric(
			exporter.Metrics["share_status_counter"].Metric,
			prometheus.GaugeValue,
			float64(count),
			status)
	}

	return nil
}

func ListShareStatus(ctx context.Context, exporter *BaseOpenStackExporter, ch chan<- prometheus.Metric) error {
	var allShares []shares.Share

	allPagesShares, err := shares.ListDetail(exporter.ClientV2, shares.ListOpts{AllTenants: true}).AllPages(ctx)
	if err != nil {
		return err
	}

	allShares, err = shares.ExtractShares(allPagesShares)
	if err != nil {
		return err
	}

	// Share status metrics
	for _, share := range allShares {
		ch <- prometheus.MustNewConstMetric(exporter.Metrics["share_status"].Metric,
			prometheus.GaugeValue, float64(mapVolumeStatus(share.Status)), share.ID, share.Name,
			share.Status, strconv.Itoa(share.Size), share.ShareType, share.ShareProto, share.ShareTypeName, share.ProjectID)
	}

	return nil
}

// manilaQuotaDetail is a single entry of the Manila quota set detail response.
type manilaQuotaDetail struct {
	InUse    int `json:"in_use"`
	Reserved int `json:"reserved"`
	Limit    int `json:"limit"`
}

// manilaQuotaSetDetail is the Manila quota set detail response at
// microversion 2.25. gophercloud has no package for the Manila quota API.
type manilaQuotaSetDetail struct {
	Shares            manilaQuotaDetail `json:"shares"`
	Gigabytes         manilaQuotaDetail `json:"gigabytes"`
	Snapshots         manilaQuotaDetail `json:"snapshots"`
	SnapshotGigabytes manilaQuotaDetail `json:"snapshot_gigabytes"`
	ShareNetworks     manilaQuotaDetail `json:"share_networks"`
}

// getManilaQuotaSetDetail returns the quotas and their usage for a project.
// The client must be set to microversion 2.25 or later.
func getManilaQuotaSetDetail(ctx context.Context, client *gophercloud.ServiceClient, projectID string) (*manilaQuotaSetDetail, error) {
	var body struct {
		QuotaSet *manilaQuotaSetDetail `json:"quota_set"`
	}

	_, err := client.Get(ctx, client.ServiceURL("quota-sets", projectID, "detail"), &body, nil)
	if err != nil {
		return nil, err
	}
	if body.QuotaSet == nil {
		return nil, fmt.Errorf("empty share quota set response for project %s", projectID)
	}

	return body.QuotaSet, nil
}

// ListShareQuotas emits the Manila quotas of every project returned by
// GetProjects.
func ListShareQuotas(ctx context.Context, exporter *BaseOpenStackExporter, ch chan<- prometheus.Metric) error {
	allProjects, err := GetProjects(ctx, exporter)
	if err != nil {
		return err
	}

	// The quota detail requires microversion 2.25; a copy of the client is
	// used to not affect other metrics.
	client := *exporter.ClientV2
	client.Microversion = manilaQuotaDetailMicroversion

	for _, p := range allProjects {
		quotaSet, err := getManilaQuotaSetDetail(ctx, &client, p.ID)
		if err != nil {
			return err
		}

		collectManilaQuotaDetail(ch, exporter, "quota_shares", quotaSet.Shares, p.Name, p.ID)
		collectManilaQuotaDetail(ch, exporter, "quota_gigabytes", quotaSet.Gigabytes, p.Name, p.ID)
		collectManilaQuotaDetail(ch, exporter, "quota_snapshots", quotaSet.Snapshots, p.Name, p.ID)
		collectManilaQuotaDetail(ch, exporter, "quota_snapshot_gigabytes", quotaSet.SnapshotGigabytes, p.Name, p.ID)
		collectManilaQuotaDetail(ch, exporter, "quota_share_networks", quotaSet.ShareNetworks, p.Name, p.ID)
	}

	return nil
}

// collectManilaQuotaDetail emits the in_use, reserved and limit values of a
// Manila quota. The metric is skipped when it has been disabled with
// --disable-metric.
func collectManilaQuotaDetail(ch chan<- prometheus.Metric, exporter *BaseOpenStackExporter, name string, q manilaQuotaDetail, projectName, projectID string) {
	metric, ok := exporter.Metrics[name]
	if !ok {
		return
	}

	ch <- prometheus.MustNewConstMetric(metric.Metric, prometheus.GaugeValue, float64(q.InUse), "in_use", projectName, projectID)
	ch <- prometheus.MustNewConstMetric(metric.Metric, prometheus.GaugeValue, float64(q.Reserved), "reserved", projectName, projectID)
	ch <- prometheus.MustNewConstMetric(metric.Metric, prometheus.GaugeValue, float64(q.Limit), "limit", projectName, projectID)
}
