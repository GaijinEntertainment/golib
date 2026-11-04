---
paths:
  - "**/CLAUDE.md"
  - "**/.claude/rules/**/*.md"
---

<agent-instructions>

- Delimit instruction blocks with XML tags, each tag on its own line. Write the rules inside a block as bullets, with a
  blank line between a tag and its bullets.
- After an edit, format the file with Prettier through a launcher such as `npx`; the repository has no Node.js
  toolchain:
  `--write --use-tabs --print-width 120 --parser markdown --prose-wrap always --bracket-same-line --html-whitespace-sensitivity strict <file-path>`.

</agent-instructions>
