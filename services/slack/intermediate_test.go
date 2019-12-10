package slack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-server/v5/model"
)

func TestIntermediateChannelSanitise(t *testing.T) {
	t.Run("Properties should respect the max length", func(t *testing.T) {
		channel := IntermediateChannel{
			Name:        strings.Repeat("a", 70),
			DisplayName: strings.Repeat("b", 70),
			Purpose:     strings.Repeat("c", 400),
			Header:      strings.Repeat("d", 1100),
		}

		expectedName := strings.Repeat("a", 64)
		expectedDisplayName := strings.Repeat("b", 64)
		expectedPurpose := strings.Repeat("c", 250)
		expectedHeader := strings.Repeat("d", 1024)

		channel.Sanitise()

		assert.Equal(t, expectedName, channel.Name)
		assert.Equal(t, expectedDisplayName, channel.DisplayName)
		assert.Equal(t, expectedPurpose, channel.Purpose)
		assert.Equal(t, expectedHeader, channel.Header)
	})

	t.Run("Name and DisplayName should be trimmed", func(t *testing.T) {
		channel := IntermediateChannel{
			Name:        "_-_channel--name-_-__",
			DisplayName: "-display_name--",
		}

		channel.Sanitise()

		assert.Equal(t, "channel--name", channel.Name)
		assert.Equal(t, "display_name", channel.DisplayName)
	})

	t.Run("Name and DisplayName should be longer than 1 character", func(t *testing.T) {
		channel := IntermediateChannel{
			Name:        "a",
			DisplayName: "-_---_--b----",
		}

		channel.Sanitise()

		assert.Equal(t, "slack-channel-a", channel.Name)
		assert.Equal(t, "slack-channel-b", channel.DisplayName)
	})

	t.Run("Name and DisplayName should contain valid characters or return id", func(t *testing.T) {
		channel := IntermediateChannel{
			Id:          "channelId1",
			Name:        "_-_chännel--name-_-__",
			DisplayName: "-døsplay_name--",
		}

		channel.Sanitise()

		assert.Equal(t, "channelid1", channel.Name)
		assert.Equal(t, "channelid1", channel.DisplayName)
	})
}

func TestTransformPublicChannels(t *testing.T) {
	users := map[string]*IntermediateUser{"m1": {}, "m2": {}, "m3": {}}
	publicChannels := []SlackChannel{
		{
			Id:      "id1",
			Name:    "channel-name-1",
			Creator: "creator1",
			Members: []string{"m1", "m2", "m3"},
			Purpose: SlackChannelSub{
				Value: "purpose1",
			},
			Topic: SlackChannelSub{
				Value: "topic1",
			},
			Type: model.CHANNEL_OPEN,
		},
		{
			Id:      "id2",
			Name:    "channel-name-2",
			Creator: "creator2",
			Members: []string{"m1", "m2", "m3"},
			Purpose: SlackChannelSub{
				Value: "purpose2",
			},
			Topic: SlackChannelSub{
				Value: "topic2",
			},
			Type: model.CHANNEL_OPEN,
		},
		{
			Id:      "id3",
			Name:    "channel-name-3",
			Creator: "creator3",
			Members: []string{"m1", "m2", "m3"},
			Purpose: SlackChannelSub{
				Value: "purpose3",
			},
			Topic: SlackChannelSub{
				Value: "topic3",
			},
			Type: model.CHANNEL_OPEN,
		},
	}

	result := TransformChannels(publicChannels, users)
	require.Len(t, result, len(publicChannels))

	for i := range result {
		assert.Equal(t, fmt.Sprintf("channel-name-%d", i+1), result[i].Name)
		assert.Equal(t, fmt.Sprintf("channel-name-%d", i+1), result[i].DisplayName)
		assert.Equal(t, []string{"m1", "m2", "m3"}, result[i].Members)
		assert.Equal(t, fmt.Sprintf("purpose%d", i+1), result[i].Purpose)
		assert.Equal(t, fmt.Sprintf("topic%d", i+1), result[i].Header)
		assert.Equal(t, model.CHANNEL_OPEN, result[i].Type)
	}
}

