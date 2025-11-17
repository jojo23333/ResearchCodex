
1. I want codex to behave like a capable undergrad, and the user would be a supervisor.


## Expected behavior of the codex agent

* When it comes to new idea, plan and figure out first. (maybe use gpt-5 for this phase)
* ideas.md and todos.md is generated at this phase for human expert to review, correct and give directional guidelines, it's like plan mode.
* Once ideas and implementation pathway is fixed, proceed with adding it to the database. 
* The rest of implemetation should be autonomous but can be stage wise. Todos are created within the database and are fullfilled by codex.

### rcodex init
'rcodex init' create the following project
```
- .exps (stores the result/output from experiements)
- threads (contains description of threads and plans to )
- artifacts (Contains papers and reference codes, storage ways TBD)
- srcs (Contains modulized tools that is potentially sharable in between experiments)
- experiments (Contains python code to do the experiments)
- .env (stores environmental parameter that maybe useful for doing experiment)
- .rcodex
  - config.yaml  
- Topics.md
- Agents.md
```

### rcodex thread -n 'thread_name'
Initialize a new thread of ideas, specifically creating following subfolder and subfiles:
```
- .exps
  - thread_name (stores the result/output from experiements on thread_1_name)
- threads (contains description of threads and plans to )
  - thread_name
    - ideas.md (To be filled with a detailed explanation of the idea1 idea2 ...)
    - todos.md (To be filled with a implementation todo plans for idea1 idea2 ...)
    - results.md (To be filled post experiment analysis..)
- experiments
  - thread_name
```
start working on the thread by setting .rcodex config.yaml to thread_name 

### rcodex thread -t 'thread_name'
Start work on the thread 'thread_name', set .rcodex config.yaml thread feild to 'thread_name'

<!-- ### rcodex status
Give a summary of research status on current thread
1. get current_thread_name from .rcodex config
2. Return @threads/thread_name/thread.md + @threads/thread_name/todo.md  -->
<!-- 3. Summerize the thread, where we are in terms of thread development, as well as if there is any result yet. -->

### rcodex get idea|todo
Get the current idea or todo stored in db

### rcodex create idea|todo|experiment "full content" 
Create the idea node, experiment, result node, set experiment and result status as null.
[todos] <- idea  --> exp1 -> result1
                 `-> exp2 -> result2 
Return unique id to idea, set in.rcodex config.yaml idea_head to idea_id

A new idea is default to depend on the privious idea on the chain, make sure previous idea todo is all finished before creating the new one.

### rcodex update idea|todo|experiment 
todos should have todo and doc field, once impelementation completed running rcodex consolidate summarize the implmentation as doc (TODO consider a better naming other than todo?)


TODO: 
1. design storage field in rcodex
2. design defualt .default_agents.md for repo that use rcodex
3. start vibe coding Rcodex