// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ram_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/ram"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-awsgps/internal/acctest"
	"terraform-provider-awsgps/internal/conns"
	tfram "terraform-provider-awsgps/internal/service/ram"
	"terraform-provider-awsgps/internal/tfresource"
	"terraform-provider-awsgps/names"
)

func TestAccRAMResourceAssociation_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var resourceShareAssociation ram.ResourceShareAssociation
	resourceName := "aws_ram_resource_association.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.RAMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAssociationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceAssociationExists(ctx, resourceName, &resourceShareAssociation),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRAMResourceAssociation_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var resourceShareAssociation ram.ResourceShareAssociation
	resourceName := "aws_ram_resource_association.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.RAMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAssociationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceAssociationExists(ctx, resourceName, &resourceShareAssociation),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfram.ResourceResourceAssociation(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRAMResourceAssociation_duplicate(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.RAMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourceAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceAssociationConfig_duplicate(rName),
				ExpectError: regexache.MustCompile(`RAM Resource Association .* already exists`),
			},
		},
	})
}

func testAccCheckResourceAssociationExists(ctx context.Context, n string, v *ram.ResourceShareAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).RAMConn(ctx)

		output, err := tfram.FindResourceAssociationByTwoPartKey(ctx, conn, rs.Primary.Attributes["resource_share_arn"], rs.Primary.Attributes["resource_arn"])

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckResourceAssociationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).RAMConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_ram_resource_association" {
				continue
			}

			_, err := tfram.FindResourceAssociationByTwoPartKey(ctx, conn, rs.Primary.Attributes["resource_share_arn"], rs.Primary.Attributes["resource_arn"])

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("RAM Resource Association %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccResourceAssociationConfig_basic(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 1), fmt.Sprintf(`
resource "aws_ram_resource_share" "test" {
  name = %[1]q
}

resource "aws_ram_resource_association" "test" {
  resource_arn       = aws_subnet.test[0].arn
  resource_share_arn = aws_ram_resource_share.test.id
}
`, rName))
}

func testAccResourceAssociationConfig_duplicate(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 1), fmt.Sprintf(`
resource "aws_ram_resource_share" "test" {
  name = %[1]q
}

resource "aws_ram_resource_association" "test1" {
  resource_arn       = aws_subnet.test[0].arn
  resource_share_arn = aws_ram_resource_share.test.id
}

resource "aws_ram_resource_association" "test2" {
  resource_arn       = aws_subnet.test[0].arn
  resource_share_arn = aws_ram_resource_association.test1.resource_share_arn
}
`, rName))
}
