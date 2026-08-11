import type { Metadata } from "next";
import { DM_Sans, JetBrains_Mono, Syne } from "next/font/google";
import "./globals.css";

const syne = Syne({
  subsets: ["latin"],
  variable: "--font-syne",
  display: "swap",
});

const dmSans = DM_Sans({
  subsets: ["latin"],
  variable: "--font-dm-sans",
  display: "swap",
});

const jetbrains = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-jetbrains",
  display: "swap",
});

export const metadata: Metadata = {
  title: "TinyBrain OS — AI Runtime Kernel for Local Hardware",
  description:
    "A hardware-aware AI runtime kernel that orchestrates small LLM agents under strict VRAM and RAM budgets. Process scheduling, paged KV memory, and bare-metal CUDA inference on consumer hardware.",
  keywords: [
    "TinyBrain",
    "local AI",
    "LLM runtime",
    "MLFQ scheduler",
    "KV cache",
    "llama.cpp",
    "open source",
  ],
  openGraph: {
    title: "TinyBrain OS",
    description: "An operating system for AI agents on local hardware.",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      className={`${syne.variable} ${dmSans.variable} ${jetbrains.variable}`}
    >
      <body className="font-sans">{children}</body>
    </html>
  );
}
