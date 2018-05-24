package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"config"

	"github.com/matrix-org/gomatrix"
	"github.com/sirupsen/logrus"
)

// Fixed variables
var (
	// ErrRoomIDEmptyOrInvalid is fired if no room ID localpart has been provided
	// or is invalid, and is followed by the command's usage.
	ErrRoomIDEmptyOrInvalid = fmt.Errorf("The room ID localpart is either empty or invalid")
	// ErrInvalidRoomMembersNb is fired if the number of joined members in the
	// room isn't 3 (i.e.: the user, their friend, and the Facebook bot).
	ErrInvalidRoomMembersNb = fmt.Errorf("Invalid number of members in the room: either the friend hasn't joined yet, or there's more than one friend in the room")
	// WarnNoAvatar is displayed if the friend doesn't have an avatar.
	WarnNoAvatar = "The friend doesn't have an avatar set"
	// WarnNoDisplayName is displayed if the friend doesn't have a display name.
	WarnNoDisplayName = "The friend doesn't have a display name set"
	// InfoAvatarUpdated is displayed if the room's avatar has been updated.
	InfoAvatarUpdated = "Room's avatar updated"
	// InfoNameUpdated is displayed if the room's name has been updated.
	InfoNameUpdated = "Room's name updated"
	// InfoProcessIsOver is displayed once the whole process is over, just before
	// exiting.
	InfoProcessIsOver = "The room has been fully updated. Don't forget to mark it as direct chat in Riot, and to edit its notification rules."
)

// Command line flags
var (
	localpart  = flag.String("room-id-localpart", "", "Room ID localpart (i.e. 'ZUFHhmRzEyUdzljKRz')")
	configFile = flag.String("config", "config.yaml", "Configuration file")
)

// MRoomAvatarContent represents the content of the "m.room.avatar" state event.
// https://matrix.org/docs/spec/client_server/r0.3.0.html#m-room-avatar
type MRoomAvatarContent struct {
	URL string `json:"url"`
}

// MRoomNameContent represents the content of the "m.room.name" state event.
// https://matrix.org/docs/spec/client_server/r0.3.0.html#m-room-name
type MRoomNameContent struct {
	Name string `json:"name"`
}

func main() {
	logConfig()

	flag.Parse()

	// We need the room ID's localpart to be non-empty and only composed of letters.
	roomIDLocalpartRgxp := regexp.MustCompile("^[a-zA-Z]+")
	if len(*localpart) == 0 || !roomIDLocalpartRgxp.Match([]byte(*localpart)) {
		logrus.Error(ErrRoomIDEmptyOrInvalid)
		flag.Usage()
		os.Exit(1)
	}

	// Load the configuration from the configuration file.
	cfg, err := config.Parse(*configFile)
	if err != nil {
		panic(err)
	}

	// Compute the room's ID along with the current user's.
	roomID := fmt.Sprintf("!%s:%s", *localpart, cfg.Matrix.ServerName)
	userID := fmt.Sprintf("@%s:%s", cfg.Matrix.Localpart, cfg.Matrix.ServerName)

	// Load the Matrix client from configuration data.
	cli, err := gomatrix.NewClient(cfg.Matrix.HomeserverURL, userID, cfg.Matrix.AccessToken)
	if err != nil {
		logrus.Panic(err)
	}

	// Retrieve the list of joined members in the room.
	membersResp, err := cli.JoinedMembers(roomID)
	if err != nil {
		logrus.Panic(err)
	}

	// Retrieve the current user's own display  name.
	displayNameResp, err := cli.GetOwnDisplayName()
	if err != nil {
		logrus.Panic(err)
	}

	// Check if the number of joined members is three, as it should be with a
	// 1:1 puppeted chat (the current user, their friend, and the AS bot).
	if len(membersResp.Joined) != 3 {
		logrus.Error(ErrInvalidRoomMembersNb)
		os.Exit(1)
	}

	// Iterate over the slice of joined members.
	var avatarURL, displayName string
	for _, member := range membersResp.Joined {
		// The friend should be the only joined member who has a display name set
		// which isn't the same as the current user's.
		if member.DisplayName != nil && *(member.DisplayName) != displayNameResp.DisplayName {
			displayName = *(member.DisplayName)
			// If there's also an avatar set for the friend, use it.
			if member.AvatarURL != nil {
				avatarURL = *(member.AvatarURL)
			}
		}
	}

	logrus.WithFields(logrus.Fields{
		"display_name": displayName,
		"avatar_url":   avatarURL,
	}).Info("Found the friend")

	// If the avatar has been found, set it as the room's avatar using a
	// m.room.avatar state event.
	if len(avatarURL) > 0 {
		if _, err := cli.SendStateEvent(
			roomID,
			"m.room.avatar",
			"",
			MRoomAvatarContent{
				URL: avatarURL,
			},
		); err != nil {
			logrus.Panic(err)
		}
		logrus.Info(InfoAvatarUpdated)
	} else {
		// Else print a warning so the user can see it clearly.
		logrus.Warn(WarnNoAvatar)
	}

	// If the display name has been found, set it as the room's name using a
	// m.room.name state event. This condition shouldn't be necessary, but heh,
	// at least that might cover a potential regression from the bridge.
	if len(displayName) > 0 {
		if _, err := cli.SendStateEvent(
			roomID,
			"m.room.name",
			"",
			MRoomNameContent{
				Name: displayName + " (Facebook)", // TODO: Allow custom suffix
			},
		); err != nil {
			logrus.Panic(err)
		}
		logrus.Info(InfoNameUpdated)
	} else {
		// Else print a warning so the user can see it clearly.
		logrus.Warn(WarnNoDisplayName)
	}

	// Print a shiny message telling the user the process is over, but it's up
	// to them to set the room as a direct chat and to update the room's push
	// notification settings, since that's not supported by gomatrix.
	logrus.Info(InfoProcessIsOver)
}
