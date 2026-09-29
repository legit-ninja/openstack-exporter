package exporters

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

type LoadbalancerTestSuite struct {
	BaseOpenStackTestSuite
}

var loadbalancerExpectedUp = `
# HELP openstack_loadbalancer_pool_status pool_status
# TYPE openstack_loadbalancer_pool_status gauge
openstack_loadbalancer_pool_status{id="7fe02a57-6fca-4578-b76f-20c3dffb2773",lb_algorithm="ROUND_ROBIN",loadbalancers="e7284bb2-f46a-42ca-8c9b-e08671255125",name="disabled_pool",operating_status="OFFLINE",project_id="8b1632d90bfe407787d9996b7f662fd7",protocol="TCP",provisioning_status="ACTIVE"} 2
openstack_loadbalancer_pool_status{id="ca00ed86-94e3-440e-95c6-ffa35531081e",lb_algorithm="ROUND_ROBIN",loadbalancers="e7284bb2-f46a-42ca-8c9b-e08671255125",name="my_test_pool",operating_status="ERROR",project_id="8b1632d90bfe407787d9996b7f662fd7",protocol="TCP",provisioning_status="ACTIVE"} 3
# HELP openstack_loadbalancer_amphora_status amphora_status
# TYPE openstack_loadbalancer_amphora_status gauge
openstack_loadbalancer_amphora_status{cert_expiration="2020-08-08T23:44:31Z",compute_id="667bb225-69aa-44b1-8908-694dc624c267",ha_ip="10.0.0.6",id="45f40289-0551-483a-b089-47214bc2a8a4",lb_network_ip="192.168.0.6",loadbalancer_id="882f2a9d-9d53-4bd0-b0e9-08e9d0de11f9",role="MASTER",status="READY"} 2
openstack_loadbalancer_amphora_status{cert_expiration="2020-08-08T23:44:30Z",compute_id="9cd0f9a2-fe12-42fc-a7e3-5b6fbbe20395",ha_ip="10.0.0.6",id="7f890893-ced0-46ed-8697-33415d070e5a",lb_network_ip="192.168.0.17",loadbalancer_id="882f2a9d-9d53-4bd0-b0e9-08e9d0de11f9",role="BACKUP",status="READY"} 2
# HELP openstack_loadbalancer_loadbalancer_status loadbalancer_status
# TYPE openstack_loadbalancer_loadbalancer_status gauge
openstack_loadbalancer_loadbalancer_status{id="607226db-27ef-4d41-ae89-f2a800e9c2db",name="best_load_balancer",operating_status="ONLINE",project_id="e3cd678b11784734bc366148aa37580e",provider="octavia",provisioning_status="ACTIVE",vip_address="203.0.113.50"} 0
openstack_loadbalancer_loadbalancer_status{id="a1f3c9e2-6b7d-4e8f-9a0b-1c2d3e4f5a01",name="admin_lb_1",operating_status="ONLINE",project_id="0c4e939acacf4376bdcd1129f1a054ad",provider="amphora",provisioning_status="ACTIVE",vip_address="203.0.113.51"} 0
openstack_loadbalancer_loadbalancer_status{id="a1f3c9e2-6b7d-4e8f-9a0b-1c2d3e4f5a02",name="admin_lb_2",operating_status="OFFLINE",project_id="0c4e939acacf4376bdcd1129f1a054ad",provider="amphora",provisioning_status="ACTIVE",vip_address="203.0.113.52"} 2
# HELP openstack_loadbalancer_quota_healthmonitor quota_healthmonitor
# TYPE openstack_loadbalancer_quota_healthmonitor gauge
openstack_loadbalancer_quota_healthmonitor{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="limit"} 10
openstack_loadbalancer_quota_healthmonitor{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="limit"} -1
openstack_loadbalancer_quota_healthmonitor{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="limit"} -1
openstack_loadbalancer_quota_healthmonitor{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="limit"} -1
openstack_loadbalancer_quota_healthmonitor{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="limit"} -1
openstack_loadbalancer_quota_healthmonitor{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="limit"} -1
openstack_loadbalancer_quota_healthmonitor{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="limit"} -1
openstack_loadbalancer_quota_healthmonitor{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="limit"} -1
# HELP openstack_loadbalancer_quota_l7policy quota_l7policy
# TYPE openstack_loadbalancer_quota_l7policy gauge
openstack_loadbalancer_quota_l7policy{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="limit"} 50
openstack_loadbalancer_quota_l7policy{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="limit"} -1
openstack_loadbalancer_quota_l7policy{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="limit"} -1
openstack_loadbalancer_quota_l7policy{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="limit"} -1
openstack_loadbalancer_quota_l7policy{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="limit"} -1
openstack_loadbalancer_quota_l7policy{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="limit"} -1
openstack_loadbalancer_quota_l7policy{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="limit"} -1
openstack_loadbalancer_quota_l7policy{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="limit"} -1
# HELP openstack_loadbalancer_quota_l7rule quota_l7rule
# TYPE openstack_loadbalancer_quota_l7rule gauge
openstack_loadbalancer_quota_l7rule{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="limit"} 100
openstack_loadbalancer_quota_l7rule{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="limit"} -1
openstack_loadbalancer_quota_l7rule{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="limit"} -1
openstack_loadbalancer_quota_l7rule{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="limit"} -1
openstack_loadbalancer_quota_l7rule{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="limit"} -1
openstack_loadbalancer_quota_l7rule{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="limit"} -1
openstack_loadbalancer_quota_l7rule{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="limit"} -1
openstack_loadbalancer_quota_l7rule{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="limit"} -1
# HELP openstack_loadbalancer_quota_listener quota_listener
# TYPE openstack_loadbalancer_quota_listener gauge
openstack_loadbalancer_quota_listener{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="in_use"} 2
openstack_loadbalancer_quota_listener{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="limit"} 20
openstack_loadbalancer_quota_listener{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="in_use"} 0
openstack_loadbalancer_quota_listener{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="limit"} -1
openstack_loadbalancer_quota_listener{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="in_use"} 0
openstack_loadbalancer_quota_listener{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="limit"} -1
openstack_loadbalancer_quota_listener{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="in_use"} 0
openstack_loadbalancer_quota_listener{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="limit"} -1
openstack_loadbalancer_quota_listener{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="in_use"} 0
openstack_loadbalancer_quota_listener{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="limit"} -1
openstack_loadbalancer_quota_listener{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="in_use"} 0
openstack_loadbalancer_quota_listener{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="limit"} -1
openstack_loadbalancer_quota_listener{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="in_use"} 0
openstack_loadbalancer_quota_listener{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="limit"} -1
openstack_loadbalancer_quota_listener{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="in_use"} 0
openstack_loadbalancer_quota_listener{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="limit"} -1
# HELP openstack_loadbalancer_quota_loadbalancer quota_loadbalancer
# TYPE openstack_loadbalancer_quota_loadbalancer gauge
openstack_loadbalancer_quota_loadbalancer{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="in_use"} 2
openstack_loadbalancer_quota_loadbalancer{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="limit"} 5
openstack_loadbalancer_quota_loadbalancer{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="in_use"} 0
openstack_loadbalancer_quota_loadbalancer{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="limit"} 10
openstack_loadbalancer_quota_loadbalancer{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="in_use"} 0
openstack_loadbalancer_quota_loadbalancer{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="limit"} 10
openstack_loadbalancer_quota_loadbalancer{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="in_use"} 0
openstack_loadbalancer_quota_loadbalancer{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="limit"} 10
openstack_loadbalancer_quota_loadbalancer{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="in_use"} 0
openstack_loadbalancer_quota_loadbalancer{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="limit"} 10
openstack_loadbalancer_quota_loadbalancer{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="in_use"} 0
openstack_loadbalancer_quota_loadbalancer{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="limit"} 10
openstack_loadbalancer_quota_loadbalancer{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="in_use"} 0
openstack_loadbalancer_quota_loadbalancer{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="limit"} 10
openstack_loadbalancer_quota_loadbalancer{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="in_use"} 0
openstack_loadbalancer_quota_loadbalancer{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="limit"} 10
# HELP openstack_loadbalancer_quota_member quota_member
# TYPE openstack_loadbalancer_quota_member gauge
openstack_loadbalancer_quota_member{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="limit"} 100
openstack_loadbalancer_quota_member{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="limit"} -1
openstack_loadbalancer_quota_member{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="limit"} -1
openstack_loadbalancer_quota_member{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="limit"} -1
openstack_loadbalancer_quota_member{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="limit"} -1
openstack_loadbalancer_quota_member{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="limit"} -1
openstack_loadbalancer_quota_member{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="limit"} -1
openstack_loadbalancer_quota_member{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="limit"} -1
# HELP openstack_loadbalancer_quota_pool quota_pool
# TYPE openstack_loadbalancer_quota_pool gauge
openstack_loadbalancer_quota_pool{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="in_use"} 2
openstack_loadbalancer_quota_pool{tenant="admin",tenant_id="0c4e939acacf4376bdcd1129f1a054ad",type="limit"} 10
openstack_loadbalancer_quota_pool{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="in_use"} 0
openstack_loadbalancer_quota_pool{tenant="alt_demo",tenant_id="fdb8424c4e4f4c0ba32c52e2de3bd80e",type="limit"} -1
openstack_loadbalancer_quota_pool{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="in_use"} 0
openstack_loadbalancer_quota_pool{tenant="demo",tenant_id="0cbd49cbf76d405d9c86562e1d579bd3",type="limit"} -1
openstack_loadbalancer_quota_pool{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="in_use"} 0
openstack_loadbalancer_quota_pool{tenant="invisible_to_admin",tenant_id="5961c443439d4fcebe42643723755e9d",type="limit"} -1
openstack_loadbalancer_quota_pool{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="in_use"} 0
openstack_loadbalancer_quota_pool{tenant="service",tenant_id="3d594eb0f04741069dbbb521635b21c7",type="limit"} -1
openstack_loadbalancer_quota_pool{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="in_use"} 0
openstack_loadbalancer_quota_pool{tenant="swifttenanttest1",tenant_id="43ebde53fc314b1c9ea2b8c5dc744927",type="limit"} -1
openstack_loadbalancer_quota_pool{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="in_use"} 0
openstack_loadbalancer_quota_pool{tenant="swifttenanttest2",tenant_id="2db68fed84324f29bb73130c6c2094fb",type="limit"} -1
openstack_loadbalancer_quota_pool{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="in_use"} 0
openstack_loadbalancer_quota_pool{tenant="swifttenanttest4",tenant_id="4b1eb781a47440acb8af9850103e537f",type="limit"} -1
# HELP openstack_loadbalancer_total_amphorae total_amphorae
# TYPE openstack_loadbalancer_total_amphorae gauge
openstack_loadbalancer_total_amphorae 2
# HELP openstack_loadbalancer_total_loadbalancers total_loadbalancers
# TYPE openstack_loadbalancer_total_loadbalancers gauge
openstack_loadbalancer_total_loadbalancers 3
# HELP openstack_loadbalancer_total_pools total_pools
# TYPE openstack_loadbalancer_total_pools gauge
openstack_loadbalancer_total_pools 2
# HELP openstack_loadbalancer_up up
# TYPE openstack_loadbalancer_up gauge
openstack_loadbalancer_up 1
`

