# snapi3

> :warning: The current version is not finished, the interface is not stable and is highly prone to change.

snapi3 is a companion app for i3wm. It aims to help you manipulating your windows and especially the floating ones.
Its main feature, inspired by [i3-grid](https://github.com/lukeshimanuki/i3-grid), is to offer a grid based positioning system for your floating windows.

Snapi3 displays a grid GUI requesting the user to select one the cell where the floating window shall be positioned and resized. The grid is completely configurable.

## Features
- Snapping floating windows to a grid
- Showing/Hiding windows
- Centering windows
- GUI or CLI execution 

Each functionalities can be applied to one window or to a group of windows. Snapi3 also offers a way to easily manage groups of windows.

## How to use
After installation, execute `snapi3 help` to obtain the detailed help.

Triggering the GUI can be done with the command `snapi3 gui`

The configuration file `snapi3.yml` of snapi3 is stored in your default `config` directory of your system (usually `~/.config`).

## How to build
Snapi3 is written in golang.

The easiest way to modify snapi3 is to use (visual studio code remote development tools)(https://code.visualstudio.com/docs/remote/remote-overview) and especially the container extension.

You'll also need a container execution tool like [docker](https://www.docker.com/) or [podman](https://podman.io/).

The `.devcontainer` located in the root of this repository contains a configuration of container to be used with visual studio code.

So to modify the repository perform the following operations:
- Clone the snapi3 repository
- Use the `Open Folder in container` function of the vscode palette to open the directory
- Wait for the container to be running
- Inside the container, execute the command:
  ```bash
  goreleaser --snapshot --skip-publish --rm-dist
  ```
- The output binary shall be located in the `dist` component