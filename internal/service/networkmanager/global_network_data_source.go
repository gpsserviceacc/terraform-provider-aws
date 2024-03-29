// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package networkmanager

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	tftags "terraform-provider-awsgps/internal/tags"
)

// @SDKDataSource("aws_networkmanager_global_network")
func DataSourceGlobalNetwork() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceGlobalNetworkRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_network_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceGlobalNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).NetworkManagerConn(ctx)
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	globalNetworkID := d.Get("global_network_id").(string)
	globalNetwork, err := FindGlobalNetworkByID(ctx, conn, globalNetworkID)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Network Manager Global Network (%s): %s", globalNetworkID, err)
	}

	d.SetId(globalNetworkID)
	d.Set("arn", globalNetwork.GlobalNetworkArn)
	d.Set("description", globalNetwork.Description)
	d.Set("global_network_id", globalNetwork.GlobalNetworkId)

	if err := d.Set("tags", KeyValueTags(ctx, globalNetwork.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	return diags
}
