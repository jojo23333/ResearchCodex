# ResearchCodex

ResearchCodex (`rcodex`) is a Go-based CLI that scaffolds lightweight research workflows for AI coding agents. See `Background.md` and `DevPlan_v1.md` for philosophy and behavior.

## Quick Install / Update (Linux/macOS, no root)

Run the installer script directly from GitHub. It automatically:
1. Caches the repo at `~/.cache/rcodex-src` (reused for updates).
2. Builds the CLI with Go.
3. Installs/updates the binary in `~/.local/bin`.

```bash
curl -fsSL https://raw.githubusercontent.com/jojo/ResearchCodex/main/install.sh | bash
```

Re-running the same command later updates to the latest `main`. No sudo required, but Go 1.21+ must be available.

### Customizing the install

Use environment variables when invoking the script (works for curl piping or local execution):

| Variable | Default | Description |
| --- | --- | --- |
| `RCODEX_INSTALL_DIR` | `~/.local/bin` | Destination directory for the `rcodex` binary |
| `RCODEX_CACHE_DIR` | `~/.cache/rcodex-src` | Location of the cached clone used for incremental updates |
| `RCODEX_REF` | `main` | Git ref or branch to install |
| `RCODEX_BINARY_NAME` | `rcodex` | Output binary name |

Example:

```bash
curl -fsSL https://raw.githubusercontent.com/jojo/ResearchCodex/main/install.sh \
  | RCODEX_INSTALL_DIR="$HOME/bin" bash
```

If the install directory is not on your `PATH`, add it to your shell profile:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

## Basic Usage

```bash
# Initialize a workspace (run once per repository)
rcodex init

# Create your first project; this sets the active project/scope
rcodex project new my-research-topic

# Switch phases as you refine the work
rcodex scope    # scoping formation (idea.md structure)
rcodex plan     # implementation planning (plans.md)
rcodex code     # implementation/execution

# Add more scopes as you refine the research thread
rcodex idea new "Better contrastive loss" --body "Try temperature annealing on video-text pairs."

# Inspect context and mode
rcodex status
```

See `DevPlan_v1.md` for the complete v1 command set and expectations.
