## CODE Mode

CODE mode turns plans into real code and documentation.

### 1 Preconditions

Only enter CODE mode when:

- `<current_project>/<current_idea>/plans.md` exist and are reasonably complete.
- Human instruction has signaled “implement” 

If not, go back to PLAN mode using following command and return a hint message
`rcodex plan`

### 2. Working Directories

- **Shared utilities** → `srcs/` (codes that are likely to be reused across ideas and projects, e.g. dataset loader, utils, common algorithms/models, general visualization tools, helper functions..).
- **Idea-specific code** → `experiments/<current_project>/` (entry script, idea-specific code)

### 3 Implementation

For each implementation session:

1. Read and internerlize `<current_project>/<current_idea>/plans.md`
2. Implement the task:
   - Create/modify files under `srcs/` and `experiments/<current_project>/`.
   - Use clear functions, docstrings, and simple entrypoints.
3. Add **basic tests** for new behavior:
   - Place them under `experiments/<current_project>/tests`.
   - Use small synthetic data or small dataset subsets.
4. If possible, run tests ande debut yourself.
   - If you cannot run them, document how to run them.
5. Update `PLANS_MD`:
   - Mark tasks as done or partially done.
   - Note any deviations:
     - `Implementation note: originally planned X, implemented Y because Z.`

### 4 Implementation doc (`implementation.md`)

Create a detailed Implelemtation doc under `<current_project>/<current_idea>/implementation.md`

Suggested structure:

```markdown
# Implementation Doc: 

Project: {{current_project}}
Idea path: {{current_idea}}
Last AI update: {{NOW_TBD}}

## Implementation Doc
A detailed description of what has been implemented

## Code Changes
Which piece of code has been changed for what

## Test & Usage
Test commands, How to use the code.. 

## Known Issues / TODO
- ...
```

This document is the **primary factual record** for future debugging and reuse.
