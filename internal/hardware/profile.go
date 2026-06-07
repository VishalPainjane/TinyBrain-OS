package hardware

// ProfileName is a hardware capability tier assigned at boot.
// See docs/architecture/hardware.md.
type ProfileName string

const (
	ProfileTiny        ProfileName = "Tiny"
	ProfileStandard    ProfileName = "Standard"
	ProfileWorkstation ProfileName = "Workstation"
)

// Backend identifies the inference backend available on the host.
type Backend string

const (
	BackendCPU   Backend = "CPU"
	BackendCUDA  Backend = "CUDA"
	BackendMetal Backend = "Metal"
	BackendROCm  Backend = "ROCm"
)

// ProbeResult holds detected hardware metrics used for profile classification.
type ProbeResult struct {
	TotalRAMBytes uint64
	VRAMBytes     uint64
	CPUInfo       string
	Backend       Backend
}

// HardwareProfile is the classified hardware profile exposed to registry and runtime.
type HardwareProfile struct {
	Name    ProfileName
	Probe   ProbeResult
}

const (
	gib = 1024 * 1024 * 1024

	ramTinyMax        = 8 * gib
	ramStandardMin    = 16 * gib
	ramWorkstationMin = 64 * gib

	vramStandardMin = 4 * gib
)

// ClassifyProfile assigns Tiny, Standard, or Workstation from probe metrics.
// CPU-only machines with constrained RAM classify as Tiny.
func ClassifyProfile(r ProbeResult) HardwareProfile {
	name := ProfileTiny

	hasDiscreteGPU := r.VRAMBytes >= vramStandardMin && r.Backend != BackendCPU

	switch {
	case r.TotalRAMBytes >= ramWorkstationMin && hasDiscreteGPU:
		name = ProfileWorkstation
	case r.TotalRAMBytes >= ramStandardMin && hasDiscreteGPU:
		name = ProfileStandard
	case r.TotalRAMBytes <= ramTinyMax && r.Backend == BackendCPU:
		name = ProfileTiny
	case r.TotalRAMBytes >= ramStandardMin:
		// RAM meets Standard tier but GPU missing or below threshold — stay Tiny per CPU-only rule.
		name = ProfileTiny
	default:
		name = ProfileTiny
	}

	return HardwareProfile{
		Name:  name,
		Probe: r,
	}
}
