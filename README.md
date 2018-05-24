# Fixer for Matrix's Facebook Messenger puppet bridge

I'm using the [Matrix↔️Facebook Messenger puppet bridge](https://github.com/matrix-hacks/matrix-puppet-facebook) to talk with my friends on Facebook using [Riot](https://riot.im) plugged to my [Matrix](https://matrix.org) homeserver.

When I'm having a 1:1 chat with a Facebook friend, the bridge creates a with the friend, the appservice's bot and myself. Because of that, and also because it doesn't automatically set the room's avatar and name, I have troubles identifying what's happening in that room.

This small tool will take the local part of the room ID created by the Matrix↔️Facebook Messenger bot once the friend has joined it, identify the friend, and grab their avatar and display name to set the room's.

**Note: I'm not shaming anyone here, as the bridge is currently still in development. This is totally a temporary solution.**

## Install

You'll first need a Go install, and [gb](https://getgb.io/):

```
go get github.com/constabulary/gb/...
```

Then clone this repo on your computer (or server), walk into it, install its dependencies and build it:

```
git clone https://github.com/babolivier/matrix-puppet-facebook-1to1-fixer
cd matrix-puppet-facebook-1to1-fixer
gb vendor restore
gb build
```

## Configure

An example configuration is located in the [config.sample.yaml](config.sample.yaml) file.

Copy it somewhere, and edit it accordingly with the key's comments.

## Run

You have to run this tool manually for each room you want to edit.

You can run this tool by typing

```
./bin/fixer --room-id-localpart ROOM_ID_LOCALPART
```

where `ROOM_ID_LOCALPART` is the room's ID's localpart (which usually looks like something like `ZUFHhmRzEyUdzljKRz`).

If you moved your configuration file to a different path than `./config.yaml`, you can specify the configuration file's path by appending the `--config CONFIG_PATH` to your command, where `CONFIG_PATH` is the path to your configuration file.

## Stuff it doesn't do

Because it's not supported by [gomatrix](https://github.com/matrix-org/gomatrix) yet, this tool won't set the room as direct chat, nor will it change its push notification settings. These features will be added as soon as gomatrix supports them.
