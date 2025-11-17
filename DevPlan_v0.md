
# ResearchCodex Development Plan (v0)

## Implementation language

- **Language**: Go
- **CLI framework**: Cobra (or a minimal custom flag parser if easier).
- **DB**: SQLite (via standard Go `database/sql` + a SQLite driver).

## Project layout (Go)

Minimal layout (v0):

```text
rcodex/
  cmd/
    rcodex/
      main.go           # CLI entrypoint
  internal/
    project/
      init.go           # rcodex init
      thread.go         # thread create/switch/list helpers
    store/
      store.go          # DB connect/open, migrations
      models.go         # structs + basic CRUD
  go.mod
  go.sum
```

We keep things as flat and straightforward as possible. No generics gymnastics or heavy abstractions.

---

## Directory layout created by `rcodex init`

Running:

```bash
rcodex init
```

in an empty directory should create:

```text
.
├── .rcodex/
│   ├── config.yaml         # current thread / idea IDs, etc.
│   └── rcodex.db           # SQLite database
├── threads/
│   └── (empty for now)
├── experiments/
│   └── (empty for now)
├── .exps/
│   └── (empty for now; future result files can go here)
├── artifacts/
│   └── (placeholder for papers, reference code, etc.)
├── srcs/
│   └── (placeholder for shared tools/utilities)
├── Topics.md
└── Agents.md
```

### `config.yaml` structure (v0)

Simple YAML:

```yaml
current_thread: null
current_idea: null
```

- `current_thread`: string (thread name) or null.
- `current_idea`: integer idea ID or null.

---

## SQLite schema (v0)

Create a single DB file: `.rcodex/rcodex.db`.

Minimal tables:

### 1. `threads`

Represents a line of research.

- `id` INTEGER PRIMARY KEY AUTOINCREMENT
- `name` TEXT UNIQUE NOT NULL
- `title` TEXT
- `created_at` TEXT (ISO timestamp)
- `status` TEXT DEFAULT 'active'  -- e.g., 'active' | 'archived'

### 2. `ideas`

Represents a research idea / hypothesis within a thread.

- `id` INTEGER PRIMARY KEY AUTOINCREMENT
- `thread_id` INTEGER NOT NULL
- `title` TEXT NOT NULL
- `body` TEXT NOT NULL           -- usually derived from `ideas.md`
- `status` TEXT NOT NULL DEFAULT 'planned'
  -- 'planned' | 'in_progress' | 'done' 
- `prev_idea_id` INTEGER DEFAULT NULL
  -- simple chain: previous idea in this thread
- `created_at` TEXT

### 3. `todos`

Represents work items / tasks related to an idea.

- `id` INTEGER PRIMARY KEY AUTOINCREMENT
- `idea_id` INTEGER NOT NULL
- `text` TEXT NOT NULL           -- short todo text
- `doc` TEXT DEFAULT ''          -- implementation summary / notes
- `status` TEXT NOT NULL DEFAULT 'open'
  -- 'open' | 'in_progress' | 'done'
- `created_at` TEXT

### 4. `experiments`

Represents a planned experiment that tests an idea.

- `id` INTEGER PRIMARY KEY AUTOINCREMENT
- `idea_id` INTEGER NOT NULL
- `description` TEXT NOT NULL    -- free-form description
- `status` TEXT NOT NULL DEFAULT 'planned'
  -- 'planned' | 'running' | 'done' 
- `code_path` TEXT DEFAULT ''    -- e.g. "experiments/thread_name/exp_3.py"
- `created_at` TEXT

For v0 we **do not** implement a separate `results` table. Results can be stored as files under `.exps` and summarized in markdown (`results.md`) by the human and/or Codex.

---

## CLI commands (v0)

### 1. `rcodex init`

**Purpose:** Initialize a new RCodex project in the current directory.

**Behavior:**

- If `.rcodex/` exists, exit with a friendly error.
- Create:
  - `.rcodex/rcodex.db`
  - `.rcodex/config.yaml`
  - `threads/`, `experiments/`, `.exps/`, `artifacts/`, `srcs/`.
  - `Topics.md` (empty placeholder).
  - `Agents.md` (empty placeholder).
- Initialize SQLite DB with v0 schema.
- Initialize `config.yaml` with all `null`.

**Example:**

```bash
rcodex init
```

---

### 2. `rcodex thread create <name>`

**Purpose:** Create a new research thread and its markdown scaffolding.

**Behavior:**

- Insert a new row into `threads`:
  - `name = <name>`
  - `title = <name>` (for now)
- Create directory:

  ```text
  threads/<name>/
    ideas.md
    todos.md
    results.md
  ```

- Create directory for experiments under this thread:

  ```text
  experiments/<name>/
  ```

- Create directory for experiment outputs for this thread:

  ```text
  .exps/<name>/
  ```

- Update `.rcodex/config.yaml`:
  - `current_thread: <name>`
  - `current_idea: null`

**Example:**

```bash
rcodex thread create topic-modeling
```

---

### 3. `rcodex thread switch <name>`

**Purpose:** Set the current active thread in `config.yaml`.

**Behavior:**

- Look up `<name>` in `threads` table.
- If not found, error.
- If found, update config:

  ```yaml
  current_thread: "<name>"
  current_idea: null
  ```

**Example:**

```bash
rcodex thread switch topic-modeling
```

