You are a coding agent responsible for implementing the first minimal version (v0) of a CLI tool called ResearchCodex (rcodex).

Context:
- There are two design documents available in this repository:
  - Background.md
  - DevPlan_v0.md
- These documents define:
  - The purpose and philosophy of ResearchCodex.
  - The exact directory structure to create.
  - The v0 SQLite schema.
  - The CLI commands and their expected behavior.
- .beads_reference_code is a issue management repo called beads, it's the source of idea and is implemented with go as well. Before planning, read and comprehend it's general architecture

Your job:
1. Read the Background.md to have a general idea about the project
2. Read and internalize DevPlan_v0.md. 
3. Reference .beads_reference_code folder for the major design of the go module.
4. Implement as specified in DevPlan_v0.md, including commands and databases
5. After implementing each command:
   - Manually test it on a small example project.
   - Make sure the directory layout and DB rows match the specification in DevPlan_v0.md.
   - Fix obvious bugs or panics.
6. After finishing the implemetation, comes up with a testing plan and started with a comprehensive test by your self, you are authorize to install anything you need just mark it down explictly.

Very important:
- Do NOT invent extra features or commands in v0.
- If any behavior is ambiguous or you are uncertain. reference .beads_reference_code for a implementation that still matches the spirit of the document 
- Delve deeper and plan more comprehensively before start
- Follow AGENTS.md's guideline on using the history/ folder.

Start by scaffolding the Go module (go mod init, main.go, etc.), then implement rcodex init, setting up data storing, and proceed command by command.

<!-- I want create a command called rcodex project organize that will merge all the ideas within the current project together into a new base.  design a way to archive previous ideas and create a new base, do not throw away the previous dependency information though -->
---


You are a coding agent responsible for implementing the first minimal version (v0) of a CLI tool called ResearchCodex (rcodex).

Context:
- There are two design documents available in this repository:
  - Background.md
  - DevPlan_v1.md
- These documents define:
  - The purpose and philosophy of ResearchCodex.
  - The directory structure to create.
  - The CLI commands and their expected behavior.
- .beads_reference_code is a issue management repo called beads, it's the source of idea and is implemented with go as well. Before planning, read and comprehend it's general architecture

Your job:
1. Read the Background.md to have a general idea about the project
2. Read and internalize DevPlan_v0.md. 
3. Reference .beads_reference_code folder for the major design of the go module.
4. Implement as specified in DevPlan_v1.md
5. After implementing each command:
   - Manually test it on a small example project.
   - Make sure the directory layout and DB rows match the specification in DevPlan_v0.md.
   - Fix obvious bugs or panics.
6. After finishing the implemetation, comes up with a testing plan and started with a comprehensive test by your self, you are authorize to install anything you need just mark it down explictly.

Very important:
- Do NOT invent extra features or commands in v0.
- If any behavior is ambiguous or you are uncertain. reference .beads_reference_code for a implementation that still matches the spirit of the document 
- Delve deeper and plan more comprehensively before start
- Follow AGENTS.md's guideline on using the history/ folder.

Start by scaffolding the Go module (go mod init, main.go, etc.), then implement rcodex init, setting up data storing, and proceed command by command.
