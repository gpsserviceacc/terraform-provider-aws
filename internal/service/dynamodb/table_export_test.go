// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynamodb_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-awsgps/internal/acctest"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/create"
	tfdynamodb "terraform-provider-awsgps/internal/service/dynamodb"
	"terraform-provider-awsgps/names"
)

func TestAccDynamoDBTableExport_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var tableexport dynamodb.DescribeExportOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_dynamodb_table_export.test"
	s3BucketResourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				Config: testAccTableExportConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableExportExists(ctx, resourceName, &tableexport),
					resource.TestCheckResourceAttr(resourceName, "export_format", "DYNAMODB_JSON"),
					resource.TestCheckResourceAttr(resourceName, "export_status", "COMPLETED"),
					resource.TestCheckResourceAttr(resourceName, "item_count", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "s3_bucket", s3BucketResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "s3_bucket_owner", ""),
					resource.TestCheckResourceAttr(resourceName, "s3_prefix", ""),
					resource.TestCheckResourceAttr(resourceName, "s3_sse_algorithm", "AES256"),
					resource.TestCheckResourceAttr(resourceName, "s3_sse_kms_key_id", ""),
					resource.TestCheckResourceAttrSet(resourceName, "manifest_files_s3_key"),
					resource.TestCheckResourceAttrSet(resourceName, "export_time"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "end_time"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "dynamodb", regexache.MustCompile(
						fmt.Sprintf("table\\/%s\\/export\\/+.", rName),
					)),
					acctest.CheckResourceAttrRegionalARN(resourceName, "table_arn", "dynamodb", fmt.Sprintf("table/%s", rName)),
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

func TestAccDynamoDBTableExport_kms(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var tableexport dynamodb.DescribeExportOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_dynamodb_table_export.test"
	s3BucketResourceName := "aws_s3_bucket.test"
	kmsKeyResourceName := "aws_kms_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.DynamoDB)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				Config: testAccTableExportConfig_kms(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableExportExists(ctx, resourceName, &tableexport),
					resource.TestCheckResourceAttr(resourceName, "export_format", "DYNAMODB_JSON"),
					resource.TestCheckResourceAttr(resourceName, "export_status", "COMPLETED"),
					resource.TestCheckResourceAttr(resourceName, "item_count", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "s3_bucket", s3BucketResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "s3_bucket_owner", ""),
					resource.TestCheckResourceAttr(resourceName, "s3_prefix", ""),
					resource.TestCheckResourceAttr(resourceName, "s3_sse_algorithm", "KMS"),
					resource.TestCheckResourceAttrPair(resourceName, "s3_sse_kms_key_id", kmsKeyResourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "manifest_files_s3_key"),
					resource.TestCheckResourceAttrSet(resourceName, "export_time"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "end_time"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "dynamodb", regexache.MustCompile(
						fmt.Sprintf("table\\/%s\\/export\\/+.", rName),
					)),
					acctest.CheckResourceAttrRegionalARN(resourceName, "table_arn", "dynamodb", fmt.Sprintf("table/%s", rName)),
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

func TestAccDynamoDBTableExport_s3Prefix(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var tableexport dynamodb.DescribeExportOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_dynamodb_table_export.test"
	s3BucketResourceName := "aws_s3_bucket.test"
	s3BucketPrefix := "test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.DynamoDB)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.DynamoDBServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             acctest.CheckDestroyNoop,
		Steps: []resource.TestStep{
			{
				Config: testAccTableExportConfig_s3Prefix(rName, s3BucketPrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTableExportExists(ctx, resourceName, &tableexport),
					resource.TestCheckResourceAttr(resourceName, "export_format", "DYNAMODB_JSON"),
					resource.TestCheckResourceAttr(resourceName, "export_status", "COMPLETED"),
					resource.TestCheckResourceAttr(resourceName, "item_count", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "s3_bucket", s3BucketResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "s3_bucket_owner", ""),
					resource.TestCheckResourceAttr(resourceName, "s3_prefix", "test"),
					resource.TestCheckResourceAttr(resourceName, "s3_sse_algorithm", "AES256"),
					resource.TestCheckResourceAttr(resourceName, "s3_sse_kms_key_id", ""),
					resource.TestCheckResourceAttrSet(resourceName, "manifest_files_s3_key"),
					resource.TestCheckResourceAttrSet(resourceName, "export_time"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "end_time"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "dynamodb", regexache.MustCompile(
						fmt.Sprintf("table\\/%s\\/export\\/+.", rName),
					)),
					acctest.CheckResourceAttrRegionalARN(resourceName, "table_arn", "dynamodb", fmt.Sprintf("table/%s", rName)),
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

func testAccCheckTableExportExists(ctx context.Context, name string, tableexport *dynamodb.DescribeExportOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.DynamoDB, create.ErrActionCheckingExistence, tfdynamodb.ResNameTableExport, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.DynamoDB, create.ErrActionCheckingExistence, tfdynamodb.ResNameTableExport, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DynamoDBConn(ctx)
		resp, err := tfdynamodb.FindTableExportByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return create.Error(names.DynamoDB, create.ErrActionCheckingExistence, tfdynamodb.ResNameTableExport, rs.Primary.ID, err)
		}

		*tableexport = *resp

		return nil
	}
}

func testAccPreCheck(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).DynamoDBConn(ctx)

	input := &dynamodb.ListExportsInput{}
	_, err := conn.ListExportsWithContext(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccTableExportConfig_baseConfig(tableName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_dynamodb_table" "test" {
  name           = %[1]q
  read_capacity  = 2
  write_capacity = 2
  hash_key       = "TestTableHashKey"

  attribute {
    name = "TestTableHashKey"
    type = "S"
  }

  point_in_time_recovery {
    enabled = true
  }
}
`, tableName)
}

func testAccTableExportConfig_basic(tableName string) string {
	return acctest.ConfigCompose(testAccTableExportConfig_baseConfig(tableName), ` 
resource "aws_dynamodb_table_export" "test" {
  s3_bucket = aws_s3_bucket.test.id
  table_arn = aws_dynamodb_table.test.arn
}
`)
}

func testAccTableExportConfig_kms(tableName string) string {
	return acctest.ConfigCompose(testAccTableExportConfig_baseConfig(tableName), `
resource "aws_kms_key" "test" {
  deletion_window_in_days = 7
}

resource "aws_dynamodb_table_export" "test" {
  s3_bucket         = aws_s3_bucket.test.id
  s3_sse_kms_key_id = aws_kms_key.test.id
  s3_sse_algorithm  = "KMS"
  table_arn         = aws_dynamodb_table.test.arn
}`)
}

func testAccTableExportConfig_s3Prefix(tableName, s3BucketPrefix string) string {
	return acctest.ConfigCompose(testAccTableExportConfig_baseConfig(tableName), fmt.Sprintf(`
resource "aws_dynamodb_table_export" "test" {
  s3_bucket        = aws_s3_bucket.test.id
  s3_prefix        = %[1]q
  s3_sse_algorithm = "AES256"
  table_arn        = aws_dynamodb_table.test.arn
}`, s3BucketPrefix))
}
