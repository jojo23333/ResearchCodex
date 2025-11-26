# PLAN Mode

PLAN Mode converts a well-formed idea into a **concrete, modular implementation plan**.
Your output target is:

`PLANS_MD = <current_project>/<current_idea>/plans.md`

Your role in PLAN Mode:
→ Translate the finalized idea (IDEA_MD) into an **actionable execution plan**  
→ Reuse and extend existing code modules in `src/`  
→ Incorporate any human-provided context under ` <current_project>/<current_idea>/context/`  

## Responsibilities

1. **Read all available sources**
   - `IDEA_MD` – the finalized idea specification  
   - `PLANS_MD` – existing plans (if any)  
   - `src/` – previously implemented modules, functions, classes  
   - ` <current_project>/<current_idea>/context/`   – human-provided implementation references or constraints  

2. **Identify Reusable Code**
   - Inspect the `src/` tree to determine:
     - which modules already implement needed features
     - which modules can be extended
     - which parts must be newly created  
   - Prefer **reuse over rewrite**.  
   - When extending code, specify *where* and *how* the changes will occur.

3. **Produce a Modular Implementation Plan**
   - Break work into **small, atomic tasks** (each implementable in Code Mode).  
   - For each major component:
     - whether it reuses an existing module  
     - whether it extends a module  
     - whether it introduces a new module  
   - Explicitly reference file paths and functions (e.g., `src/processor.py:FeatureExtractor`).

4. **Use Context for Guidance**
   - Incorporate patterns, conventions, and examples from `context/` if available.  
   - Respect all human-specified constraints in the context documents.

5. **Write/Update PLANS_MD**
   - PLANS_MD must contain:
     - A short summary of the implementation strategy  
     - A hierarchical task list (`- [ ] Task`)  
     - Notes on module reuse or extension  
     - Any open questions that require human clarification  

6. **Plan for Future Execution**
   - Write tasks so that Code Mode can directly follow them.  
   - Tasks should be:
     - unambiguous  
     - sequential  
     - small enough to complete in one coding step  

7. **Do not generate code in PLAN Mode.**
   - Only produce *plans*, never code.
