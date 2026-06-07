Day 1 Deliverable
Only:
README.md
Structure:
What is TinyBrain?
Why does it exist?
High level architecture.
Current roadmap.
Current status.
One page.
Not 20 pages.
I’m turning your spec into a tight one-page README and keeping it strictly to the six sections you listed.
Thought for 10s
Done: README.md
Built as a one-page Day 1 README from your TinyBrain OS spec.
Create:
vision.md
This answers:
Why does TinyBrain exist?
Example topics:
Problem
Local AI is constrained
Monolithic models waste resources
Need dynamic runtime scheduling
I’m drafting it as a short vision statement, centered on the problem, the constraint, and the runtime idea.
Thought for a couple of seconds
Analyzed
from pathlib import Path
content = """# TinyBrain Vision
TinyBrain exists because local AI is still treated like a single-model problem.
Today’s common approach is simple: load one large model, send every task to it, and hope it can do everything well enough. That works for demos, but it is inefficient on real hardware. A monolithic model consumes too much VRAM and RAM, wastes compute on tasks that do not need deep reasoning, and becomes slow or impossible to run on consumer machines.
TinyBrain is built around a different idea: AI should behave like a runtime, not just a chat interface.
## The problem
Local AI has three major constraints:
- hardware is limited
- memory is expensive
- not every task deserves the same model
A single model is forced to handle planning, browsing, coding, reasoning, and retrieval at once. That leads to poor specialization and poor resource usage. On constrained devices, this becomes even worse: large models cannot stay loaded comfortably, context becomes expensive, and switching between tasks is inefficient.
## Why TinyBrain exists
TinyBrain exists to make local AI practical on everyday hardware.
Instead of one model trying to do everything, TinyBrain treats agents like scheduled processes. Different tasks can be routed to different models, loaded only when needed, and swapped out when they are no longer active. That means the system can stay lean, responsive, and hardware-aware.
The goal is not just to run AI locally. The goal is to run it intelligently.
## The vision
TinyBrain is a local AI runtime that:
- schedules work dynamically
- allocates resources based on task needs
- loads and unloads models as needed
- keeps memory usage under control
- makes specialized agents work together
This is the core vision: move from a static model-centric workflow to a dynamic runtime-centric system.
## The outcome
If TinyBrain works, then local AI becomes:
- more efficient
- more modular
- more scalable on weak hardware
- more transparent to inspect and debug
- more usable as a real system, not just a prompt wrapper
TinyBrain exists to prove that a smart runtime can make small and mid-sized models feel far more capable than a single monolithic model.
That is the direction.
"""
path = Path("/mnt/data/vision.md")
path.write_text(content, encoding="utf-8")
print(f"Wrote {path} ({len(content.splitlines())} lines)")
STDOUT/STDERR
Wrote /mnt/data/vision.md (51 lines)
Done: vision.md
I’m writing the event catalog as a tight source-of-truth file with just the event names you listed.
Thought for a second
Analyzed
from pathlib import Path
base = Path("/mnt/data/docs")
base.mkdir(parents=True, exist_ok=True)
content = """# TinyBrain Events
This document defines the core events used by TinyBrain.
## Events
- TaskCreated
- TaskAssigned
- AgentStarted
- AgentStopped
- SwapStarted
- SwapCompleted
- KVStored
- KVLoaded
"""
path = base / "events.md"
path.write_text(content, encoding="utf-8")
print(f"Wrote {path}")
STDOUT/STDERR
Wrote /mnt/data/docs/events.md
Done: events.md
Day 5
Registry Design
Create:
docs/registry.md
This is important.
Remember our newer architecture.
We removed:
Planner
Browser
Coder
Reasoner
from core.
Instead define:
Agent
conceptually.
Registry owns:
Agent Definitions
Model Definitions
Tool Definitions
I’m defining the registry around generic agents, models, and tools so it matches the newer architecture cleanly.
Thought for a second
Analyzed
from pathlib import Path
base = Path("/mnt/data/docs")
base.mkdir(parents=True, exist_ok=True)
content = """# TinyBrain Registry
The registry is the system of record for all configurable capabilities in TinyBrain.
It does not define fixed roles like Planner, Browser, Coder, or Reasoner at the core level.
Instead, TinyBrain treats everything as an Agent with a configuration.
## What the Registry Owns
### Agent Definitions
Agent definitions describe what an agent is allowed to do, what it is called, and how it should be treated by the runtime.
### Model Definitions
Model definitions describe available models, their capabilities, and their resource requirements.
### Tool Definitions
Tool definitions describe external actions or abilities that agents may request or use through the runtime.
## Purpose
The registry keeps TinyBrain modular.
It allows the runtime to:
- discover available agents
- select compatible models
- attach the right tools
- adapt to different hardware profiles
## Core Rule
The registry describes capabilities.
The scheduler and runtime consume those capabilities.
The core system should not hardcode specialized agent roles.
"""
path = base / "registry.md"
path.write_text(content, encoding="utf-8")
print(f"Wrote {path}")
STDOUT/STDERR
Wrote /mnt/data/docs/registry.md
Done: registry.md
Day 6
Hardware Profiles
Create:
docs/hardware-profiles.md
Define:
Tiny
8GB RAM
CPU only
Standard
16GB RAM
RTX 3050
Workstation
64GB+
This document becomes extremely important later.
I’m adding the hardware profile tiers as a compact reference doc for later model selection and scheduling decisions.
Thought for a second
Analyzed
from pathlib import Path
base = Path("/mnt/data/docs")
base.mkdir(parents=True, exist_ok=True)
content = """# TinyBrain Hardware Profiles
TinyBrain uses hardware profiles to understand what kind of local environment it is running on.
These profiles guide model selection, scheduling, and resource planning.
## Tiny
- 8GB RAM
- CPU only
This profile represents the most constrained target environment.
## Standard
- 16GB RAM
- RTX 3050
This is the baseline profile for the main TinyBrain design.
## Workstation
- 64GB+ RAM
This profile represents a high-capability system with room for larger workloads and more aggressive orchestration.
## Core Rule
Hardware profiles are not just labels.
They are inputs to the runtime, scheduler, and registry when deciding how TinyBrain should behave.
"""
path = base / "hardware-profiles.md"
path.write_text(content, encoding="utf-8")
print(f"Wrote {path}")
STDOUT/STDERR
Wrote /mnt/data/docs/hardware-profiles.md
Done: hardware-profiles.md
I'll treat the document you uploaded as the source of truth for the project specification.
Part 1 — What Has Already Been Decided
This is not just an AI project.
The project has already evolved into:
A Local Multi-Agent Operating System for LLMs
A resource-aware runtime that schedules multiple tiny AI models like OS processes on constrained hardware.
Think:
Linux scheduler
Kubernetes control plane
LLM orchestration
Memory manager
Agent runtime
combined into one system.
Project Identity
Final Name
TinyBrain OS
Alternative naming:
TinyBrain Runtime
TinyBrain Core
TinyBrain Agent OS
But the document consistently converges toward:
TinyBrain OS.
Core Thesis
The entire project is built on one hypothesis:
Instead of:
1 × 14B model
use:
1 × 2B Planner
1 × 2B Browser
1 × 3B Coder
1 × 4B Reasoner
and orchestrate them intelligently.
The swarm should:
consume less VRAM
consume less RAM
improve task specialization
improve throughput
improve cost efficiency
compared to a monolithic model. (arXiv)
Target Hardware
The entire architecture is optimized for:
Baseline Device
16 GB RAM
RTX 3050
4 GB VRAM
This is a hard design constraint.
Everything revolves around fitting inside this envelope.
Not 4090.
Not A100.
Not H100.
Consumer laptop.
The System Is NOT Cloud First
A major decision already made:
Local First
Edge First
Offline Capable
No dependency on:
OpenAI
Anthropic
Gemini
for core operation.
Everything should run locally.
Fundamental Architecture
The architecture is:
User
↓
Proxy Router
↓
Orchestrator
↓
Scheduler
↓
Agent Runtime
↓
Models
Where:
Planner
Browser
Coder
Reasoner
are independent processes.
Not prompts.
Processes.
Agent Roles Are Already Defined
Planner
Responsible for:
Task decomposition
Workflow generation
Routing
Planning
Expected size:
2B
Browser
Responsible for:
Search
Documentation lookup
Retrieval
Expected size:
2B
Coder
Responsible for:
Code generation
Code modification
Refactoring
Expected size:
3B
Reasoner
Responsible for:
Deep reasoning
Verification
Synthesis
Expected size:
4B
Quantization Strategy
Already decided.
Not optional.
Primary Format
GGUF
Reason:
mmap support
llama.cpp compatibility
local inference
Quantization
Q4_K_M
Preferred default.
Reasons:
Good quality
Good memory savings
Fast
Optional
AWQ
For future optimization.
Runtime Engine
Already largely selected.
Primary
llama.cpp
Reason:
GGUF
KV handling
Memory mapping
GPU offloading
Possible secondary:
vLLM
Ollama
but llama.cpp is clearly the primary engine.
Most Important Innovation
The project's biggest differentiator:
Dynamic Model Swapping
Instead of:
Load all models
you:
Load Planner
Run
Unload
Load Browser
Run
Unload
Load Coder
Run
Unload
This is the core novelty.
Scheduler Type
Already decided.
MLFQ
Multi-Level Feedback Queue
Same concept used in operating systems.
It handles:
priority
fairness
preemption
responsiveness
Queues:
P0 → highest
P1
P2
P3 → lowest
Interactive tasks jump to top.
Background tasks move downward.
Preemption Exists
Critical decision.
Agents can be interrupted.
Example:
Coder running
User asks something.
Scheduler:
Pause coder
Run planner
Resume coder later
Exactly like OS process scheduling.
KV Cache Is First-Class
This is huge.
Most agent systems ignore it.
TinyBrain does not.
When an agent is paused:
Save KV cache
instead of:
Recompute prompt
This becomes:
Context Hibernation
KV Compression Strategy
Chosen:
FP16
↓
Q4 KV
Goal:
~72% reduction.
This is one of the biggest performance wins in the whole architecture.
Persistent Block Pool
Already decided.
Storage hierarchy:
VRAM
↓
RAM
↓
NVMe
KV caches live here.
Not inside model memory.
Memory Mapping Strategy
Very important.
The project will use:
mmap()
for model loading.
Reason:
Avoid loading entire model into RAM.
Only load pages when required.
Model Router
Already defined.
API style:
POST /models/load
POST /models/unload
GET /models
POST /models/swap
The scheduler controls models via APIs.
Not shell scripts.
LRU Eviction
Chosen.
When VRAM full:
Least Recently Used model
↓
Evict
This frees VRAM automatically.
Communication Protocol
A massive design decision.
Agents DO NOT communicate using English.
Bad:
Hey coder, write a function...
Good:
{
"task": "write_code",
"language": "python"
}
Everything becomes structured.
IPC Format
Chosen:
JSON
Everywhere.
No exceptions.
Structured Outputs
Chosen technology:
JSON Schema
GBNF Grammar
or
XGrammar
This guarantees:
Valid JSON
100% parseable
between agents.
Kubernetes Is Included
One of the strongest resume points.
Kubernetes is not decorative.
It becomes:
Control Plane
for the whole runtime.
The host machine becomes:
Single Node Cluster
This is a brilliant talking point for recruiters.
This is only about 30-40% of the decisions already made.
The next part would cover:
exact repository structure
service architecture
telemetry system
brain-top TUI
APIs
database decisions
scheduler internals
Kubernetes CRDs
deployment architecture
benchmark suite
observability stack
development phases
portfolio/recruiter positioning
MVP vs V1 vs V2 roadmap
what should actually be built first vs what can be mocked
Reply with:
"continue with part 2"
and I'll go even deeper into the engineering architecture.
Part 2 — Repository Architecture, Services, APIs, Telemetry & What You Actually Build
Now we're moving from ideas into real engineering decisions.
This is where recruiters stop seeing "AI project" and start seeing:
Systems Engineer
Distributed Systems
Infrastructure Engineer
Runtime Engineer
ML Systems Engineer
territory.
The architecture you've described is surprisingly close to a miniature operating system. The concepts of MLFQ scheduling, KV cache offloading, memory mapping, preemption, and multi-tier memory management are all established systems ideas being adapted to LLM workloads. (Wikipedia)
Final Repository Structure
This should NOT become:
backend/
frontend/
models/
That's too ordinary.
Instead:
tinybrain-os/
├── cmd/
│ ├── tinybrain
│ └── brain-top
├── orchestrator/
│ ├── scheduler/
│ ├── router/
│ ├── planner/
│ └── policies/
├── runtime/
│ ├── model-loader/
│ ├── kv-manager/
│ ├── swap-manager/
│ └── process-manager/
├── agents/
│ ├── planner/
│ ├── browser/
│ ├── coder/
│ └── reasoner/
├── api/
│ ├── rest/
│ ├── grpc/
│ └── websocket/
├── telemetry/
│ ├── metrics/
│ ├── tracing/
│ └── logs/
├── dashboard/
│ └── web-ui/
├── tui/
│ └── brain-top/
├── benchmarks/
├── deployments/
│ ├── k8s/
│ ├── docker/
│ └── local/
└── docs/
Reason:
You are building an OS runtime.
Not a chatbot.
Major Services
The project is actually 7 separate systems.
Service 1
API Gateway
Purpose:
Receive requests
Authenticate
Stream output
Endpoints:
POST /v1/chat
POST /v1/task
GET /v1/status
GET /v1/metrics
Think:
TinyBrain API = OpenAI API replacement.
Service 2
Router
Purpose:
Decide:
Planner?
Coder?
Browser?
Reasoner?
before scheduler sees anything.
Example:
User:
Write a Next.js middleware
Router:
{
"agent": "coder",
"priority": 0
}
Scheduler Service
Most important component.
Responsibilities:
Queue management
Priorities
Preemption
Aging
Starvation prevention
Everything flows through scheduler.
Inputs:
{
"agent":"coder",
"priority":1,
"memory":"2.4GB",
"context":"..."
}
Outputs:
{
"state":"RUNNING",
"gpu":"allocated"
}
Scheduler State Machine
Every agent exists in one state.
NEW
READY
RUNNING
WAITING
SUSPENDED
TERMINATED
Exactly like Linux.
Transitions:
READY
↓
RUNNING
↓
WAITING
↓
RUNNING
↓
TERMINATED
or
RUNNING
↓
PREEMPTED
↓
SUSPENDED
Process Manager
Hidden but critical.
Purpose:
Track every model.
Example:
PID 101 Planner
PID 102 Browser
PID 103 Coder
PID 104 Reasoner
Displayed inside brain-top.
KV Manager
This becomes a separate subsystem.
Responsibilities:
Save KV
Load KV
Compress KV
Restore KV
Delete KV
The document repeatedly assumes KV cache persistence/offloading is a first-class primitive, which aligns with LMCache/KVSwap-style architectures. (arXiv)
Model Loader
Purpose:
Load model
Unload model
Warm model
Prefetch model
Uses:
mmap()
so weights load lazily rather than fully copying into RAM. Modern llama.cpp tooling explicitly supports mmap-backed model loading. (manpages.debian.org)
Model Registry
Local database.
Stores:
{
"planner": {
"path":"models/planner.gguf",
"size":"1.3GB"
}
}
Can later evolve into:
HuggingFace integration
Agent Runtime
Each agent gets:
Memory quota
Priority
Context
Tools
Example:
Planner:
memory: 1.5GB
priority: HIGH
Coder:
memory: 2.5GB
priority: MEDIUM
Tool Execution Layer
A separate service.
Not part of agents.
Tools:
Search
Filesystem
Terminal
Python
Git
Agents request tools.
Agents never execute tools directly.
Example:
Planner:
{
"tool":"search",
"query":"Next.js 15 routing"
}
IPC Bus
One of the smartest decisions.
Never send:
human language
between agents.
Send:
{
"action":"write_code",
"language":"typescript"
}
This reduces agent drift dramatically.
Structured Output Engine
Dedicated subsystem.
Responsible for:
JSON schemas
GBNF
Grammar constraints
Validation
The idea is that every agent output is machine-consumable rather than conversational.
Telemetry Stack
This is recruiter gold.
Most student projects skip observability.
Yours shouldn't.
Metrics:
VRAM usage
RAM usage
SSD usage
TPS
TTFT
Queue depth
Swap latency
brain-top (TUI)
This becomes your signature feature.
Think:
htop
btop
nvidia-smi
agent scheduler
Layout:
┌─────────────────────┐
│ AGENTS │
├─────────────────────┤
│ Planner Sleeping │
│ Browser Waiting │
│ Coder Running │
│ Reasoner Suspended │
└─────────────────────┘
Second Panel
GPU
███████░░░ 72%
RAM
█████░░░░░ 55%
SSD
██░░░░░░░░ 21%
Third Panel
Queue monitor.
P0 [2]
P1 [4]
P2 [1]
P3 [0]
Fourth Panel
Swap monitor.
Planner
VRAM → RAM
Coder
RAM → VRAM
Browser
RAM → SSD
Benchmark Suite
Mandatory.
Without benchmarks:
Interesting idea
With benchmarks:
Research-grade system
Benchmark A
Monolith:
Llama 8B
Benchmark B
Swarm:
Planner
Browser
Coder
Reasoner
Measure:
TTFT
TPS
Peak VRAM
Peak RAM
Completion time
The entire project's claim depends on proving that orchestrated small models can compete under tight hardware budgets. KV offloading and hierarchical cache systems are already showing significant throughput gains in current research. (arXiv)
Kubernetes Objects
You'll eventually create CRDs.
Example:
kind: Agent
kind: Model
kind: Task
kind: SwapPolicy
This is where the Kubernetes angle becomes real instead of decorative.
What Actually Gets Built First
This is extremely important.
Do NOT start with:
MLFQ
KV compression
Cross-phase injection
Start with:
MVP
Load model
Unload model
Switch model
Only.
Then:
V0.2
Planner
Coder
Only two agents.
Then:
V0.3
KV save/load
Then:
V0.4
brain-top
Then:
V1
Real scheduler
If you try building everything immediately:
Project dies.
What Makes This Portfolio-Worthy
When someone asks:
"What did you build?"
The answer is NOT:
A multi-agent AI system
The answer becomes:
I built a local operating system for AI agents.
It implements process scheduling,
resource-aware model orchestration,
KV-cache hibernation,
dynamic model paging,
multi-tier memory management,
structured IPC,
and Kubernetes-based control planes
for running specialized LLM swarms on
consumer hardware.
That's a completely different level of project positioning.
Part 3 would go even deeper into:
exact MVP architecture
technologies (Go vs Rust vs TypeScript)
database choices
CRD design
event-driven architecture
scheduler algorithms
swap heuristics
deployment pipeline
GitHub roadmap
how to make it look like a startup-grade project instead of a college project.
Part 3 — Exact MVP, Tech Stack Decisions, CRDs, Event Architecture, and How to Make This Look Like a Real Company Product
This is the part most people get wrong.
The document describes a future vision.
But investors, recruiters, and engineers judge you on:
What actually works today?
The difference between a cool architecture diagram and a real systems project is whether you can define the first working version.
What Is The Actual MVP?
Not:
Multi-agent OS
Not:
AI Kubernetes Runtime
Not:
KV Compression
The MVP is much smaller.
MVP Goal
Prove ONE thing:
Multiple models can be treated like operating system processes and dynamically loaded/unloaded under a resource budget.
That's it.
Nothing more.
MVP Workflow
User enters:
Create a FastAPI CRUD application
System:
Planner
↓
creates task list
Coder
↓
generates code
Return result
Only:
2 agents
No Browser.
No Reasoner.
No parallelism.
Technology Stack Decisions
Many choices exist.
Most are wrong.
Backend
Choose:
Go
Reason:
Kubernetes ecosystem
Controller-runtime
Operators
Concurrency
Performance
Single binary
Kubernetes itself is written in Go. The CRD/controller/operator ecosystem is built around Go-based tooling and controller-runtime patterns. (Kubernetes)
Avoid:
Node.js
for orchestration.
Use Node only for frontend.
Why Not Rust?
Rust is excellent.
But:
Kubernetes Operator SDK
controller-runtime
client-go
are all Go-first.
You would spend months fighting tooling.
Recommendation:
Go = Core
Rust = Optional future runtime optimizations
Frontend
Choose:
Next.js
Reason:
Dashboard
Landing page
Metrics UI
Docs
No reason to overcomplicate.
API Layer
Use:
gRPC internally
Between:
Scheduler
Runtime
KV Manager
Model Loader
Use:
REST externally
For users.
POST /chat
POST /task
GET /metrics
This mirrors how many large distributed systems are structured.
Database Decision
This surprises people.
Initially:
NO DATABASE
Use:
BoltDB
SQLite
only.
Store:
Task metadata
Model metadata
Execution history
Do not deploy:
Postgres
Redis
Kafka
in MVP.
Zero value.
Huge complexity.
Event Driven Architecture
Everything becomes events.
Bad:
Planner directly calls Coder
Good:
PlannerFinished
↓
TaskCreated
↓
SchedulerEvent
↓
CoderStarted
Now every component is decoupled.
Events:
AgentLoaded
AgentUnloaded
TaskCreated
TaskCompleted
SwapStarted
SwapFinished
KVStored
KVRestored
Internal Event Bus
For MVP:
Go Channels
Later:
NATS
Eventually:
Kafka
Do NOT start with Kafka.
The Core CRDs
This is where Kubernetes becomes meaningful.
CRDs + Controllers are the standard way to extend Kubernetes with your own resources and automation. (Kubernetes)
Agent CRD
kind: Agent
Represents:
Planner
Coder
Browser
Reasoner
Example:
apiVersion: tinybrain.io/v1
kind: Agent
metadata:
name: planner
spec:
model: stablelm-2b
memoryLimit: 2Gi
priority: high
Task CRD
kind: Task
Represents:
Work
Example:
spec:
prompt: "Build FastAPI CRUD"
assignedAgent: planner
SwapPolicy CRD
Unique feature.
kind: SwapPolicy
Example:
spec:
vramThreshold: 80
swapTarget: ram
Now scheduler behavior becomes declarative.
Huge Kubernetes flex.
KVCache CRD
This one is powerful.
kind: KVCache
Tracks:
Location
Size
Owner Agent
Compression Type
Example:
spec:
owner: planner
location: ram
size: 700mb
Scheduler Architecture
Real architecture:
API
↓
Router
↓
Scheduler
↓
Runtime
Scheduler never talks to models directly.
Instead:
Scheduler
↓
Runtime Manager
↓
Model Process
This separation prevents future rewrites.
brain-top Architecture
brain-top should become its own product.
Input Sources:
Scheduler
GPU
RAM
Swap Manager
KV Manager
Output:
Terminal Dashboard
Framework:
Bubble Tea (Go)
or
Charm Bracelet ecosystem
Reason:
Beautiful TUIs.
Fast.
Modern.
Metrics Pipeline
Every event emits metrics.
Example:
TaskCreated
increments:
tasks_total
Example:
SwapFinished
records:
swap_duration_ms
Metrics:
agent_load_time
agent_unload_time
swap_time
queue_wait_time
vram_usage
ram_usage
tokens_per_second
Benchmark Architecture
Benchmarks become a first-class module.
Not a script.
Create:
benchmarks/
Benchmarks:
Memory
Measure:
Peak RAM
Peak VRAM
Scheduling
Measure:
Queue latency
Inference
Measure:
TTFT
TPS
Swapping
Measure:
Swap latency
Restore latency
Deployment Strategy
Three modes.
Mode 1
Developer Mode
tinybrain run
Everything local.
Mode 2
Docker
docker compose up
Mode 3
Kubernetes
helm install tinybrain
This makes the project feel enterprise-grade.
GitHub Roadmap
The repository should expose milestones.
Milestone 1
Single model runtime
Milestone 2
Model switching
Milestone 3
Planner + Coder
Milestone 4
KV persistence
Milestone 5
MLFQ Scheduler
Milestone 6
brain-top
Milestone 7
Kubernetes Operator
Milestone 8
Full Swarm
What Recruiters Will Actually Care About
Not:
Phi-3
Gemma
Qwen
They care about:
Scheduling
Memory Management
Resource Allocation
Systems Design
Observability
Kubernetes
Runtime Engineering
The strongest pitch is:
Built a local AI operating system
that dynamically schedules and
swaps specialized LLM agents under
strict GPU memory constraints.
Implemented Kubernetes-native
controllers, custom schedulers,
KV-cache hibernation,
memory-mapped model loading,
and real-time runtime telemetry.
That sounds like a systems engineer built it, not a prompt engineer.
The One Feature That Can Make It Viral
Not the agents.
Not the models.
brain-top
When people run:
brain-top
and see:
Planner RUNNING
Coder SUSPENDED
Browser SWAPPED
Reasoner WAITING
VRAM 82%
RAM 54%
Swap Latency 32ms
TPS 64
they instantly understand the project.
The dashboard becomes the visual proof that TinyBrain OS is behaving like an operating system for AI agents.
Part 4 would go into the deepest layer:
exact scheduler algorithm implementation
MLFQ data structures
swap heuristics
KV block pool architecture
memory allocator design
model lifecycle manager
process table design
how to implement this on an RTX 3050 without the architecture collapsing
what can be faked for demo purposes vs what must be genuinely engineered.
Part 4 — The Real Scheduler, Memory System, Swap Heuristics & What Must Actually Be Engineered
This is where TinyBrain OS stops being a "multi-agent project" and becomes a runtime system.
Most people can build:
Agent A
→ Agent B
→ Agent C
Very few can build:
Scheduler
Memory Manager
Process Table
Context Hibernation
The scheduler is the heart of the entire system.
The Scheduler Is NOT Scheduling Agents
This is the first major realization.
The scheduler is actually scheduling:
VRAM
Not agents.
Not tasks.
Not prompts.
VRAM.
Every design decision revolves around:
4 GB GPU memory
and keeping it alive.
TinyBrain's Process Table
You need a real process table.
Like Linux.
Structure:
type AgentProcess struct {
PID            string
AgentType      string
State          State
Priority       int
MemoryUsage    uint64
VRAMUsage      uint64
KVCacheID      string
LastExecution  time.Time
TokensProduced int
TaskID         string
}
Example:
PID 101 Planner
PID 102 Browser
PID 103 Coder
PID 104 Reasoner
brain-top reads directly from this table.
Agent States
Final state machine:
NEW
READY
RUNNING
WAITING
PREEMPTED
HIBERNATED
TERMINATED
Meaning:
READY
Waiting for GPU.
RUNNING
Currently owns GPU.
WAITING
Waiting on tool.
Example:
Browser waiting on HTTP request
PREEMPTED
Removed from GPU.
KV still exists.
HIBERNATED
KV moved.
Weights unloaded.
TERMINATED
Finished.
The MLFQ Structure
Actual queues:
Q0 → Highest
Q1
Q2
Q3 → Lowest
This follows classic MLFQ behavior where interactive or short jobs stay higher and CPU-heavy workloads gradually move lower. (GeeksforGeeks)
Data structure:
type Queue struct {
Level int
Quantum int
Processes []*AgentProcess
}
Queue Quantums
For TinyBrain:
Q0 = 32 tokens
Q1 = 64 tokens
Q2 = 128 tokens
Q3 = 256 tokens
Notice:
Not milliseconds.
Not CPU cycles.
Tokens.
This is a major innovation.
Traditional OS:
10ms quantum
TinyBrain:
32 token quantum
Scheduling Loop
Every token generated:
for {
process := HighestPriorityProcess()
Run(process)
if QuantumExceeded {
    Demote(process)
}
}
Meaning:
Coder generates:
32 tokens
Scheduler interrupts.
Checks queue.
Potentially switches model.
Priority Boosting
Classic MLFQ problem:
Starvation
Lower queues never run.
Solution:
Every:
30 seconds
or
500 generated tokens
perform:
BoostAll()
Move everyone to:
Q0
Again.
This is straight from operating-system scheduling practice. (GeeksforGeeks)
Context Switch Cost
This is where most designs fail.
Switching models is expensive.
Cost Components:
Unload weights
Save KV
Load weights
Restore KV
If:
Switch Cost
Execution Benefit
scheduler loses.
Therefore:
TinyBrain needs:
Swap heuristic
Swap Heuristic
Never swap immediately.
Bad:
Coder paused
→ unload
Good:
Coder paused
If idle > 10 seconds:
swap
else:
keep warm
This single heuristic massively reduces thrashing.
The Three Memory Layers
TinyBrain memory hierarchy:
VRAM
↓
RAM
↓
NVMe
Exactly like:
L1
L2
L3
cache hierarchy.
VRAM Layer
Stores:
Active weights
Active KV
Only.
Target:
< 90% utilization
Never 100%.
RAM Layer
Stores:
Compressed KV
Warm Models
Target:
< 70%
Always leave headroom.
NVMe Layer
Stores:
Cold KV
Old Sessions
Archived Contexts
Think:
Swap partition
for AI.
KV Block Pool
This is probably the most unique component.
Structure:
type KVBlock struct {
ID string
Agent string
Size uint64
Location string
LastAccess time.Time
}
Locations:
VRAM
RAM
SSD
Every context becomes a block.
KV Block Manager
Responsibilities:
Allocate
Move
Compress
Restore
Delete
Only.
Nothing else.
This should be its own package.
Compression Layer
Future versions:
FP16
↓
INT8
↓
Q4
for KV storage.
The document's proposed advantage comes from reducing memory pressure enough that more suspended contexts can coexist simultaneously.
Model Lifecycle Manager
Hidden hero.
States:
NOT_LOADED
LOADING
ACTIVE
WARM
UNLOADING
UNLOADED
Example:
Planner
ACTIVE
↓
WARM
↓
UNLOADED
This prevents duplicate loads.
Predictive Loading
Very important.
Planner usually followed by:
Coder
Scheduler learns:
Planner
→ Coder
pattern.
While Planner runs:
Start loading Coder
in background.
Result:
Near-zero wait.
Router Intelligence
Router should evolve.
Version 1:
Rules
Example:
contains "write"
→ coder
Version 2:
Tiny router model
~500M parameters.
Version 3:
Learned routing
based on telemetry.
What Can Be Faked For Demo?
A lot.
And that's okay.
Can fake:
KV Compression
Store JSON metadata.
Pretend compression happened.
Can fake:
SSD tier
Just move files locally.
Can fake:
Cross-phase injection
Simple restore.
Can fake:
Predictive scheduling
Static rules.
Recruiters won't know.
Engineers won't care initially.
What MUST Be Real
Absolutely real:
Dynamic model loading
Planner unload
Coder load
Process table
Live states
Scheduler
Queues
Priority
Preemption
brain-top
Real metrics
Telemetry
GPU
RAM
Task tracking
Without these:
It's a chatbot workflow.
Not an operating system.
The Actual Killer Demo
The strongest demonstration isn't code generation.
It's this:
brain-top
Window 1:
Planner RUNNING
Coder WAITING
Reasoner HIBERNATED
Window 2:
VRAM 78%
RAM 52%
SSD 12%
Window 3:
Q0 [2]
Q1 [3]
Q2 [1]
Q3 [0]
Window 4:
SWAPS
Planner
VRAM → RAM
Coder
RAM → VRAM
Then execute:
tinybrain run \
"Research Next.js 15 and generate migration script"
and watch agents move through states live.
At that moment people immediately understand:
"This isn't another AI wrapper. This is a runtime system."
The next part (Part 5) gets into the deepest engineering layer:
exact Go package structure
interfaces and abstractions
scheduler pseudocode
agent runtime contracts
model loader implementation
llama.cpp integration
Kubernetes operator design
MVP timeline (week-by-week)
what you, as a pre-final-year student targeting backend/systems roles, should prioritize and what to deliberately postpone.
Part 5 — Actual Code Architecture, Interfaces, Kubernetes Operator Design & The "Don't Ruin This Project" Roadmap
This is the level where you stop thinking:
What features should I build?
and start thinking:
What abstractions should exist?
The best systems projects are basically:
Good abstractions
Good boundaries
Good observability
Everything else becomes implementation details.
The Biggest Mistake You Can Make
Do NOT build:
Planner
Browser
Coder
Reasoner
first.
That's a trap.
Most students build agents.
The actual hard problem is:
Scheduler
Memory Manager
Runtime
Agents are replaceable.
Runtime is not.
TinyBrain's Real Core
The entire system should revolve around:
type Runtime interface {
Start()
Stop()
Submit(task Task)
Pause(agentID string)
Resume(agentID string)
Metrics() Metrics
}
Everything else becomes plugins.
The Package Hierarchy
Final Go structure:
internal/
├── scheduler/
├── runtime/
├── process/
├── kv/
├── memory/
├── loader/
├── router/
├── telemetry/
├── api/
├── operator/
└── agents/
Each package should have:
interface.go
service.go
types.go
This keeps things clean.
Runtime Interface
The scheduler should NEVER know:
llama.cpp
Ollama
vLLM
exists.
Instead:
type ModelRuntime interface {
LoadModel(model string)
UnloadModel(model string)
Generate(req Request)
SaveContext(id string)
RestoreContext(id string)
}
This abstraction is critical.
Why This Matters
Future:
Today:
llama.cpp
Tomorrow:
vLLM
Next year:
Custom CUDA Engine
No code changes elsewhere.
Agent Contract
Every agent follows one contract.
type Agent interface {
Name() string
Execute(
    ctx Context,
    task Task,
) Result
}
Planner.
Coder.
Reasoner.
Everything.
Same interface.
Process Manager Design
Process manager owns:
map[PID]*AgentProcess
Nothing else.
Responsibilities:
Create Process
Destroy Process
Track State
Track Metrics
Scheduler Interface
type Scheduler interface {
Enqueue(*AgentProcess)
Schedule()
Preempt(pid PID)
Boost()
}
Clean.
Simple.
Expandable.
Event-Driven Runtime
Everything becomes events.
Example:
Task Submitted
↓
TaskCreated
↓
RouterSelectedAgent
↓
ProcessSpawned
↓
ModelLoaded
↓
ExecutionStarted
↓
TaskCompleted
This architecture scales forever.
The Event Bus
Version 1:
chan Event
Version 2:
NATS
NATS is ideal because it is lightweight, event-driven, and commonly used for distributed control-plane style architectures.
Why Not Kafka?
Kafka solves:
millions of events
You solve:
few hundred events
Huge difference.
Kubernetes Operator Design
This is where the project becomes genuinely impressive.
Kubernetes Operators combine:
CRDs
Controllers
to create new platform primitives. (Wikipedia)
TinyBrain should become:
TinyBrain Operator
TinyBrain CRDs
Four initial CRDs:
Agent
Task
KVCache
SwapPolicy
That's enough.
Do not create 20 CRDs.
Agent Controller
Controller watches:
kind: Agent
Whenever:
spec.state: Running
controller ensures:
Model Loaded
Resources Assigned
Exactly like Kubernetes reconciliation loops. (Wikipedia)
Task Controller
Task created:
kind: Task
↓
Controller sees event
↓
Creates Process
↓
Scheduler receives process
Beautiful architecture.
KVCache Controller
This is your flex.
Most projects don't have this.
Example:
kind: KVCache
spec:
owner: coder
location: ram
Controller guarantees:
Cache exists
Cache tracked
Cache recoverable
SwapPolicy Controller
Allows users to write:
spec:
vramThreshold: 85
ramThreshold: 75
Scheduler obeys automatically.
This makes TinyBrain configurable.
Reconciliation Loop
The entire operator runs on:
for {
desired := Spec
actual := CurrentState
reconcile(desired, actual)
}
This is literally how Kubernetes controllers operate. (Wikipedia)
The Model Loader
One of the most important components.
Responsibilities:
Load
Unload
Warm
Prefetch
Evict
Only.
Never let scheduler touch:
llama.cpp
directly.
The Model Registry
Create:
models.yaml
Example:
planner:
path: planner.gguf
memory: 1.3GB
coder:
path: coder.gguf
memory: 2.3GB
Later:
Hugging Face
Model Hub
integration.
brain-top Architecture
brain-top should never query services directly.
Instead:
Metrics Stream
↓
brain-top
Like:
htop
for TinyBrain.
Metrics Source
brain-top consumes:
Scheduler Metrics
Runtime Metrics
Memory Metrics
Swap Metrics
Unified format:
type Metric struct {
Name string
Value float64
Timestamp time.Time
}
OpenTelemetry
When you reach V2:
Use:
OpenTelemetry
for traces.
This gives:
Task Timeline
Example:
Task Submitted
↓ 20ms
Planner Loaded
↓ 45ms
Planner Finished
↓ 120ms
Coder Loaded
↓ 2.1s
Task Completed
This is exactly the type of observability modern control planes rely on.
What The MVP Should ACTUALLY Do
Not:
Research Next.js
Not:
Autonomous Coding
The MVP should do:
Planner
↓
Coder
only.
One workflow.
One success case.
Nothing more.
The MVP Success Criteria
If these work:
✅ Dynamic loading
✅ Dynamic unloading
✅ Scheduler
✅ Process table
✅ brain-top
✅ Metrics
Then MVP succeeds.
Everything else:
Reasoning
RAG
Browser
Agents
can come later.
The Hidden Opportunity
This project can actually become:
A Local AI Runtime
not just:
A Portfolio Project
There is growing interest in orchestrating constrained AI workloads and extending Kubernetes with custom controllers/operators rather than building monolithic applications. Research is even exploring lightweight controller execution and custom scheduling strategies for specialized workloads. (arXiv)
If I Were Building This
I would prioritize in this exact order:
Month 1
Runtime
Process Table
Model Loader
Month 2
Scheduler
MLFQ
brain-top
Month 3
KV Manager
Swap Manager
Month 4
Planner + Coder
Month 5
Kubernetes Operator
CRDs
Controllers
Month 6
Benchmarks
Website
Documentation
Only after that would I touch:
Reasoner
Browser
Autonomous Agents
Because the runtime is the moat.
The agents are replaceable. The runtime is what makes TinyBrain OS a systems project rather than another AI wrapper.
This changes the project fundamentally.
After reading your additions, I would actually rewrite the project's core thesis.
The original thesis:
4 Fixed Agents
(Planner, Browser, Coder, Reasoner)
Dynamic Swapping
is not ambitious enough anymore.
New Core Thesis
TinyBrain OS is NOT:
An operating system for 4 AI agents
TinyBrain OS is:
A hardware-aware runtime kernel
for dynamically orchestrating
arbitrary AI agents under
resource constraints.
That's a much stronger statement.
The changes you proposed are exactly what would push this from a cool project to a genuine systems architecture.
What I Would Change In The Architecture
Remove Fixed Agent Types
Currently:
Planner
Browser
Coder
Reasoner
should NOT exist in source code.
Instead:
type Agent struct {
Name string
ModelProfile string
ToolSet []Tool
ResourceProfile ResourceProfile
Priority Priority
}
Meaning:
Planner
becomes:
name: planner
not:
PlannerAgent{}
This is a huge architectural improvement.
The Real OS Analogy
Linux doesn't know:
Chrome
VSCode
Discord
exist.
Linux only knows:
Process
Memory
CPU
Priority
Similarly:
TinyBrain should only know:
Agent
Resources
Priority
State
Everything else becomes plugins.
New Agent Architecture
Instead of:
Planner
Browser
Coder
Reasoner
you create:
Agent Registry
agents:
planner:
model: qwen-2b
tools:
  - task_split
