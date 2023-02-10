---
id: cli
title: CLI options
---

### Global options

|Option|Type|Default|Description|
|------|:--:|:-----:|-----------|
|`-E, --env`|`stringToString`||set env variable for running command KEY=VALUE (default [])|
|`--all`|`bool`|false|show all commands (include hidden commands which start with `_`)|
|`--init`|`bool`|false|creates a new lets.yaml in the current folder|
|`--only`|`stringArray`||run only specified command(s) described in cmd as map|
|`--exclude`|`stringArray`||run all but excluded command(s) described in cmd as map|
|`--upgrade`|`bool`|false|upgrade lets to latest version|
|`--no-depends`|`bool`|false|skip 'depends' for running command|
|`-c, --config`|`string`|lets.yaml|specify config|
|`-d, --debug`|`bool`|false|verbose logs|
|`-h, --help`|||help for lets|
|`-v, --version`|||version for lets|
