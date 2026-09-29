package exporters

import (
	"net/http"
	"os"
	"strings"

	"github.com/jarcoal/httpmock"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

type ManilaTestSuite struct {
	BaseOpenStackTestSuite
}

var manilaExpectedUp = `
# HELP openstack_sharev2_quota_gigabytes quota_gigabytes
# TYPE openstack_sharev2_quota_gigabytes gauge
openstack_sharev2_quota_gigabytes{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="in_use"} 120
openstack_sharev2_quota_gigabytes{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="limit"} 2048
openstack_sharev2_quota_gigabytes{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="reserved"} 10
openstack_sharev2_quota_gigabytes{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="in_use"} 0
openstack_sharev2_quota_gigabytes{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="limit"} 1000
openstack_sharev2_quota_gigabytes{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="reserved"} 0
openstack_sharev2_quota_gigabytes{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="in_use"} 0
openstack_sharev2_quota_gigabytes{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="limit"} 1000
openstack_sharev2_quota_gigabytes{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="reserved"} 0
openstack_sharev2_quota_gigabytes{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="in_use"} 0
openstack_sharev2_quota_gigabytes{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="limit"} 1000
openstack_sharev2_quota_gigabytes{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="reserved"} 0
openstack_sharev2_quota_gigabytes{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="in_use"} 0
openstack_sharev2_quota_gigabytes{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="limit"} 1000
openstack_sharev2_quota_gigabytes{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="reserved"} 0
openstack_sharev2_quota_gigabytes{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="in_use"} 0
openstack_sharev2_quota_gigabytes{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="limit"} 1000
openstack_sharev2_quota_gigabytes{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="reserved"} 0
openstack_sharev2_quota_gigabytes{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="in_use"} 0
openstack_sharev2_quota_gigabytes{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="limit"} 1000
openstack_sharev2_quota_gigabytes{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="reserved"} 0
openstack_sharev2_quota_gigabytes{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="in_use"} 0
openstack_sharev2_quota_gigabytes{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="limit"} 1000
openstack_sharev2_quota_gigabytes{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="reserved"} 0
# HELP openstack_sharev2_quota_share_networks quota_share_networks
# TYPE openstack_sharev2_quota_share_networks gauge
openstack_sharev2_quota_share_networks{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="in_use"} 1
openstack_sharev2_quota_share_networks{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="limit"} 20
openstack_sharev2_quota_share_networks{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="reserved"} 0
openstack_sharev2_quota_share_networks{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="in_use"} 0
openstack_sharev2_quota_share_networks{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="limit"} 10
openstack_sharev2_quota_share_networks{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="reserved"} 0
openstack_sharev2_quota_share_networks{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="in_use"} 0
openstack_sharev2_quota_share_networks{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="limit"} 10
openstack_sharev2_quota_share_networks{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="reserved"} 0
openstack_sharev2_quota_share_networks{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="in_use"} 0
openstack_sharev2_quota_share_networks{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="limit"} 10
openstack_sharev2_quota_share_networks{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="reserved"} 0
openstack_sharev2_quota_share_networks{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="in_use"} 0
openstack_sharev2_quota_share_networks{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="limit"} 10
openstack_sharev2_quota_share_networks{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="reserved"} 0
openstack_sharev2_quota_share_networks{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="in_use"} 0
openstack_sharev2_quota_share_networks{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="limit"} 10
openstack_sharev2_quota_share_networks{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="reserved"} 0
openstack_sharev2_quota_share_networks{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="in_use"} 0
openstack_sharev2_quota_share_networks{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="limit"} 10
openstack_sharev2_quota_share_networks{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="reserved"} 0
openstack_sharev2_quota_share_networks{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="in_use"} 0
openstack_sharev2_quota_share_networks{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="limit"} 10
openstack_sharev2_quota_share_networks{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="reserved"} 0
# HELP openstack_sharev2_quota_shares quota_shares
# TYPE openstack_sharev2_quota_shares gauge
openstack_sharev2_quota_shares{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="in_use"} 3
openstack_sharev2_quota_shares{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="limit"} 50
openstack_sharev2_quota_shares{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="reserved"} 1
openstack_sharev2_quota_shares{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="in_use"} 0
openstack_sharev2_quota_shares{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="limit"} 50
openstack_sharev2_quota_shares{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="reserved"} 0
openstack_sharev2_quota_shares{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="in_use"} 0
openstack_sharev2_quota_shares{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="limit"} 50
openstack_sharev2_quota_shares{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="reserved"} 0
openstack_sharev2_quota_shares{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="in_use"} 0
openstack_sharev2_quota_shares{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="limit"} 50
openstack_sharev2_quota_shares{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="reserved"} 0
openstack_sharev2_quota_shares{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="in_use"} 0
openstack_sharev2_quota_shares{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="limit"} 50
openstack_sharev2_quota_shares{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="reserved"} 0
openstack_sharev2_quota_shares{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="in_use"} 0
openstack_sharev2_quota_shares{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="limit"} 50
openstack_sharev2_quota_shares{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="reserved"} 0
openstack_sharev2_quota_shares{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="in_use"} 0
openstack_sharev2_quota_shares{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="limit"} 50
openstack_sharev2_quota_shares{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="reserved"} 0
openstack_sharev2_quota_shares{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="in_use"} 0
openstack_sharev2_quota_shares{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="limit"} 50
openstack_sharev2_quota_shares{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="reserved"} 0
# HELP openstack_sharev2_quota_snapshot_gigabytes quota_snapshot_gigabytes
# TYPE openstack_sharev2_quota_snapshot_gigabytes gauge
openstack_sharev2_quota_snapshot_gigabytes{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="in_use"} 20
openstack_sharev2_quota_snapshot_gigabytes{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="limit"} 1000
openstack_sharev2_quota_snapshot_gigabytes{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="reserved"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="in_use"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="limit"} 1000
openstack_sharev2_quota_snapshot_gigabytes{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="reserved"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="in_use"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="limit"} 1000
openstack_sharev2_quota_snapshot_gigabytes{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="reserved"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="in_use"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="limit"} 1000
openstack_sharev2_quota_snapshot_gigabytes{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="reserved"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="in_use"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="limit"} 1000
openstack_sharev2_quota_snapshot_gigabytes{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="reserved"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="in_use"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="limit"} 1000
openstack_sharev2_quota_snapshot_gigabytes{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="reserved"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="in_use"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="limit"} 1000
openstack_sharev2_quota_snapshot_gigabytes{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="reserved"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="in_use"} 0
openstack_sharev2_quota_snapshot_gigabytes{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="limit"} 1000
openstack_sharev2_quota_snapshot_gigabytes{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="reserved"} 0
# HELP openstack_sharev2_quota_snapshots quota_snapshots
# TYPE openstack_sharev2_quota_snapshots gauge
openstack_sharev2_quota_snapshots{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="in_use"} 2
openstack_sharev2_quota_snapshots{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="limit"} -1
openstack_sharev2_quota_snapshots{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="reserved"} 0
openstack_sharev2_quota_snapshots{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="in_use"} 0
openstack_sharev2_quota_snapshots{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="limit"} 50
openstack_sharev2_quota_snapshots{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="reserved"} 0
openstack_sharev2_quota_snapshots{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="in_use"} 0
openstack_sharev2_quota_snapshots{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="limit"} 50
openstack_sharev2_quota_snapshots{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="reserved"} 0
openstack_sharev2_quota_snapshots{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="in_use"} 0
openstack_sharev2_quota_snapshots{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="limit"} 50
openstack_sharev2_quota_snapshots{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="reserved"} 0
openstack_sharev2_quota_snapshots{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="in_use"} 0
openstack_sharev2_quota_snapshots{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="limit"} 50
openstack_sharev2_quota_snapshots{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="reserved"} 0
openstack_sharev2_quota_snapshots{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="in_use"} 0
openstack_sharev2_quota_snapshots{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="limit"} 50
openstack_sharev2_quota_snapshots{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="reserved"} 0
openstack_sharev2_quota_snapshots{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="in_use"} 0
openstack_sharev2_quota_snapshots{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="limit"} 50
openstack_sharev2_quota_snapshots{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="reserved"} 0
openstack_sharev2_quota_snapshots{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="in_use"} 0
openstack_sharev2_quota_snapshots{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="limit"} 50
openstack_sharev2_quota_snapshots{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="reserved"} 0
# HELP openstack_sharev2_share_gb share_gb
# TYPE openstack_sharev2_share_gb gauge
openstack_sharev2_share_gb{availability_zone="az1",id="4be93e2e-ffff-ffff-ffff-603e3ec2a5d6",name="share-test",project_id="ffff8fa0ca1a468db8ad00970c1effff",share_proto="NFS",share_type="az1",share_type_name="",status="available"} 1
# HELP openstack_sharev2_share_status share_status
# TYPE openstack_sharev2_share_status gauge
openstack_sharev2_share_status{id="4be93e2e-ffff-ffff-ffff-603e3ec2a5d6",name="share-test",project_id="ffff8fa0ca1a468db8ad00970c1effff",share_proto="NFS",share_type="az1",share_type_name="",size="1",status="available"} 1
# HELP openstack_sharev2_share_status_counter share_status_counter
# TYPE openstack_sharev2_share_status_counter gauge
openstack_sharev2_share_status_counter{status="available"} 1
openstack_sharev2_share_status_counter{status="creating"} 0
openstack_sharev2_share_status_counter{status="deleting"} 0
openstack_sharev2_share_status_counter{status="error"} 0
openstack_sharev2_share_status_counter{status="error_deleting"} 0
openstack_sharev2_share_status_counter{status="extending"} 0
openstack_sharev2_share_status_counter{status="inactive"} 0
openstack_sharev2_share_status_counter{status="managing"} 0
openstack_sharev2_share_status_counter{status="migrating"} 0
openstack_sharev2_share_status_counter{status="migration_error"} 0
openstack_sharev2_share_status_counter{status="restoring"} 0
openstack_sharev2_share_status_counter{status="reverting"} 0
openstack_sharev2_share_status_counter{status="reverting_error"} 0
openstack_sharev2_share_status_counter{status="reverting_to_snapshot"} 0
openstack_sharev2_share_status_counter{status="shrinking"} 0
openstack_sharev2_share_status_counter{status="shrinking_error"} 0
openstack_sharev2_share_status_counter{status="soft_deleting"} 0
openstack_sharev2_share_status_counter{status="unmanaging"} 0
openstack_sharev2_share_status_counter{status="updating"} 0
# HELP openstack_sharev2_shares_counter shares_counter
# TYPE openstack_sharev2_shares_counter gauge
openstack_sharev2_shares_counter 1
# HELP openstack_sharev2_up up
# TYPE openstack_sharev2_up gauge
openstack_sharev2_up 1
`

func (suite *ManilaTestSuite) TestManilaExporter() {
	err := testutil.CollectAndCompare(*suite.Exporter, strings.NewReader(manilaExpectedUp))
	assert.NoError(suite.T(), err)
}

// The quota detail must be requested at microversion 2.25, while the share
// listing must keep using the exporter's default client.
func (suite *ManilaTestSuite) TestManilaQuotaMicroversion() {
	const adminQuotaURL = "/shares/v2/quota-sets/0c4e939acacf4376bdcd1129f1a054ad/detail"

	quotaVersions := make(chan string, 2)
	quotaFixture, err := os.ReadFile(suite.FixturePath("manila_quota_sets_detail_admin"))
	suite.Require().NoError(err)
	httpmock.RegisterResponder("GET", suite.MakeURL(adminQuotaURL, ""),
		func(req *http.Request) (*http.Response, error) {
			quotaVersions <- req.Header.Get("X-OpenStack-Manila-API-Version")
			return jsonFixtureResponse(req, quotaFixture), nil
		})

	shareVersions := make(chan string, 4)
	sharesFixture, err := os.ReadFile(suite.FixturePath("manila_shares"))
	suite.Require().NoError(err)
	httpmock.RegisterResponder("GET", suite.MakeURL("/shares/v2/shares/detail?all_tenants=true", ""),
		func(req *http.Request) (*http.Response, error) {
			shareVersions <- req.Header.Get("X-OpenStack-Manila-API-Version")
			return jsonFixtureResponse(req, sharesFixture), nil
		})

	families := suite.GatherFamilies(*suite.Exporter)
	assert.Contains(suite.T(), families, "openstack_sharev2_quota_shares")
	assert.Contains(suite.T(), families, "openstack_sharev2_shares_counter")

	// One quota request for the admin project, two share listings
	// (shares_counter and share_status).
	assert.Len(suite.T(), quotaVersions, 1)
	assert.Len(suite.T(), shareVersions, 2)
	close(quotaVersions)
	close(shareVersions)
	for v := range quotaVersions {
		assert.Equal(suite.T(), manilaQuotaDetailMicroversion, v)
	}
	for v := range shareVersions {
		assert.Empty(suite.T(), v)
	}
}

// The quota metrics are slow and must be dropped by --disable-slow-metrics.
func (suite *ManilaTestSuite) TestManilaQuotaSlowMetricsDisabled() {
	exporter := suite.NewTestExporter(nil, true)

	families := suite.GatherFamilies(exporter)
	for name := range families {
		assert.NotContains(suite.T(), name, "openstack_sharev2_quota_")
	}
	assert.Contains(suite.T(), families, "openstack_sharev2_shares_counter")
}

// A secondary quota metric disabled with --disable-metric must not be emitted
// and must not break the collection of the other quota metrics.
func (suite *ManilaTestSuite) TestManilaQuotaMetricDisabled() {
	exporter := suite.NewTestExporter([]string{"sharev2-quota_share_networks"}, false)

	families := suite.GatherFamilies(exporter)
	assert.NotContains(suite.T(), families, "openstack_sharev2_quota_share_networks")
	assert.Len(suite.T(), families["openstack_sharev2_quota_gigabytes"].GetMetric(), 8*3)
}

// jsonFixtureResponse builds a JSON response to req, like the responses
// installed by SetResponseFromFixture.
func jsonFixtureResponse(req *http.Request, data []byte) *http.Response {
	resp := httpmock.NewBytesResponse(200, data)
	resp.Header = http.Header{"Content-Type": []string{"application/json"}}
	resp.Request = req
	return resp
}

// A failing project must be skipped instead of dropping the metrics of all
// projects.
func (suite *ManilaTestSuite) TestManilaQuotasSkipFailingProject() {
	suite.FailRequest("/shares/v2/quota-sets/"+failingProjectID+"/detail", 404)

	families := suite.GatherFamilies(*suite.Exporter)
	suite.AssertProjectSkipped(families, "openstack_sharev2_quota_shares", "openstack_sharev2_quota_share_networks")
}
