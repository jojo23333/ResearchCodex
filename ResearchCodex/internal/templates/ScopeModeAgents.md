# SCOPE Mode

Scope mode refines rough idea into structured, reviewable ideas before planning and coding.

You will mainly edit:
- `IDEA_MD=<current_project>/<current_idea>/idea.md`

## Goal
**Turn a rough idea into a detailed, reviewable **AI-elaborated idea**.

`idea.md` should contain two parts:
1) Original human description (organized lightly for readability)  
2) AI-Elaborated idea

```markdown
## Original Idea
...
## AI-Elaborated Idea
...
```

## Steps:
1. Ensure `IDEA_MD` exists. Incorporate any new human content into **Original Idea**.
2. Read helpful context in `projects/<current_project>/context/` (optional).
3. Itemize datasets, models, hyperparams, evaluation plan, and risks. When choices need human input, add a block under AI-Elaborated Idea:
```markdown
***********************
> - <Clarification Question>
> [ ] Human Response: TODO
************************
```
4. When the human fills in `Human Response:` or checks boxes, absorb those decisions into the AI-Elaborated Idea, replacing placeholders.
5. Keep AI-Elaborated Idea structured and coherent.

**Completion rule**
- If no unresolved `> [ ]` blocks remain (or human says to proceed), move to PLAN mode.
- Otherwise, stop and ask the human to review/fill `idea.md`.
