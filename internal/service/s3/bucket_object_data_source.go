// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package s3

// WARNING: This code is DEPRECATED and will be removed in a future release!!
// DO NOT apply fixes or enhancements to the data source in this file.
// INSTEAD, apply fixes and enhancements to the data source in "object_data_source.go".

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	tftags "terraform-provider-awsgps/internal/tags"
)

// @SDKDataSource("aws_s3_bucket_object", name="Bucket Object")
func dataSourceBucketObject() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceBucketObjectRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"body": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bucket": {
				Deprecated: "Use the aws_s3_object data source instead",
				Type:       schema.TypeString,
				Required:   true,
			},
			"bucket_key_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cache_control": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_disposition": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_encoding": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_language": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"object_lock_legal_hold_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_lock_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_lock_retain_until_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"range": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_side_encryption": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sse_kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_class": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"version_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"website_redirect_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		DeprecationMessage: `use the aws_s3_object data source instead`,
	}
}

func dataSourceBucketObjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).S3Client(ctx)
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	bucket := d.Get("bucket").(string)
	key := sdkv1CompatibleCleanKey(d.Get("key").(string))
	input := &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	if v, ok := d.GetOk("range"); ok {
		input.Range = aws.String(v.(string))
	}
	if v, ok := d.GetOk("version_id"); ok {
		input.VersionId = aws.String(v.(string))
	}

	out, err := findObject(ctx, conn, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading S3 Bucket (%s) Object (%s): %s", bucket, key, err)
	}

	if aws.ToBool(out.DeleteMarker) {
		return sdkdiag.AppendErrorf(diags, "S3 Bucket (%s) Object (%s) has been deleted", bucket, key)
	}

	id := bucket + "/" + d.Get("key").(string)
	if v, ok := d.GetOk("version_id"); ok {
		id += "@" + v.(string)
	}
	d.SetId(id)

	arn, err := newObjectARN(meta.(*conns.AWSClient).Partition, bucket, key)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading S3 Bucket (%s) Object (%s): %s", bucket, key, err)
	}
	d.Set("arn", arn.String())

	d.Set("bucket_key_enabled", out.BucketKeyEnabled)
	d.Set("cache_control", out.CacheControl)
	d.Set("content_disposition", out.ContentDisposition)
	d.Set("content_encoding", out.ContentEncoding)
	d.Set("content_language", out.ContentLanguage)
	d.Set("content_length", out.ContentLength)
	d.Set("content_type", out.ContentType)
	// See https://forums.aws.amazon.com/thread.jspa?threadID=44003
	d.Set("etag", strings.Trim(aws.ToString(out.ETag), `"`))
	d.Set("expiration", out.Expiration)
	if out.Expires != nil {
		d.Set("expires", out.Expires.Format(time.RFC1123))
	} else {
		d.Set("expires", nil)
	}
	if out.LastModified != nil {
		d.Set("last_modified", out.LastModified.Format(time.RFC1123))
	} else {
		d.Set("last_modified", "")
	}
	d.Set("metadata", out.Metadata)
	d.Set("object_lock_legal_hold_status", out.ObjectLockLegalHoldStatus)
	d.Set("object_lock_mode", out.ObjectLockMode)
	d.Set("object_lock_retain_until_date", flattenObjectDate(out.ObjectLockRetainUntilDate))
	d.Set("server_side_encryption", out.ServerSideEncryption)
	d.Set("sse_kms_key_id", out.SSEKMSKeyId)
	// The "STANDARD" (which is also the default) storage
	// class when set would not be included in the results.
	d.Set("storage_class", types.ObjectStorageClassStandard)
	if out.StorageClass != "" {
		d.Set("storage_class", out.StorageClass)
	}
	d.Set("version_id", out.VersionId)
	d.Set("website_redirect_location", out.WebsiteRedirectLocation)

	if isContentTypeAllowed(out.ContentType) {
		input := &s3.GetObjectInput{
			Bucket:    aws.String(bucket),
			Key:       aws.String(key),
			VersionId: out.VersionId,
		}
		if v, ok := d.GetOk("range"); ok {
			input.Range = aws.String(v.(string))
		}

		out, err := conn.GetObject(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "downloading S3 Bucket (%s) Object (%s): %s", bucket, key, err)
		}

		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(out.Body); err != nil {
			return sdkdiag.AppendFromErr(diags, err)
		}

		d.Set("body", buf.String())
	}

	tags, err := objectListTags(ctx, conn, bucket, key)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "listing tags for S3 Bucket (%s) Object (%s): %s", bucket, key, err)
	}

	if err := d.Set("tags", tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	return diags
}
