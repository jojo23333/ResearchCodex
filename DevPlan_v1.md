
# ResearchCodex Development Plan (v0)

## Implementation language

- **Language**: Go
- **CLI framework**: Cobra (or a minimal custom flag parser if easier).
- **Storage**: Filesystem + Markdown (no database in v0).


## Directory layout created by `rcodex init`

Running:

```bash
rcodex init
```

in an empty directory should create:

```text
.
├── .rcodex/
│   ├── config.yaml           # current project / idea path, etc.
│   ├── idea_deps.jsonl       # per-line JSON dependency records (created empty)
│   ├── .plan_mode_agents.md  # default agents used in PLAN mode
│   └── .code_mode_agents.md  # default agents used in CODE mode
├── projects/
│   └── (empty for now)
├── experiments/
│   └── (empty for now)
├── srcs/
│   └── (placeholder for shared tools/utilities)
└── AGENTS.md
```

### `config.yaml` structure (v0)

Simple YAML (paths are relative to repo root):

```yaml
current_project: null
current_idea: null
mode: null
```

- `current_project`: string (project name) or null.
- `current_idea`: string path to the current idea folder (e.g., "projects/topic-modeling/20251115_103012_better_contrastive_loss") or null.
- `mode`: either `plan` or `code` indicating the active focus; defaults to `plan` when a project is created or an idea is created.

---

## CLI commands (v0)

### 1. `rcodex init`
**Purpose:** Initialize a new RCodex project in the current directory.

**Behavior:**

- If `.rcodex/` exists, exit with a friendly error.
- Create:
  - `.rcodex/config.yaml`
  - `.rcodex/idea_deps.jsonl` (empty file)
  - `.rcodex/.plan_mode_agents.md` (initialized with default PLAN mode agent notes below)
  - `.rcodex/.code_mode_agents.md` (initialized with default CODE mode agent notes below)
  - `projects/`, `experiments/`, `srcs/`.
  - `projects/AGENTS.md` (should be initialized as the content in `rcodex/.default_project_agents.md`)
- Initialize `config.yaml` with all `null`.
**Example:**

```bash
rcodex init
```

#### Default mode agent files

These files are created under `.rcodex/` with brief starter content you can edit anytime.

- `.plan_mode_agents.md`

```markdown
# Plan Mode Agents

- Planner: clarify problem statements, assumptions, success criteria
- Decomposer: break goals into steps, experiments, and checkpoints
- Researcher: collect references, related work, and notes

Usage: referenced by `rcodex status` when in PLAN mode.
```

- `.code_mode_agents.md`

```markdown
# Code Mode Agents

- Implementer: write code tied to the current idea chain
- Tester: create minimal test harness and run smoke tests
- Reviewer: review diffs, enforce style, and validate run instructions

Usage: referenced by `rcodex status` when in CODE mode.
```

---
### 2. `rcodex project create <name>`

**Purpose:** Create a new research project and its markdown scaffolding.

**Behavior:**

- Create directories and files:

  ```text
  projects/<name>/
    AGENTS.md           # project-specific agents
    context/             # project context files (empty for now)
  experiments/<name>/    # project-specific experiments folder
  ```

- Create a base idea for the project:

  - Title: `base`
  - Slug: `YYYYMMDD_HHMMSS_base`
  - Path: `projects/<name>/<slug>/`
  - Files:
    - `idea.md` (title: base, created_at, body empty)
    - `plans.md` (empty)
  - Append a record to `.rcodex/idea_deps.jsonl` with `depends_on: null` (base point)

- Update `.rcodex/config.yaml`:
  - `current_project: <name>`
  - `current_idea: "projects/<name>/<slug>"`
  - `mode: plan`

**Example:**

```bash
rcodex project create topic-modeling
```

---

### 3. `rcodex project switch <name>`

**Purpose:** Set the current active project in `config.yaml`.

**Behavior:**

- Verify `projects/<name>` exists; if not, error.
- If found, update config:

  ```yaml
  current_project: "<name>"
  current_idea: null
  ```

**Example:**

```bash
rcodex project switch topic-modeling
```

---

### 4. `rcodex idea create "title" --body "full body text" [--depends-on <slug|path|none>]`

**Purpose:** Create a new idea in the current project.

**Behavior:**

- Read `current_project` from `config.yaml`.
  - If null → error and ask user to run `rcodex project create` or `rcodex project switch`.
- Create a new idea folder under the project using a timestamped slug:
  - Slug format: `YYYYMMDD_HHMMSS_<slugified-title>` (lowercase, spaces to `_`, keep alphanumerics and `_`)
  - Path: `projects/<project_name>/<slug>/`
- Inside the idea folder, create:
  - `idea.md` containing title, created_at, and body
  - `plans.md` (empty scaffold)
