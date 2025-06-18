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

(*) Generated journal files are stored in `~/Documents/tyn/journal/{year}/{month}/YYYYMMDD.md`. This path will be OS-sensitive and eventually configurable.

## Usage

### Build

```
make build
```

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
```

Note: The double quotes in the examples above are used for clarity, but they are not required to capture notes. You can omit them if your input doesn't contain special characters that need escaping in your shell.

The following special symbols are used when capturing nodes to provide additional metadata:

```
#tag       - Tags your note with a category (e.g., #projectX, #reading)
@place     - Associates your note with a location (e.g., @home, @office)
:status    - Sets the status of a task (e.g., :todo, :done, :wip)
^date      - Sets a due date for a task (e.g., ^2025-06-17)
URL        - Any valid URL is automatically recognized (e.g., https://example.com)
```

### List Nodes

List all nodes:

```
tn list
```

List only tasks:

```
tn list task
```

List only notes with a tag:

```
tn list note --tag projectX
```

List only tasks at a place with a status:

```
tn list task --place home --status todo
```

## WIP
This project is a work in progress. Output formatting and features are basic and intended as a starting point for further development.

## License
MIT
