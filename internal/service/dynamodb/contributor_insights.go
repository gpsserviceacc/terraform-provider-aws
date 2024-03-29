// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dynamodb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	"terraform-provider-awsgps/internal/tfresource"
)

// @SDKResource("aws_dynamodb_contributor_insights")
func ResourceContributorInsights() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceContributorInsightsCreate,
		ReadWithoutTimeout:   resourceContributorInsightsRead,
		DeleteWithoutTimeout: resourceContributorInsightsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"index_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"table_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceContributorInsightsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).DynamoDBConn(ctx)

	input := &dynamodb.UpdateContributorInsightsInput{
		ContributorInsightsAction: aws.String(dynamodb.ContributorInsightsActionEnable),
	}

	if v, ok := d.GetOk("table_name"); ok {
		input.TableName = aws.String(v.(string))
	}

	var indexName string
	if v, ok := d.GetOk("index_name"); ok {
		input.IndexName = aws.String(v.(string))
		indexName = v.(string)
	}

	output, err := conn.UpdateContributorInsightsWithContext(ctx, input)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating DynamoDB ContributorInsights for table (%s): %s", d.Get("table_name").(string), err)
	}

	id := EncodeContributorInsightsID(aws.StringValue(output.TableName), indexName, meta.(*conns.AWSClient).AccountID)
	d.SetId(id)

	if err := waitContributorInsightsCreated(ctx, conn, aws.StringValue(output.TableName), indexName, d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for DynamoDB ContributorInsights (%s) create: %s", d.Id(), err)
	}

	return append(diags, resourceContributorInsightsRead(ctx, d, meta)...)
}

func resourceContributorInsightsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).DynamoDBConn(ctx)

	tableName, indexName, err := DecodeContributorInsightsID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "unable to decode DynamoDB ContributorInsights ID (%s): %s", d.Id(), err)
	}

	out, err := FindContributorInsights(ctx, conn, tableName, indexName)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] DynamoDB ContributorInsights (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading DynamoDB ContributorInsights (%s): %s", d.Id(), err)
	}

	d.Set("index_name", out.IndexName)
	d.Set("table_name", out.TableName)

	return diags
}

func resourceContributorInsightsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).DynamoDBConn(ctx)

	log.Printf("[INFO] Deleting DynamoDB ContributorInsights %s", d.Id())

	tableName, indexName, err := DecodeContributorInsightsID(d.Id())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "unable to decode DynamoDB ContributorInsights ID (%s): %s", d.Id(), err)
	}

	input := &dynamodb.UpdateContributorInsightsInput{
		ContributorInsightsAction: aws.String(dynamodb.ContributorInsightsActionDisable),
		TableName:                 aws.String(tableName),
	}

	if indexName != "" {
		input.IndexName = aws.String(indexName)
	}

	_, err = conn.UpdateContributorInsightsWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, dynamodb.ErrCodeResourceNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting DynamoDB ContributorInsights (%s): %s", d.Id(), err)
	}

	if err := waitContributorInsightsDeleted(ctx, conn, tableName, indexName, d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for DynamoDB ContributorInsights (%s) to be deleted: %s", d.Id(), err)
	}

	return diags
}

func EncodeContributorInsightsID(tableName, indexName, accountID string) string {
	return fmt.Sprintf("name:%s/index:%s/%s", tableName, indexName, accountID)
}

func DecodeContributorInsightsID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 3 || idParts[0] == "" || idParts[2] == "" {
		return "", "", fmt.Errorf("expected ID in the form of table_name/account_id, given: %q", id)
	}

	tableName := strings.TrimPrefix(idParts[0], "name:")
	indexName := strings.TrimPrefix(idParts[1], "index:")

	return tableName, indexName, nil
}
