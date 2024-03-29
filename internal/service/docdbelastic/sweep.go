// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package docdbelastic

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/docdbelastic"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-awsgps/internal/sweep"
	"terraform-provider-awsgps/internal/sweep/awsv2"
	"terraform-provider-awsgps/internal/sweep/framework"
	"terraform-provider-awsgps/names"
)

func RegisterSweepers() {
	resource.AddTestSweepers("aws_docdbelastic_cluster", &resource.Sweeper{
		Name: "aws_docdbelastic_cluster",
		F:    sweepClusters,
	})
}

func sweepClusters(region string) error {
	ctx := sweep.Context(region)
	if region == names.USWest1RegionID {
		log.Printf("[WARN] Skipping DocDB Elastic Cluster sweep for region: %s", region)
		return nil
	}
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.DocDBElasticClient(ctx)
	input := &docdbelastic.ListClustersInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	pages := docdbelastic.NewListClustersPaginator(conn, input)

	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)
		if awsv2.SkipSweepError(err) {
			log.Printf("[WARN] Skipping DocDB Elastic Clusters sweep for %s: %s", region, err)
			return nil
		}
		if err != nil {
			return fmt.Errorf("error retrieving DocDB Elastic Clusters: %w", err)
		}

		for _, cluster := range page.Clusters {
			arn := aws.ToString(cluster.ClusterArn)

			log.Printf("[INFO] Deleting DocDB Elastic Cluster: %s", aws.ToString(cluster.ClusterName))
			sweepResources = append(sweepResources, framework.NewSweepResource(newResourceCluster, client,
				framework.NewAttribute("id", arn),
			))
		}
	}

	if err := sweep.SweepOrchestrator(ctx, sweepResources); err != nil {
		return fmt.Errorf("error sweeping DocDB Elastic Clusters for %s: %w", region, err)
	}

	return nil
}
