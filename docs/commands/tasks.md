# Tasks Command

The `tasks` command lets you list, filter, and manage your tasks. You can update any property of a task using the flexible `update` subcommand with flags, or use dedicated subcommands for common actions like changing status, tags, places, or text.

## Usage

```
tn tasks [subcommand] [options]
```

### Main Subcommands
- `list`         List all tasks or filter by status, tag, or place
- `update`       Update any property of a task using flags (status, text, tags, places, due date)

### Dedicated Subcommands (Shortcuts)
- `status`       Set or cycle the status of a task
- `text`         Update only the text/content of a task
- `tag`          Add, remove, or clear tags on a task
- `place`        Add, remove, or clear places on a task
- `date`         Set or remove a due date for a task

## Examples

### Using `update` with flags
```
# Update multiple fields at once
 tn tasks update 1234 --text "Update docs" --tags "docs,urgent" --due "2025-07-01" --status done --places "home,office"

# Update only the status
 tn tasks update 1234 --status blocked

# Update only the tags
 tn tasks update 1234 --tags "bug,high-priority"

# Update only the place
 tn tasks update 1234 --places "backend-team"
```

### Using dedicated subcommands
```
# Set a task's status directly
 tn tasks status set 1234 done

# Cycle to the next status
 tn tasks status next 1234

# Update only the text
 tn tasks text 1234 "Refactor login handler"

# Add a tag
 tn tasks tag add 1234 bug

# Remove a tag
 tn tasks tag remove 1234 bug

# Clear all tags
 tn tasks tag clear 1234

# Add a place
 tn tasks place add 1234 home

# Remove a place
 tn tasks place remove 1234 home

# Clear all places
 tn tasks place clear 1234

# Set a due date
 tn tasks date set 1234 2025-07-15

# Remove a due date
 tn tasks date remove 1234
```

- You can use short IDs for tasks (e.g., `1234` instead of the full UUID).
- Status cycling follows the configured status sequence.

For more details, see the [Command Reference](index.md).
