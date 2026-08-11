"use client";

import { motion, useReducedMotion } from "framer-motion";

type NodeDef = {
  id: string;
  label: string;
  sub: string;
  x: number;
  y: number;
  w: number;
  h: number;
  accent: string;
  layer: "control" | "data";
};

const nodes: NodeDef[] = [
  { id: "router", label: "Router", sub: "HTTP · SSE", x: 48, y: 28, w: 100, h: 52, accent: "#FFD60A", layer: "control" },
  { id: "scheduler", label: "Scheduler", sub: "MLFQ Q0–Q3", x: 178, y: 28, w: 110, h: 52, accent: "#FF6B35", layer: "control" },
  { id: "process", label: "Process Table", sub: "PID · states", x: 318, y: 28, w: 118, h: 52, accent: "#E63946", layer: "control" },
  { id: "runtime", label: "Runtime", sub: "load · generate", x: 108, y: 118, w: 110, h: 52, accent: "#FF6B35", layer: "control" },
  { id: "registry", label: "Registry", sub: "agents · models", x: 268, y: 118, w: 110, h: 52, accent: "#FFD60A", layer: "control" },
  { id: "inference", label: "Inference", sub: "llama.cpp · CUDA", x: 48, y: 228, w: 110, h: 52, accent: "#E63946", layer: "data" },
  { id: "memory", label: "KV / Memory", sub: "paged · swap", x: 188, y: 228, w: 118, h: 52, accent: "#FF6B35", layer: "data" },
  { id: "telemetry", label: "Telemetry", sub: "brain-top", x: 338, y: 228, w: 100, h: 52, accent: "#FFD60A", layer: "data" },
];

const edges: { from: string; to: string; curved?: boolean }[] = [
  { from: "router", to: "scheduler" },
  { from: "scheduler", to: "process" },
  { from: "scheduler", to: "runtime", curved: true },
  { from: "runtime", to: "registry" },
  { from: "runtime", to: "inference", curved: true },
  { from: "runtime", to: "memory", curved: true },
  { from: "memory", to: "telemetry" },
  { from: "inference", to: "memory" },
];

function center(n: NodeDef) {
  return { cx: n.x + n.w / 2, cy: n.y + n.h / 2 };
}

function pathBetween(a: NodeDef, b: NodeDef, curved?: boolean) {
  const from = center(a);
  const to = center(b);
  if (curved) {
    const midY = (from.cy + to.cy) / 2;
    return `M ${from.cx} ${from.cy} C ${from.cx} ${midY}, ${to.cx} ${midY}, ${to.cx} ${to.cy}`;
  }
  return `M ${from.cx} ${from.cy} L ${to.cx} ${to.cy}`;
}

function NodeCard({
  node,
  delay,
  reduceMotion,
}: {
  node: NodeDef;
  delay: number;
  reduceMotion: boolean | null;
}) {
  return (
    <motion.g
      initial={reduceMotion ? {} : { opacity: 0, scale: 0.92 }}
      whileInView={{ opacity: 1, scale: 1 }}
      viewport={{ once: true }}
      transition={{ delay, duration: 0.45, ease: "easeOut" }}
    >
      <defs>
        <filter id={`glow-${node.id}`} x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur stdDeviation="4" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
        <clipPath id={`clip-${node.id}`}>
          <rect x={node.x} y={node.y} width={node.w} height={node.h} rx="14" />
        </clipPath>
      </defs>
      <rect
        x={node.x}
        y={node.y}
        width={node.w}
        height={node.h}
        rx="14"
        fill="white"
        stroke={node.accent}
        strokeWidth="2"
        filter={`url(#glow-${node.id})`}
        opacity="0.95"
      />
      <rect
        x={node.x}
        y={node.y}
        width={node.w}
        height={6}
        fill={node.accent}
        clipPath={`url(#clip-${node.id})`}
        opacity="0.9"
      />
      <rect
        x={node.x}
        y={node.y}
        width={node.w}
        height={node.h}
        rx="14"
        fill="none"
        stroke={node.accent}
        strokeWidth="2"
        opacity="0.95"
      />
      <text
        x={node.x + node.w / 2}
        y={node.y + 30}
        textAnchor="middle"
        fill="#171717"
        fontSize="13"
        fontWeight="700"
        fontFamily="var(--font-syne), system-ui"
      >
        {node.label}
      </text>
      <text
        x={node.x + node.w / 2}
        y={node.y + 44}
        textAnchor="middle"
        fill="#737373"
        fontSize="10"
        fontFamily="var(--font-jetbrains), monospace"
      >
        {node.sub}
      </text>
    </motion.g>
  );
}