coder:
model: deepseek-coder
tools:
  - filesystem
  - git
reviewer:
model: phi
tools:
  - static_analysis
shell:
model: stablelm
tools:
  - bash
Now users create fleets.
Not you.
Exactly as you suggested.
Biggest Architectural Upgrade
The project should become:
Hardware Aware
Current design:
Planner = 2B
Coder = 3B
Reasoner = 4B
Wrong.
Your observation is completely correct.
A runtime should NEVER assume hardware.
Instead:
Hardware Probe
↓
Capability Profile
↓
Model Assignment
New Boot Sequence
At startup:
TinyBrain Bootloader
runs:
Phase 1
Hardware Detection
Collect:
RAM
VRAM
CPU
Backend
CUDA
Metal
ROCm
Phase 2
Capability Classification
Example:
Tier 0
CPU only
Tier 1
8GB RAM
Tier 2
16GB RAM + RTX 3050
Tier 3
32GB RAM + RTX 4070
Tier 4
64GB+ Workstation
Phase 3
Agent Mapping
Instead of:
Planner → Qwen 2B
Use:
Planner → Best available model
Example:
Tier 1
Planner
Browser
↓
Same model
Tier 4
Planner → 8B
Browser → 8B
Coder → 32B
Reasoner → 70B
without changing any code.
This is the real operating-system approach.
The Single Biggest Insight
I think the actual innovation is no longer:
Dynamic Model Swapping
The actual innovation becomes:
Dynamic Capability Allocation
Meaning:
The scheduler decides:
Which agent?
Which model?
Which precision?
Which device?
Which memory tier?
for every task.
Example:
User asks:
List files in project
Scheduler:
Shell Agent
1.5B model
CPU
User asks:
Refactor microservices architecture
Scheduler:
Reasoner Agent
14B model
GPU
User asks:
Generate SQL query
Scheduler:
SQL Agent
SQLCoder
CPU
This is far more interesting than fixed agents.
The Next Major Component
I think the architecture is missing something extremely important:
Capability Manager
Current:
Scheduler
↓
Agent
Future:
Scheduler
↓
Capability Manager
↓
Agent
Capability Manager owns:
Agent Discovery
Agent Registration
Tool Registration
Model Assignment
Hardware Mapping
Like:
tinybrain agent install reviewer
or
tinybrain agent install sql-agent
The Future Cloud Design
I completely agree with your hexagonal architecture idea.
The scheduler should NEVER know:
Local
Cloud
API
It should only know:
type InferenceProvider interface {
Generate(...)
}
Then:
Local GGUF Provider
implements it.
Later:
OpenAI Provider
Anthropic Provider
vLLM Provider
implement the same interface.
Your scheduler remains unchanged.
What Becomes The Killer Feature Now?
Initially I thought:
KV Swapping
was the killer feature.
After reading your additions:
I disagree.
The real killer feature becomes:
Hardware-Adaptive Agent Runtime
Imagine:
tinybrain boot
on:
Laptop
Detected:
16GB RAM
RTX 3050
Recommended Fleet:
Planner → Qwen 2B
Coder → DeepSeek 7B Q4
Reviewer → Phi 3B
Same command on:
Workstation
Detected:
128GB RAM
2× RTX 4090
Recommended Fleet:
Planner → Llama 8B
Coder → DeepSeek 33B
Reviewer → Command-R 35B
Reasoner → Llama 70B
No code changes.
No user configuration.
Just runtime adaptation.
That is something very few local AI projects attempt.
If This Were My Project
I would officially redefine TinyBrain OS as:
A Hardware-Aware AI Runtime Kernel
that dynamically discovers hardware,
selects optimal model fleets,
allocates resources,
schedules agent execution,
and performs adaptive model swapping
across heterogeneous memory tiers.
Notice what disappeared:
Planner
Browser
Coder
Reasoner
Those are now examples.
Not architecture.
That shift alone makes the project much more scalable, much more systems-oriented, and much harder for someone else to dismiss as "just another multi-agent wrapper."
Only:
README.md
Structure:
What is TinyBrain?
Why does it exist?
High level architecture.
Current roadmap.
Current status.
One page.
