package slack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlert_importBasic(t *testing.T) {
	random := acctest.RandInt()

	testAccCheckSlackChannelImporterConfig := fmt.Sprintf(`
		resource "slack_channel" "channel" {
			name	= "test_import-%d"
		}
	`, random)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlackChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSlackChannelImporterConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlackChannelExists("slack_channel.channel"),
				),
			},
			{
				ResourceName:      "slack_channel.channel",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
