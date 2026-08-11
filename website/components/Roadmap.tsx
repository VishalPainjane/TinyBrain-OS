"use client";

import { motion, useReducedMotion } from "framer-motion";
import { releases } from "@/lib/content";

export function RoadmapSection() {
  const reduceMotion = useReducedMotion();

  return (
    <section className="section-padding">
      <div className="mx-auto max-w-7xl">
        <h2 className="font-display text-3xl font-bold tracking-tight md:text-4xl">
          Release evolution
          <span className="gradient-text"> v0.3 → v1.1</span>
        </h2>
        <p className="mt-4 max-w-2xl text-lg text-neutral-600">
          Monthly sprints with tagged releases — kernel first, inference engine,
          scheduler and memory, agents, Kubernetes control plane, then production
          maturity.
        </p>

        <div className="mt-12 relative">
          <div
            className="absolute left-0 top-0 bottom-0 w-1 rounded-full bg-gradient-to-b from-flame-red via-flame-orange to-flame-yellow md:left-8"
            aria-hidden
          />
          <div className="space-y-6 pl-8 md:pl-20">
            {releases.map((release, i) => (
              <motion.div
                key={release.version}
                className={`relative rounded-2xl border p-6 ${
                  release.current
                    ? "border-flame-orange bg-white shadow-lg shadow-flame-orange/10"
                    : "border-neutral-200 bg-white/80"
                }`}
                initial={reduceMotion ? {} : { opacity: 0, x: 20 }}
                whileInView={{ opacity: 1, x: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.06 }}
              >
                <div
                  className="absolute -left-[33px] top-8 hidden h-4 w-4 rounded-full border-2 border-white bg-flame-orange md:block"
                  aria-hidden
                />
                <div className="flex flex-wrap items-center gap-3">
                  <span className="font-mono text-lg font-bold text-flame-red">
                    {release.version}
                  </span>
                  <span className="font-display font-semibold">
                    {release.label}
                  </span>
                  {release.current && (
                    <span className="rounded-full bg-flame-yellow px-2 py-0.5 text-xs font-bold uppercase">
                      Current
                    </span>
                  )}
                </div>
                <p className="mt-2 text-neutral-600">{release.detail}</p>
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