func TestTransformPublicChannelsWithAnInvalidMember(t *testing.T) {
	users := map[string]*IntermediateUser{"m1": {}, "m2": {}}
	publicChannels := []SlackChannel{
		{
			Id:      "id1",
			Name:    "channel-name-1",
			Creator: "creator1",
			Members: []string{"m1", "m2", "m3"},
			Purpose: SlackChannelSub{
				Value: "purpose1",
			},
			Topic: SlackChannelSub{
				Value: "topic1",
			},
			Type: model.CHANNEL_OPEN,
		},
		{
			Id:      "id2",
			Name:    "channel-name-2",
			Creator: "creator2",
			Members: []string{"m1", "m2", "m3"},
			Purpose: SlackChannelSub{
				Value: "purpose2",
			},
			Topic: SlackChannelSub{
				Value: "topic2",
			},
			Type: model.CHANNEL_OPEN,
		},
		{
			Id:      "id3",
			Name:    "channel-name-3",
			Creator: "creator3",
			Members: []string{"m1", "m2", "m3"},
			Purpose: SlackChannelSub{
				Value: "purpose3",
			},
			Topic: SlackChannelSub{
				Value: "topic3",
			},
			Type: model.CHANNEL_OPEN,
		},
	}

	result := TransformChannels(publicChannels, users)
	require.Len(t, result, len(publicChannels))

	for i := range result {
		assert.Equal(t, fmt.Sprintf("channel-name-%d", i+1), result[i].Name)
		assert.Equal(t, fmt.Sprintf("channel-name-%d", i+1), result[i].DisplayName)
		assert.Equal(t, []string{"m1", "m2"}, result[i].Members)
		assert.Equal(t, fmt.Sprintf("purpose%d", i+1), result[i].Purpose)
		assert.Equal(t, fmt.Sprintf("topic%d", i+1), result[i].Header)
		assert.Equal(t, model.CHANNEL_OPEN, result[i].Type)
	}
}

func TestTransformPrivateChannels(t *testing.T) {
	users := map[string]*IntermediateUser{"m1": {}, "m2": {}, "m3": {}}
	privateChannels := []SlackChannel{
		{
			Id:      "id1",
			Name:    "channel-name-1",
			Creator: "creator1",
			Members: []string{"m1", "m2", "m3"},
			Purpose: SlackChannelSub{
				Value: "purpose1",
			},
			Topic: SlackChannelSub{
				Value: "topic1",
			},
			Type: model.CHANNEL_PRIVATE,
		},
		{
			Id:      "id2",
			Name:    "channel-name-2",
			Creator: "creator2",
			Members: []string{"m1", "m2", "m3"},
			Purpose: SlackChannelSub{
				Value: "purpose2",
			},
			Topic: SlackChannelSub{
				Value: "topic2",
			},
			Type: model.CHANNEL_PRIVATE,
		},
		{
			Id:      "id3",
			Name:    "channel-name-3",
			Creator: "creator3",
			Members: []string{"m1", "m2", "m3"},
			Purpose: SlackChannelSub{
				Value: "purpose3",
			},
			Topic: SlackChannelSub{
				Value: "topic3",
			},
			Type: model.CHANNEL_PRIVATE,
		},
	}

	result := TransformChannels(privateChannels, users)
	require.Len(t, result, len(privateChannels))

	for i := range result {
		assert.Equal(t, fmt.Sprintf("channel-name-%d", i+1), result[i].Name)
		assert.Equal(t, fmt.Sprintf("channel-name-%d", i+1), result[i].DisplayName)
		assert.Equal(t, []string{"m1", "m2", "m3"}, result[i].Members)
		assert.Equal(t, fmt.Sprintf("purpose%d", i+1), result[i].Purpose)
		assert.Equal(t, fmt.Sprintf("topic%d", i+1), result[i].Header)
		assert.Equal(t, model.CHANNEL_PRIVATE, result[i].Type)
	}
}

