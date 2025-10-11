# Nots

Nots is a CLI tool for organizing your notes with markdown documents,
note-series and simple templating.

Nots main use is in the terminal, and is configurable to work with your
favourite pager and editor. If it can be with a command, it will work with nots.

## Usage

Open note based on config. (note-series|previous)

```bash
nots
```

Open my-note.md from the notes root directory. create it if one does not exist.
Or, you can select a different template from your template folder instead of the
default one.

```bash
nots open my-note.md
nots open my-note.md --template my-template.md
```

View the note in your preferred pager, or print it straigt to the stdout.

```bash
nots view my-note.md
nots view my-note.md --std-out
```

You can evaluate you note templates with `nots template`. This is a good way to test your templates and control their output before using the for your notes.

```bash
nots template my-template.go
```

use the `-h` flag for any of the commands to see what you can do.

## Configuration

Nots is configurable to allow you to control most things. You can choose
which editor to use, where to store your notes.

Create you config at `~/.config/nots.nots.toml`.

```toml
# set the root directory for all notes. defaults to ~/nots
notes-dir = "~/nots"

# sets which editor to use for notes. defaults to $EDITOR
editor = "hx" # code, nano, vim

# sets the pager for viewing notes. defaults to $PAGER
viewer = 'glow -p' # less, bat

# sets the default open mode for the base nots command. default to "previous"
# valid values: "series" | "previous"
default-open-mode="previous"

# if you want to use the series, you need to define the template
# default-open-mode="series"
# open-mode-series="daily"

# set the default template when creating new notes.
# paths a relative to the `templates/` directory
default-template = "daily.md"

[[ note-series ]]
# required. the name of the series.
name = "daily"

# required. takes a single template expression that we evaluate with the template package,
# and check to see of a file with the evaluated name exists in the series directory. (filename will be appended with .md)
filename-expression = 'time_now("YYYY-MM-DD")'

# optional. if set, controls what template to use for new notes in the series.
template = "daily.md"

# optional. where to put the note series. defaults to the series-name.
directory = ""
```

## Templates

_templates are still WIP. there are more features i would like to add._

Nots has its own template parsing and evaluation. The templates allow for simple expression evaluation.

To define a template, go to `~/.config/nots/templates` and create a .md file.

template:

```
{{ "string"|upper_case }}
{{ filename }}
{{ time_now("YYYY-MM-DD") }}
```

output:

```
STRING
my-file.md
2025-10-01
```

The templates support variables, filters, and functions. You can at any time
evaluate a template with `print_env` to see the available "symbols" for the
evaluation environment
