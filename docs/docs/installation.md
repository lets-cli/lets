---
id: installation
title: Installation
sidebar_label: Installation
---

### **Shell script**

This will install `lets` binary to `/usr/local/bin` directory and will require `sudo`. But you can change install location to any directory you want

```bash
sudo curl -sfL https://raw.githubusercontent.com/lets-cli/lets/master/install.sh | sudo sh -s -- -b /usr/local/bin
```

Alternatively, if you do not want to install `lets` with a `sudo`, you can do next:

1. Create (if do not have any) some bin directory somewhere in your home dir 

```bash
cd ~
mkdir bin
```

2. Add this dir to your `PATH`

```bash
vim ~/.bashrc # or ~/.zshrc

Add the folowing line, save file and restart the shell.

```bash
export PATH=$PATH:$HOME/bin
```

Now you can add any binary files to ~/bin and you system will see and executable globally.

3. Install lets to ~/bin

```bash
curl -sfL https://raw.githubusercontent.com/lets-cli/lets/master/install.sh | sh -s -- -b ~/bin
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
