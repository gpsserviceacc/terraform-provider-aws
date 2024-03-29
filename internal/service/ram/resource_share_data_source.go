// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ram

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ram"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	"terraform-provider-awsgps/internal/flex"
	tfslices "terraform-provider-awsgps/internal/slices"
	tftags "terraform-provider-awsgps/internal/tags"
	"terraform-provider-awsgps/internal/tfresource"
)

// @SDKDataSource("aws_ram_resource_share", name="Resource Shared")
// @Tags
func dataSourceResourceShare() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceResourceShareRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"owning_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_arns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_owner": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(ram.ResourceOwner_Values(), false),
			},
			"resource_share_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(ram.ResourceShareStatus_Values(), false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceResourceShareRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).RAMConn(ctx)

	resourceOwner := d.Get("resource_owner").(string)
	inputG := &ram.GetResourceSharesInput{
		ResourceOwner: aws.String(resourceOwner),
	}

	if v, ok := d.GetOk("name"); ok {
		inputG.Name = aws.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok && v.(*schema.Set).Len() > 0 {
		inputG.TagFilters = expandTagFilters(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("resource_share_status"); ok {
		inputG.ResourceShareStatus = aws.String(v.(string))
	}

	share, err := findResourceShare(ctx, conn, inputG)

	if err != nil {
		return sdkdiag.AppendFromErr(diags, tfresource.SingularDataSourceFindError("RAM Resource Share", err))
	}

	arn := aws.StringValue(share.ResourceShareArn)
	d.SetId(arn)
	d.Set("arn", arn)
	d.Set("name", share.Name)
	d.Set("owning_account_id", share.OwningAccountId)
	d.Set("status", share.Status)

	setTagsOut(ctx, share.Tags)

	inputL := &ram.ListResourcesInput{
		ResourceOwner:     aws.String(resourceOwner),
		ResourceShareArns: aws.StringSlice([]string{arn}),
	}
	resources, err := findResources(ctx, conn, inputL)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading RAM Resource Share (%s) resources: %s", arn, err)
	}

	resourceARNs := tfslices.ApplyToAll(resources, func(r *ram.Resource) string {
		return aws.StringValue(r.Arn)
	})
	d.Set("resource_arns", resourceARNs)

	return diags
}

func expandTagFilter(tfMap map[string]interface{}) *ram.TagFilter {
	if tfMap == nil {
		return nil
	}

	apiObject := &ram.TagFilter{}

	if v, ok := tfMap["name"].(string); ok && v != "" {
		apiObject.TagKey = aws.String(v)
	}

	if v, ok := tfMap["values"].([]interface{}); ok && len(v) > 0 {
		apiObject.TagValues = flex.ExpandStringList(v)
	}

	return apiObject
}

func expandTagFilters(tfList []interface{}) []*ram.TagFilter {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*ram.TagFilter

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandTagFilter(tfMap)

		if apiObject == nil {
			continue
		}

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}
