# PLAN Mode

Plan mode turns coarse,under-prepared research ideas to detailed implementation plans.
You will be mainly working on two files:
`IDEA_MD=<current_project>/<current_idea>/idea.md` and `PLANS_MD=<current_project>/<current_idea>/plans.md`

## Stage 1 – Idea Formation 

**Goal:** turn a rough idea into a detailed, reviewable **AI-elaborated idea**.

The IDEA_MD should contain two part
1. The original human description of the idea, lightly organized/rewritten into readable format
2. AI-Elaborated idea 

```markdown
## Original Idea
...
## AI-Elaborated Idea
...
```

Steps:

1. Ensure `IDEA_MD` exists.
   - **If there is new content by user added them into original idea section**
   - **If There is still no content in the original idea section, stop here and ask for human input("Please input a idea or modify `<current_project>/<current_idea>/idea.md`")**
2. Read:
   * `IDEA_MD` (The original Idea section as well as the AI-elaborated idea section)
   * Any useful context under `projects/<current_project>/context/` (optional).
3. Reason carefully and itemize all plausible details (datasets, models, hyperparams, evaluation plan, etc.).
    * When a concrete choice needs human input, insert an explicit block inside ## AI-Elaborated Idea:
```markdown
    ***********************
    > - <Clarification Question>
    > [ ] Human Response: TODO
    ************************
```
4. When the human has filled in `Human Response:` fields and/or updated checkboxes to `[x]`, read those decisions and obsorb them into the main **AI-Elaborated Idea** text (replace vague placeholders with the chosen options)
5. Rewrite and update AI-Elaborated Idea as a structured, coherent document, using context if helpful.

**Completion rule**

* **If IDEA_MD contains no > [ ] Human input needed blocks, or the human explicitly says to proceed with planning, go to the next stage.**
* **Otherwise, stop and ask the human to review and revise idea.md (including filling in Human Response: fields).**


## Stage 2 – Plan Formation (`plans.md`)

**Goal:** convert `IDEA_MD` into a concrete **implementation plan**.

Steps:

1. Read `IDEA_MD` and existing `PLANS_MD` (if any).
2. Write/update `PLANS_MD` as a concise but comprehensibe task lists.
3. Keep tasks **small and actionable**.  
4. You will later mark tasks as done / modified during implementation.
