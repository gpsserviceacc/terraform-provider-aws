// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elasticsearch

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	elasticsearch "github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"terraform-provider-awsgps/internal/tfresource"
)

func FindDomainByName(ctx context.Context, conn *elasticsearch.ElasticsearchService, name string) (*elasticsearch.ElasticsearchDomainStatus, error) {
	input := &elasticsearch.DescribeElasticsearchDomainInput{
		DomainName: aws.String(name),
	}

	output, err := conn.DescribeElasticsearchDomainWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, elasticsearch.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.DomainStatus == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.DomainStatus, nil
}
