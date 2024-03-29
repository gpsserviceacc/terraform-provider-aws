// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package codecommit

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/aws/aws-sdk-go-v2/service/codecommit/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/errs"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	tftags "terraform-provider-awsgps/internal/tags"
	"terraform-provider-awsgps/internal/tfresource"
	"terraform-provider-awsgps/internal/verify"
	"terraform-provider-awsgps/names"
)

// @SDKResource("aws_codecommit_repository", name="Repository")
// @Tags(identifierAttribute="arn")
func resourceRepository() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceRepositoryCreate,
		UpdateWithoutTimeout: resourceRepositoryUpdate,
		ReadWithoutTimeout:   resourceRepositoryRead,
		DeleteWithoutTimeout: resourceRepositoryDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clone_url_http": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"clone_url_ssh": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_branch": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"kms_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidARN,
			},
			"repository_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repository_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 100),
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceRepositoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CodeCommitClient(ctx)

	name := d.Get("repository_name").(string)
	input := &codecommit.CreateRepositoryInput{
		RepositoryName: aws.String(name),
		Tags:           getTagsIn(ctx),
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		input.KmsKeyId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		input.RepositoryDescription = aws.String(v.(string))
	}

	_, err := conn.CreateRepository(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating CodeCommit Repository (%s): %s", name, err)
	}

	d.SetId(name)

	if v, ok := d.GetOk("default_branch"); ok {
		if err := updateRepositoryDefaultBranch(ctx, conn, d.Id(), v.(string)); err != nil {
			return sdkdiag.AppendFromErr(diags, err)
		}
	}

	return append(diags, resourceRepositoryRead(ctx, d, meta)...)
}

func resourceRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CodeCommitClient(ctx)

	repository, err := findRepositoryByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] CodeCommit Repository %s not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading CodeCommit Repository (%s): %s", d.Id(), err)
	}

	d.Set("arn", repository.Arn)
	d.Set("clone_url_http", repository.CloneUrlHttp)
	d.Set("clone_url_ssh", repository.CloneUrlSsh)
	if _, ok := d.GetOk("default_branch"); ok {
		// The default branch can only be set when there is code in the repository.
		// Preserve the configured value.
		if v := repository.DefaultBranch; v != nil { // nosemgrep:ci.helper-schema-ResourceData-Set-extraneous-nil-check
			d.Set("default_branch", v)
		}
	}
	d.Set("description", repository.RepositoryDescription)
	d.Set("kms_key_id", repository.KmsKeyId)
	d.Set("repository_id", repository.RepositoryId)
	d.Set("repository_name", repository.RepositoryName)

	return diags
}

func resourceRepositoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CodeCommitClient(ctx)

	if d.HasChange("repository_name") {
		newName := d.Get("repository_name").(string)
		input := &codecommit.UpdateRepositoryNameInput{
			NewName: aws.String(newName),
			OldName: aws.String(d.Id()),
		}

		_, err := conn.UpdateRepositoryName(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating CodeCommit Repository (%s) name: %s", d.Id(), err)
		}

		d.SetId(newName)
	}

	if d.HasChange("default_branch") {
		if err := updateRepositoryDefaultBranch(ctx, conn, d.Id(), d.Get("default_branch").(string)); err != nil {
			return sdkdiag.AppendFromErr(diags, err)
		}
	}

	if d.HasChange("description") {
		input := &codecommit.UpdateRepositoryDescriptionInput{
			RepositoryDescription: aws.String(d.Get("description").(string)),
			RepositoryName:        aws.String(d.Id()),
		}

		_, err := conn.UpdateRepositoryDescription(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating CodeCommit Repository (%s) description: %s", d.Id(), err)
		}
	}

	if d.HasChange("kms_key_id") {
		input := &codecommit.UpdateRepositoryEncryptionKeyInput{
			KmsKeyId:       aws.String((d.Get("kms_key_id").(string))),
			RepositoryName: aws.String(d.Id()),
		}

		_, err := conn.UpdateRepositoryEncryptionKey(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating CodeCommit Repository (%s) encryption key: %s", d.Id(), err)
		}
	}

	return append(diags, resourceRepositoryRead(ctx, d, meta)...)
}

func resourceRepositoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CodeCommitClient(ctx)

	log.Printf("[INFO] Deleting CodeCommit Repository: %s", d.Id())
	_, err := conn.DeleteRepository(ctx, &codecommit.DeleteRepositoryInput{
		RepositoryName: aws.String(d.Id()),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting CodeCommit Repository (%s): %s", d.Id(), err)
	}

	return diags
}

func updateRepositoryDefaultBranch(ctx context.Context, conn *codecommit.Client, name, defaultBranch string) error {
	inputL := &codecommit.ListBranchesInput{
		RepositoryName: aws.String(name),
	}

	output, err := conn.ListBranches(ctx, inputL)

	if err != nil {
		return fmt.Errorf("listing CodeCommit Repository (%s) branches: %s", name, err)
	}

	if len(output.Branches) == 0 {
		return nil
	}

	inputU := &codecommit.UpdateDefaultBranchInput{
		DefaultBranchName: aws.String(defaultBranch),
		RepositoryName:    aws.String(name),
	}

	_, err = conn.UpdateDefaultBranch(ctx, inputU)

	if err != nil {
		return fmt.Errorf("updating CodeCommit Repository (%s) default branch: %w", name, err)
	}

	return nil
}

func findRepositoryByName(ctx context.Context, conn *codecommit.Client, name string) (*types.RepositoryMetadata, error) {
	input := &codecommit.GetRepositoryInput{
		RepositoryName: aws.String(name),
	}

	output, err := conn.GetRepository(ctx, input)

	if errs.IsA[*types.RepositoryDoesNotExistException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.RepositoryMetadata == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.RepositoryMetadata, nil
}
