# TinyBrain OS — Showcase Website

Marketing and documentation hub for the TinyBrain OS project. Static export suitable for GitHub Pages, Vercel, or Netlify.

## Design

- **Palette:** rich red (`#E63946`), orange (`#FF6B35`), yellow (`#FFD60A`), cream white (`#FFF8F0`)
- **Typography:** Syne (display), DM Sans (body), JetBrains Mono (code)
- **Motion:** Framer Motion with `prefers-reduced-motion` support
- **Inspiration:** Awwwards product storytelling, llama.app hardware grid, Ollama CLI clarity

## Develop

```bash
cd website
npm install
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

## Build static site

```bash
npm run build
```

Output is in `out/` (Next.js static export).

## Deploy

### Vercel

```bash
npx vercel --cwd website
```

Set root directory to `website` in project settings.

### GitHub Pages

```bash
npm run build
# Publish contents of website/out to gh-pages branch or GitHub Pages artifact
```

For project sites at `username.github.io/TinyBrain-OS`, set `basePath` in `next.config.ts` if needed.

## Structure

| Path | Purpose |
|------|---------|
| `app/page.tsx` | Landing — vision, architecture, lifecycle, research |
| `app/docs/page.tsx` | Documentation hub linking to repo markdown |
| `components/` | Hero orbit illustration, sections, header/footer |
| `lib/content.ts` | ADRs, research links, features, release timeline |

Documentation links point to GitHub blob URLs so the site stays synchronized with the repository without duplicating content.
