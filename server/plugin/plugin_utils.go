package plugin

import (
	"github.com/mattermost/mattermost/server/public/model"
	"google.golang.org/api/drive/v3"
)

// CreateBotDMPost posts a direct message using the bot account.
// Any error are not returned and instead logged.
func (p *Plugin) createBotDMPost(userID, message string, props map[string]any) {
	channel, err := p.client.Channel.GetDirect(userID, p.BotUserID)
	if err != nil {
		p.client.Log.Warn("Couldn't get bot's DM channel", "userID", userID, "error", err.Error())
		return
	}

	post := &model.Post{
		UserId:    p.BotUserID,
		ChannelId: channel.Id,
		Message:   message,
		Props:     props,
	}

	if err = p.client.Post.CreatePost(post); err != nil {
		p.client.Log.Warn("Failed to create DM post", "userID", userID, "post", post, "error", err.Error())
		return
	}
}

func (p *Plugin) getAllChannelUsers(channelId string) []*model.User {
	page := 0
	perPage := 100
	allUsers := make([]*model.User, 0)
	for {
		users, appErr := p.API.GetUsers(&model.UserGetOptions{InChannelId: channelId, Page: page, PerPage: perPage})
		if appErr != nil || len(users) == 0 {
			break
		}
		allUsers = append(allUsers, users...)
		page += 1
	}
	return allUsers
}

func (p *Plugin) getUserDisplayName(user *drive.User) string {
	userDisplay := ""
	if user != nil {
		if user.DisplayName != "" {
			userDisplay += user.DisplayName
		}
		if user.EmailAddress != "" {
			userDisplay += "(" + user.EmailAddress + ")"
			user, _ := p.API.GetUserByEmail(user.EmailAddress)

			if user != nil {
				userDisplay += "@" + user.Username
			}
		}
	}
	if userDisplay == "" {
		userDisplay = "Someone"
	}
	return userDisplay
}
