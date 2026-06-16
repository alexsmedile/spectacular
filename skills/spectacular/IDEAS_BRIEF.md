# IDEAS_BRIEF — reusable design-brief contract

A copy-paste **creative-director prompt** for handing a product's interface to an
AI design studio (v0, Figma AI/Make, Lovable, Midjourney for hero shots) or a
human team. The output is a single self-contained prompt — not wireframes, not a
spec. This file is the **contract**: fill every `{{slot}}`, delete any section
that genuinely doesn't apply, and ship the result as plain text.

---

## How to use

1. Copy the **Template** block below.
2. Replace every `{{double-brace slot}}` with project specifics. Slots with a
   trailing `?` are optional — drop the line if empty.
3. Keep the section headings (`ROLE`, `THE PRODUCT`, …) — they are the contract's
   load-bearing structure. A studio reads them in order.
4. Strip any guidance in `‹angle-quote comments›` — those are notes to you, not
   the designer.
5. Deliver as plain text. No markdown fences in the final handoff unless the
   target tool wants them.

## Filling rules (what makes a brief good vs. generic)

- **PRODUCT SOUL** is the most important section. Three feelings, max. Each must
  be falsifiable ("must feel instant and calm", not "must be modern"). If a
  design choice can't be argued against one of these, the soul is too vague.
- **Anti-requirements** earn their keep — say what it must *not* look like. "Never
  corporate-dark by default" steers harder than ten adjectives.
- **CONSTRAINTS** are where briefs save the studio's time. State the host platform
  (native macOS / iOS / web / Electron), what's achievable in it, and what content
  the shell does *not* own (embedded third-party UI, web views, ad slots).
- **Every state, not just the happy path.** Loading, empty/first-run, error,
  offline, and any signature behavior (e.g. a sleep/wake cycle) get explicit
  screens. The happy path is the easy 20%.
- Name real **brand values** — hex codes, corner radius as a %, the actual
  typeface — not "use brand colors."
- One **GOAL** sentence at the end: the single first-glance impression. If the
  finished design nails this, it's done.

---

## Template

```
ROLE
You are a senior product designer + creative director specializing in
{{platform, e.g. native macOS apps}}. Design the complete interface for a
{{platform}} app/product called "{{Product Name}}". Deliver an opinionated,
polished, cohesive UI — not wireframes.

THE PRODUCT
{{2–4 sentences: what it is, who it's for, the one job it does better than the
incumbent. Name the incumbent and the pain it causes — that contrast is the pitch.
End with the core technical shape if it affects the UI (e.g. "one web view per
service, each sandboxed").}}

PRODUCT SOUL — let every pixel reinforce these
‹Exactly 3 feelings. Each falsifiable. Number them.›
1. {{Feeling #1}}. {{One line on what it implies for the UI.}}
2. {{Feeling #2}}. {{…}}
3. {{Feeling #3}}. {{…}}

BRAND / VISUAL LANGUAGE
- Identity/mascot?: {{the signature mark or character, and how sparingly to use it}}
- Palette: {{named colors with hex — e.g. lime #BCE66B → green #619129}}. Pair with
  {{surface treatment — near-white, frosted/vibrancy, dark, etc.}}.
- Shape language: {{corner radius as % of width, roundness, shadow style}}.
- Typography: {{typeface — favor the platform system font unless there's a reason}}.
- Tone: {{3–5 adjectives}} — never {{anti-tone, e.g. corporate or techy-dark}}.
  Must respect {{platform}} conventions and feel native, just more delightful.

SCREENS & SURFACES TO DESIGN (full {{scope: app / flow / single screen}})
‹One numbered block per surface. Every product has more states than screens —
list the states explicitly. Delete blocks that don't apply.›
1. {{Primary screen}} — {{layout, key regions, the main interaction}}
2. {{Navigation / chrome}} — {{rail, tabs, sidebar, window controls}}
3. {{Signature-behavior states}} — {{the thing that makes this product unique;
   show every state it passes through}}
4. {{Create / add / onboarding flow}} — {{2–3 steps max, light and inviting}}
5. {{Settings / per-item config}} — {{the proper surface, not a hidden menu}}
6. {{Empty / first-run state}} — {{what a brand-new user sees; teach the core idea
   in one glance; let the brand shine here}}
7. {{Loading / error / offline}} — {{calm, never alarming}}

DELIVERABLES
- A clear visual direction: typography, color usage, materials, spacing system,
  corner radii, shadow/elevation, and motion principles (how the signature
  transitions feel).
- High-fidelity layouts for each surface above, in BOTH light and dark mode.
- A small component set: {{list the 4–6 reusable components and their states}}.
- {{Any product-specific deliverable, e.g. an active/selected-state treatment}}.
- Annotations on anything non-obvious — especially {{the signature behavior}}.

CONSTRAINTS
- {{Host platform}} — design within what's achievable natively
  ({{platform primitives — e.g. NSVisualEffectView vibrancy, SF Symbols, system
  fonts}}). No {{out-of-platform}} flourishes that wouldn't feel native.
- Respect {{platform}} HIG: {{window chrome, light/dark, accessibility contrast}}.
- {{What the shell does NOT own}} — design only {{Product}}'s UI AROUND it, never
  restyle {{the embedded/third-party content}}.
- {{Anti-requirement — if a choice adds weight/complexity without value, cut it.}}

GOAL
{{One sentence: the single first-glance impression the finished design must
deliver — while feeling 100% at home on {{platform}}.}}
```

