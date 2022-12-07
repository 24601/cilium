# Copyright Authors of Cilium
# SPDX-License-Identifier: Apache-2.0
include ../../Makefile.defs

# Update this via images/scripts/update-cilium-builder-image.sh
CILIUM_BUILDER_IMAGE=quay.io/cilium/cilium-builder:7b3868a9157f7b2a1da391793b72b268151d140a@sha256:17edab226533ba26622978ba192f9bffd01538a0b66c31a2dcb8bc82f523f62a

.PHONY: proto
proto:
	$(QUIET)$(CONTAINER_ENGINE) container run --rm \
		--volume $(CURDIR):/src \
		--user "$(shell id -u):$(shell id -g)" \
		$(CILIUM_BUILDER_IMAGE) \
		make -C /src -f Makefile.protoc
