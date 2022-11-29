// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package cmd

import (
	"github.com/cilium/cilium/pkg/datapath"
	"github.com/cilium/cilium/pkg/datapath/types"
	"github.com/cilium/cilium/pkg/defaults"
	"github.com/cilium/cilium/pkg/gops"
	"github.com/cilium/cilium/pkg/hive/cell"
	"github.com/cilium/cilium/pkg/k8s"
	k8sClient "github.com/cilium/cilium/pkg/k8s/client"
	"github.com/cilium/cilium/pkg/node"
	"github.com/cilium/cilium/pkg/readiness"
	serviceManager "github.com/cilium/cilium/pkg/service"
	serviceCache "github.com/cilium/cilium/pkg/service/cache"
)

var (
	Agent = cell.Module(
		"agent",
		"Cilium Agent",

		Infrastructure,
		ControlPlane,
		Datapath,
	)

	// Infrastructure provides access and services to the outside.
	// A cell should live here instead of ControlPlane if it is not needed by
	// integrations tests, or needs to be mocked.
	Infrastructure = cell.Module(
		"infra",
		"Infrastructure",

		// Runs the gops agent, a tool to diagnose Go processes.
		gops.Cell(defaults.GopsPortAgent),

		// Provides Clientset, API for accessing Kubernetes objects.
		k8sClient.Cell,
	)

	// ControlPlane implement the per-node control functions. These are pure
	// business logic and depend on datapath or infrastructure to perform
	// actions. This separation enables non-privileged integration testing of
	// the control-plane.
	ControlPlane = cell.Module(
		"controlplane",
		"Control Plane",

		// Readiness allows modules to register as readiness signal providers.
		// Daemon waits for the signal before finishing initialization and telling
		// Kubernetes that the agent is ready for CNI requests.
		readiness.Cell,

		// LocalNodeStore holds onto the information about the local node and allows
		// observing changes to it.
		node.LocalNodeStoreCell,

		// Shared resources provide access to k8s resources as event streams or as
		// read-only stores.
		k8s.SharedResourcesCell,

		// ServiceCache provides an API for accessing services and their associated
		// endpoints.
		serviceCache.Cell,

		// ServiceManager manages the datapath resources for services and backends.
		serviceManager.Cell,

		// Provide NodeAddressing for ServiceCache.
		cell.Provide(
			func(dp datapath.Datapath) types.NodeAddressing {
				return dp.LocalNodeAddressing()
			},
		),

		// daemonCell wraps the legacy daemon initialization and provides Promise[*Daemon].
		daemonCell,
	)

	// Datapath provides the privileged operations to apply control-plane
	// decision to the kernel.
	Datapath = cell.Module(
		"datapath",
		"Datapath",

		cell.Provide(
			newWireguardAgent,
			newDatapath,
		),
	)
)