package hardware_test

import (
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/hardware"
)

func TestClassifyProfile_CPUOnlyTiny(t *testing.T) {
	r := hardware.ProbeResult{
		TotalRAMBytes: 8 * 1024 * 1024 * 1024,
		VRAMBytes:     0,
		Backend:       hardware.BackendCPU,
	}

	profile := hardware.ClassifyProfile(r)
	if profile.Name != hardware.ProfileTiny {
		t.Errorf("Name = %q, want Tiny", profile.Name)
	}
}

func TestClassifyProfile_Standard(t *testing.T) {
	r := hardware.ProbeResult{
		TotalRAMBytes: 16 * 1024 * 1024 * 1024,
		VRAMBytes:     4 * 1024 * 1024 * 1024,
		Backend:       hardware.BackendCUDA,
	}

	profile := hardware.ClassifyProfile(r)
	if profile.Name != hardware.ProfileStandard {
		t.Errorf("Name = %q, want Standard", profile.Name)
	}
}

func TestClassifyProfile_Workstation(t *testing.T) {
	r := hardware.ProbeResult{
		TotalRAMBytes: 128 * 1024 * 1024 * 1024,
		VRAMBytes:     24 * 1024 * 1024 * 1024,
		Backend:       hardware.BackendCUDA,
	}

	profile := hardware.ClassifyProfile(r)
	if profile.Name != hardware.ProfileWorkstation {
		t.Errorf("Name = %q, want Workstation", profile.Name)
	}
}

func TestClassifyProfile_HighRAMCPUOnlyStaysTiny(t *testing.T) {
	r := hardware.ProbeResult{
		TotalRAMBytes: 32 * 1024 * 1024 * 1024,
		VRAMBytes:     0,
		Backend:       hardware.BackendCPU,
	}

	profile := hardware.ClassifyProfile(r)
	if profile.Name != hardware.ProfileTiny {
		t.Errorf("Name = %q, want Tiny for CPU-only", profile.Name)
	}
}

func TestOSProber_Probe_ReturnsRAM(t *testing.T) {
	var p hardware.Prober = hardware.OSProber{}
	r, err := p.Probe()
	if err != nil {
		t.Fatalf("Probe() error = %v", err)
	}
	if r.TotalRAMBytes == 0 {
		t.Error("TotalRAMBytes should be non-zero on a real machine")
	}
	if r.CPUInfo == "" {
		t.Error("CPUInfo should be populated")
	}
}
