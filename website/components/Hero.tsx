"use client";

import { motion, useReducedMotion } from "framer-motion";
import { GITHUB_URL } from "@/lib/content";
import { AgentOrbit } from "@/components/AgentOrbit";

export function Hero() {
  const reduceMotion = useReducedMotion();

  return (
    <section className="section-padding relative overflow-hidden pt-32">
      <div className="pointer-events-none absolute -right-32 top-20 h-96 w-96 rounded-full bg-flame-yellow/20 blur-3xl" />
      <div className="pointer-events-none absolute -left-20 bottom-0 h-80 w-80 rounded-full bg-flame-red/10 blur-3xl" />

      <div className="mx-auto grid max-w-7xl items-center gap-16 lg:grid-cols-2">
        <motion.div
          initial={reduceMotion ? {} : { opacity: 0, y: 24 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7 }}
        >
          <p className="font-mono text-sm font-medium uppercase tracking-widest text-flame-orange">
            v1.1 — Advanced Subsystems
          </p>
          <h1 className="mt-4 font-display text-4xl font-bold leading-[1.1] tracking-tight md:text-5xl lg:text-6xl">
            An operating system
            <br />
            <span className="gradient-text">for AI agents</span>
          </h1>
          <p className="mt-6 max-w-xl text-lg leading-relaxed text-neutral-600">
            TinyBrain OS schedules, swaps, and orchestrates small LLMs on
            consumer hardware — with process tables, MLFQ preemption, paged KV
            memory, and bare-metal CUDA inference. Not a chatbot wrapper. A
            runtime kernel.
          </p>

          <div className="mt-8 flex flex-wrap gap-4">
            <a
              href={GITHUB_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="rounded-full bg-neutral-900 px-6 py-3 text-sm font-semibold text-white transition hover:bg-flame-red focus-ring"
            >
              View on GitHub
            </a>
            <a
              href="#lifecycle"
              className="rounded-full border border-neutral-300 bg-white px-6 py-3 text-sm font-semibold text-neutral-800 transition hover:border-flame-orange focus-ring"
            >
              See how it works
            </a>
          </div>

          <div className="mt-10 rounded-2xl border border-neutral-200 bg-white/60 p-4 font-mono text-sm">
            <span className="text-neutral-400">$ </span>
            <span className="text-neutral-800">git clone </span>
            <span className="text-flame-orange">TinyBrain-OS</span>
            <span className="text-neutral-800"> && go build -o tinybrain ./cmd/tinybrain</span>
          </div>
        </motion.div>

        <motion.div
          initial={reduceMotion ? {} : { opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.8, delay: 0.2 }}
        >
          <AgentOrbit />
        </motion.div>
      </div>
    </section>
  );
}
