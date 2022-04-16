package mem

import (
	"log"
	"runtime"

	"flashcat.cloud/categraf/inputs"
	"flashcat.cloud/categraf/inputs/system"
	"flashcat.cloud/categraf/types"
)

const InputName = "mem"

type MemStats struct {
	ps       system.PS
	platform string

	PrintConfigs          bool  `toml:"print_configs"`
	IntervalSeconds       int64 `toml:"interval_seconds"`
	CollectPlatformFields bool  `toml:"collect_platform_fields"`
}

func init() {
	ps := system.NewSystemPS()
	inputs.Add(InputName, func() inputs.Input {
		return &MemStats{
			ps: ps,
		}
	})
}

func (s *MemStats) GetInputName() string {
	return InputName
}

func (s *MemStats) GetIntervalSeconds() int64 {
	return s.IntervalSeconds
}

// overwrite func
func (s *MemStats) Init() error {
	s.platform = runtime.GOOS
	return nil
}

func (s *MemStats) Gather() []*types.Sample {
	vm, err := s.ps.VMStat()
	if err != nil {
		log.Println("E! failed to get vmstat:", err)
		return nil
	}

	fields := map[string]interface{}{
		"total":             vm.Total,     // bytes
		"available":         vm.Available, // bytes
		"used":              vm.Used,      // bytes
		"used_percent":      100 * float64(vm.Used) / float64(vm.Total),
		"available_percent": 100 * float64(vm.Available) / float64(vm.Total),
	}

	if s.CollectPlatformFields {
		switch s.platform {
		case "darwin":
			fields["active"] = vm.Active
			fields["free"] = vm.Free
			fields["inactive"] = vm.Inactive
			fields["wired"] = vm.Wired
		case "openbsd":
			fields["active"] = vm.Active
			fields["cached"] = vm.Cached
			fields["free"] = vm.Free
			fields["inactive"] = vm.Inactive
			fields["wired"] = vm.Wired
		case "freebsd":
			fields["active"] = vm.Active
			fields["buffered"] = vm.Buffers
			fields["cached"] = vm.Cached
			fields["free"] = vm.Free
			fields["inactive"] = vm.Inactive
			fields["laundry"] = vm.Laundry
			fields["wired"] = vm.Wired
		case "linux":
			fields["active"] = vm.Active
			fields["buffered"] = vm.Buffers
			fields["cached"] = vm.Cached
			fields["commit_limit"] = vm.CommitLimit
			fields["committed_as"] = vm.CommittedAS
			fields["dirty"] = vm.Dirty
			fields["free"] = vm.Free
			fields["high_free"] = vm.HighFree
			fields["high_total"] = vm.HighTotal
			fields["huge_pages_free"] = vm.HugePagesFree
			fields["huge_page_size"] = vm.HugePageSize
			fields["huge_pages_total"] = vm.HugePagesTotal
			fields["inactive"] = vm.Inactive
			fields["low_free"] = vm.LowFree
			fields["low_total"] = vm.LowTotal
			fields["mapped"] = vm.Mapped
			fields["page_tables"] = vm.PageTables
			fields["shared"] = vm.Shared
			fields["slab"] = vm.Slab
			fields["sreclaimable"] = vm.Sreclaimable
			fields["sunreclaim"] = vm.Sunreclaim
			fields["swap_cached"] = vm.SwapCached
			fields["swap_free"] = vm.SwapFree
			fields["swap_total"] = vm.SwapTotal
			fields["vmalloc_chunk"] = vm.VmallocChunk
			fields["vmalloc_total"] = vm.VmallocTotal
			fields["vmalloc_used"] = vm.VmallocUsed
			fields["write_back_tmp"] = vm.WriteBackTmp
			fields["write_back"] = vm.WriteBack
		}
	}

	return inputs.NewSamples(fields)
}