// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kms

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	"terraform-provider-awsgps/internal/flex"
	itypes "terraform-provider-awsgps/internal/types"
)

// @SDKDataSource("aws_kms_ciphertext")
func DataSourceCiphertext() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceCiphertextRead,

		Schema: map[string]*schema.Schema{
			"ciphertext_blob": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"context": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"plaintext": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceCiphertextRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).KMSConn(ctx)

	keyID := d.Get("key_id").(string)
	input := &kms.EncryptInput{
		KeyId:     aws.String(keyID),
		Plaintext: []byte(d.Get("plaintext").(string)),
	}

	if v, ok := d.GetOk("context"); ok && len(v.(map[string]interface{})) > 0 {
		input.EncryptionContext = flex.ExpandStringMap(v.(map[string]interface{}))
	}

	output, err := conn.EncryptWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "encrypting with KMS Key (%s): %s", keyID, err)
	}

	d.SetId(aws.StringValue(output.KeyId))
	d.Set("ciphertext_blob", itypes.Base64Encode(output.CiphertextBlob))

	return diags
}
