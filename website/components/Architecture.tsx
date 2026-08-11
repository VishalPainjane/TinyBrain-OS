"use client";

import { motion, useReducedMotion } from "framer-motion";
import { ArchitectureDiagram } from "@/components/ArchitectureDiagram";

const layers = [
  {
    name: "API / Router",
    desc: "HTTP, OpenAI-compatible SSE",
    width: "100%",
    color: "#FFD60A",
  },
  {
    name: "Scheduler",
    desc: "MLFQ — never sees models",
    width: "88%",
    color: "#FF6B35",
  },
  {
    name: "Runtime",
    desc: "Load, warm, generate",
    width: "76%",
    color: "#FF6B35",
  },
  {
    name: "Inference + Memory",
    desc: "llama.cpp · KV tiers · swap",
    width: "64%",
    color: "#E63946",
  },
  {
    name: "Models (GGUF)",
    desc: "Q4_K_M on mmap",
    width: "52%",
    color: "#E63946",
  },
];

export function VisionSection() {
  const reduceMotion = useReducedMotion();

  return (
    <section id="vision" className="section-padding">
      <div className="mx-auto max-w-7xl">
        <div className="grid gap-16 lg:grid-cols-2 lg:items-center">
          <div>
            <h2 className="font-display text-3xl font-bold tracking-tight md:text-4xl">
              Why local AI breaks
              <span className="gradient-text"> on one GPU</span>
            </h2>
            <p className="mt-6 text-lg leading-relaxed text-neutral-600">
              A single 7B model can consume an entire RTX 3050&apos;s 4 GB VRAM.
              Running a second agent, a code tool, or a retrieval pipeline becomes
              impossible without evicting the first. Cloud APIs and bigger GPUs
              defeat the purpose of sovereign, offline, local AI.
            </p>
            <p className="mt-4 text-lg leading-relaxed text-neutral-600">
              TinyBrain treats agents like OS processes: scheduled, preempted,
              hibernated, and coordinated under strict memory budgets. Small
              specialized models — orchestrated smartly — beat one monolithic
              giant on the hardware you already own.
            </p>

            <dl className="mt-8 grid gap-4 sm:grid-cols-2">
              {[
                { k: "Target VRAM", v: "4 GB (RTX 3050)" },
                { k: "Target RAM", v: "16 GB consumer" },
                { k: "Scheduling", v: "Token-boundary MLFQ" },
                { k: "IPC", v: "Structured JSON only" },
              ].map((item) => (
                <div
                  key={item.k}
                  className="rounded-xl border border-neutral-200 bg-white p-4"
                >
                  <dt className="font-mono text-xs uppercase text-neutral-500">
                    {item.k}
                  </dt>
                  <dd className="mt-1 font-display font-semibold">{item.v}</dd>
                </div>
              ))}
            </dl>
          </div>

          <div className="card-surface p-8">
            <p className="font-mono text-xs uppercase tracking-widest text-flame-orange">
              Constraint envelope
            </p>
            <div className="mt-6 space-y-3">
              {layers.map((layer, i) => (
                <motion.div
                  key={layer.name}
                  initial={reduceMotion ? {} : { opacity: 0, x: -20 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: i * 0.1 }}
                >
                  <div
                    className="rounded-xl px-4 py-3 text-white transition-all"
                    style={{
                      width: layer.width,
                      backgroundColor: layer.color,
                      marginLeft: "auto",
                    }}
                  >
                    <div className="font-display text-sm font-bold">
                      {layer.name}
                    </div>
                    <div className="text-xs opacity-90">{layer.desc}</div>
                  </div>
                </motion.div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export function ArchitectureSection() {
  const reduceMotion = useReducedMotion();

  return (
    <section id="architecture" className="section-padding bg-white">
      <div className="mx-auto max-w-7xl">
        <div className="max-w-2xl">
          <h2 className="font-display text-3xl font-bold tracking-tight md:text-4xl">
            Hexagonal architecture
            <span className="gradient-text"> with hard boundaries</span>
          </h2>
          <p className="mt-4 text-lg text-neutral-600">
            Scheduler never imports inference. Runtime never imports scheduler.
            Agents never call models directly. Eight invariants — tested in CI
            import boundary checks.
          </p>
        </div>

        <div className="mt-16 grid gap-12 lg:grid-cols-2 lg:items-start">
          <ArchitectureDiagram />

          <div className="space-y-3">
            {[
              "INV-001: Scheduler never imports runtime",
              "INV-002: Runtime never imports scheduler",
              "INV-003: Agents never invoke models directly",
              "INV-004: Registry owns capability discovery",
              "INV-005: Hardware determines model selection",
              "INV-006: Structured JSON IPC only",
              "INV-007: Core never hardcodes agent roles",
              "INV-008: Inference via Provider interface only",
            ].map((inv, i) => (
              <motion.div
                key={inv}
                className="flex items-start gap-3 rounded-xl border border-neutral-100 bg-flame-cream/50 p-4"
                initial={reduceMotion ? {} : { opacity: 0, x: 12 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.05 }}
              >
                <span className="mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-flame-orange text-xs font-bold text-white">
                  {i + 1}
                </span>
                <span className="font-mono text-sm text-neutral-700">{inv}</span>
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
