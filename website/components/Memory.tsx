"use client";

import { motion, useReducedMotion } from "framer-motion";
import { memoryTiers } from "@/lib/content";

export function MemorySection() {
  const reduceMotion = useReducedMotion();

  return (
    <section id="memory" className="section-padding">
      <div className="mx-auto max-w-7xl">
        <div className="grid gap-16 lg:grid-cols-2 lg:items-center">
          <div>
            <h2 className="font-display text-3xl font-bold tracking-tight md:text-4xl">
              Paged virtual memory
              <span className="gradient-text"> for attention state</span>
            </h2>
            <p className="mt-4 text-lg text-neutral-600">
              Inspired by PagedAttention and CPU cache hierarchies. KV blocks map
              through a page table; spill to RAM or NVMe without reloading
              multi-gigabyte weights. Zstandard compression on lower tiers.
            </p>

            <div className="mt-8 space-y-4">
              {memoryTiers.map((tier, i) => (
                <motion.div
                  key={tier.tier}
                  className="flex items-center gap-4 rounded-xl border border-neutral-200 bg-white p-4"
                  initial={reduceMotion ? {} : { opacity: 0, x: -16 }}
                  whileInView={{ opacity: 1, x: 0 }}
                  viewport={{ once: true }}
                  transition={{ delay: i * 0.1 }}
                >
                  <div
                    className="h-12 w-12 shrink-0 rounded-lg"
                    style={{ backgroundColor: tier.color }}
                  />
                  <div>
                    <div className="font-display font-bold">{tier.tier}</div>
                    <div className="text-sm text-neutral-600">{tier.label}</div>
                    <div className="font-mono text-xs text-neutral-500">
                      {tier.target}
                    </div>
                  </div>
                </motion.div>
              ))}
            </div>
          </div>

          <div className="card-surface p-6 font-mono text-xs md:p-8">
            <p className="text-neutral-500">brain-top snapshot</p>
            <pre className="mt-4 overflow-x-auto leading-relaxed text-neutral-800">
{`┌─────────────────────────────────────┐
│            AGENT PROCESSES          │
├─────────────────────────────────────┤
│  PID 101  Planner     RUNNING       │
│  PID 102  Browser     WAITING       │
│  PID 103  Coder       PREEMPTED     │
│  PID 104  Reasoner    HIBERNATED    │
├─────────────────────────────────────┤
│  GPU   ███████░░░  72%              │
│  RAM   █████░░░░░  55%              │
│  SSD   ██░░░░░░░░  21%              │
├─────────────────────────────────────┤
│  Q0 [2]  Q1 [4]  Q2 [1]  Q3 [0]     │
├─────────────────────────────────────┤
│  SWAPS                              │
│  Planner   VRAM → RAM               │
│  Coder     RAM → VRAM               │
└─────────────────────────────────────┘`}
            </pre>
            <p className="mt-4 text-sm text-neutral-500">
              htop for AI agents — real-time process states, MLFQ queues, and
              tier movement.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
