// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package drop

import (
	"context"
	"strings"
	"testing"

	flowpb "github.com/cilium/cilium/api/v1/flow"
	"github.com/cilium/cilium/pkg/hubble/metrics/api"
	monitorAPI "github.com/cilium/cilium/pkg/monitor/api"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDropHandler(t *testing.T) {
	registry := prometheus.NewRegistry()
	opts := api.Options{}
	h := &dropHandler{}
	assert.NoError(t, h.Init(registry, opts))
	flow := flowpb.Flow{
		Verdict:        flowpb.Verdict_DROPPED,
		DropReasonDesc: flowpb.DropReason_POLICY_DENIED,
		EventType:      &flowpb.CiliumEventType{Type: monitorAPI.MessageTypeDrop},
		L4: &flowpb.Layer4{
			Protocol: &flowpb.Layer4_TCP{
				TCP: &flowpb.TCP{DestinationPort: 80},
			},
		},
		Source:      &flowpb.Endpoint{Namespace: "src-a"},
		Destination: &flowpb.Endpoint{Namespace: "src-b"},
	}
	assert.NoError(t, h.ProcessFlow(context.Background(), &flow))
	flow.DropReasonDesc = flowpb.DropReason_DROP_REASON_UNKNOWN
	assert.NoError(t, h.ProcessFlow(context.Background(), &flow))
	expected := strings.NewReader(`# HELP hubble_drop_total Number of drops
# TYPE hubble_drop_total counter
hubble_drop_total{protocol="TCP",reason="DROP_REASON_UNKNOWN"} 1
hubble_drop_total{protocol="TCP",reason="POLICY_DENIED"} 1
`)
	assert.NoError(t, testutil.CollectAndCompare(h.drops, expected))
}
