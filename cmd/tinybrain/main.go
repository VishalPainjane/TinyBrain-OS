// Command tinybrain is the composition root for TinyBrain OS.
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "doctor":
		return runDoctor(stdout)
	case "probe":
		jsonOut := false
		if len(args) > 1 && args[1] == "--json" {
			jsonOut = true
		}
		return runProbe(stdout, jsonOut)
	case "models":
		if len(args) < 2 {
			fmt.Fprintln(stderr, "usage: tinybrain models <list|pull>")
			return 2
		}
		switch args[1] {
		case "list":
			return runModelsList(stdout)
		case "pull":
			return runModelsPull(args[2:], stdout, stderr)
		default:
			fmt.Fprintf(stderr, "unknown models subcommand: %s\n", args[1])
			fmt.Fprintln(stderr, "usage: tinybrain models <list|pull>")
			return 2
		}
	case "agents":
		if len(args) < 2 || args[1] != "list" {
			fmt.Fprintln(stderr, "usage: tinybrain agents list")
			return 2
		}
		return runAgentsList(stdout)
	case "run":
		return runRun(args[1:], stdout, stderr)
	case "daemon":
		return runDaemon(args[1:], stdout, stderr)
	case "submit":
		return runSubmit(args[1:], stdout, stderr)
	case "workflow":
		return runWorkflow(args[1:], stdout, stderr)
	case "status":
		return runStatus(stdout)
	case "version", "--version", "-v":
		fmt.Fprintf(stdout, "tinybrain %s\n", Version)
		return 0
	case "help", "--help", "-h":
		printUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		printUsage(stderr)
		return 2
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, `TinyBrain OS CLI %s

Usage:
  tinybrain doctor
  tinybrain probe [--json]
  tinybrain models list
  tinybrain models pull --id ID --url URL
  tinybrain agents list
  tinybrain run --agent ID --prompt TEXT
  tinybrain daemon [--port 8080] [--agent ID]
  tinybrain submit [--port 8080] --agent ID --prompt TEXT
  tinybrain workflow --prompt TEXT
  tinybrain status
  tinybrain version

Environment:
  TB_MODELS_DB       Registry database (default ~/.tinybrain/models.db)
  TB_MODELS_SEED     YAML seed when DB is empty
  TB_NGLAYERS        GPU layer offload (CUDA builds)
  TB_LLAMA_LIB_DIR   llama.cpp library directory
`, Version)
}