---

## Worked example (the brief this contract was extracted from)

> Filled for **Wasabi** — a featherweight native-macOS multi-service chat browser
> shell. Use as a reference for how dense/specific each slot should get; do not
> copy verbatim.

```
ROLE
You are a senior product designer + creative director specializing in native
macOS apps. Design the complete interface for a shipping macOS app called
"Wasabi." Deliver an opinionated, polished, cohesive UI — not wireframes.

THE PRODUCT
Wasabi is a featherweight native macOS browser shell that hosts web messaging
apps (WhatsApp Web, Telegram, more later) behind a single vertical icon rail.
It's the calm, low-RAM alternative to bloated Electron chat clients (which eat
2–5 GB each). Wasabi keeps the active service live and quietly sleeps idle ones,
targeting under 500 MB. Native: Swift + AppKit + WebKit, one sandboxed,
persistent WKWebView per service.

PRODUCT SOUL — let every pixel reinforce these
1. Featherweight. Must FEEL light, instant, calm. Generous space. Quiet by default.
2. Many apps, one window. A frictionless home for all your chat web-apps.
3. Private & isolated. Each service is its own encrypted sandbox; logins persist.
   Convey trust and tidiness without shouting about security.

BRAND / VISUAL LANGUAGE
- Mascot: a kawaii white wasabi-paste dollop with a tiny smiling face — playful,
  clean. Use sparingly as a signature, not wallpaper.
- Palette: wasabi-greens on a mesh gradient — lime #BCE66B → mid #85B545 →
  wasabi #619129 → teal #4DBC8C. Pair with near-white surfaces and frosted
  macOS vibrancy.
- Shape language: macOS squircle, corner radius ~22.5% of width, soft shadows.
- Typography: SF Pro / system.
- Tone: friendly, minimal, fast, slightly kawaii — never corporate or techy-dark.
  Must respect macOS conventions (traffic lights, vibrancy, light/dark).

SCREENS & SURFACES TO DESIGN (full app UI)
1. Main window — slim (~48pt) frosted vertical icon rail (real site favicons,
   rounded) + a "+" add tile at the bottom; main content area right of a hairline
   divider holds the active service's web view.
2. Active vs inactive rail states — distinctive but restrained active indicator;
   quiet unread badges/dots.
3. Sleeping / waking states — the signature behavior: idle services are torn down
   to save RAM. Show a calm SLEEPING rail state and a graceful WAKING/loading
   state on reselect. Lean into "featherweight" — intentional, not broken.
4. Add-service flow — from "+": URL or pick a suggestion, 2–3 steps, inviting.
5. Per-service settings — sleep policy (Keep running / Sleep after 5 min / Smart
   sleep) as a proper popover, not a right-click menu.
6. Empty / first-run state — before any service is added; mascot + brand shine,
   teach the core idea at a glance.

DELIVERABLES
- Visual direction: type, color, materials (vibrancy/frost), spacing, radii,
  elevation, motion (sleep/wake, selection, add transitions).
- Hi-fi layouts per surface, light AND dark mode.
- Components: rail tile (active/inactive/unread/sleeping), add tile, sleep-policy
  control, badges, divider, empty-state hero.
- A distinctive-yet-restrained active + unread treatment.
- Annotations, especially on the sleep/wake states.

CONSTRAINTS
- Native macOS (AppKit/WebKit) — design within native primitives
  (NSVisualEffectView vibrancy, standard window chrome, SF Symbols, system fonts).
  No web-only flourishes.
- Respect macOS HIG: traffic lights, menu bar, light/dark, accessibility contrast.
- The hosted web content (WhatsApp/Telegram) is owned by those sites — design only
  Wasabi's shell AROUND it, never restyle the web content.
- It must look as light as it runs. If a choice adds visual weight without value,
  cut it.

GOAL
A distinctive, delightful, unmistakably-Wasabi interface that makes "a calm,
featherweight home for all my chat apps" obvious and desirable at first glance —
while feeling 100% at home on macOS.
```
