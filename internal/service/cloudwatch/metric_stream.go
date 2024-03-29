// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudwatch

import (
	"context"
	"log"
	"time"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/create"
	"terraform-provider-awsgps/internal/enum"
	"terraform-provider-awsgps/internal/errs"
	"terraform-provider-awsgps/internal/errs/sdkdiag"
	"terraform-provider-awsgps/internal/flex"
	tftags "terraform-provider-awsgps/internal/tags"
	"terraform-provider-awsgps/internal/tfresource"
	"terraform-provider-awsgps/internal/verify"
	"terraform-provider-awsgps/names"
)

// @SDKResource("aws_cloudwatch_metric_stream", name="Metric Stream")
// @Tags(identifierAttribute="arn")
func resourceMetricStream() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceMetricStreamCreate,
		ReadWithoutTimeout:   resourceMetricStreamRead,
		UpdateWithoutTimeout: resourceMetricStreamUpdate,
		DeleteWithoutTimeout: resourceMetricStreamDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"exclude_filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_names": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringLenBetween(1, 255),
							},
						},
						"namespace": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
					},
				},
				ConflictsWith: []string{"include_filter"},
			},
			"firehose_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
			"include_filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_names": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringLenBetween(1, 255),
							},
						},
						"namespace": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 255),
						},
					},
				},
				ConflictsWith: []string{"exclude_filter"},
			},
			"include_linked_accounts_metrics": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"last_update_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  validateMetricStreamName,
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validateMetricStreamName,
			},
			"output_format": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: enum.Validate[types.MetricStreamOutputFormat](),
			},
			"role_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"statistics_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_statistics": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.All(
									validation.Any(
										validation.StringMatch(
											regexache.MustCompile(`(^IQM$)|(^(p|tc|tm|ts|wm)(100|\d{1,2})(\.\d{0,10})?$)|(^[ou]\d+(\.\d*)?$)`),
											"invalid statistic, see: https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Statistics-definitions.html",
										),
										validation.StringMatch(
											regexache.MustCompile(`^(TM|TC|TS|WM)\(((((\d{1,2})(\.\d{0,10})?|100(\.0{0,10})?)%)?:((\d{1,2})(\.\d{0,10})?|100(\.0{0,10})?)%|((\d{1,2})(\.\d{0,10})?|100(\.0{0,10})?)%:(((\d{1,2})(\.\d{0,10})?|100(\.0{0,10})?)%)?)\)|(TM|TC|TS|WM|PR)\(((\d+(\.\d{0,10})?|(\d+(\.\d{0,10})?[Ee][+-]?\d+)):((\d+(\.\d{0,10})?|(\d+(\.\d{0,10})?[Ee][+-]?\d+)))?|((\d+(\.\d{0,10})?|(\d+(\.\d{0,10})?[Ee][+-]?\d+)))?:(\d+(\.\d{0,10})?|(\d+(\.\d{0,10})?[Ee][+-]?\d+)))\)$`),
											"invalid statistic, see: https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Statistics-definitions.html",
										),
									),
									validation.StringDoesNotMatch(
										regexache.MustCompile(`^p0(\.0{0,10})?|p100(\.\d{0,10})?$`),
										"invalid statistic, see: https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Statistics-definitions.html",
									),
								),
							},
						},
						"include_metric": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringLenBetween(1, 255),
									},
									"namespace": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringLenBetween(1, 255),
									},
								},
							},
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},
	}
}

func resourceMetricStreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CloudWatchClient(ctx)

	name := create.Name(d.Get("name").(string), d.Get("name_prefix").(string))
	input := &cloudwatch.PutMetricStreamInput{
		FirehoseArn:                  aws.String(d.Get("firehose_arn").(string)),
		IncludeLinkedAccountsMetrics: aws.Bool(d.Get("include_linked_accounts_metrics").(bool)),
		Name:                         aws.String(name),
		OutputFormat:                 types.MetricStreamOutputFormat(d.Get("output_format").(string)),
		RoleArn:                      aws.String(d.Get("role_arn").(string)),
		Tags:                         getTagsIn(ctx),
	}

	if v, ok := d.GetOk("exclude_filter"); ok && v.(*schema.Set).Len() > 0 {
		input.ExcludeFilters = expandMetricStreamFilters(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("include_filter"); ok && v.(*schema.Set).Len() > 0 {
		input.IncludeFilters = expandMetricStreamFilters(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("statistics_configuration"); ok && v.(*schema.Set).Len() > 0 {
		input.StatisticsConfigurations = expandMetricStreamStatisticsConfigurations(v.(*schema.Set).List())
	}

	output, err := conn.PutMetricStream(ctx, input)

	// Some partitions (e.g. ISO) may not support tag-on-create.
	if input.Tags != nil && errs.IsUnsupportedOperationInPartitionError(meta.(*conns.AWSClient).Partition, err) {
		input.Tags = nil

		output, err = conn.PutMetricStream(ctx, input)
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating CloudWatch Metric Stream (%s): %s", name, err)
	}

	d.SetId(name)

	if _, err := waitMetricStreamRunning(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for CloudWatch Metric Stream (%s) create: %s", d.Id(), err)
	}

	// For partitions not supporting tag-on-create, attempt tag after create.
	if tags := getTagsIn(ctx); input.Tags == nil && len(tags) > 0 {
		err := createTags(ctx, conn, aws.ToString(output.Arn), tags)

		// If default tags only, continue. Otherwise, error.
		if v, ok := d.GetOk(names.AttrTags); (!ok || len(v.(map[string]interface{})) == 0) && errs.IsUnsupportedOperationInPartitionError(meta.(*conns.AWSClient).Partition, err) {
			return append(diags, resourceMetricStreamRead(ctx, d, meta)...)
		}

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "setting CloudWatch Metric Stream (%s) tags: %s", d.Id(), err)
		}
	}

	return append(diags, resourceMetricStreamRead(ctx, d, meta)...)
}

func resourceMetricStreamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CloudWatchClient(ctx)

	output, err := findMetricStreamByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] CloudWatch Metric Stream (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading CloudWatch Metric Stream (%s): %s", d.Id(), err)
	}

	d.Set("arn", output.Arn)
	d.Set("creation_date", output.CreationDate.Format(time.RFC3339))
	if output.ExcludeFilters != nil {
		if err := d.Set("exclude_filter", flattenMetricStreamFilters(output.ExcludeFilters)); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting exclude_filter: %s", err)
		}
	}
	d.Set("firehose_arn", output.FirehoseArn)
	if output.IncludeFilters != nil {
		if err := d.Set("include_filter", flattenMetricStreamFilters(output.IncludeFilters)); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting include_filter: %s", err)
		}
	}
	d.Set("include_linked_accounts_metrics", output.IncludeLinkedAccountsMetrics)
	d.Set("last_update_date", output.CreationDate.Format(time.RFC3339))
	d.Set("name", output.Name)
	d.Set("name_prefix", create.NamePrefixFromName(aws.ToString(output.Name)))
	d.Set("output_format", output.OutputFormat)
	d.Set("role_arn", output.RoleArn)
	d.Set("state", output.State)
	if output.StatisticsConfigurations != nil {
		if err := d.Set("statistics_configuration", flattenMetricStreamStatisticsConfigurations(output.StatisticsConfigurations)); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting statistics_configuration: %s", err)
		}
	}

	return diags
}

func resourceMetricStreamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CloudWatchClient(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		input := &cloudwatch.PutMetricStreamInput{
			FirehoseArn:                  aws.String(d.Get("firehose_arn").(string)),
			IncludeLinkedAccountsMetrics: aws.Bool(d.Get("include_linked_accounts_metrics").(bool)),
			Name:                         aws.String(d.Id()),
			OutputFormat:                 types.MetricStreamOutputFormat(d.Get("output_format").(string)),
			RoleArn:                      aws.String(d.Get("role_arn").(string)),
		}

		if v, ok := d.GetOk("exclude_filter"); ok && v.(*schema.Set).Len() > 0 {
			input.ExcludeFilters = expandMetricStreamFilters(v.(*schema.Set).List())
		}

		if v, ok := d.GetOk("include_filter"); ok && v.(*schema.Set).Len() > 0 {
			input.IncludeFilters = expandMetricStreamFilters(v.(*schema.Set).List())
		}

		if v, ok := d.GetOk("statistics_configuration"); ok && v.(*schema.Set).Len() > 0 {
			input.StatisticsConfigurations = expandMetricStreamStatisticsConfigurations(v.(*schema.Set).List())
		}

		_, err := conn.PutMetricStream(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating CloudWatch Metric Stream (%s): %s", d.Id(), err)
		}

		if _, err := waitMetricStreamRunning(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for CloudWatch Metric Stream (%s) update: %s", d.Id(), err)
		}
	}

	return append(diags, resourceMetricStreamRead(ctx, d, meta)...)
}

func resourceMetricStreamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CloudWatchClient(ctx)

	log.Printf("[INFO] Deleting CloudWatch Metric Stream: %s", d.Id())
	_, err := conn.DeleteMetricStream(ctx, &cloudwatch.DeleteMetricStreamInput{
		Name: aws.String(d.Id()),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting CloudWatch Metric Stream (%s): %s", d.Id(), err)
	}

	if _, err := waitMetricStreamDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for CloudWatch Metric Stream (%s) delete: %s", d.Id(), err)
	}

	return diags
}

func findMetricStreamByName(ctx context.Context, conn *cloudwatch.Client, name string) (*cloudwatch.GetMetricStreamOutput, error) {
	input := &cloudwatch.GetMetricStreamInput{
		Name: aws.String(name),
	}

	output, err := conn.GetMetricStream(ctx, input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func statusMetricStream(ctx context.Context, conn *cloudwatch.Client, name string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := findMetricStreamByName(ctx, conn, name)

		if tfresource.NotFound(err) {
			return nil, "", nil
		}

		if err != nil {
			return nil, "", err
		}

		return output, aws.ToString(output.State), nil
	}
}

const (
	metricStreamStateRunning = "running"
	metricStreamStateStopped = "stopped"
)

func waitMetricStreamDeleted(ctx context.Context, conn *cloudwatch.Client, name string, timeout time.Duration) (*cloudwatch.GetMetricStreamOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{metricStreamStateRunning, metricStreamStateStopped},
		Target:  []string{},
		Refresh: statusMetricStream(ctx, conn, name),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*cloudwatch.GetMetricStreamOutput); ok {
		return output, err
	}

	return nil, err
}

