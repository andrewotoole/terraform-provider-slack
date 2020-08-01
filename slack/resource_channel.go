package slack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/slack-go/slack"
)

func resourceSlackChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateSlackChannel,
		Read:   resourceReadSlackChannel,
		Update: resourceUpdateSlackChannel,
		Delete: resourceDeleteSlackChannel,
		Exists: resourceExistsSlackChannel,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the Slack channel",
				ForceNew:     true,
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
		},
	}
}

func resourceCreateSlackChannel(d *schema.ResourceData, m interface{}) error {
	apiClient := slack.New(m.(*Config).Token)

	channel, err := apiClient.CreateConversation(d.Get("name").(string), false)

	if err != nil {
		return err
	}

	d.SetId(channel.ID)

	if _, err := apiClient.SetTopicOfConversation(d.Id(), d.Get("topic").(string)); err != nil {
		return err
	}

	if _, err := apiClient.SetPurposeOfConversation(d.Id(), d.Get("purpose").(string)); err != nil {
		return err
	}

	return nil
}

func resourceReadSlackChannel(d *schema.ResourceData, m interface{}) error {
	apiClient := slack.New(m.(*Config).Token)

	channel, err := apiClient.GetConversationInfo(d.Id(), false)
	if err != nil {
		d.SetId("")
		return nil
	}

	d.SetId(channel.ID)
	d.Set("name", channel.Name)
	d.Set("topic", channel.Topic)
	d.Set("purpose", channel.Purpose)

	return nil
}

func resourceUpdateSlackChannel(d *schema.ResourceData, m interface{}) error {
	apiClient := slack.New(m.(*Config).Token)

	if _, err := apiClient.RenameConversation(d.Id(), d.Get("name").(string)); err != nil {
		return err
	}

	if _, err := apiClient.SetTopicOfConversation(d.Id(), d.Get("topic").(string)); err != nil {
		return err
	}

	if _, err := apiClient.SetPurposeOfConversation(d.Id(), d.Get("purpose").(string)); err != nil {
		return err
	}

	return nil
}

func resourceDeleteSlackChannel(d *schema.ResourceData, m interface{}) error {
	apiClient := slack.New(m.(*Config).Token)

	if err := apiClient.ArchiveConversation(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func resourceExistsSlackChannel(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := slack.New(m.(*Config).Token)

	_, err := apiClient.GetConversationInfo(d.Id(), false)
	if err != nil {
		return false, nil
	}
	return true, nil
}
