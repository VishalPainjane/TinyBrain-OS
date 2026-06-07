package hardware

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// Prober detects host hardware metrics.
type Prober interface {
	Probe() (ProbeResult, error)
}

// OSProber detects hardware using OS facilities and optional CLI tools.
type OSProber struct{}

// Probe collects RAM, VRAM, CPU info, and inference backend on the current machine.
func (OSProber) Probe() (ProbeResult, error) {
	ram, err := detectRAMBytes()
	if err != nil {
		return ProbeResult{}, fmt.Errorf("detect RAM: %w", err)
	}

	vram, backend := detectGPU()
	cpuInfo := fmt.Sprintf("%s/%s cores=%d", runtime.GOOS, runtime.GOARCH, runtime.NumCPU())

	return ProbeResult{
		TotalRAMBytes: ram,
		VRAMBytes:     vram,
		CPUInfo:       cpuInfo,
		Backend:       backend,
	}, nil
}

// ProbeAndClassify runs OSProber and returns the classified hardware profile.
func ProbeAndClassify() (HardwareProfile, error) {
	var p Prober = OSProber{}
	r, err := p.Probe()
	if err != nil {
		return HardwareProfile{}, err
	}
	return ClassifyProfile(r), nil
}

func detectGPU() (vram uint64, backend Backend) {
	backend = BackendCPU

	if path, err := exec.LookPath("nvidia-smi"); err == nil {
		out, err := exec.Command(path, "--query-gpu=memory.total", "--format=csv,noheader,nounits").Output()
		if err == nil {
			vram = parseNvidiaSMIVRAM(out)
			if vram > 0 {
				backend = BackendCUDA
			}
		}
	}

	if backend == BackendCPU && runtime.GOOS == "darwin" {
		// Metal presence is inferred on Darwin without CGO bindings in v0.3.
		backend = BackendMetal
	}

	return vram, backend
}

func parseNvidiaSMIVRAM(out []byte) uint64 {
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimSuffix(line, " MiB")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		mib, err := strconv.ParseUint(line, 10, 64)
		if err != nil {
			continue
		}
		return mib * 1024 * 1024
	}
	return 0
}
