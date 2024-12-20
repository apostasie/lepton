/*
   Copyright The containerd Authors.

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

/*
   Portions from https://github.com/moby/moby/blob/cff4f20c44a3a7c882ed73934dec6a77246c6323/pkg/sysinfo/sysinfo.go
   Copyright (C) Docker/Moby authors.
   Licensed under the Apache License, Version 2.0
   NOTICE: https://github.com/moby/moby/blob/cff4f20c44a3a7c882ed73934dec6a77246c6323/NOTICE
*/

// Package sysinfo stores information about which features a kernel supports.
package sysinfo // import "github.com/docker/docker/pkg/sysinfo"

import "github.com/docker/docker/pkg/parsers"

// Opt for New().
type Opt func(info *SysInfo)

// SysInfo stores information about which features a kernel supports.
// TODO Windows: Factor out platform specific capabilities.
type SysInfo struct {
	// Whether the kernel supports AppArmor or not
	AppArmor bool
	// Whether the kernel supports Seccomp or not
	Seccomp bool

	cgroupMemInfo
	cgroupCPUInfo
	cgroupBlkioInfo
	cgroupCpusetInfo
	cgroupPids

	// Whether the kernel supports cgroup namespaces or not
	CgroupNamespaces bool

	// Whether IPv4 forwarding is supported or not, if this was disabled, networking will not work
	IPv4ForwardingDisabled bool

	// Whether the cgroup has the mountpoint of "devices" or not
	CgroupDevicesEnabled bool

	// Whether the cgroup is in unified mode (v2).
	CgroupUnified bool

	// Warnings contains a slice of warnings that occurred  while collecting
	// system information. These warnings are intended to be informational
	// messages for the user, and can either be logged or returned to the
	// client; they are not intended to be parsed / used for other purposes,
	// and do not have a fixed format.
	Warnings []string

	// cg2GroupPath is the cgroup v2 group path to inspect availability of the controllers.
	cg2GroupPath string //nolint:unused,nolintlint

	// cg2Controllers is an index of available cgroup v2 controllers.
	cg2Controllers map[string]struct{} //nolint:unused,nolintlint
}

type cgroupMemInfo struct {
	// Whether memory limit is supported or not
	MemoryLimit bool

	// Whether swap limit is supported or not
	SwapLimit bool

	// Whether soft limit is supported or not
	MemoryReservation bool

	// Whether OOM killer disable is supported or not
	OomKillDisable bool

	// Whether memory swappiness is supported or not
	MemorySwappiness bool

	// Whether kernel memory limit is supported or not. This option is used to
	// detect support for kernel-memory limits on API < v1.42. Kernel memory
	// limit (`kmem.limit_in_bytes`) is not supported on cgroups v2, and has been
	// removed in kernel 5.4.
	KernelMemory bool

	// Whether kernel memory TCP limit is supported or not. Kernel memory TCP
	// limit (`memory.kmem.tcp.limit_in_bytes`) is not supported on cgroups v2.
	KernelMemoryTCP bool
}

type cgroupCPUInfo struct {
	// Whether CPU shares is supported or not
	CPUShares bool

	// Whether CPU CFS (Completely Fair Scheduler) is supported
	CPUCfs bool

	// Whether CPU real-time scheduler is supported
	CPURealtime bool
}

type cgroupBlkioInfo struct {
	// Whether Block IO weight is supported or not
	BlkioWeight bool

	// Whether Block IO weight_device is supported or not
	BlkioWeightDevice bool

	// Whether Block IO read limit in bytes per second is supported or not
	BlkioReadBpsDevice bool

	// Whether Block IO write limit in bytes per second is supported or not
	BlkioWriteBpsDevice bool

	// Whether Block IO read limit in IO per second is supported or not
	BlkioReadIOpsDevice bool

	// Whether Block IO write limit in IO per second is supported or not
	BlkioWriteIOpsDevice bool
}

type cgroupCpusetInfo struct {
	// Whether Cpuset is supported or not
	Cpuset bool

	// Available Cpuset's cpus
	Cpus string

	// Available Cpuset's memory nodes
	Mems string
}

type cgroupPids struct {
	// Whether Pids Limit is supported or not
	PidsLimit bool
}

// IsCpusetCpusAvailable returns `true` if the provided string set is contained
// in cgroup's cpuset.cpus set, `false` otherwise.
// If error is not nil a parsing error occurred.
func (c cgroupCpusetInfo) IsCpusetCpusAvailable(provided string) (bool, error) {
	return isCpusetListAvailable(provided, c.Cpus)
}

// IsCpusetMemsAvailable returns `true` if the provided string set is contained
// in cgroup's cpuset.mems set, `false` otherwise.
// If error is not nil a parsing error occurred.
func (c cgroupCpusetInfo) IsCpusetMemsAvailable(provided string) (bool, error) {
	return isCpusetListAvailable(provided, c.Mems)
}

func isCpusetListAvailable(provided, available string) (bool, error) {
	parsedAvailable, err := parsers.ParseUintList(available)
	if err != nil {
		return false, err
	}
	// 8192 is the normal maximum number of CPUs in Linux, so accept numbers up to this
	// or more if we actually have more CPUs.
	maxCPUs := 8192
	for m := range parsedAvailable {
		if m > maxCPUs {
			maxCPUs = m
		}
	}
	parsedProvided, err := parsers.ParseUintListMaximum(provided, maxCPUs)
	if err != nil {
		return false, err
	}
	for k := range parsedProvided {
		if !parsedAvailable[k] {
			return false, nil
		}
	}
	return true, nil
}
