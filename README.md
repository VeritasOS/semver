# semver

A simple tool that helps you menage semantic versions, by bumping up its components
accordingly. If semantic versions are maintained as git tags, semver will use old
tags and optional lines in your commit message to bump up the semantic versions and
tag the repo. By default, patch version is bumped up by 1, if no other options are
specified. See the usage section for more details. 

<!--
# <a name="install"></a>Installation

## Download

- [MacOSX][dist-darwin]
- [Linux][dist-linux]
- [Windows][dist-windows]

## Installing on MacOSX and Linux

To install, put the `semver` binary in your `PATH`. On MacOSX and Linux, we
recommend `$HOME/bin/`. You may need to [update your `PATH`][home-bin-path] in
your `$HOME/.bashrc` or other shell's config file to include the directory
where you put `semver`.

```sh
# abbreviated install instructions (MacOSX/Linux)
mkdir -p ~/bin
# use wget or curl to download semver
SEMVER_URL="https://<update-me-please>/semver/<version>/semver_$(uname -s | tr '[:upper:]' '[:lower:]')_amd64"
wget -O ~/bin/semver "$SEMVER_URL" || curl -o ~/bin/semver "$SEMVER_URL"
chmod +x ~/bin/semver
```

## Installing on Windows

```powershell
# abbreviated install instructions (Windows Powershell)
mkdir "C:\Program Files\semver"
Invoke-WebRequest -Uri https://<update-me-please>/semver/<version>/semver_windows_amd64 -OutFile "C:\Program Files\semver\semver.exe"
```

In other words:

- Rename `semver_windows_amd64` to `semver.exe`
- Move `semver.exe` to `C:\Program Files\semver\semver.exe`

## Upgrading

To upgrade, run this:

```sh
semver upgrade
```
-->

# Usage

### When the current semver is specified in the command line

```sh
semver 0.1.2         # will return 0.1.3 (patch is the default bump type)
semver 0.1.2 --minor # will return 0.2.0
semver 0.1.2 --major # will return 1.0.0
```

### When the semversion is stored as git tags

If `--git` is specified at runtime, the new semver is determined by using:
- The latest tag (on the current branch of the repo where it is being ran) AND
- The bump type(s) found in the commit messgaes from the last tag on the current branch up to HEAD. 
 - If anywhere in the messages, a string of types `+major`, `+minor` or `+patch` is found, "major", 
 "minor" or "patch" will be used as the bump type, respectively. 
 - If more than one bump type is found in the commit messages, the highest one will be used as the 
bump type (major > minor > patch)

Similarly, if the commit message contains `+meta=<some_meta_data>` the latest meta data will be 
attached to the new sematic version when specifying `--git`

For example, assume that the latest tag was `1.3.4+6.8` and that the last commit message was

```
Adding new functionality in a backwards-compatible manner

New feature to do something cool!

+minor
+meta=7.3

That's all folks!
```

In this case:

```sh
semver --git         # will print 1.4.0+6.8
```

If you want to create a tag using this semver and push it:

```sh
semver --git --push  # will print 1.4.0+7.3 and push the tag to the repo.
```

[home-bin-path]: https://askubuntu.com/questions/402353/how-to-add-home-username-bin-to-path#402356
