# SSSsssHhHingg. — SSH Terminal Portfolio

A minimal portfolio you visit over SSH. It's a Go program built on the
[Charmbracelet](https://charm.sh) stack — Bubble Tea, Bubbles, Lip Gloss,
Wish, Harmonica, and Log.

```
ssh localhost -p 2222
```

## What it does

- A **preloader** types out `Pr3thiv` with a blinking red cursor.
- Then a **tabbed UI**: About, Projects, Blogs, Links, Playground.
- **About** shows your bio next to a randomly-picked ASCII art, with a
  "last updated" line at the bottom.
- **Blogs** is a selectable list; open a post to read it in a scrollable
  in-terminal reader.
- **Links** are clickable (OSC 8) and also shown as raw URLs so any terminal
  can Ctrl/Cmd+click them.
- **Playground** is a colorful particle sandbox driven by real physics
  (Harmonica projectiles with gravity + bouncing).
- Every visitor connection is **logged** (time, and terminal) to `visits.log`.

## Controls (vim-style)

| Key           | Action                    |
| ------------- | ------------------------- |
| `h` / `l`     | switch tabs               |
| `j` / `k`     | scroll / move selection   |
| `g` / `G`     | jump to top / bottom      |
| `enter`       | open a blog post          |
| `q` / `esc`   | back (in reader) / quit   |
| `space` `f` `r` | Playground: launch / drop / reset |

Arrow keys work too.

## How it works

```
main.go                 Wish SSH server + visitor logging
internal/config         loads content.yaml + ASCII art
internal/ui/root.go     switches preloader -> portfolio
internal/ui/preloader.go  the Pr3thiv splash
internal/ui/portfolio.go  tabs, blog reader, links
internal/ui/playground.go Harmonica physics sandbox
internal/ui/styles.go   the theme (colors, spacing)
```

Each SSH session gets its own Bubble Tea program wired to the connection's
terminal. Wish handles the SSH side; Bubble Tea handles the UI loop
(state -> update -> view). Content is re-read from `content.yaml` on every
connection, so edits show up without a rebuild.

## Customizing (edit `content.yaml`)

Everything you see is driven by `content.yaml` — no code changes needed.

```yaml
name: "Pr3thiv"
tagline: "researcher · developer · terminal enjoyer"
last_updated: "2026-07-05"
about: |
  Your bio here. Blank lines become paragraphs.

projects:
  - name: "Project name"
    description: "What it does."
    url: "https://github.com/you/project"

blogs:
  - title: "Post title"
    date: "2026-05-01"
    url: "https://example.com/post"
    body: |
      Optional inline text shown in the reader.

links:
  - label: "GitHub"
    url: "https://github.com/you"
```

**ASCII art:** drop `.txt` files into `assets/ascii/`. One is picked at
random each session and shown on the About page.

## Run

```powershell
# Windows PowerShell
$env:GOFLAGS="-mod=mod"; go run .
```

```bash
# macOS / Linux
go run .
```

A host key is generated automatically at `.ssh/portfolio_ed25519` on first
run. Then connect from another terminal:

```
ssh localhost -p 2222
```

### Config via environment variables

| Variable        | Default                    | Purpose              |
| --------------- | -------------------------- | -------------------- |
| `PORT_ADDR`     | `:2222`                    | listen address       |
| `HOST_KEY_PATH` | `.ssh/portfolio_ed25519`   | SSH host key         |
| `CONTENT_PATH`  | `content.yaml`             | content file         |
| `VISITS_LOG`    | `visits.log`               | visitor log file     |

## Built with

Bubble Tea · Bubbles · Lip Gloss · Wish · Harmonica · Log
