# wyvern

wyvern is a work-in-progress workflow engine for orchestration and automation of tasks. it organized tasks into a DAG (directed acyclic graph) and executes them in parallel.

wyvern supports atomic tasks plugins, which are written in golang and can be extended to support any task type.

wyvern also supports custom store implementation, which can be used to persist and retrieve task state.

## Usage