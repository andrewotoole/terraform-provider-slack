package slack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/slack-go/slack"
	"strings"
)

func resourceSlackChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateSlackChannel,
		Read:   resourceReadSlackChannel,
		Update: resourceUpdateSlackChannel,
		Delete: resourceDeleteSlackChannel,
		Exists: resourceExistsSlackChannel,
		Importer: &schema.ResourceImporter{
			State: resourceSlackChannelImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the Slack channel",
				ValidateFunc: validation.StringLenBetween(1, 80),
			},
			"topic": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The topic of the Slack channel",
			},
			"purpose": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The purpose of the Slack channel",
			},
			"is_archived": {
				Type: 		 schema.TypeBool,
				Optional: 	 true,
				Default: 	 false,
				Description: "Determines if a channel is archived",
			},
		},
	}
}

func resourceCreateSlackChannel(d *schema.ResourceData, m interface{}) error {
	apiClient := slack.New(m.(*Config).Token)

	channel, err := apiClient.CreateConversation(d.Get("name").(string), false)
	if err != nil {
		return err
	}

	terraformId, err := createTerraformId(channel.ID, channel.Name)
	if err != nil {
		return err
	}
	d.SetId(terraformId)

	if _, err := apiClient.SetTopicOfConversation(channel.ID, d.Get("topic").(string)); err != nil {
		return err
	}

	if _, err := apiClient.SetPurposeOfConversation(channel.ID, d.Get("purpose").(string)); err != nil {
		return err
	}

	if d.Get("is_archived").(bool) {
		if err := apiClient.ArchiveConversation(channel.ID); err != nil {
			return err
		}
	}

	return nil
}

func resourceReadSlackChannel(d *schema.ResourceData, m interface{}) error {
	channel, err := getSlackChannelFromTerraformId(d.Id(), m); if err != nil {
		d.SetId("")
		return err
	}

	terraformId, err := createTerraformId(channel.ID, channel.Name)
	if err != nil {
		return err
	}
	d.SetId(terraformId)
	d.Set("name", channel.Name)
	d.Set("topic", channel.Topic)
	d.Set("purpose", channel.Purpose)
	d.Set("is_archived", channel.IsArchived)

	return nil
}

func resourceUpdateSlackChannel(d *schema.ResourceData, m interface{}) error {
	apiClient := slack.New(m.(*Config).Token)

	channel, err := getSlackChannelFromTerraformId(d.Id(), m); if err != nil {
		d.SetId("")
		return err
	}

	if channel.Name != d.Get("name").(string) {
		if _, err := apiClient.RenameConversation(channel.ID, d.Get("name").(string)); err != nil {
			return err
		}

		// Update terraform ID if name changes
		newId, err := createTerraformId(channel.ID, d.Get("name").(string))
		if err != nil {
			return err
		}
		d.SetId(newId)
	}

	if channel.Topic.Value != d.Get("topic").(string) {
		if _, err := apiClient.SetTopicOfConversation(channel.ID, d.Get("topic").(string)); err != nil {
			return err
		}
	}

	if channel.Purpose.Value != d.Get("purpose").(string) {
		if _, err := apiClient.SetPurposeOfConversation(channel.ID, d.Get("purpose").(string)); err != nil {
			return err
		}
	}

	if channel.IsArchived != d.Get("is_archived").(bool) {
		if d.Get("is_archived").(bool) {
			if err := apiClient.ArchiveConversation(channel.ID); err != nil {
				return err
			}
		} else {
			if err := apiClient.UnArchiveConversation(channel.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceDeleteSlackChannel(d *schema.ResourceData, m interface{}) error {
	apiClient := slack.New(m.(*Config).Token)

	channel, err := getSlackChannelFromTerraformId(d.Id(), m); if err != nil {
		return err
	}

	// no way to truly "delete" a conversation, so archiving it is the next best thing
	if !d.Get("is_archived").(bool) {
		err = apiClient.ArchiveConversation(channel.ID)
		return err
	}
	return nil
}

func resourceExistsSlackChannel(d *schema.ResourceData, m interface{}) (bool, error) {
	_, err := getSlackChannelFromTerraformId(d.Id(), m); if err != nil {
		return false, nil
	}
	return true, nil
}

/***************************************
// These function is needed due to the issue outlined here: https://api.slack.com/docs/conversations-api
//
// 		🚧 Channel IDs can become unstable in certain situations
// 		There are a few circumstances channel IDs might change within a workspace. If you
// 		can operate without depending on their stability, you'll be well-prepared for
// 		unfortunate hijinks.
//
// 		In the future, we'll mitigate this unexpected transition with appropriate Events API
// 		events or other solutions.
//
// 		In the meantime, be aware this might happen and use conversations.list regularly to
// 		monitor change for known #channel names if ID stability is important to you.
//
// To solve for this, the terraform ID is composed of slackId:channelName
****************************************/
func getSlackChannelFromTerraformId(terraformId string, m interface{}) (slack.Channel, error) {
	_, name, err := parseTerraformId(terraformId)
	if err != nil {
		return slack.Channel{}, err
	}

	return searchForChannelByName(name, m)
}

func searchForChannelByName(name string, m interface{}) (slack.Channel, error) {
	apiClient := slack.New(m.(*Config).Token)

	params := slack.GetConversationsParameters{
		ExcludeArchived: "false",
	}
	cursor := "start"
	for cursor != "" {
		channels, nextCursor, err := apiClient.GetConversations(&params)
		if err != nil {
			return slack.Channel{}, err
		}
		for _, c := range channels {
			if c.Name == name {
				return c, nil
			}
		}
		cursor = nextCursor
		params.Cursor = cursor
	}
	return slack.Channel{}, fmt.Errorf("not_found")
}

func parseTerraformId(terraformId string) (string, string, error) {
	parts := strings.Split(terraformId, ":")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed_id")
	}

	return parts[0], parts[1], nil
}

func createTerraformId(slackId string, name string) (string, error) {
	if slackId == "" || name == "" {
		return "", fmt.Errorf("malformed_properties - slack_id:%s | name:%s", slackId, name)
	}
	return fmt.Sprintf("%s:%s", slackId, name), nil
}
