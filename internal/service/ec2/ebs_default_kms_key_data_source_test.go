// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ec2_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-awsgps/internal/acctest"
	"terraform-provider-awsgps/names"
)

func TestAccEC2EBSDefaultKMSKeyDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEBSDefaultKMSKeyDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEBSDefaultKMSKey(ctx, "data.aws_ebs_default_kms_key.current"),
				),
			},
		},
	})
}

const testAccEBSDefaultKMSKeyDataSourceConfig_basic = `
data "aws_ebs_default_kms_key" "current" {}
`
