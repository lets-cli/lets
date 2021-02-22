---
id: installation
title: Installation
sidebar_label: Installation
---

## Install

### **Shell script**

```bash
curl --proto '=https' --tlsv1.2 -sSf https://lets-cli.org/install.sh | sh -- -b ~/bin
```

This will install **latest** `lets` binary to `~/bin` directory.

To be able to run `lets` from any place in system add `$HOME/bin` to your `PATH`

Open one of these files

```bash
vim ~/.profile # or vim ~/.bashrc or ~/.zshrc
```

Add the following line at the end of file, save file and restart the shell.

```bash
export PATH=$PATH:$HOME/bin
```

You can change install location to any directory you want, probably to directory that is in your $PATH

To install a specific version of `lets` (for example `v0.0.21`):

```bash
curl --proto '=https' --tlsv1.2 -sSf https://lets-cli.org/install.sh | sh -s -- v0.0.21
```

To use `lets` globally in system you may want to install `lets` to `/usr/local/bin`

> May require `sudo`

```bash
curl --proto '=https' --tlsv1.2 -sSf https://lets-cli.org/install.sh | sh -s -- -b /usr/local/bin
```

### **Binary (Cross-platform)**

Download the version you need for your platform from [Lets Releases](https://github.com/lets-cli/lets/releases). 

Once downloaded, the binary can be run from anywhere.

Ideally, you should install it somewhere in your PATH for easy use. `/usr/local/bin` is the most probable location.

### **Arch Linux**

You can get binary release from https://aur.archlinux.org/packages/lets-bin/

If you are using `yay` as AUR helper:

```bash
yay -S lets-bin
```

Also you can get bleeding edge version from https://aur.archlinux.org/packages/lets-git/

```bash
yay -S lets-git
```

## Update

### Self upgrade

```bash
lets --upgrade
```

### Shell script 

To update `lets` you can use the same [shell script](#shell-script) from above.

Running script will update `lets` to **latest** version.

### Binary 

You can download latest version from [Lets Releases](https://github.com/lets-cli/lets/releases). 

### Arch Linux

AUR repository always provides **latest** version.