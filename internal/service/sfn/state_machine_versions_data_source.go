// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sfn

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	"terraform-provider-awsgps/internal/verify"
)

// @SDKDataSource("aws_sfn_state_machine_versions")
func DataSourceStateMachineVersions() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceStateMachineVersionsRead,

		Schema: map[string]*schema.Schema{
			"statemachine_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
			"statemachine_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceStateMachineVersionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SFNConn(ctx)

	smARN := d.Get("statemachine_arn").(string)
	input := &sfn.ListStateMachineVersionsInput{
		StateMachineArn: aws.String(smARN),
	}
	var smvARNs []string

	err := listStateMachineVersionsPages(ctx, conn, input, func(page *sfn.ListStateMachineVersionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.StateMachineVersions {
			if v != nil {
				smvARNs = append(smvARNs, aws.StringValue(v.StateMachineVersionArn))
			}
		}

		return !lastPage
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "listing Step Functions State Machine (%s) Versions: %s", smARN, err)
	}

	d.SetId(smARN)
	d.Set("statemachine_versions", smvARNs)

	return diags
}
