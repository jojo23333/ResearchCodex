# SCOPE Mode

Scope mode refines rough concepts into structured, reviewable scopes before planning and coding.

You will mainly edit:
- `IDEA_MD=<current_project>/<current_scope>/idea.md`
- `PLANS_MD=<current_project>/<current_scope>/plans.md`

## Stage 1 – Scope Formation

**Goal:** turn a rough scope into a detailed, reviewable **AI-elaborated scope**.

`idea.md` should contain two parts:
1) Original human description (organized lightly for readability)  
2) AI-Elaborated Scope

```markdown
## Original Scope
...
## AI-Elaborated Scope
...
```

Steps:
1. Ensure `IDEA_MD` exists. Incorporate any new human content into **Original Scope**.
2. Read helpful context in `projects/<current_project>/context/` (optional).
3. Itemize datasets, models, hyperparams, evaluation plan, and risks. When choices need human input, add a block under AI-Elaborated Scope:
```markdown
***********************
> - <Clarification Question>
> [ ] Human Response: TODO
************************
```
4. When the human fills in `Human Response:` or checks boxes, absorb those decisions into the AI-Elaborated Scope, replacing placeholders.
5. Keep AI-Elaborated Scope structured and coherent.

**Completion rule**
- If no unresolved `> [ ]` blocks remain (or human says to proceed), move to PLAN mode.
- Otherwise, stop and ask the human to review/fill `idea.md`.

## Stage 2 – Plan Formation (`plans.md`)

**Goal:** convert `IDEA_MD` into a concrete **implementation plan** in `plans.md`.

Steps:
1. Read `IDEA_MD` and existing `PLANS_MD`.
2. Update `PLANS_MD` as concise, actionable tasks.
3. Keep tasks small and testable; you will mark them done/modified later.

