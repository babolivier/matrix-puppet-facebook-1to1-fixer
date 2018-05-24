package main

import (
	"flag"
	"fmt"
	"os"

	"config"

	"github.com/matrix-org/gomatrix"
	"github.com/sirupsen/logrus"
)

// Fixed variables
var (
	// ErrEmptyRoomID is fired if no room ID localpart has been provided, and is
	// followed by the command's usage.
	ErrEmptyRoomID = fmt.Errorf("The room ID localpart cannot be empty")
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
	InfoProcessIsOver = "The room has been fully updated. Don't forget to mark it as direct chat in Riot, and to edit its push rules."
)

// Command line flags
var (
	localpart  = flag.String("room-id-localpart", "", "Room ID localpart")
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

	if len(*localpart) == 0 {
		logrus.Error(ErrEmptyRoomID)
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.Parse(*configFile)
	if err != nil {
		panic(err)
	}

	roomID := fmt.Sprintf("!%s:%s", *localpart, cfg.Matrix.ServerName)
	userID := fmt.Sprintf("@%s:%s", cfg.Matrix.Localpart, cfg.Matrix.ServerName)

	cli, err := gomatrix.NewClient(cfg.Matrix.HomeserverURL, userID, cfg.Matrix.AccessToken)
	if err != nil {
		logrus.Panic(err)
	}

	membersResp, err := cli.JoinedMembers(roomID)
	if err != nil {
		logrus.Panic(err)
	}

	displayNameResp, err := cli.GetOwnDisplayName()
	if err != nil {
		logrus.Panic(err)
	}

	if len(membersResp.Joined) != 3 {
		logrus.Error(ErrInvalidRoomMembersNb)
		os.Exit(1)
	}

	var avatarURL, displayName string
	for _, member := range membersResp.Joined {
		if member.DisplayName != nil && *(member.DisplayName) != displayNameResp.DisplayName {
			displayName = *(member.DisplayName)
			if member.AvatarURL != nil {
				avatarURL = *(member.AvatarURL)
			}
		}
	}

	logrus.WithFields(logrus.Fields{
		"display_name": displayName,
		"avatar_url":   avatarURL,
	}).Info("Found the friend")

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
		logrus.Warn(WarnNoAvatar)
	}

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
		logrus.Warn(WarnNoDisplayName)
	}

	logrus.Info(InfoProcessIsOver)
}
