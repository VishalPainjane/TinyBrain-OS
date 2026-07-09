package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
)

func runModelsPull(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("pull", flag.ContinueOnError)
	fs.SetOutput(stderr)

	id := fs.String("id", "", "Model ID (required)")
	urlStr := fs.String("url", "", "URL to download GGUF from (required)")
	quant := fs.String("quant", "Q4_K_M", "Quantization level")
	caps := fs.String("caps", "chat,code", "Comma-separated capabilities")
	budget := fs.Int("budget", 1000000000, "Memory budget in bytes")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *id == "" || *urlStr == "" {
		fmt.Fprintln(stderr, "error: --id and --url are required")
		fs.Usage()
		return 2
	}

	filename := filepath.Base(*urlStr)
	// If the url doesn't have a clean filename, fallback to ID + .gguf
	if !strings.HasSuffix(filename, ".gguf") {
		filename = *id + ".gguf"
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(stderr, "error getting home dir: %v\n", err)
		return 1
	}

	destDir := filepath.Join(home, ".tinybrain")
	destPath := filepath.Join(destDir, filename)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		fmt.Fprintf(stderr, "error creating directory: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Starting download for model %s from:\n%s\n", *id, *urlStr)
	fmt.Fprintf(stdout, "Saving to: %s\n\n", destPath)

	resp, err := http.Get(*urlStr)
	if err != nil {
		fmt.Fprintf(stderr, "HTTP get error: %v\n", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(stderr, "HTTP status error: %s\n", resp.Status)
		return 1
	}

	out, err := os.Create(destPath)
	if err != nil {
		fmt.Fprintf(stderr, "File creation error: %v\n", err)
		return 1
	}
	defer out.Close()

	size := resp.ContentLength
	var downloaded int64
	buffer := make([]byte, 32*1024)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if size > 0 {
					pct := float64(downloaded) / float64(size) * 100
					fmt.Fprintf(stdout, "\rProgress: %.2f MB / %.2f MB (%.1f%%)", float64(downloaded)/(1024*1024), float64(size)/(1024*1024), pct)
				} else {
					fmt.Fprintf(stdout, "\rProgress: %.2f MB", float64(downloaded)/(1024*1024))
				}
			}
		}
	}()

	for {
		n, readErr := resp.Body.Read(buffer)
		if n > 0 {
			_, writeErr := out.Write(buffer[:n])
			if writeErr != nil {
				done <- true
				fmt.Fprintf(stderr, "\nWrite error: %v\n", writeErr)
				return 1
			}
			downloaded += int64(n)
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			done <- true
			fmt.Fprintf(stderr, "\nRead error: %v\n", readErr)
			return 1
		}
	}
	done <- true

	fmt.Fprintf(stdout, "\n\nSuccessfully downloaded file to: %s (Total size: %.2f MB)\n\n", destPath, float64(downloaded)/(1024*1024))

	// Register in DB
	reg, err := openRegistry()
	if err != nil {
		fmt.Fprintf(stderr, "Warning: downloaded successfully but failed to open registry: %v\n", err)
		fmt.Fprintln(stderr, "Please stop your running tinybrain daemon (Ctrl+C) and run the registration manually or pull again.")
		return 1
	}
	defer reg.Close()

	capsList := strings.Split(*caps, ",")
	for i, c := range capsList {
		capsList[i] = strings.TrimSpace(c)
	}

	m := registry.ModelDefinition{
		ID:           *id,
		Path:         destPath,
		SizeBytes:    uint64(downloaded),
		MemoryBudget: uint64(*budget),
		Quantization: *quant,
		Capabilities: capsList,
	}

	if err := reg.RegisterModel(m); err == nil {
		fmt.Fprintf(stdout, "Successfully registered %s in registry!\n", *id)
	} else if err.Error() == "duplicate ID" || err.Error() == "ErrDuplicateID" {
		fmt.Fprintf(stdout, "%s is already registered in registry.\n", *id)
	} else {
		fmt.Fprintf(stderr, "Error registering %s: %v\n", *id, err)
		if strings.Contains(err.Error(), "read-only") {
			fmt.Fprintln(stderr, "This is likely because the daemon is running and holds the lock.")
			fmt.Fprintln(stderr, "Please stop the daemon and try registering again.")
		}
		return 1
	}

	return 0
}
