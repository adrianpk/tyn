# tyn

Tyn is a command-line tool for capturing notes quickly and precisely. A minimalist companion for your terminal workflow. The name comes from the Welsh word tyn, meaning to tighten or draw in, evoking the act of focusing thought just before release, like drawing a bowstring. It's also a play on the letters T and N, a mnemonic for tight notes, a nod to the tool’s command (`tn`), and a pair that’s easily accessible on a Colemak layout.

A simple CLI for capturing, listing, and managing notes, tasks, and links, fast, focused, and distraction-free.

## Features
- Capture notes, tasks, and links from the command line
- List all nodes or filter by type, tag, place, or status
- Automatic daily journal generation from captured nodes (*)
- System notifications for tasks with due dates
- Pretty-printed output for easy inspection
- More to come

(*) Generated journal files are stored in `~/Documents/tyn/journal/{year}/{month}/YYYYMMDD.md`. This path will be OS-sensitive and eventually configurable. For a sample of what the generated output looks like, check out our [example journal entry](docs/examples/20250619.md). The system also maintains a [rotating index](docs/examples/index.md) accessible at `~/Documents/tyn/index.md` with links to journal entries.

## Installation

```
go install github.com/adrianpk/tyn@latest
```

If you prefer to compile and test locally, you can clone the repository and use make:

```
git clone https://github.com/adrianpk/tyn.git
cd tyn
make build
```

## Usage

### Capture a Node

You can capture a node directly from the command line. Here are some real-world examples:

```
# Capture a quick meeting note
 tn capture "Sync with Alice and Bob #projectX :done ^2025-06-17 Discussed Q3 roadmap https://company.com/roadmap"

# Add a task for today
 tn capture "Write project summary :todo #writing @home"

# Save a useful link for later
 tn capture "Read about Go generics https://go.dev/doc/tutorial/generics #reading"

# Log a place-based note
 tn capture "Coffee with Carol #networking @cafe"

# Capture a completed task with a specific date
 tn capture "Submit tax report :done ^2025-04-15 #finance"

# Capture a draft snippet for a future document
 tn capture +code-echo "Network Security Alert: Identifying Echo-Pattern Vulnerabilities #security Our team recently discovered a critical vulnerability in proxy services."
```

Note: The double quotes in the examples above are used for clarity, but they are not required to capture notes. You can omit them if your input doesn't contain special characters that need escaping in your shell.

For convenience, the `capture` command can also be invoked using the shorter aliases `cap` or `c` (e.g., `tn cap`, `tn c`).

The following special symbols are used when capturing nodes to provide additional metadata:

```
#tag       - Tags your note with a category (e.g., #projectX, #reading)
@place     - Associates your note with a location (e.g., @home, @office)
:status    - Sets the status of a task (e.g., :todo, :done, :wip)
^date      - Sets a due date for a task (e.g., ^2025-06-17)
URL        - Any valid URL is automatically recognized (e.g., https://example.com)
```

### Managing Tasks

Tyn provides specialized commands to manage tasks with more efficiency:

```
tn tasks            # List all tasks (shorthand for 'tn tasks list')
tn tasks list       # List all tasks explicitly
tn tasks list :todo # List tasks with todo status
tn tasks list #urgent # List tasks with urgent tag
tn tasks list @home # List tasks at home location
```

You can also combine filters:

```
tn tasks list :wip #urgent      # In-progress urgent tasks
tn tasks list @office :todo     # Todo tasks at the office
tn tasks list #project @home    # Project tasks at home
```

Each task is displayed with a short ID that can be used to reference it in other commands:

```
ID     STATUS     CONTENT                                            TAGS/PLACES
--------------------------------------------------------------------------------
e0e9   [wip]      Fix critical bug    Need to fix memory leak issue  #urgent
```

To change task status, you can use the following commands:

```
tn tasks status set e0e9 done   # Set specific status
tn tasks status next e0e9       # Move to next status in cycle
tn tasks status prev e0e9       # Move to previous status in cycle
```

For convenience, these commands can also be used with shorter aliases:

```
tn t l                # Short for "tn tasks list"
tn t s next e0e9      # Short for "tn tasks status next e0e9"
```

Example output of `tn tasks list`:

```
ID     STATUS     CONTENT                                       TAGS/PLACES          !
------------------------------------------------------------------------------------------
5f5c   [done]     Sync with Alice and Bob Discussed Q3 roadmap  #projectx             
cf41   [todo]     Write project summary Due by end of week      #writing @home        
5701   [done]     Submit tax report Filed electronically        #finance @office      
f5bb   [wip]      Research cloud providers Comparing AWS, GC... #infrastructure       
a5b7   [todo]     Schedule dentist appointment                  #health @personal     
cdce   [todo]     Review pull request                           #23 #coding           
e991   [done]     Design database schema Finalized user and ... #projectX @home       
b3a2   [wip]      Order new laptop Looking at developer-focu... #shopping @online     
993f   [blocked]  Fix critical bug Need to fix memory leak i... #urgent              ⌛
```

When changing a task's status with `tn tasks status next`, you'll see:

```
Task status updated: 'wip' → 'blocked'
todo → ready → <wip> → [blocked] → on-hold → review → done → canceled → waiting
```

This visualization shows the status cycle, with:
- `<wip>` indicating the original status
- `[blocked]` highlighting the new status

When we executed `tn tasks status next 993f` on the task "Fix critical bug", it moved from 'wip' (work in progress) to 'blocked' status in the cycle. The task is still overdue as indicated by the ⌛ symbol in the list view.

After the command is executed, this task appears as:

```
993f   [blocked]  Fix critical bug Need to fix memory leak i... #urgent              ⌛
```

You can also set a specific status directly without cycling through states:

```
tn tasks status set 993f done
```

This produces:

```
Task status updated: 'blocked' → 'done'
todo → ready → wip → <blocked> → on-hold → review → [done] → canceled → waiting
```

And the task will appear as:

```
993f   [done]     Fix critical bug Need to fix memory leak i... #urgent
```

Note that the overdue indicator (⌛) disappears when a task is marked as done, even if its due date has passed.

### Update Tasks

Tyn provides convenient commands to update various aspects of your tasks:

```
# Update task text content
tn tasks update d356 --text "Fix memory leak in authentication service"

# Set or update a task's due date
tn tasks update d356 --due "2025-07-01"

# Change task tags
tn tasks update d356 --tags "bug,high-priority"

# Update task location/places
tn tasks update d356 --places "office,backend-team"

# Combine multiple updates in one command
tn tasks update d356 --text "Prepare Q2 financial report" --due "2025-06-30" --tags "finance,quarterly" --places "accounting"
```

You can also use dedicated commands for specific updates:

```
# Text updates
tn tasks text d356 "Call vendor about license renewal issues"

# Tag management
tn tasks tag add d356 blocker
tn tasks tag remove d356 blocker
tn tasks tag clear d356

# Place management
tn tasks place add d356 meeting-room
tn tasks place remove d356 meeting-room
tn tasks place clear d356

# Due date management
tn tasks date set d356 2025-07-15
tn tasks date remove d356
```

For a full list and detailed explanation of all commands, see [docs/commands/index.md](docs/commands/index.md).

## Roadmap

While additional features, new commands, and integrations are being considered for the future, the next efforts will be dedicated to:

- Polishing and strengthening the current features.
- Testing and covering edge cases to ensure reliable behavior.
- Building a robust test suite to guarantee that future refactoring does not introduce regressions or new bugs.

The current priority is to ensure the existing functionality is solid and well-tested.

## License
MIT
