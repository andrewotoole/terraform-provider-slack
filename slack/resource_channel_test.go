package slack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/slack-go/slack"
)

func TestAccSlackChannel_Basic(t *testing.T) {
	random := acctest.RandInt()

	testAccSlackChannelConfig := fmt.Sprintf(`
		resource "slack_channel" "test_channel" {
			name	= "test-%d"
		}
	`, random)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlackChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlackChannelConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlackChannelExists("slack_channel.test_channel"),
					resource.TestCheckResourceAttr(
						"slack_channel.test_channel", "name", fmt.Sprintf("test-%d", random)),
				),
			},
		},
	})
}

func TestAccSlackChannel_Update(t *testing.T) {
	random := acctest.RandInt()

	testAccSlackChannelUpdatePre := fmt.Sprintf(`
		resource "slack_channel" "test_update" {
			name    	= "name_original-%d"
			topic		= "topic_original"
			purpose 	= "purpose_original"
		}
	`, random)

	testAccSlackChannelUpdatePost := fmt.Sprintf(`
		resource "slack_channel" "test_update" {
			name    	= "name_updated-%d"
			topic		= "topic_updated"
			purpose 	= "purpose_updated"
		}
	`, random)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSlackChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlackChannelUpdatePre,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlackChannelExists("slack_channel.test_update"),
					resource.TestCheckResourceAttr(
						"slack_channel.test_update", "name", fmt.Sprintf("name_original-%d", random)),
					resource.TestCheckResourceAttr(
						"slack_channel.test_update", "topic", "topic_original"),
					resource.TestCheckResourceAttr(
						"slack_channel.test_update", "purpose", "purpose_original"),
				),
			},
			{
				Config: testAccSlackChannelUpdatePost,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlackChannelExists("slack_channel.test_update"),
					resource.TestCheckResourceAttr(
						"slack_channel.test_update", "name", fmt.Sprintf("name_updated-%d", random)),
					resource.TestCheckResourceAttr(
						"slack_channel.test_update", "topic", "topic_updated"),
					resource.TestCheckResourceAttr(
						"slack_channel.test_update", "purpose", "purpose_updated"),
				),
			},
		},
	})
}

func testAccCheckSlackChannelDestroy(s *terraform.State) error {
	apiClient := slack.New(testAccProvider.Meta().(*Config).Token)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "slack_channel" {
			continue
		}

		channel, err := apiClient.GetConversationInfo(rs.Primary.ID, false)
		if err == nil {
			if !channel.IsArchived {
				return fmt.Errorf("channel still exists")
			}
		}
		return nil
	}
	return nil
}

func testAccCheckSlackChannelExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := slack.New(testAccProvider.Meta().(*Config).Token)
		_, err := apiClient.GetConversationInfo(rs.Primary.ID, false)
		if err != nil {
			return fmt.Errorf("error fetching channel with resource %s. %s", resource, err)
		}
		return nil
	}
}