func TestTransformBigGroupChannels(t *testing.T) {
	channelMembers := []string{"m1", "m2", "m3", "m4", "m5", "m6", "m7", "m8", "m9"}
	users := map[string]*IntermediateUser{"m1": {}, "m2": {}, "m3": {}, "m4": {}, "m5": {}, "m6": {}, "m7": {}, "m8": {}, "m9": {}}
	bigGroupChannels := []SlackChannel{
		{
			Id:      "id1",
			Creator: "creator1",
			Members: channelMembers,
			Purpose: SlackChannelSub{
				Value: "purpose1",
			},
			Topic: SlackChannelSub{
				Value: "topic1",
			},
			Type: model.CHANNEL_GROUP,
		},
		{
			Id:      "id2",
			Name:    "invalid",
			Creator: "creator2",
			Members: channelMembers,
			Purpose: SlackChannelSub{
				Value: "purpose2",
			},
			Topic: SlackChannelSub{
				Value: "topic2",
			},
			Type: model.CHANNEL_GROUP,
		},
		{
			Id:      "id3",
			Creator: "creator3",
			Members: channelMembers,
			Purpose: SlackChannelSub{
				Value: "purpose3",
			},
			Topic: SlackChannelSub{
				Value: "topic3",
			},
			Type: model.CHANNEL_GROUP,
		},
	}

	result := TransformChannels(bigGroupChannels, users)
	require.Len(t, result, len(bigGroupChannels))

	for i := range result {
		assert.Equal(t, fmt.Sprintf("purpose%d", i+1), result[i].Name)
		assert.Equal(t, fmt.Sprintf("purpose%d", i+1), result[i].DisplayName)
		assert.Equal(t, channelMembers, result[i].Members)
		assert.Equal(t, fmt.Sprintf("purpose%d", i+1), result[i].Purpose)
		assert.Equal(t, fmt.Sprintf("topic%d", i+1), result[i].Header)
		assert.Equal(t, model.CHANNEL_PRIVATE, result[i].Type)
	}
}

func TestTransformRegularGroupChannels(t *testing.T) {
	users := map[string]*IntermediateUser{"m1": {}, "m2": {}, "m3": {}}
	regularGroupChannels := []SlackChannel{
		{
			Id:      "id1",
			Name:    "channel-name-1",
			Members: []string{"m1", "m2", "m3"},
			Purpose: SlackChannelSub{
				Value: "purpose1",
			},
			Topic: SlackChannelSub{
				Value: "topic1",
			},
			Type: model.CHANNEL_GROUP,
		},
		{
			Id:      "id2",
			Name:    "channel-name-2",
			Creator: "creator2",
			Members: []string{"m1", "m2", "m3"},
			Purpose: SlackChannelSub{
				Value: "purpose2",
			},
			Topic: SlackChannelSub{
				Value: "topic2",
			},
			Type: model.CHANNEL_GROUP,
		},
		{
			Id:      "id3",
			Name:    "channel-name-3",
			Creator: "creator3",
			Members: []string{"m1", "m2", "m3"},
			Purpose: SlackChannelSub{
				Value: "purpose3",
			},
			Topic: SlackChannelSub{
				Value: "topic3",
			},
			Type: model.CHANNEL_GROUP,
		},
	}

	result := TransformChannels(regularGroupChannels, users)
	require.Len(t, result, len(regularGroupChannels))

	for i := range result {
		assert.Equal(t, fmt.Sprintf("channel-name-%d", i+1), result[i].Name)
		assert.Equal(t, fmt.Sprintf("channel-name-%d", i+1), result[i].DisplayName)
		assert.Equal(t, []string{"m1", "m2", "m3"}, result[i].Members)
		assert.Equal(t, fmt.Sprintf("purpose%d", i+1), result[i].Purpose)
		assert.Equal(t, fmt.Sprintf("topic%d", i+1), result[i].Header)
		assert.Equal(t, model.CHANNEL_GROUP, result[i].Type)
	}
}

func TestTransformDirectChannels(t *testing.T) {
	users := map[string]*IntermediateUser{"m1": {}, "m2": {}, "m3": {}}
	directChannels := []SlackChannel{
		{
			Id:      "id1",
			Creator: "creator1",
			Members: []string{"m1", "m2", "m3"},
			Type:    model.CHANNEL_DIRECT,
		},
		{
			Id:      "id2",
			Creator: "creator2",
			Members: []string{"m1", "m2", "m3"},
			Type:    model.CHANNEL_DIRECT,
		},
		{
			Id:      "id2",
			Creator: "creator2",
			Members: []string{"m1", "m2", "m3"},
			Type:    model.CHANNEL_DIRECT,
		},
	}

	result := TransformChannels(directChannels, users)
	require.Len(t, result, len(directChannels))

	for i := range result {
		assert.Equal(t, []string{"m1", "m2", "m3"}, result[i].Members)
		assert.Equal(t, model.CHANNEL_DIRECT, result[i].Type)
	}
}

