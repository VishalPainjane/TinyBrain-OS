"use client";

import Link from "next/link";
import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { GITHUB_URL, navLinks } from "@/lib/content";

export function Header() {
  const [open, setOpen] = useState(false);

  return (
    <header className="fixed top-0 z-50 w-full border-b border-neutral-200/60 bg-flame-cream/90 backdrop-blur-md">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-6 md:px-12">
        <Link
          href="/"
          className="font-display text-lg font-bold tracking-tight focus-ring rounded-lg"
        >
          <span className="gradient-text">TinyBrain</span>
          <span className="text-neutral-500"> OS</span>
        </Link>

        <nav className="hidden items-center gap-8 md:flex" aria-label="Main">
          {navLinks.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              className="text-sm font-medium text-neutral-600 transition-colors hover:text-flame-red focus-ring rounded px-1"
            >
              {link.label}
            </Link>
          ))}
          <a
            href={GITHUB_URL}
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-full bg-neutral-900 px-4 py-2 text-sm font-medium text-white transition hover:bg-flame-orange focus-ring"
          >
            GitHub
          </a>
        </nav>

        <button
          type="button"
          className="md:hidden focus-ring rounded-lg p-2"
          onClick={() => setOpen(!open)}
          aria-expanded={open}
          aria-label="Toggle menu"
        >
          <span className="block h-0.5 w-6 bg-neutral-800" />
          <span className="mt-1.5 block h-0.5 w-6 bg-neutral-800" />
        </button>
      </div>

      <AnimatePresence>
        {open && (
          <motion.nav
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: "auto" }}
            exit={{ opacity: 0, height: 0 }}
            className="border-t border-neutral-200 bg-white md:hidden"
            aria-label="Mobile"
          >
            <div className="flex flex-col gap-1 p-4">
              {navLinks.map((link) => (
                <Link
                  key={link.href}
                  href={link.href}
                  className="rounded-lg px-3 py-2 text-neutral-700 hover:bg-neutral-100"
                  onClick={() => setOpen(false)}
                >
                  {link.label}
                </Link>
              ))}
              <a
                href={GITHUB_URL}
                className="mt-2 rounded-lg bg-neutral-900 px-3 py-2 text-center text-white"
                target="_blank"
                rel="noopener noreferrer"
              >
                GitHub
              </a>
            </div>
          </motion.nav>
        )}
      </AnimatePresence>
    </header>
  );
}
