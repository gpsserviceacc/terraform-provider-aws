// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package codepipeline

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/enum"
	"terraform-provider-awsgps/internal/errs"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	tfslices "terraform-provider-awsgps/internal/slices"
	tftags "terraform-provider-awsgps/internal/tags"
	"terraform-provider-awsgps/internal/tfresource"
	"terraform-provider-awsgps/internal/verify"
	"terraform-provider-awsgps/names"
)

// @SDKResource("aws_codepipeline_custom_action_type", name="Custom Action Type")
// @Tags(identifierAttribute="arn")
func resourceCustomActionType() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCustomActionTypeCreate,
		ReadWithoutTimeout:   resourceCustomActionTypeRead,
		UpdateWithoutTimeout: resourceCustomActionTypeUpdate,
		DeleteWithoutTimeout: resourceCustomActionTypeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: enum.Validate[types.ActionCategory](),
			},
			"configuration_property": {
				Type:     schema.TypeList,
				MaxItems: 10,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"queryable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"required": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"secret": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"type": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: enum.Validate[types.ActionConfigurationPropertyType](),
						},
					},
				},
			},
			"input_artifact_details": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maximum_count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 5),
						},
						"minimum_count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 5),
						},
					},
				},
			},
			"output_artifact_details": {
				Type:     schema.TypeList,
				ForceNew: true,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"maximum_count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 5),
						},
						"minimum_count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 5),
						},
					},
				},
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_name": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 35),
			},
			"settings": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_url_template": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"execution_url_template": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"revision_url_template": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"third_party_configuration_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"version": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 9),
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceCustomActionTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CodePipelineClient(ctx)

	category := d.Get("category").(string)
	provider := d.Get("provider_name").(string)
	version := d.Get("version").(string)
	id := CustomActionTypeCreateResourceID(category, provider, version)
	input := &codepipeline.CreateCustomActionTypeInput{
		Category: types.ActionCategory(category),
		Provider: aws.String(provider),
		Tags:     getTagsIn(ctx),
		Version:  aws.String(version),
	}

	if v, ok := d.GetOk("configuration_property"); ok && len(v.([]interface{})) > 0 {
		input.ConfigurationProperties = expandActionConfigurationProperties(v.([]interface{}))
	}

	if v, ok := d.GetOk("input_artifact_details"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.InputArtifactDetails = expandArtifactDetails(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("output_artifact_details"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.OutputArtifactDetails = expandArtifactDetails(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("settings"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.Settings = expandActionTypeSettings(v.([]interface{})[0].(map[string]interface{}))
	}

	_, err := conn.CreateCustomActionType(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating CodePipeline Custom Action Type (%s): %s", id, err)
	}

	d.SetId(id)

	return append(diags, resourceCustomActionTypeRead(ctx, d, meta)...)
}

func resourceCustomActionTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CodePipelineClient(ctx)

	category, provider, version, err := CustomActionTypeParseResourceID(d.Id())
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	actionType, err := findCustomActionTypeByThreePartKey(ctx, conn, category, provider, version)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] CodePipeline Custom Action Type %s not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading CodePipeline Custom Action Type (%s): %s", d.Id(), err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   "codepipeline",
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("actiontype:%s/%s/%s/%s", types.ActionOwnerCustom, category, provider, version),
	}.String()
	d.Set("arn", arn)
	d.Set("category", actionType.Id.Category)
	if err := d.Set("configuration_property", flattenActionConfigurationProperties(d, actionType.ActionConfigurationProperties)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting configuration_property: %s", err)
	}
	if actionType.InputArtifactDetails != nil {
		if err := d.Set("input_artifact_details", []interface{}{flattenArtifactDetails(actionType.InputArtifactDetails)}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting input_artifact_details: %s", err)
		}
	} else {
		d.Set("input_artifact_details", nil)
	}
	if actionType.OutputArtifactDetails != nil {
		if err := d.Set("output_artifact_details", []interface{}{flattenArtifactDetails(actionType.OutputArtifactDetails)}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting output_artifact_details: %s", err)
		}
	} else {
		d.Set("output_artifact_details", nil)
	}
	d.Set("owner", actionType.Id.Owner)
	d.Set("provider_name", actionType.Id.Provider)
	if actionType.Settings != nil &&
		// Service can return empty ({}) Settings.
		(actionType.Settings.EntityUrlTemplate != nil || actionType.Settings.ExecutionUrlTemplate != nil || actionType.Settings.RevisionUrlTemplate != nil || actionType.Settings.ThirdPartyConfigurationUrl != nil) {
		if err := d.Set("settings", []interface{}{flattenActionTypeSettings(actionType.Settings)}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting settings: %s", err)
		}
	} else {
		d.Set("settings", nil)
	}
	d.Set("version", actionType.Id.Version)

	return diags
}

func resourceCustomActionTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Tags only.

	return append(diags, resourceCustomActionTypeRead(ctx, d, meta)...)
}

func resourceCustomActionTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CodePipelineClient(ctx)

	category, provider, version, err := CustomActionTypeParseResourceID(d.Id())
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	log.Printf("[INFO] Deleting CodePipeline Custom Action Type: %s", d.Id())
	_, err = conn.DeleteCustomActionType(ctx, &codepipeline.DeleteCustomActionTypeInput{
		Category: category,
		Provider: aws.String(provider),
		Version:  aws.String(version),
	})

	if errs.IsA[*types.ActionNotFoundException](err) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting CodePipeline Custom Action Type (%s): %s", d.Id(), err)
	}

	return diags
}

const customActionTypeResourceIDSeparator = "/"

func CustomActionTypeCreateResourceID(category, provider, version string) string {
	parts := []string{category, provider, version}
	id := strings.Join(parts, customActionTypeResourceIDSeparator)

	return id
}

func CustomActionTypeParseResourceID(id string) (types.ActionCategory, string, string, error) {
	parts := strings.Split(id, customActionTypeResourceIDSeparator)

	if len(parts) == 3 && parts[0] != "" && parts[1] != "" && parts[2] != "" {
		return types.ActionCategory(parts[0]), parts[1], parts[2], nil
	}

	return "", "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected category%[2]sprovider%[2]sversion", id, customActionTypeResourceIDSeparator)
}

func findCustomActionTypeByThreePartKey(ctx context.Context, conn *codepipeline.Client, category types.ActionCategory, provider, version string) (*types.ActionType, error) {
	input := &codepipeline.ListActionTypesInput{
		ActionOwnerFilter: types.ActionOwnerCustom,
	}

	return findActionType(ctx, conn, input, func(v *types.ActionType) bool {
		return v.Id.Category == category && aws.ToString(v.Id.Provider) == provider && aws.ToString(v.Id.Version) == version
	})
}

func findActionType(ctx context.Context, conn *codepipeline.Client, input *codepipeline.ListActionTypesInput, filter tfslices.Predicate[*types.ActionType]) (*types.ActionType, error) {
	output, err := findActionTypes(ctx, conn, input, filter)

	if err != nil {
		return nil, err
	}

	return tfresource.AssertSinglePtrResult(output)
}

func findActionTypes(ctx context.Context, conn *codepipeline.Client, input *codepipeline.ListActionTypesInput, filter tfslices.Predicate[*types.ActionType]) ([]*types.ActionType, error) {
	var output []*types.ActionType

	pages := codepipeline.NewListActionTypesPaginator(conn, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if err != nil {
			return nil, err
		}

		for _, v := range page.ActionTypes {
			v := v
			if v := &v; filter(v) {
				output = append(output, v)
			}
		}
	}

	return output, nil
}

