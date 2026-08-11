import { GITHUB_URL } from "@/lib/content";

export function Footer() {
  return (
    <footer className="border-t border-neutral-200 bg-white px-6 py-12 md:px-12">
      <div className="mx-auto flex max-w-7xl flex-col gap-8 md:flex-row md:items-center md:justify-between">
        <div>
          <p className="font-display text-xl font-bold">
            <span className="gradient-text">TinyBrain</span> OS
          </p>
          <p className="mt-2 max-w-md text-sm text-neutral-600">
            A hardware-aware AI runtime kernel. Built with discipline. Engineered
            from first principles. MIT License.
          </p>
        </div>
        <div className="flex flex-wrap gap-6 text-sm">
          <a
            href={GITHUB_URL}
            target="_blank"
            rel="noopener noreferrer"
            className="text-neutral-600 hover:text-flame-red focus-ring rounded"
          >
            GitHub
          </a>
          <a href="/docs" className="text-neutral-600 hover:text-flame-red focus-ring rounded">
            Documentation
          </a>
          <a
            href={`${GITHUB_URL}/blob/main/SECURITY.md`}
            target="_blank"
            rel="noopener noreferrer"
            className="text-neutral-600 hover:text-flame-red focus-ring rounded"
          >
            Security
          </a>
          <a
            href={`${GITHUB_URL}/blob/main/CHANGELOG.md`}
            target="_blank"
            rel="noopener noreferrer"
            className="text-neutral-600 hover:text-flame-red focus-ring rounded"
          >
            Changelog
          </a>
        </div>
      </div>
      <p className="mx-auto mt-8 max-w-7xl text-center text-xs text-neutral-400">
        This is not a chatbot wrapper. This is a runtime system.
      </p>
    </footer>
  );
}