func TestTransformChannelWithOneValidMember(t *testing.T) {
	users := map[string]*IntermediateUser{"m1": {}}
	t.Run("A direct channel with only one valid member should not be transformed", func(t *testing.T) {
		directChannels := []SlackChannel{
			{
				Id:      "id1",
				Creator: "creator1",
				Members: []string{"m1", "m2", "m3"},
				Type:    model.CHANNEL_DIRECT,
			},
		}

		result := TransformChannels(directChannels, users)
		require.Len(t, result, 0)
	})

	t.Run("A group channel with only one valid member should not be transformed", func(t *testing.T) {
		groupChannels := []SlackChannel{
			{
				Id:      "id1",
				Name:    "channel-name-1",
				Members: []string{"m1", "m2", "m3"},
				Purpose: SlackChannelSub{
					Value: "purpose1",
				},
				Topic: SlackChannelSub{
					Value: "topic1",
				},
				Type: model.CHANNEL_GROUP,
			},
		}

		result := TransformChannels(groupChannels, users)
		require.Len(t, result, 0)
	})
}

func TestIntermediateUserSanitise(t *testing.T) {
	t.Run("If there is no email, a placeholder should be used", func(t *testing.T) {
		user := IntermediateUser{
			Username: "test-username",
			Email:    "",
		}

		user.Sanitise()

		assert.Equal(t, "test-username@example.com", user.Email)
	})
}

func TestTransformUsers(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	id3 := "id3"

	intermediate := &Intermediate{}
	users := []SlackUser{
		{
			Id:       id1,
			Username: "username1",
			Profile: SlackProfile{
				FirstName: "firstname1",
				LastName:  "lastname1",
				Email:     "email1@example.com",
			},
		},
		{
			Id:       id2,
			Username: "username2",
			Profile: SlackProfile{
				FirstName: "firstname2",
				LastName:  "lastname2",
				Email:     "email2@example.com",
			},
		},
		{
			Id:       id3,
			Username: "username3",
			Profile: SlackProfile{
				FirstName: "firstname3",
				LastName:  "lastname3",
				Email:     "email3@example.com",
			},
		},
	}

	TransformUsers(users, intermediate)
	require.Len(t, intermediate.UsersById, len(users))

	for i, id := range []string{id1, id2, id3} {
		assert.Equal(t, fmt.Sprintf("id%d", i+1), intermediate.UsersById[id].Id)
		assert.Equal(t, fmt.Sprintf("username%d", i+1), intermediate.UsersById[id].Username)
		assert.Equal(t, fmt.Sprintf("firstname%d", i+1), intermediate.UsersById[id].FirstName)
		assert.Equal(t, fmt.Sprintf("lastname%d", i+1), intermediate.UsersById[id].LastName)
		assert.Equal(t, fmt.Sprintf("email%d@example.com", i+1), intermediate.UsersById[id].Email)
	}
}

func TestPopulateUserMemberships(t *testing.T) {
	intermediate := &Intermediate{
		UsersById: map[string]*IntermediateUser{"id1": {}, "id2": {}, "id3": {}},
		PublicChannels: []*IntermediateChannel{
			{
				Name:    "c1",
				Members: []string{"id1", "id3"},
			},
			{
				Name:    "c2",
				Members: []string{"id1", "id2"},
			},
		},
		PrivateChannels: []*IntermediateChannel{
			{
				Name:    "c3",
				Members: []string{"id3"},
			},
		},
	}

	PopulateUserMemberships(intermediate)

	assert.Equal(t, []string{"c1", "c2"}, intermediate.UsersById["id1"].Memberships)
	assert.Equal(t, []string{"c2"}, intermediate.UsersById["id2"].Memberships)
	assert.Equal(t, []string{"c1", "c3"}, intermediate.UsersById["id3"].Memberships)
}

func TestPopulateChannelMemberships(t *testing.T) {
	c1 := IntermediateChannel{
		Name:    "c1",
		Members: []string{"id1", "id3"},
	}
	c2 := IntermediateChannel{
		Name:    "c2",
		Members: []string{"id1", "id2"},
	}
	c3 := IntermediateChannel{
		Name:    "c3",
		Members: []string{"id3"},
	}

	intermediate := &Intermediate{
		UsersById: map[string]*IntermediateUser{
			"id1": {Username: "u1"},
			"id2": {Username: "u2"},
			"id3": {Username: "u3"},
		},
		GroupChannels:  []*IntermediateChannel{&c1, &c2},
		DirectChannels: []*IntermediateChannel{&c3},
	}

	PopulateChannelMemberships(intermediate)

	assert.Equal(t, []string{"u1", "u3"}, c1.MembersUsernames)
	assert.Equal(t, []string{"u1", "u2"}, c2.MembersUsernames)
	assert.Equal(t, []string{"u3"}, c3.MembersUsernames)
}
