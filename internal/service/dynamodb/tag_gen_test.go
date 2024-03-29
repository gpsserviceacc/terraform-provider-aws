// Code generated by internal/generate/tagresource/main.go; DO NOT EDIT.

package dynamodb_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-awsgps/internal/acctest"
	"terraform-provider-awsgps/internal/conns"
	tfdynamodb "terraform-provider-awsgps/internal/service/dynamodb"
	tftags "terraform-provider-awsgps/internal/tags"
	"terraform-provider-awsgps/internal/tfresource"
	"terraform-provider-awsgps/names"
)

func testAccCheckTagDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DynamoDBClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_dynamodb_tag" {
				continue
			}

			identifier, key, err := tftags.GetResourceID(rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = tfdynamodb.GetTag(ctx, conn, identifier, key)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("%s resource (%s) tag (%s) still exists", names.DynamoDB, identifier, key)
		}

		return nil
	}
}

func testAccCheckTagExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		identifier, key, err := tftags.GetResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DynamoDBClient(ctx)

		_, err = tfdynamodb.GetTag(ctx, conn, identifier, key)

		return err
	}
}
