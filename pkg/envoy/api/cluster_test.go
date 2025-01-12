/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package envoy

import (
	"testing"
	"time"

	v3Cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

func TestNewCluster(t *testing.T) {
	name := "myTestCluster_12345"
	connectTimeout := 5 * time.Second

	endpoint1 := NewLBEndpoint("127.0.0.1", 1234)
	endpoint2 := NewLBEndpoint("127.0.0.2", 1234)
	endpoints := []*endpoint.LbEndpoint{endpoint1, endpoint2}

	// With HTTP2
	c := NewCluster(name, connectTimeout, endpoints, true, v3Cluster.Cluster_STATIC)
	assert.Equal(t, c.GetConnectTimeout().Seconds, int64(connectTimeout.Seconds()))
	//nolint: staticcheck // TODO: GetHttp2ProtocolOptions() is deprecated.
	assert.Assert(t, c.GetHttp2ProtocolOptions() != nil)
	assert.Equal(t, c.GetName(), name)
	assert.DeepEqual(t, c.LoadAssignment.Endpoints[0].LbEndpoints, endpoints, protocmp.Transform())

	// Without HTTP2
	c = NewCluster(name, connectTimeout, endpoints, false, v3Cluster.Cluster_STATIC)
	//nolint: staticcheck // TODO: GetHttp2ProtocolOptions() is deprecated.
	assert.Assert(t, c.GetHttp2ProtocolOptions() == nil)
}