- Determine dependency (within the same project):
  - If `--depends-on` is provided:
    - `none` → no dependency (explicit base point)
    - otherwise resolve to an existing idea by slug or path under `projects/<project_name>/`
  - If `--depends-on` is omitted:
    - default to the most recently created idea in this project (lexicographic max of `YYYYMMDD_HHMMSS_*`), if any
    - if none exists, this is the base point idea
- Append one JSON object to `.rcodex/idea_deps.jsonl` recording the relationship (idea 2 -> idea 1).
- Update `config.yaml`:
  - `current_idea: "projects/<project_name>/<slug>"`
  - `mode: plan`

**Examples:**

```bash
rcodex idea create "Better contrastive loss" --body "Try temperature annealing on video-text pairs."
```

Optional (nice-to-have but still within v0 scope): if `--body` is omitted, read from stdin.

---

### 5. `rcodex idea status`

**Purpose:** Fetch and display the current idea and its dependency chain.

**Behavior:**

- Read `current_idea` (path) from `config.yaml`.
- If null → error suggesting to create or select an idea.
- Read and print `idea.md` from that folder in a simple human-readable summary. Optionally show the folder path.
- Resolve and print the full dependency chain for this idea up to the base point by walking `.rcodex/idea_deps.jsonl` (project-local predecessors), oldest→newest.

  ```text
  Idea: Better contrastive loss
  Project: topic-modeling
  Created: 2025-11-16T10:30:12Z
  Path: projects/topic-modeling/20251116_103012_better_contrastive_loss
  Dependency chain:
    projects/topic-modeling/20251115_083000_base
    -> projects/topic-modeling/20251115_153045_initial_baseline
    -> projects/topic-modeling/20251116_103012_better_contrastive_loss

  Body:
    Try temperature annealing on video-text pairs.
  ```

---


### 6. `rcodex status`

**Purpose:** Get current general status

**Behaviour:**

- Read `.rcodex/config.yaml`:

  ```yaml
  current_project: <string or null>
  current_idea: <string or null>
  mode: plan|code
  ```

- If `current_project` is `null`, respond:
  "No project setted. please response by 'Please setup project by running rcodex project create <project_name>'"

- Otherwise, show which project and idea (by folder name) you're on, and the current `mode`.

- If mode is `plan`, show:

  """
  You are currently in PLAN mode
  - For PLAN mode, you will work in `projects/<current_project>/<current_idea>`
  <.plan_mode_agents.md>
  """
  <.plan_mode_agents.md> is the context in .rcodex/.plan_mode_agents.md
  use actual project name and idea name to replace <current_project> and <current_idea>


- If mode is `code`, show:

  """
  You are currently in CODE mode 
  you will reference `projects/<current_project>/<current_idea>` work on 
    1. Shared code folder `srcs/` 
    2. Idea-specific scripts go in `experiments/<current_project>/`.
  <.code_mode_agents.md>
  """
  <.code_mode_agents.md> is the context in .rcodex/.code_mode_agents.md
  use actual project name and idea name to replace <current_project> and <current_idea>


**Example**
```bash
rcodex status
```

---

## Data model and persistence (v0)

- No database is used. All data lives in the filesystem as folders and Markdown files.
- `config.yaml` tracks the current project and idea path for convenience.
  - Also tracks `mode` which toggles between planning and coding focus.
- Projects live under `projects/<name>/` with a `context/` folder and `AGENTS.md`.
- Ideas live under `projects/<name>/<YYYYMMDD_HHMMSS_slug>/` with `idea.md` and `plans.md`.
- Dependency between ideas lies in `.rcodex/idea_deps.jsonl`
- Experiments are organized under `experiments/<project_name>/` (created on project creation) and can contain scripts, notebooks, and outputs as you see fit.

### Dependencies between ideas

- Stored append-only in `.rcodex/idea_deps.jsonl` as one JSON object per line.
- Fields per record:
  - `project`: project name (string)
  - `idea_path`: relative path to the idea folder (string)
  - `depends_on`: relative path to the prerequisite idea folder, or `null` if base point
  - `created_at`: ISO-8601 timestamp (string)
  - `note` (optional): free-form text
- On first idea in a project: write `depends_on: null` (base point).
- If `--depends-on` is omitted when creating a new idea: set dependency to the latest existing idea in the same project (by slug timestamp); otherwise use the resolved idea indicated by `--depends-on`.
- `rcodex idea status` should reconstruct and print the chain from base → ... → current using these records.

### 7. `rcodex plan`

**Purpose:** Switch the active working mode to planning.

**Behavior:**
- Verify a `current_project` is set; if not, error.
- Set `mode: plan` in `config.yaml`.
- Optionally print a short confirmation.

**Example:**

```bash
rcodex plan
```

### 8. `rcodex code`

**Purpose:** Switch the active working mode to coding.

**Behavior:**
- Verify a `current_project` is set; if not, error.
- Set `mode: code` in `config.yaml`.
- Optionally print a short confirmation.

**Example:**

```bash
rcodex code
```


