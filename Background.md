# ResearchCodex Background

## What is ResearchCodex?

ResearchCodex (RCodex) is a small CLI + library that turns a coding agent (like Codecs/Codex) into a **research assistant** instead of a generic coding bot.

The core idea:

- Treat the **human** as a PI/supervisor.
- Treat **Codex** as a capable undergrad / RA who:
  - Helps plan experiments.
  - Writes and runs code.
  - Logs what was tried, what worked, and what failed.

Instead of focusing on code quality, lint, or bug triage (like Beads), ResearchCodex focuses on:

- **Tracking research ideas** (what we’re trying to test).
- **Tracking experiments** (how we test it).
- **Tracking outcomes** (results and interpretations).
- Creating a **lightweight memory** of the project that’s queryable by both humans and agents.

## Design philosophy

1. **Plan first, implement second**

   - New ideas are first written in human-readable markdown:
     - `ideas.md` — problem statements, hypotheses, approaches.
     - `todos.md` — concrete tasks and experiment plans.
   - A human supervisor reviews and adjusts these.
   - Once the direction is clear, ideas and tasks are **committed into a structured store** via RCodex commands.

2. **Structured store as the source of truth**

   - Markdown is for humans.
   - A **SQLite database** in `.rcodex/` is the source of truth for:
     - Threads (research lines).
     - Ideas.
     - Todos (tasks).
     - Experiments.
   - The agent (Codex) should **read/write the DB only through the `rcodex` CLI**, never by editing the DB directly.

3. **Codex behaves like a capable undergrad**

   - Codex:
     - Reads `Agents.md` for instructions.
     - Uses RCodex CLI commands to:
       - Create new ideas/tasks/experiments.
       - Update statuses.
       - Fetch current work items.
   - Human:
     - Sets direction.
     - Reviews and edits `ideas.md`, `todos.md`, and `results.md`.
     - Approves or rejects bigger plan changes.