func waitMetricStreamRunning(ctx context.Context, conn *cloudwatch.Client, name string, timeout time.Duration) (*cloudwatch.GetMetricStreamOutput, error) { //nolint:unparam
	stateConf := &retry.StateChangeConf{
		Pending: []string{metricStreamStateStopped},
		Target:  []string{metricStreamStateRunning},
		Refresh: statusMetricStream(ctx, conn, name),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*cloudwatch.GetMetricStreamOutput); ok {
		return output, err
	}

	return nil, err
}

func validateMetricStreamName(v interface{}, k string) (ws []string, errors []error) {
	return validation.All(
		validation.StringLenBetween(1, 255),
		validation.StringMatch(regexache.MustCompile(`^[0-9A-Za-z_-]*$`), "must match [0-9A-Za-z_-]"),
	)(v, k)
}

func expandMetricStreamFilters(tfList []interface{}) []types.MetricStreamFilter {
	var apiObjects []types.MetricStreamFilter

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})
		if !ok {
			continue
		}

		apiObject := types.MetricStreamFilter{}

		if v, ok := tfMap["metric_names"].(*schema.Set); ok && v.Len() > 0 {
			apiObject.MetricNames = flex.ExpandStringValueSet(v)
		}

		if v, ok := tfMap["namespace"].(string); ok && v != "" {
			apiObject.Namespace = aws.String(v)
		}

		apiObjects = append(apiObjects, apiObject)
	}

	if len(apiObjects) == 0 {
		return nil
	}

	return apiObjects
}

func flattenMetricStreamFilters(apiObjects []types.MetricStreamFilter) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject.Namespace != nil {
			tfMap := map[string]interface{}{
				"metric_names": apiObject.MetricNames,
			}

			if v := apiObject.Namespace; v != nil {
				tfMap["namespace"] = aws.ToString(v)
			}

			tfList = append(tfList, tfMap)
		}
	}

	return tfList
}

func expandMetricStreamStatisticsConfigurations(tfList []interface{}) []types.MetricStreamStatisticsConfiguration {
	var apiObjects []types.MetricStreamStatisticsConfiguration

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})
		if !ok {
			continue
		}

		apiObject := types.MetricStreamStatisticsConfiguration{}

		if v, ok := tfMap["additional_statistics"].(*schema.Set); ok && v.Len() > 0 {
			apiObject.AdditionalStatistics = flex.ExpandStringValueSet(v)
		}

		if v, ok := tfMap["include_metric"].(*schema.Set); ok && v.Len() > 0 {
			apiObject.IncludeMetrics = expandMetricStreamStatisticsConfigurationsIncludeMetrics(v.List())
		}

		apiObjects = append(apiObjects, apiObject)
	}

	if len(apiObjects) == 0 {
		return nil
	}

	return apiObjects
}

func expandMetricStreamStatisticsConfigurationsIncludeMetrics(tfList []interface{}) []types.MetricStreamStatisticsMetric {
	var apiObjects []types.MetricStreamStatisticsMetric

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})
		if !ok {
			continue
		}

		apiObject := types.MetricStreamStatisticsMetric{}

		if v, ok := tfMap["metric_name"].(string); ok && v != "" {
			apiObject.MetricName = aws.String(v)
		}

		if v, ok := tfMap["namespace"].(string); ok && v != "" {
			apiObject.Namespace = aws.String(v)
		}

		apiObjects = append(apiObjects, apiObject)
	}

	if len(apiObjects) == 0 {
		return nil
	}

	return apiObjects
}

func flattenMetricStreamStatisticsConfigurations(apiObjects []types.MetricStreamStatisticsConfiguration) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		tfMap := map[string]interface{}{}

		if v := apiObject.AdditionalStatistics; v != nil {
			tfMap["additional_statistics"] = flex.FlattenStringValueSet(v)
		}

		if v := apiObject.IncludeMetrics; v != nil {
			tfMap["include_metric"] = flattenMetricStreamStatisticsConfigurationsIncludeMetrics(v)
		}

		tfList = append(tfList, tfMap)
	}

	return tfList
}

func flattenMetricStreamStatisticsConfigurationsIncludeMetrics(apiObjects []types.MetricStreamStatisticsMetric) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		tfMap := map[string]interface{}{}

		if v := apiObject.MetricName; v != nil {
			tfMap["metric_name"] = aws.ToString(v)
		}

		if v := apiObject.Namespace; v != nil {
			tfMap["namespace"] = aws.ToString(v)
		}

		tfList = append(tfList, tfMap)
	}

	return tfList
}
