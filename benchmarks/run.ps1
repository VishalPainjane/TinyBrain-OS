param(
    [string]$Prompt = "Explain the theory of relativity in 5 paragraphs."
)

Write-Host "Starting TinyBrain v1.0 Benchmark Suite" -ForegroundColor Cyan
Write-Host "Hardware Profile: Standard"
Write-Host ""

Write-Host "Running Monolith Benchmark..."
$monolithOut = go run ./cmd/benchmark --mode=monolith --prompt="$Prompt" 2>&1
$monolithOut | Out-String | Write-Host

Write-Host "Running Swarm Benchmark..."
$swarmOut = go run ./cmd/benchmark --mode=swarm --prompt="$Prompt" 2>&1
$swarmOut | Out-String | Write-Host

Write-Host "Benchmarks complete. Please update report_template.md with the results." -ForegroundColor Green
