"use client";

import { motion, useReducedMotion, useAnimationFrame } from "framer-motion";
import { useRef, useState } from "react";

// ═══════════════════════════════════════════════════════════════════
//  COORDINATE SYSTEM
//
//  Everything is positioned in CSS % relative to the card container.
//  The orbit ring is an HTML div (border-radius:50%), the kernel is
//  an HTML div, and the agent pills are HTML divs — all sharing the
//  same CSS positioning context.  SVG is used only for lightweight
//  decorative animations (pulse rings, spokes, orbiting dots) and
//  sits at inset-0 on the same square container, so viewBox 0-100
//  maps 1:1 to CSS percentages.
//
//  CY = 38% shifts the orbit upward so the bottom clears the
//  telemetry dashboard.
// ═══════════════════════════════════════════════════════════════════

const CX = 50;       // center X (%)
const CY = 38;       // center Y (%)
const ORBIT_R = 26;  // orbit radius (%)
const ORBIT_D = ORBIT_R * 2; // orbit diameter (%)
const KERNEL_D = 20; // kernel diameter (%)

// Evenly spaced at 72° — produces a clean pentagonal shape
const agents = [
  { label: "Planner",  angle: -90,  bg: "#E63946", icon: "◈" },
  { label: "Reasoner", angle: -18,  bg: "#D62839", icon: "◎" },
  { label: "Coder",    angle:  54,  bg: "#FF6B35", icon: "⌘" },
  { label: "Shell",    angle: 126,  bg: "#F4A024", icon: "▸", dark: true },
  { label: "Browser",  angle: 198,  bg: "#FF6B35", icon: "◉" },
];

function pos(angleDeg: number, r = ORBIT_R) {
  const rad = (angleDeg * Math.PI) / 180;
  return { x: CX + r * Math.cos(rad), y: CY + r * Math.sin(rad) };
}

/* ─── Orbiting dot (SVG) ─── */
function OrbitDot({
  offset,
  color,
  size = 1.3,
}: {
  offset: number;
  color: string;
  size?: number;
}) {
  const [p, setP] = useState(() => {
    const a = offset * 2 * Math.PI - Math.PI / 2;
    return { x: CX + ORBIT_R * Math.cos(a), y: CY + ORBIT_R * Math.sin(a) };
  });
  const t0 = useRef<number | null>(null);

  useAnimationFrame((t) => {
    if (!t0.current) t0.current = t;
    const f = ((t - t0.current) / 10000) % 1;
    const a = (f + offset) * 2 * Math.PI - Math.PI / 2;
    setP({ x: CX + ORBIT_R * Math.cos(a), y: CY + ORBIT_R * Math.sin(a) });
  });

  return (
    <circle
      cx={p.x}
      cy={p.y}
      r={size}
      fill={color}
      style={{ filter: `drop-shadow(0 0 4px ${color})` }}
    />
  );
}

// ═══════════════════════════════════════════════════════════════════

