# terraform-provider-slack 

Terraform provider for [slack](https://slack.com/)

## Requirements
-	[Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.13

## Installation

See the [the Provider Configuration page of the Terraform documentation](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for instructions.

Pre-compiled binaries are available from the [Releases](https://github.com/andrewotoole/terraform-provider-slack/releases) page.

## Development

### Test

Test the provider by running `make test`.

Make sure to set the following environment variables:

- `SENTRY_TEST_ORGANIZATION`
- `SENTRY_TOKEN`

### Build

See the [Writing Custom Providers page of the Terraform documentation](https://www.terraform.io/docs/extend/writing-custom-providers.html#building-the-plugin) for instructions.