---

### 4. `rcodex idea create "title" --body "full body text"`

**Purpose:** Create a new idea in the current thread.

**Behavior:**

- Read `current_thread` from `config.yaml`.
  - If null → error and ask user to run `rcodex thread create` or `rcodex thread switch`.
- Find the `thread` row by `name = current_thread`.
- Determine `prev_idea_id`:
  - Find the most recently created `idea` in this thread (if any).
- Insert a new row into `ideas`:

  - `title` = `"title"`
  - `body` = `"full body text"` (or from stdin if easier)
  - `thread_id` = thread.id
  - `prev_idea_id` = latest idea.id or NULL
  - `status` = `'planned'`
  - `created_at` = timestamp

- Append a small section to `threads/<thread_name>/ideas.md` (if present), e.g.:

  ```markdown
  ## Idea <id>: <title>

  <full body text>

  Status: planned
  ```

- Update `config.yaml`:

  - `current_idea: <new idea id>`

**Examples:**

```bash
rcodex idea create "Better contrastive loss" --body "Try temperature annealing on video-text pairs."
```

Optional (nice-to-have but still within v0 scope): if `--body` is omitted, read from stdin.

---

### 5. `rcodex todo create "task text" [--idea <id>]`

**Purpose:** Create a todo / task associated with an idea.

**Behavior:**

- Determine `idea_id`:
  - If `--idea` is given, use that.
  - Otherwise, read `current_idea` from `config.yaml`.
    - If still null → error.
- Insert into `todos`:
  - `text` = `"task text"`
  - `doc` = empty string
  - `status` = `'open'`
- Optionally append this todo to `threads/<thread_name>/todos.md`:

  ```markdown
  - [ ] (Todo <id>) task text
  ```

**Example:**

```bash
rcodex todo create "Implement baseline model for Idea 3" --idea 3
```

---

### 6. `rcodex experiment create "description" [--idea <id>]`

**Purpose:** Create an experiment associated with an idea and scaffold a script file.

**Behavior:**

- Determine `idea_id` (same logic as `todo`).
- Insert into `experiments`:
  - `description` = `"description"`
  - `status` = `'planned'`
  - `code_path` = e.g. `"experiments/<thread_name>/exp_<id>.py"` (or `.ipynb`, but pick one consistently).
- Create a minimal experiment file at `code_path`:

  Example Python stub:

  ```python
  # Experiment <id>
  # Description: <description>

  def main():
      # TODO: implement experiment
      pass

  if __name__ == "__main__":
      main()
  ```

- Optionally append to `threads/<thread_name>/todos.md`:

  ```markdown
  - [ ] (Experiment <id>) Implement and run: description
  ```

**Example:**

```bash
rcodex experiment create "Train baseline ResNet on small subset" --idea 3
```

---

### 7. `rcodex get idea`

**Purpose:** Fetch and display the *current idea*.

**Behavior:**

- Read `current_idea` from `config.yaml`.
- If null → error suggesting to create or select an idea.
- Query `ideas` table for that `id`.
- Print a simple human-readable summary:

  ```text
  Idea 5 (Thread: topic-modeling)
  Status: in_progress
  Title: Better contrastive loss

  Body:
    Try temperature annealing on video-text pairs.
  ```

For v0 we do **not** implement `--json`. (We can add it in later versions.)

---

### 8. `rcodex get todo`

**Purpose:** Fetch and display all todos for the current ideas.

**Behavior:**

- Fetch todos for the current ideas 
- default by fetching todos with status `open` and `in_progress`

**Example:**

```bash
rcodex get todo 
rcodex get todo --all 
```

---

### 9. `rcodex update todo <id> --status <status> [--doc "text"]`

**Purpose:** Update the status and/or doc field of a todo.

**Behavior:**

- Supported statuses: `open`, `in_progress`, `done`.
- If `--doc` is provided:
  - Replace `doc` with this content.
  - (In v0 we don't do append vs replace; future versions can add `--append-doc`.)
- No markdown updates are strictly required in v0 (but it’s okay to optionally mark `[x]` in `todos.md` if trivial to implement).

**Example:**

```bash
rcodex update todo 12 --status done --doc "Implemented baseline, results saved under .exps/topic-modeling/baseline_run1/"
```

---

### 10. `rcodex update idea <id> --status <status>`

**Purpose:** Change an idea’s status.

**Behavior:**

- Supported statuses: `planned`, `in_progress`, `done`.
- For v0, no constraints (e.g., we don’t enforce “all todos must be done” yet).
- Optional: update `threads/<thread_name>/ideas.md` status line if existing.

**Example:**

```bash
rcodex update idea 5 --status in_progress
```

### 11. `rcodex status`

**Purpose:** Get current general status

**Behaviour:**

- Display information on what is the current thread name and idea name
- For the current idea show what you will get with 'rcodex get todo' and 'rcodex get idea'

**Example**
```bash
rcodex status
```

---

## Minimal migration / DB initialization logic

On first open of `.rcodex/rcodex.db`:

- Check if tables exist.
- If not, create them with the schema above.
- For v0, support **only** schema version 1.
- Store schema version in a small meta table:

  - `meta(key TEXT PRIMARY KEY, value TEXT)`
  - Insert row: `('schema_version', '1')`

We do **not** implement schema migrations yet.


