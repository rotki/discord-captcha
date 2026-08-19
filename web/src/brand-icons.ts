/**
 * Brand logos the app registers itself.
 *
 * `@rotki/ui-library` went brand-free when lucide v1 dropped all brand/logo
 * icons: brand marks are third-party trademarks and belong in the consuming
 * app. Anything listed here must also be allow-listed via `customIcons` in
 * vite.config.ts, or the icon plugin flags it as unknown.
 *
 * The data is byte-identical to what the library rendered before it went
 * brand-free, so there is no visual change. `lu-github` is lucide 0.577's
 * outline glyph; `lu-x-twitter` is the mark rotki already shipped.
 */
interface BrandIcon {
  name: string;
  components: [string, Record<string, string>][];
}

export const github: BrandIcon = {
  name: 'lu-github',
  components: [
    ['path', { d: 'M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4' }],
    ['path', { d: 'M9 18c-4.51 2-5-2-7-2' }],
  ],
};

export const xTwitter: BrandIcon = {
  name: 'lu-x-twitter',
  components: [
    ['path', { 'd': 'M3 21L10.5484 13.4516M10.5484 13.4516L3 3H8L13.4516 10.5484M10.5484 13.4516L16 21H21L13.4516 10.5484M21 3L13.4516 10.5484', 'stroke': 'currentColor', 'stroke-width': '2', 'stroke-linecap': 'round', 'stroke-linejoin': 'round' }],
  ],
};

export const brandIcons = [github, xTwitter];
