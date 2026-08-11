import Link from "next/link";
import { Header } from "@/components/Header";
import { Footer } from "@/components/Footer";
import { docLinks, GITHUB_URL, REPO_DOCS_BASE } from "@/lib/content";

const categories = [...new Set(docLinks.map((d) => d.category))];

export default function DocsPage() {
  return (
    <>
      <Header />
      <main className="section-padding pt-28">
        <div className="mx-auto max-w-4xl">
          <Link
            href="/"
            className="text-sm text-neutral-500 hover:text-flame-orange focus-ring rounded"
          >
            ← Back to home
          </Link>
          <h1 className="mt-6 font-display text-4xl font-bold tracking-tight">
            Documentation
            <span className="gradient-text"> hub</span>
          </h1>
          <p className="mt-4 text-lg text-neutral-600">
            Links to guides, architecture, governance, and reference material in
            the repository. TinyBrain runs in your terminal; this hub points to
            the source-of-truth markdown.
          </p>

          <div className="mt-8 rounded-2xl border border-neutral-200 bg-white p-6">
            <h2 className="font-display font-bold">Quick commands</h2>
            <pre className="mt-4 overflow-x-auto rounded-xl bg-neutral-900 p-4 font-mono text-sm text-neutral-100">
{`go build -o tinybrain ./cmd/tinybrain
go build -o brain-top ./cmd/brain-top
./tinybrain doctor
./tinybrain probe --json
./tinybrain models list
./brain-top snapshot`}
            </pre>
          </div>

          {categories.map((category) => (
            <section key={category} className="mt-12">
              <h2 className="font-mono text-sm font-bold uppercase tracking-widest text-flame-orange">
                {category}
              </h2>
              <ul className="mt-4 space-y-3">
                {docLinks
                  .filter((d) => d.category === category)
                  .map((doc) => (
                    <li key={doc.href}>
                      <a
                        href={doc.href}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="group block rounded-xl border border-neutral-200 bg-white p-5 transition hover:border-flame-orange focus-ring"
                      >
                        <span className="font-display font-bold group-hover:text-flame-red">
                          {doc.title}
                        </span>
                        <p className="mt-1 text-sm text-neutral-600">
                          {doc.description}
                        </p>
                      </a>
                    </li>
                  ))}
              </ul>
            </section>
          ))}

          <div className="mt-16 flex flex-wrap gap-4">
            <a
              href={GITHUB_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="rounded-full bg-neutral-900 px-6 py-3 text-sm font-semibold text-white hover:bg-flame-red focus-ring"
            >
              Open repository
            </a>
            <a
              href={`${REPO_DOCS_BASE}/docs/`}
              target="_blank"
              rel="noopener noreferrer"
              className="rounded-full border border-neutral-300 px-6 py-3 text-sm font-semibold hover:border-flame-orange focus-ring"
            >
              Browse docs/ on GitHub
            </a>
          </div>
        </div>
      </main>
      <Footer />
    </>
  );
}
