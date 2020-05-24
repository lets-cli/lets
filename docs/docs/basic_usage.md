---
id: basic_usage
title: Basic usage
---

We will start with a simple example here. More advanced usage you will find in the [Advanced section](advanced_usage.md).

Assume you have a `node.js` project.

### Create config

Go to your project repo and create `lets.yaml`.

**`touch lets.yaml`**

Now add `.lets` to `.gitignore`. `.lets` is a lets directory where it stores some internal metadata. You do not need to commit this directory.

### Write first command

First of all you want to be able to run your project.

You have your `package.json` with all dependencies and scripts in it.

Lets create first command:

```yaml
shell: bash

command:
  run:
    description: Run nodejs server
    cmd: npm run server
```

That's it. You've just created your first `lets` command.

Run `lets` in terminal to see all available commands.

### Run first command

Now you can use this command to start your server.

**`lets run`**
