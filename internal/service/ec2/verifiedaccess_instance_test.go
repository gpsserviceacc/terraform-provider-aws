// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ec2_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-awsgps/internal/acctest"
	"terraform-provider-awsgps/internal/conns"
	tfsync "terraform-provider-awsgps/internal/experimental/sync"
	tfec2 "terraform-provider-awsgps/internal/service/ec2"
	"terraform-provider-awsgps/internal/tfresource"
	"terraform-provider-awsgps/names"
)

func testAccVerifiedAccessInstance_basic(t *testing.T, semaphore tfsync.Semaphore) {
	ctx := acctest.Context(t)
	var v types.VerifiedAccessInstance
	resourceName := "aws_verifiedaccess_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVerifiedAccessSynchronize(t, semaphore)
			acctest.PreCheck(ctx, t)
			testAccPreCheckVerifiedAccessInstance(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVerifiedAccessInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifiedAccessInstanceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVerifiedAccessInstanceExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrSet(resourceName, "last_updated_time"),
					resource.TestCheckResourceAttr(resourceName, "verified_access_trust_providers.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccVerifiedAccessInstance_description(t *testing.T, semaphore tfsync.Semaphore) {
	ctx := acctest.Context(t)
	var v1, v2 types.VerifiedAccessInstance
	resourceName := "aws_verifiedaccess_instance.test"
	originalDescription := "original description"
	updatedDescription := "updated description"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVerifiedAccessSynchronize(t, semaphore)
			acctest.PreCheck(ctx, t)
			testAccPreCheckVerifiedAccessInstance(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVerifiedAccessInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifiedAccessInstanceConfig_description(originalDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVerifiedAccessInstanceExists(ctx, resourceName, &v1),
					resource.TestCheckResourceAttr(resourceName, "description", originalDescription),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccVerifiedAccessInstanceConfig_description(updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVerifiedAccessInstanceExists(ctx, resourceName, &v2),
					testAccCheckVerifiedAccessInstanceNotRecreated(&v1, &v2),
					resource.TestCheckResourceAttr(resourceName, "description", updatedDescription),
				),
			},
		},
	})
}

func testAccVerifiedAccessInstance_fipsEnabled(t *testing.T, semaphore tfsync.Semaphore) {
	ctx := acctest.Context(t)
	var v1, v2 types.VerifiedAccessInstance
	resourceName := "aws_verifiedaccess_instance.test"
	originalFipsEnabled := true
	updatedFipsEnabled := false

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVerifiedAccessSynchronize(t, semaphore)
			acctest.PreCheck(ctx, t)
			testAccPreCheckVerifiedAccessInstance(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVerifiedAccessInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifiedAccessInstanceConfig_fipsEnabled(originalFipsEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVerifiedAccessInstanceExists(ctx, resourceName, &v1),
					resource.TestCheckResourceAttr(resourceName, "fips_enabled", strconv.FormatBool(originalFipsEnabled)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccVerifiedAccessInstanceConfig_fipsEnabled(updatedFipsEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVerifiedAccessInstanceExists(ctx, resourceName, &v2),
					testAccCheckVerifiedAccessInstanceRecreated(&v1, &v2),
					resource.TestCheckResourceAttr(resourceName, "fips_enabled", strconv.FormatBool(updatedFipsEnabled)),
				),
			},
		},
	})
}

func testAccVerifiedAccessInstance_disappears(t *testing.T, semaphore tfsync.Semaphore) {
	ctx := acctest.Context(t)
	var v types.VerifiedAccessInstance
	resourceName := "aws_verifiedaccess_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVerifiedAccessSynchronize(t, semaphore)
			acctest.PreCheck(ctx, t)
			testAccPreCheckVerifiedAccessInstance(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVerifiedAccessInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifiedAccessInstanceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVerifiedAccessInstanceExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfec2.ResourceVerifiedAccessInstance(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccVerifiedAccessInstance_tags(t *testing.T, semaphore tfsync.Semaphore) {
	ctx := acctest.Context(t)
	var v1, v2, v3 types.VerifiedAccessInstance
	resourceName := "aws_verifiedaccess_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVerifiedAccessSynchronize(t, semaphore)
			acctest.PreCheck(ctx, t)
			testAccPreCheckVerifiedAccessInstance(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVerifiedAccessInstanceDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVerifiedAccessInstanceConfig_tags1("key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVerifiedAccessInstanceExists(ctx, resourceName, &v1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccVerifiedAccessInstanceConfig_tags2("key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVerifiedAccessInstanceExists(ctx, resourceName, &v2),
					testAccCheckVerifiedAccessInstanceNotRecreated(&v1, &v2),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccVerifiedAccessInstanceConfig_tags1("key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVerifiedAccessInstanceExists(ctx, resourceName, &v3),
					testAccCheckVerifiedAccessInstanceNotRecreated(&v2, &v3),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccCheckVerifiedAccessInstanceNotRecreated(before, after *types.VerifiedAccessInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.ToString(before.VerifiedAccessInstanceId), aws.ToString(after.VerifiedAccessInstanceId); before != after {
			return fmt.Errorf("Verified Access Instance (%s/%s) recreated", before, after)
		}

		return nil
	}
}

func testAccCheckVerifiedAccessInstanceRecreated(before, after *types.VerifiedAccessInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.ToString(before.VerifiedAccessInstanceId), aws.ToString(after.VerifiedAccessInstanceId); before == after {
			return fmt.Errorf("Verified Access Instance (%s) not recreated", before)
		}

		return nil
	}
}

func testAccCheckVerifiedAccessInstanceExists(ctx context.Context, n string, v *types.VerifiedAccessInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Client(ctx)

		output, err := tfec2.FindVerifiedAccessInstanceByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckVerifiedAccessInstanceDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Client(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_verifiedaccess_instance" {
				continue
			}

			_, err := tfec2.FindVerifiedAccessInstanceByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Verified Access Instance %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccPreCheckVerifiedAccessInstance(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Client(ctx)

	input := &ec2.DescribeVerifiedAccessInstancesInput{}
	_, err := conn.DescribeVerifiedAccessInstances(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}
	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccVerifiedAccessInstanceConfig_basic() string {
	return `
resource "aws_verifiedaccess_instance" "test" {}
`
}

func testAccVerifiedAccessInstanceConfig_description(description string) string {
	return fmt.Sprintf(`
resource "aws_verifiedaccess_instance" "test" {
  description = %[1]q
}
`, description)
}

func testAccVerifiedAccessInstanceConfig_fipsEnabled(fipsEnabled bool) string {
	return fmt.Sprintf(`
resource "aws_verifiedaccess_instance" "test" {
  fips_enabled = %[1]t
}
`, fipsEnabled)
}

func testAccVerifiedAccessInstanceConfig_tags1(tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_verifiedaccess_instance" "test" {
  tags = {
    %[1]q = %[2]q
  }
}
`, tagKey1, tagValue1)
}

func testAccVerifiedAccessInstanceConfig_tags2(tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_verifiedaccess_instance" "test" {
  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }
}
`, tagKey1, tagValue1, tagKey2, tagValue2)
}