export function AgentOrbit() {
  const rm = useReducedMotion();

  return (
    <div className="relative mx-auto w-full max-w-lg">
      {/* Ambient glow behind card */}
      <div
        className="pointer-events-none absolute inset-0 rounded-full opacity-50 blur-3xl"
        style={{
          background:
            "radial-gradient(circle at 50% 36%, rgba(255,107,53,0.4) 0%, rgba(255,214,10,0.15) 50%, transparent 72%)",
        }}
      />

      {/* ── Card ── */}
      <div className="relative aspect-square rounded-3xl border border-white/50 bg-white/45 shadow-2xl shadow-flame-orange/8 backdrop-blur-lg">
        {/* Dot grid texture */}
        <div
          className="pointer-events-none absolute inset-0 rounded-3xl opacity-[0.15]"
          style={{
            backgroundImage:
              "radial-gradient(circle, #FF6B35 0.6px, transparent 0.6px)",
            backgroundSize: "16px 16px",
          }}
        />

        {/* ═══════ SVG decoration layer ═══════ */}
        <svg
          className="absolute inset-0 h-full w-full"
          viewBox="0 0 100 100"
          aria-hidden
        >
          <defs>
            <radialGradient id="kglow" cx="50%" cy="50%" r="50%">
              <stop offset="0%" stopColor="#FF6B35" stopOpacity="0.35" />
              <stop offset="100%" stopColor="#FF6B35" stopOpacity="0" />
            </radialGradient>
          </defs>

          {/* Pulse rings */}
          {!rm && (
            <>
              <motion.circle
                cx={CX}
                cy={CY}
                r={6}
                fill="none"
                stroke="#FF6B35"
                strokeWidth="0.2"
                animate={{ r: [6, 22], opacity: [0.45, 0] }}
                transition={{
                  duration: 2.8,
                  repeat: Infinity,
                  ease: "easeOut",
                }}
              />
              <motion.circle
                cx={CX}
                cy={CY}
                r={6}
                fill="none"
                stroke="#FFD60A"
                strokeWidth="0.15"
                animate={{ r: [6, 28], opacity: [0.3, 0] }}
                transition={{
                  duration: 2.8,
                  repeat: Infinity,
                  ease: "easeOut",
                  delay: 1,
                }}
              />
            </>
          )}

          {/* Spoke lines (kernel → agent) */}
          {agents.map((agent, i) => {
            const p = pos(agent.angle);
            return (
              <motion.path
                key={agent.label}
                d={`M ${CX} ${CY} L ${p.x} ${p.y}`}
                stroke={agent.bg}
                strokeWidth="0.35"
                fill="none"
                initial={rm ? { opacity: 0.15 } : { pathLength: 0, opacity: 0 }}
                animate={{ pathLength: 1, opacity: [0.08, 0.25, 0.08] }}
                transition={{
                  pathLength: { duration: 0.8, delay: 0.4 + i * 0.12 },
                  opacity: {
                    duration: 4,
                    repeat: Infinity,
                    delay: i * 0.5,
                  },
                }}
              />
            );
          })}

          {/* Kernel ambient glow */}
          <circle cx={CX} cy={CY} r={12} fill="url(#kglow)" />

          {/* Orbiting tokens */}
          {!rm && (
            <>
              <OrbitDot offset={0} color="#FFD60A" size={1.4} />
              <OrbitDot offset={0.33} color="#FF6B35" size={1.1} />
              <OrbitDot offset={0.66} color="#E63946" size={1.0} />
            </>
          )}
        </svg>

        {/* ═══════ HTML layer (same % coord system) ═══════ */}

        {/* Orbit ring — CSS div with dashed border, centered at (CX, CY) */}
        <div
          className="absolute z-10"
          style={{
            width: `${ORBIT_D}%`,
            height: `${ORBIT_D}%`,
            left: `${CX}%`,
            top: `${CY}%`,
            transform: "translate(-50%, -50%)",
          }}
        >
          <motion.div
            className="h-full w-full rounded-full"
            style={{ border: "1.5px dashed rgba(255, 107, 53, 0.3)" }}
            animate={rm ? {} : { rotate: 360 }}
            transition={{ duration: 90, repeat: Infinity, ease: "linear" }}
          />
        </div>

        {/* Kernel chip */}
        <div
          className="absolute z-20"
          style={{
            width: `${KERNEL_D}%`,
            height: `${KERNEL_D}%`,
            left: `${CX}%`,
            top: `${CY}%`,
            transform: "translate(-50%, -50%)",
          }}
        >
          <motion.div
            className="flex h-full w-full items-center justify-center rounded-full border border-white/15 bg-neutral-900"
            style={{
              boxShadow:
                "0 0 40px rgba(255,107,53,0.4), 0 0 80px rgba(255,107,53,0.12), inset 0 1px 0 rgba(255,255,255,0.08)",
            }}
            animate={rm ? {} : { scale: [1, 1.035, 1] }}
            transition={{ duration: 3, repeat: Infinity, ease: "easeInOut" }}
          >
            <div className="absolute inset-0 rounded-full bg-gradient-to-br from-flame-orange/20 to-transparent" />
            <div className="relative text-center">
              <span className="block font-mono text-[6.5px] uppercase tracking-[0.22em] text-flame-yellow/80 sm:text-[7.5px] md:text-[8px]">
                TinyBrain
              </span>
              <span className="block font-display text-[11px] font-bold uppercase tracking-widest text-white sm:text-sm md:text-base">
                Kernel
              </span>
            </div>
          </motion.div>
        </div>

        {/* Agent pills */}
        {agents.map((agent, i) => {
          const p = pos(agent.angle);
          const isDark = "dark" in agent && agent.dark;
          return (
            <div
              key={agent.label}
              className="absolute z-30"
              style={{
                left: `${p.x}%`,
                top: `${p.y}%`,
                transform: "translate(-50%, -50%)",
              }}
            >
              <motion.div
                initial={rm ? {} : { opacity: 0, scale: 0.65 }}
                animate={{
                  opacity: 1,
                  scale: 1,
                  y: rm ? 0 : [0, -3, 0],
                }}
                transition={{
                  opacity: { duration: 0.45, delay: 0.5 + i * 0.1 },
                  scale: { duration: 0.45, delay: 0.5 + i * 0.1 },
                  y: {
                    duration: 3.2 + i * 0.3,
                    repeat: Infinity,
                    ease: "easeInOut",
                    delay: i * 0.35,
                  },
                }}
              >
                <div
                  className="flex items-center gap-1 whitespace-nowrap rounded-xl px-2.5 py-1 shadow-md transition-transform hover:scale-110 sm:rounded-2xl sm:px-3 sm:py-1.5"
                  style={{
                    backgroundColor: agent.bg,
                    boxShadow: `0 4px 14px ${agent.bg}50, 0 1px 3px rgba(0,0,0,0.1)`,
                  }}
                >
                <span
                  className={`text-[8px] sm:text-[9px] ${isDark ? "text-neutral-800/60" : "text-white/75"}`}
                >
                  {agent.icon}
                </span>
                <span
                  className={`font-mono text-[8px] font-bold uppercase tracking-wide sm:text-[9px] md:text-[10px] ${isDark ? "text-neutral-800" : "text-white"}`}
                >
                  {agent.label}
                </span>
              </div>
            </motion.div>
            </div>
          );
        })}

        {/* ── Telemetry dashboard ── */}
        <motion.div
          className="absolute bottom-2.5 left-2.5 right-2.5 z-40 rounded-xl border border-white/50 bg-white/92 p-2.5 shadow-lg backdrop-blur-md sm:bottom-3.5 sm:left-3.5 sm:right-3.5 sm:rounded-2xl sm:p-3 md:bottom-4 md:left-4 md:right-4 md:p-3.5"
          initial={rm ? {} : { opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.7 }}
        >
          <div className="mb-1.5 flex items-center justify-between sm:mb-2">
            <span className="font-mono text-[7px] font-bold uppercase tracking-[0.15em] text-neutral-400 sm:text-[8px] md:text-[9px]">
              Live Telemetry
            </span>
            <span className="flex items-center gap-1">
              <span className="relative flex h-1.5 w-1.5 sm:h-2 sm:w-2">
                <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-50" />
                <span className="relative inline-flex h-full w-full rounded-full bg-emerald-500" />
              </span>
              <span className="font-mono text-[7px] font-semibold text-emerald-600 sm:text-[8px] md:text-[9px]">
                RUNNING
              </span>
            </span>
          </div>
          <div className="grid grid-cols-3 gap-1.5 sm:gap-2.5 md:gap-3">
            {[
              { label: "VRAM", pct: 72, color: "#E63946" },
              { label: "RAM", pct: 55, color: "#FF6B35" },
              { label: "Q0 queue", pct: 40, color: "#F4A024" },
            ].map((bar) => (
              <div key={bar.label}>
                <div className="mb-0.5 flex justify-between font-mono text-[7px] text-neutral-400 sm:mb-1 sm:text-[8px] md:text-[9px]">
                  <span>{bar.label}</span>
                  <span>{bar.pct}%</span>
                </div>
                <div className="h-1 overflow-hidden rounded-full bg-neutral-100 sm:h-1.5 md:h-2">
                  <motion.div
                    className="h-full rounded-full"
                    style={{ backgroundColor: bar.color }}
                    initial={{ width: 0 }}
                    animate={{ width: `${bar.pct}%` }}
                    transition={{ duration: 1, delay: 0.9, ease: "easeOut" }}
                  />
                </div>
              </div>
            ))}
          </div>
        </motion.div>
      </div>
    </div>
  );
}
