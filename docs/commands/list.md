# List Command

The `list` command displays all nodes (notes, tasks, links, drafts) or filters them by type, tag, place, or status. It is the main way to view your captured data in Tyn.

## Usage

```
tn list [type] [--tag TAG] [--place PLACE] [--status STATUS]
```

- `[type]` can be `note`, `task`, `link`, or `draft` to filter by node type.
- `--tag` (`-t`) filters by tag.
- `--place` (`-p`) filters by place.
- `--status` (`-s`) filters by status (for tasks).

## Examples

```
# List all nodes
 tn list

# List only tasks
 tn list task

# List only notes
 tn list note

# List only drafts
 tn list draft

# List nodes with a specific tag
 tn list --tag urgent

# List tasks with a specific status
 tn list task --status done

# List nodes at a specific place
 tn list --place home

# Combine filters
 tn list task --tag projectX --place office --status wip
```

- The output includes all relevant fields, including draft name for drafts.
- You can use short or long flags for filters (e.g., `-t` or `--tag`).
- Filtering is case-sensitive for tags and places.

For more details, see the [Command Reference](index.md).