func (suite *LoadbalancerTestSuite) TestLoadbalancerExporter() {
	err := testutil.CollectAndCompare(*suite.Exporter, strings.NewReader(loadbalancerExpectedUp))
	assert.NoError(suite.T(), err)
}

// The quota metrics are slow and must be dropped by --disable-slow-metrics,
// which also avoids the per-project quota requests.
func (suite *LoadbalancerTestSuite) TestLoadbalancerQuotaSlowMetricsDisabled() {
	exporter := suite.NewTestExporter(nil, true)

	families := suite.GatherFamilies(exporter)
	for name := range families {
		assert.NotContains(suite.T(), name, "openstack_loadbalancer_quota_")
	}
	assert.Contains(suite.T(), families, "openstack_loadbalancer_total_loadbalancers")
}

// A secondary quota metric disabled with --disable-metric must not be emitted
// and must not break the collection of the other quota metrics.
func (suite *LoadbalancerTestSuite) TestLoadbalancerQuotaMetricDisabled() {
	exporter := suite.NewTestExporter([]string{"loadbalancer-quota_listener"}, false)

	families := suite.GatherFamilies(exporter)
	assert.NotContains(suite.T(), families, "openstack_loadbalancer_quota_listener")
	assert.Contains(suite.T(), families, "openstack_loadbalancer_quota_pool")
	// in_use and limit for each of the 8 projects.
	assert.Len(suite.T(), families["openstack_loadbalancer_quota_loadbalancer"].GetMetric(), 8*2)
}

// A failing project must be skipped instead of dropping the metrics of all
// projects.
func (suite *LoadbalancerTestSuite) TestLoadbalancerQuotasSkipFailingProject() {
	suite.FailRequest("/loadbalancer/v2.0/quotas/"+failingProjectID, 404)

	families := suite.GatherFamilies(*suite.Exporter)
	suite.AssertProjectSkipped(families, "openstack_loadbalancer_quota_loadbalancer", "openstack_loadbalancer_quota_l7rule")
}
