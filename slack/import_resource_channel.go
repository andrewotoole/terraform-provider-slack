package slack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSlackChannelImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	_, _, err := parseTerraformId(d.Id())
	if err != nil {
		// Import by channelName
		channel, err := searchForChannelByName(d.Id(), m)
		if err != nil {
			return nil, err
		}
		terraformId, _ := createTerraformId(channel.ID, channel.Name)
		d.SetId(terraformId)

		return []*schema.ResourceData{d}, nil
	} else{
		// Import by slackId:channelName
		d.SetId(d.Id())

		return []*schema.ResourceData{d}, nil
	}

}