func expandActionConfigurationProperty(tfMap map[string]interface{}) *types.ActionConfigurationProperty {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.ActionConfigurationProperty{}

	if v, ok := tfMap["description"].(string); ok && v != "" {
		apiObject.Description = aws.String(v)
	}

	if v, ok := tfMap["key"].(bool); ok {
		apiObject.Key = v
	}

	if v, ok := tfMap["name"].(string); ok && v != "" {
		apiObject.Name = aws.String(v)
	}

	if v, ok := tfMap["queryable"].(bool); ok && v {
		apiObject.Queryable = v
	}

	if v, ok := tfMap["required"].(bool); ok {
		apiObject.Required = v
	}

	if v, ok := tfMap["secret"].(bool); ok {
		apiObject.Secret = v
	}

	if v, ok := tfMap["type"].(string); ok && v != "" {
		apiObject.Type = types.ActionConfigurationPropertyType(v)
	}

	return apiObject
}

func expandActionConfigurationProperties(tfList []interface{}) []types.ActionConfigurationProperty {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []types.ActionConfigurationProperty

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandActionConfigurationProperty(tfMap)

		if apiObject == nil {
			continue
		}

		apiObjects = append(apiObjects, *apiObject)
	}

	return apiObjects
}

func expandArtifactDetails(tfMap map[string]interface{}) *types.ArtifactDetails {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.ArtifactDetails{}

	if v, ok := tfMap["maximum_count"].(int); ok {
		apiObject.MaximumCount = int32(v)
	}

	if v, ok := tfMap["minimum_count"].(int); ok {
		apiObject.MinimumCount = int32(v)
	}

	return apiObject
}

func expandActionTypeSettings(tfMap map[string]interface{}) *types.ActionTypeSettings {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.ActionTypeSettings{}

	if v, ok := tfMap["entity_url_template"].(string); ok && v != "" {
		apiObject.EntityUrlTemplate = aws.String(v)
	}

	if v, ok := tfMap["execution_url_template"].(string); ok && v != "" {
		apiObject.ExecutionUrlTemplate = aws.String(v)
	}

	if v, ok := tfMap["revision_url_template"].(string); ok && v != "" {
		apiObject.RevisionUrlTemplate = aws.String(v)
	}

	if v, ok := tfMap["third_party_configuration_url"].(string); ok && v != "" {
		apiObject.ThirdPartyConfigurationUrl = aws.String(v)
	}

	return apiObject
}

func flattenActionConfigurationProperty(d *schema.ResourceData, i int, apiObject types.ActionConfigurationProperty) map[string]interface{} {
	tfMap := map[string]interface{}{
		"key":       apiObject.Key,
		"queryable": apiObject.Queryable,
		"required":  apiObject.Required,
		"secret":    apiObject.Secret,
	}

	if v := apiObject.Description; v != nil {
		tfMap["description"] = aws.ToString(v)
	}

	if v := apiObject.Name; v != nil {
		tfMap["name"] = aws.ToString(v)
	}

	if v := apiObject.Type; v != "" {
		tfMap["type"] = v
	} else {
		// The AWS API does not return Type.
		key := fmt.Sprintf("configuration_property.%d.type", i)
		tfMap["type"] = d.Get(key).(string)
	}

	return tfMap
}

func flattenActionConfigurationProperties(d *schema.ResourceData, apiObjects []types.ActionConfigurationProperty) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for i, apiObject := range apiObjects {
		tfList = append(tfList, flattenActionConfigurationProperty(d, i, apiObject))
	}

	return tfList
}

func flattenArtifactDetails(apiObject *types.ArtifactDetails) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{
		"maximum_count": apiObject.MaximumCount,
		"minimum_count": apiObject.MinimumCount,
	}

	return tfMap
}

func flattenActionTypeSettings(apiObject *types.ActionTypeSettings) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.EntityUrlTemplate; v != nil {
		tfMap["entity_url_template"] = aws.ToString(v)
	}

	if v := apiObject.ExecutionUrlTemplate; v != nil {
		tfMap["execution_url_template"] = aws.ToString(v)
	}

	if v := apiObject.RevisionUrlTemplate; v != nil {
		tfMap["revision_url_template"] = aws.ToString(v)
	}

	if v := apiObject.ThirdPartyConfigurationUrl; v != nil {
		tfMap["third_party_configuration_url"] = aws.ToString(v)
	}

	return tfMap
}
