"use client";

import { motion, useReducedMotion } from "framer-motion";
import { features } from "@/lib/content";

const accentMap = {
  red: "border-flame-red/30 bg-flame-red/5",
  orange: "border-flame-orange/30 bg-flame-orange/5",
  yellow: "border-flame-yellow/50 bg-flame-yellow/10",
};

export function FeaturesSection() {
  const reduceMotion = useReducedMotion();

  return (
    <section className="section-padding">
      <div className="mx-auto max-w-7xl">
        <h2 className="font-display text-3xl font-bold tracking-tight md:text-4xl">
          Five layers of
          <span className="gradient-text"> novel capability</span>
        </h2>
        <p className="mt-4 max-w-2xl text-lg text-neutral-600">
          Dynamic allocation, token preemption, paged KV memory, context
          hibernation, and hardware-adaptive fleets — engineered from first
          principles, not bolted onto a Python wrapper.
        </p>

        <div className="mt-12 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {features.map((feature, i) => (
            <motion.article
              key={feature.title}
              className={`card-surface p-6 border ${accentMap[feature.accent as keyof typeof accentMap]}`}
              initial={reduceMotion ? {} : { opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true, margin: "-40px" }}
              transition={{ delay: i * 0.08 }}
              whileHover={{ y: -4 }}
            >
              <h3 className="font-display text-lg font-bold">{feature.title}</h3>
              <p className="mt-3 text-sm leading-relaxed text-neutral-600">
                {feature.description}
              </p>
            </motion.article>
          ))}
        </div>
      </div>
    </section>
  );
}