export function ArchitectureDiagram() {
  const reduceMotion = useReducedMotion();
  const nodeMap = Object.fromEntries(nodes.map((n) => [n.id, n]));

  return (
    <div className="relative overflow-hidden rounded-3xl border border-neutral-200/80 bg-gradient-to-br from-white via-flame-cream to-white p-4 shadow-xl shadow-flame-orange/5 md:p-8">
      {/* Ambient grid */}
      <div
        className="pointer-events-none absolute inset-0 opacity-[0.35]"
        style={{
          backgroundImage:
            "radial-gradient(circle at 1px 1px, #FF6B35 1px, transparent 0)",
          backgroundSize: "24px 24px",
        }}
      />
      <div className="pointer-events-none absolute -right-20 -top-20 h-64 w-64 rounded-full bg-flame-yellow/20 blur-3xl" />
      <div className="pointer-events-none absolute -bottom-16 -left-16 h-56 w-56 rounded-full bg-flame-red/10 blur-3xl" />

      <div className="relative flex items-center justify-between px-2 pb-4">
        <p className="font-mono text-[10px] uppercase tracking-[0.2em] text-flame-orange">
          System map
        </p>
        <div className="flex gap-2">
          <span className="rounded-full border border-flame-orange/30 bg-flame-orange/10 px-3 py-1 text-[10px] font-bold uppercase tracking-wider text-flame-orange">
            Control plane
          </span>
          <span className="rounded-full border border-flame-red/30 bg-flame-red/10 px-3 py-1 text-[10px] font-bold uppercase tracking-wider text-flame-red">
            Data plane
          </span>
        </div>
      </div>

      <svg
        viewBox="0 0 480 300"
        className="relative z-10 w-full h-auto"
        role="img"
        aria-label="TinyBrain OS architecture diagram showing Router, Scheduler, Process Table, Runtime, Registry, Inference, KV Memory, and Telemetry"
      >
        <defs>
          <linearGradient id="flowGrad" x1="0%" y1="0%" x2="100%" y2="0%">
            <stop offset="0%" stopColor="#E63946" />
            <stop offset="50%" stopColor="#FF6B35" />
            <stop offset="100%" stopColor="#FFD60A" />
          </linearGradient>
          <linearGradient id="zoneControl" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#FF6B35" stopOpacity="0.08" />
            <stop offset="100%" stopColor="#FF6B35" stopOpacity="0" />
          </linearGradient>
          <linearGradient id="zoneData" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stopColor="#E63946" stopOpacity="0.06" />
            <stop offset="100%" stopColor="#E63946" stopOpacity="0" />
          </linearGradient>
        </defs>

        {/* Zone backgrounds */}
        <rect x="20" y="12" width="440" height="168" rx="16" fill="url(#zoneControl)" />
        <rect x="20" y="200" width="440" height="88" rx="16" fill="url(#zoneData)" />

        <text x="32" y="22" fill="#FF6B35" fontSize="9" fontWeight="600" fontFamily="monospace" opacity="0.8">
          CONTROL PLANE
        </text>
        <text x="32" y="212" fill="#E63946" fontSize="9" fontWeight="600" fontFamily="monospace" opacity="0.8">
          DATA PLANE
        </text>

        {/* Connection paths */}
        {edges.map((edge, i) => {
          const a = nodeMap[edge.from];
          const b = nodeMap[edge.to];
          if (!a || !b) return null;
          const d = pathBetween(a, b, edge.curved);
          return (
            <g key={`${edge.from}-${edge.to}`}>
              <motion.path
                d={d}
                fill="none"
                stroke="url(#flowGrad)"
                strokeWidth="2"
                strokeLinecap="round"
                opacity="0.35"
                initial={reduceMotion ? {} : { pathLength: 0 }}
                whileInView={{ pathLength: 1 }}
                viewport={{ once: true }}
                transition={{ duration: 0.8, delay: i * 0.08 }}
              />
              <motion.path
                d={d}
                fill="none"
                stroke="url(#flowGrad)"
                strokeWidth="2"
                strokeLinecap="round"
                strokeDasharray="6 10"
                initial={reduceMotion ? {} : { pathLength: 0, strokeDashoffset: 0 }}
                whileInView={{ pathLength: 1 }}
                viewport={{ once: true }}
                animate={
                  reduceMotion
                    ? {}
                    : { strokeDashoffset: [0, -32] }
                }
                transition={{
                  pathLength: { duration: 0.8, delay: i * 0.08 },
                  strokeDashoffset: {
                    duration: 2.5,
                    repeat: Infinity,
                    ease: "linear",
                    delay: i * 0.15,
                  },
                }}
              />
            </g>
          );
        })}

        {nodes.map((node, i) => (
          <NodeCard key={node.id} node={node} delay={i * 0.06} reduceMotion={reduceMotion} />
        ))}

        {/* Boundary badge */}
        <motion.g
          initial={reduceMotion ? {} : { opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
          transition={{ delay: 0.6 }}
        >
          <rect x="110" y="178" width="260" height="22" rx="11" fill="#171717" />
          <text
            x="240"
            y="192"
            textAnchor="middle"
            fill="white"
            fontSize="8"
            fontFamily="monospace"
            fontWeight="600"
          >
            InferenceProvider port · no scheduler imports
          </text>
        </motion.g>
      </svg>
    </div>
  );
}
