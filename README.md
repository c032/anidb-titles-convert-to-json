# anidb-titles-convert-to-json

Convert AniDB's XML dumps to JSON.

This program is meant to be executed from cron, and is currently not
configurable.

## How it works

* Get a list of all files matching the
  `${HOME}/backups/anidb/animetitles/*.xml.xz` pattern.
* For each file name:
  * Replaces the `.xml.xz` extension with `.json` and looks if that file exists.
  * If it doesn't exist:
    * Parse the `.xml.xz` file and create a `.json` file.

## Requirements

* `xz` command.
* `xzcat` command.
* AniDB XML dump files matching path `${HOME}/backups/anidb/animetitles/*.xml.xz`
