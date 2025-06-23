# Capture Command

The `capture` command is the core of Tyn. It lets you quickly add notes, tasks, links, and drafts from the command line, using a simple, expressive syntax.

## Usage

```
tn capture [content]
```

You can also use the aliases `cap` or `c`.

## Examples

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

### Special Syntax

| Syntax    | Description                                    |
|-----------|------------------------------------------------|
| `#tag`    | Add tags to any node                           |
| `@place`  | Add a place/location                           |
| `:status` | Set a status (for tasks)                       |
| `^date`   | Set a due date (for tasks)                     |
| `+draft`  | Start a draft capture (always type `draft`)    |
| URLs      | Automatically recognized as links              |

Drafts are grouped by their draft name and can be combined later. A future command will allow you to combine all entries with the same draft name into a single markdown document.

For more details, see the [Command Reference](index.md).
