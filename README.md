# dotman

Super simple dotfiles manager. It will create symlinks between your dotfiles repository and the system files.

## Install

Via go get:

```
go get github.com/erizocosmico/dotman
```

Manual build:

```
cd /path/to/dotman
go build -o /usr/local/bin/dotman .
```

Or just grab one of the releases from the [releases page](https://github.com/erizocosmico/dotman/releases).

## Usage

### Config

Create a `config.yaml` file with the mappings between files:

```yaml
config/i3: ~/.config/i3
config.json: ~/.config/some/config.json
```

`dotman` admits only a very limited form of yaml: `SRC_PATH: DST_PATH`.
The key is the source file or directory, and the value is the destination.

For example, `config/i3: ~/.config/i3` will create a symlink between `./config/i3` and `~/.config/i3`, no matter if it's a file or a directory.

### Run

```
dotman
```
Just that.

There's a couple options you can specify:
- `-force` if the destination path already exists, delete it first and then symlink. This is specially dangerous if the destination is a folder! Use with caution!
- `-config PATH` path to the config file. `config.yaml` by default.

By default, if the destination exists, dotman will ask you if you want to replace it, so if you didn't use `-force` dotman will not mess anything unless you specifically confirm whether you want to overwrite some file or not.

## License

MIT License, see [LICENSE](/LICENSE)