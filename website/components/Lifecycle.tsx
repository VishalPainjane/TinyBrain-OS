"use client";

import { motion, useReducedMotion } from "framer-motion";
import { lifecycleSteps } from "@/lib/content";

export function LifecycleSection() {
  const reduceMotion = useReducedMotion();

  return (
    <section id="lifecycle" className="section-padding bg-white">
      <div className="mx-auto max-w-7xl">
        <div className="max-w-2xl">
          <h2 className="font-display text-3xl font-bold tracking-tight md:text-4xl">
            Request lifecycle
            <span className="gradient-text"> end to end</span>
          </h2>
          <p className="mt-4 text-lg text-neutral-600">
            From HTTP submission to token stream and brain-top telemetry — every
            hop respects architectural boundaries.
          </p>
        </div>

        <div className="mt-16 relative">
          <div
            className="absolute left-4 top-0 bottom-0 w-px bg-gradient-to-b from-flame-red via-flame-orange to-flame-yellow md:left-1/2 md:-translate-x-px"
            aria-hidden
          />

          <div className="space-y-8">
            {lifecycleSteps.map((step, i) => (
              <motion.div
                key={step.step}
                className={`relative flex flex-col gap-4 md:flex-row md:items-center ${
                  i % 2 === 0 ? "md:flex-row" : "md:flex-row-reverse"
                }`}
                initial={reduceMotion ? {} : { opacity: 0, y: 24 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true, margin: "-60px" }}
                transition={{ delay: 0.05 }}
              >
                <div className="md:w-1/2 md:px-12">
                  <div
                    className={`card-surface p-6 ${
                      i % 2 === 0 ? "md:mr-8" : "md:ml-8"
                    }`}
                  >
                    <span className="font-mono text-sm font-bold text-flame-orange">
                      {step.step}
                    </span>
                    <h3 className="mt-2 font-display text-xl font-bold">
                      {step.title}
                    </h3>
                    <p className="mt-2 text-neutral-600">{step.detail}</p>
                  </div>
                </div>

                <div
                  className="absolute left-4 flex h-8 w-8 -translate-x-1/2 items-center justify-center rounded-full border-2 border-white bg-neutral-900 text-xs font-bold text-white md:left-1/2"
                  aria-hidden
                >
                  {i + 1}
                </div>

                <div className="hidden md:block md:w-1/2" />
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
