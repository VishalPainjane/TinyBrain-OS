"use client";

import { motion, useReducedMotion } from "framer-motion";
import {
  adrs,
  externalResearch,
  inspirationSites,
  researchNotes,
  rfcs,
  REPO_DOCS_BASE,
} from "@/lib/content";

export function ResearchSection() {
  const reduceMotion = useReducedMotion();

  return (
    <section id="research" className="section-padding bg-white">
      <div className="mx-auto max-w-7xl">
        <h2 className="font-display text-3xl font-bold tracking-tight md:text-4xl">
          Research foundations
          <span className="gradient-text"> and references</span>
        </h2>
        <p className="mt-4 max-w-2xl text-lg text-neutral-600">
          Architecture drawn from OS theory, vLLM serving systems, FlashAttention,
          Kubernetes operators, and bare-metal GPU programming — documented in
          ADRs, RFCs, and research notes in the repository.
        </p>

        <div className="mt-12 grid gap-8 lg:grid-cols-3">
          <div>
            <h3 className="font-display text-lg font-bold">Repository notes</h3>
            <ul className="mt-4 space-y-2">
              {researchNotes.map((note) => (
                <li key={note.file}>
                  <a
                    href={`${REPO_DOCS_BASE}/${note.file}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="group flex flex-col rounded-lg border border-neutral-100 p-3 transition hover:border-flame-orange focus-ring"
                  >
                    <span className="font-medium group-hover:text-flame-red">
                      {note.title}
                    </span>
                    <span className="text-xs text-neutral-500">{note.topic}</span>
                  </a>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h3 className="font-display text-lg font-bold">
              Architecture decisions
            </h3>
            <ul className="mt-4 space-y-2">
              {adrs.map((adr) => (
                <li key={adr.id}>
                  <a
                    href={`${REPO_DOCS_BASE}/${adr.file}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-2 rounded-lg border border-neutral-100 p-3 text-sm transition hover:border-flame-orange focus-ring"
                  >
                    <span className="font-mono text-flame-orange">{adr.id}</span>
                    <span>{adr.title}</span>
                  </a>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h3 className="font-display text-lg font-bold">RFCs & literature</h3>
            <ul className="mt-4 space-y-2">
              {rfcs.map((rfc) => (
                <li key={rfc.id}>
                  <a
                    href={`${REPO_DOCS_BASE}/${rfc.file}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-2 rounded-lg border border-neutral-100 p-3 text-sm transition hover:border-flame-orange focus-ring"
                  >
                    <span className="font-mono text-flame-orange">{rfc.id}</span>
                    <span>{rfc.title}</span>
                  </a>
                </li>
              ))}
            </ul>
            <h4 className="mt-6 text-sm font-semibold uppercase tracking-wide text-neutral-500">
              External
            </h4>
            <ul className="mt-2 space-y-2">
              {externalResearch.map((ref) => (
                <li key={ref.url}>
                  <a
                    href={ref.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-sm text-neutral-700 underline decoration-flame-yellow decoration-2 underline-offset-2 hover:text-flame-red focus-ring rounded"
                  >
                    {ref.name}
                    <span className="ml-2 text-xs text-neutral-400">
                      {ref.type}
                    </span>
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <motion.div
          className="mt-16 rounded-2xl border border-neutral-200 bg-flame-cream p-6 md:p-8"
          initial={reduceMotion ? {} : { opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
        >
          <h3 className="font-display text-lg font-bold">Design inspiration</h3>
          <p className="mt-2 text-sm text-neutral-600">
            Showcase patterns referenced from award-winning and open-source
            product sites — editorial grids, scroll storytelling, hardware-first
            heroes, and developer-minimal CLI aesthetics.
          </p>
          <div className="mt-4 flex flex-wrap gap-3">
            {inspirationSites.map((site) => (
              <a
                key={site.url}
                href={site.url}
                target="_blank"
                rel="noopener noreferrer"
                className="rounded-full border border-neutral-300 bg-white px-4 py-2 text-sm transition hover:border-flame-orange focus-ring"
              >
                {site.name}
              </a>
            ))}
          </div>
        </motion.div>
      </div>
    </section>
  );
}
