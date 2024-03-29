// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package medialive_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/medialive"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-awsgps/internal/acctest"
	"terraform-provider-awsgps/names"
)

func TestAccMediaLiveInputDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var input medialive.DescribeInputOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_medialive_input.test"
	dataSourceName := "data.aws_medialive_input.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.MediaLiveEndpointID)
			testAccInputsPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.MediaLiveServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckInputDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccInputDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInputExists(ctx, dataSourceName, &input),
					resource.TestCheckResourceAttrPair(resourceName, "arn", dataSourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttr(dataSourceName, "destinations.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "input_class", dataSourceName, "input_class"),
					resource.TestCheckResourceAttrPair(resourceName, "input_devices", dataSourceName, "input_devices"),
					resource.TestCheckResourceAttrPair(resourceName, "input_partner_ids", dataSourceName, "input_partner_ids"),
					resource.TestCheckResourceAttrPair(resourceName, "input_source_type", dataSourceName, "input_source_type"),
					resource.TestCheckResourceAttrPair(resourceName, "security_groups", dataSourceName, "security_groups"),
					resource.TestCheckResourceAttrPair(resourceName, "sources", dataSourceName, "sources"),
					resource.TestCheckResourceAttrSet(dataSourceName, "state"),
					resource.TestCheckResourceAttrPair(resourceName, "tag_all", dataSourceName, "tags"),
					resource.TestCheckResourceAttrPair(resourceName, "type", dataSourceName, "type"),
				),
			},
		},
	})
}

func testAccInputDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_medialive_input_security_group" "test" {
  whitelist_rules {
    cidr = "10.0.0.8/32"
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_medialive_input" "test" {
  name                  = %[1]q
  input_security_groups = [aws_medialive_input_security_group.test.id]
  type                  = "UDP_PUSH"

  tags = {
    Name = %[1]q
  }
}

data "aws_medialive_input" "test" {
  id = aws_medialive_input.test.id
}
`, rName)
}
