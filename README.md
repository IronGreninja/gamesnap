# gamesnap(-cli) - simple & minimal game save snapshots

**Step 1:** Find out where your favourite game stores its saves  
**Step 2:** Put it in da config  
**Step 3:** Snapshot  
**Step 4:** Profit?

## Usage

### Config File

```toml
## gamesnap.toml
# search precedence: current-dir (.) -> home-dir (~)

[options]

# minimum backups to keep during prune
backup_max = 7

# path to put snapshots
# variables: true
destination = "." # current-dir

# name common paths and later refer to them
# by their name as substitute by prepending a dot
# and surrounding with double curly brackets
# see example below
#
[variables]

AppData = "C:/Users/igreninja/AppData/Roaming"
localappdata = "C:/Users/igreninja/AppData/Local"

[games]

Elden_Ring = [
    "{{.AppData}}/EldenRing/7688510801809128",
    "{{.AppData}}/EldenRing/GraphicsConfig.xml",
]

Nightreign = [
    "{{.AppData}}/Nightreign"
]

Expedition33 = [
    "{{.localappdata}}/Sandfall/Saved/SaveGames",
    "{{.localappdata}}/Sandfall/Saved/Config/Windows",
]
```

### Snapshots

```sh
gamesnap snap   # all games
gamesnap snap Nightreign Expedition33   # only specified
```

### Restoring Snapshots

Snapshots are restored by pointing to a particular snapshot directory. Here `dir` refers to `destination`, as specified in config

```sh
gamesnap resnap dir/Expedition33/2025-11-07_23.25.29
gamesnap resnap dir/Expedition33/latest    # points to the latest snapshot (except windows)
```

#### Snapshot Architecture

Snapshots for each game, `game_name` is stored in `destination` directory as: `game_name/timestamp`.
Inside each are directories in the form `path*`, and a file `info.toml`
`info.toml` maps `path*`s to `name, type` pairs.
If `type=dir`, the contents of `path*` is a tree of files, and should be restored inside `name`.
If `type=file`, the contents of `path*` is a single file, and should be restored as `name`.

## Credits

[Game-Save-Manager](https://github.com/dyang886/Game-Save-Manager) - Automatically detects save locations, gui-only, windows-only.
