// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package sweep_test

import (
	"context"
	"slices"

	"terraform-provider-awsgps/internal/conns"
	"terraform-provider-awsgps/internal/service/accessanalyzer"
	"terraform-provider-awsgps/internal/service/account"
	"terraform-provider-awsgps/internal/service/acm"
	"terraform-provider-awsgps/internal/service/acmpca"
	"terraform-provider-awsgps/internal/service/amp"
	"terraform-provider-awsgps/internal/service/amplify"
	"terraform-provider-awsgps/internal/service/apigateway"
	"terraform-provider-awsgps/internal/service/apigatewayv2"
	"terraform-provider-awsgps/internal/service/appautoscaling"
	"terraform-provider-awsgps/internal/service/appconfig"
	"terraform-provider-awsgps/internal/service/appfabric"
	"terraform-provider-awsgps/internal/service/appflow"
	"terraform-provider-awsgps/internal/service/appintegrations"
	"terraform-provider-awsgps/internal/service/applicationinsights"
	"terraform-provider-awsgps/internal/service/appmesh"
	"terraform-provider-awsgps/internal/service/apprunner"
	"terraform-provider-awsgps/internal/service/appstream"
	"terraform-provider-awsgps/internal/service/appsync"
	"terraform-provider-awsgps/internal/service/athena"
	"terraform-provider-awsgps/internal/service/auditmanager"
	"terraform-provider-awsgps/internal/service/autoscaling"
	"terraform-provider-awsgps/internal/service/autoscalingplans"
	"terraform-provider-awsgps/internal/service/backup"
	"terraform-provider-awsgps/internal/service/batch"
	"terraform-provider-awsgps/internal/service/bedrock"
	"terraform-provider-awsgps/internal/service/bedrockagent"
	"terraform-provider-awsgps/internal/service/budgets"
	"terraform-provider-awsgps/internal/service/ce"
	"terraform-provider-awsgps/internal/service/chime"
	"terraform-provider-awsgps/internal/service/chimesdkmediapipelines"
	"terraform-provider-awsgps/internal/service/chimesdkvoice"
	"terraform-provider-awsgps/internal/service/cleanrooms"
	"terraform-provider-awsgps/internal/service/cloud9"
	"terraform-provider-awsgps/internal/service/cloudcontrol"
	"terraform-provider-awsgps/internal/service/cloudformation"
	"terraform-provider-awsgps/internal/service/cloudfront"
	"terraform-provider-awsgps/internal/service/cloudfrontkeyvaluestore"
	"terraform-provider-awsgps/internal/service/cloudhsmv2"
	"terraform-provider-awsgps/internal/service/cloudsearch"
	"terraform-provider-awsgps/internal/service/cloudtrail"
	"terraform-provider-awsgps/internal/service/cloudwatch"
	"terraform-provider-awsgps/internal/service/codeartifact"
	"terraform-provider-awsgps/internal/service/codebuild"
	"terraform-provider-awsgps/internal/service/codecatalyst"
	"terraform-provider-awsgps/internal/service/codecommit"
	"terraform-provider-awsgps/internal/service/codeguruprofiler"
	"terraform-provider-awsgps/internal/service/codegurureviewer"
	"terraform-provider-awsgps/internal/service/codepipeline"
	"terraform-provider-awsgps/internal/service/codestarconnections"
	"terraform-provider-awsgps/internal/service/codestarnotifications"
	"terraform-provider-awsgps/internal/service/cognitoidentity"
	"terraform-provider-awsgps/internal/service/cognitoidp"
	"terraform-provider-awsgps/internal/service/comprehend"
	"terraform-provider-awsgps/internal/service/computeoptimizer"
	"terraform-provider-awsgps/internal/service/configservice"
	"terraform-provider-awsgps/internal/service/connect"
	"terraform-provider-awsgps/internal/service/connectcases"
	"terraform-provider-awsgps/internal/service/controltower"
	"terraform-provider-awsgps/internal/service/costoptimizationhub"
	"terraform-provider-awsgps/internal/service/cur"
	"terraform-provider-awsgps/internal/service/customerprofiles"
	"terraform-provider-awsgps/internal/service/dataexchange"
	"terraform-provider-awsgps/internal/service/datapipeline"
	"terraform-provider-awsgps/internal/service/datasync"
	"terraform-provider-awsgps/internal/service/dax"
	"terraform-provider-awsgps/internal/service/deploy"
	"terraform-provider-awsgps/internal/service/detective"
	"terraform-provider-awsgps/internal/service/devicefarm"
	"terraform-provider-awsgps/internal/service/devopsguru"
	"terraform-provider-awsgps/internal/service/directconnect"
	"terraform-provider-awsgps/internal/service/dlm"
	"terraform-provider-awsgps/internal/service/dms"
	"terraform-provider-awsgps/internal/service/docdb"
	"terraform-provider-awsgps/internal/service/docdbelastic"
	"terraform-provider-awsgps/internal/service/ds"
	"terraform-provider-awsgps/internal/service/dynamodb"
	"terraform-provider-awsgps/internal/service/ec2"
	"terraform-provider-awsgps/internal/service/ecr"
	"terraform-provider-awsgps/internal/service/ecrpublic"
	"terraform-provider-awsgps/internal/service/ecs"
	"terraform-provider-awsgps/internal/service/efs"
	"terraform-provider-awsgps/internal/service/eks"
	"terraform-provider-awsgps/internal/service/elasticache"
	"terraform-provider-awsgps/internal/service/elasticbeanstalk"
	"terraform-provider-awsgps/internal/service/elasticsearch"
	"terraform-provider-awsgps/internal/service/elastictranscoder"
	"terraform-provider-awsgps/internal/service/elb"
	"terraform-provider-awsgps/internal/service/elbv2"
	"terraform-provider-awsgps/internal/service/emr"
	"terraform-provider-awsgps/internal/service/emrcontainers"
	"terraform-provider-awsgps/internal/service/emrserverless"
	"terraform-provider-awsgps/internal/service/events"
	"terraform-provider-awsgps/internal/service/evidently"
	"terraform-provider-awsgps/internal/service/finspace"
	"terraform-provider-awsgps/internal/service/firehose"
	"terraform-provider-awsgps/internal/service/fis"
	"terraform-provider-awsgps/internal/service/fms"
	"terraform-provider-awsgps/internal/service/fsx"
	"terraform-provider-awsgps/internal/service/gamelift"
	"terraform-provider-awsgps/internal/service/glacier"
	"terraform-provider-awsgps/internal/service/globalaccelerator"
	"terraform-provider-awsgps/internal/service/glue"
	"terraform-provider-awsgps/internal/service/grafana"
	"terraform-provider-awsgps/internal/service/greengrass"
	"terraform-provider-awsgps/internal/service/groundstation"
	"terraform-provider-awsgps/internal/service/guardduty"
	"terraform-provider-awsgps/internal/service/healthlake"
	"terraform-provider-awsgps/internal/service/iam"
	"terraform-provider-awsgps/internal/service/identitystore"
	"terraform-provider-awsgps/internal/service/imagebuilder"
	"terraform-provider-awsgps/internal/service/inspector"
	"terraform-provider-awsgps/internal/service/inspector2"
	"terraform-provider-awsgps/internal/service/internetmonitor"
	"terraform-provider-awsgps/internal/service/iot"
	"terraform-provider-awsgps/internal/service/iotanalytics"
	"terraform-provider-awsgps/internal/service/iotevents"
	"terraform-provider-awsgps/internal/service/ivs"
	"terraform-provider-awsgps/internal/service/ivschat"
	"terraform-provider-awsgps/internal/service/kafka"
	"terraform-provider-awsgps/internal/service/kafkaconnect"
	"terraform-provider-awsgps/internal/service/kendra"
	"terraform-provider-awsgps/internal/service/keyspaces"
	"terraform-provider-awsgps/internal/service/kinesis"
	"terraform-provider-awsgps/internal/service/kinesisanalytics"
	"terraform-provider-awsgps/internal/service/kinesisanalyticsv2"
	"terraform-provider-awsgps/internal/service/kinesisvideo"
	"terraform-provider-awsgps/internal/service/kms"
	"terraform-provider-awsgps/internal/service/lakeformation"
	"terraform-provider-awsgps/internal/service/lambda"
	"terraform-provider-awsgps/internal/service/launchwizard"
	"terraform-provider-awsgps/internal/service/lexmodels"
	"terraform-provider-awsgps/internal/service/lexv2models"
	"terraform-provider-awsgps/internal/service/licensemanager"
	"terraform-provider-awsgps/internal/service/lightsail"
	"terraform-provider-awsgps/internal/service/location"
	"terraform-provider-awsgps/internal/service/logs"
	"terraform-provider-awsgps/internal/service/lookoutmetrics"
	"terraform-provider-awsgps/internal/service/m2"
	"terraform-provider-awsgps/internal/service/macie2"
	"terraform-provider-awsgps/internal/service/mediaconnect"
	"terraform-provider-awsgps/internal/service/mediaconvert"
	"terraform-provider-awsgps/internal/service/medialive"
	"terraform-provider-awsgps/internal/service/mediapackage"
	"terraform-provider-awsgps/internal/service/mediapackagev2"
	"terraform-provider-awsgps/internal/service/mediastore"
	"terraform-provider-awsgps/internal/service/memorydb"
	"terraform-provider-awsgps/internal/service/meta"
	"terraform-provider-awsgps/internal/service/mq"
	"terraform-provider-awsgps/internal/service/mwaa"
	"terraform-provider-awsgps/internal/service/neptune"
	"terraform-provider-awsgps/internal/service/networkfirewall"
	"terraform-provider-awsgps/internal/service/networkmanager"
	"terraform-provider-awsgps/internal/service/oam"
	"terraform-provider-awsgps/internal/service/opensearch"
	"terraform-provider-awsgps/internal/service/opensearchserverless"
	"terraform-provider-awsgps/internal/service/opsworks"
	"terraform-provider-awsgps/internal/service/organizations"
	"terraform-provider-awsgps/internal/service/osis"
	"terraform-provider-awsgps/internal/service/outposts"
	"terraform-provider-awsgps/internal/service/pcaconnectorad"
	"terraform-provider-awsgps/internal/service/pinpoint"
	"terraform-provider-awsgps/internal/service/pipes"
	"terraform-provider-awsgps/internal/service/polly"
	"terraform-provider-awsgps/internal/service/pricing"
	"terraform-provider-awsgps/internal/service/qbusiness"
	"terraform-provider-awsgps/internal/service/qldb"
	"terraform-provider-awsgps/internal/service/quicksight"
	"terraform-provider-awsgps/internal/service/ram"
	"terraform-provider-awsgps/internal/service/rbin"
	"terraform-provider-awsgps/internal/service/rds"
	"terraform-provider-awsgps/internal/service/redshift"
	"terraform-provider-awsgps/internal/service/redshiftdata"
	"terraform-provider-awsgps/internal/service/redshiftserverless"
	"terraform-provider-awsgps/internal/service/rekognition"
	"terraform-provider-awsgps/internal/service/resourceexplorer2"
	"terraform-provider-awsgps/internal/service/resourcegroups"
	"terraform-provider-awsgps/internal/service/resourcegroupstaggingapi"
	"terraform-provider-awsgps/internal/service/rolesanywhere"
	"terraform-provider-awsgps/internal/service/route53"
	"terraform-provider-awsgps/internal/service/route53domains"
	"terraform-provider-awsgps/internal/service/route53recoverycontrolconfig"
	"terraform-provider-awsgps/internal/service/route53recoveryreadiness"
	"terraform-provider-awsgps/internal/service/route53resolver"
	"terraform-provider-awsgps/internal/service/rum"
	"terraform-provider-awsgps/internal/service/s3"
	"terraform-provider-awsgps/internal/service/s3control"
	"terraform-provider-awsgps/internal/service/s3outposts"
	"terraform-provider-awsgps/internal/service/sagemaker"
	"terraform-provider-awsgps/internal/service/scheduler"
	"terraform-provider-awsgps/internal/service/schemas"
	"terraform-provider-awsgps/internal/service/secretsmanager"
	"terraform-provider-awsgps/internal/service/securityhub"
	"terraform-provider-awsgps/internal/service/securitylake"
	"terraform-provider-awsgps/internal/service/serverlessrepo"
	"terraform-provider-awsgps/internal/service/servicecatalog"
	"terraform-provider-awsgps/internal/service/servicecatalogappregistry"
	"terraform-provider-awsgps/internal/service/servicediscovery"
	"terraform-provider-awsgps/internal/service/servicequotas"
	"terraform-provider-awsgps/internal/service/ses"
	"terraform-provider-awsgps/internal/service/sesv2"
	"terraform-provider-awsgps/internal/service/sfn"
	"terraform-provider-awsgps/internal/service/shield"
	"terraform-provider-awsgps/internal/service/signer"
	"terraform-provider-awsgps/internal/service/simpledb"
	"terraform-provider-awsgps/internal/service/sns"
	"terraform-provider-awsgps/internal/service/sqs"
	"terraform-provider-awsgps/internal/service/ssm"
	"terraform-provider-awsgps/internal/service/ssmcontacts"
	"terraform-provider-awsgps/internal/service/ssmincidents"
	"terraform-provider-awsgps/internal/service/ssmsap"
	"terraform-provider-awsgps/internal/service/sso"
	"terraform-provider-awsgps/internal/service/ssoadmin"
	"terraform-provider-awsgps/internal/service/storagegateway"
	"terraform-provider-awsgps/internal/service/sts"
	"terraform-provider-awsgps/internal/service/swf"
	"terraform-provider-awsgps/internal/service/synthetics"
	"terraform-provider-awsgps/internal/service/timestreamwrite"
	"terraform-provider-awsgps/internal/service/transcribe"
	"terraform-provider-awsgps/internal/service/transfer"
	"terraform-provider-awsgps/internal/service/verifiedpermissions"
	"terraform-provider-awsgps/internal/service/vpclattice"
	"terraform-provider-awsgps/internal/service/waf"
	"terraform-provider-awsgps/internal/service/wafregional"
	"terraform-provider-awsgps/internal/service/wafv2"
	"terraform-provider-awsgps/internal/service/wellarchitected"
	"terraform-provider-awsgps/internal/service/worklink"
	"terraform-provider-awsgps/internal/service/workspaces"
	"terraform-provider-awsgps/internal/service/xray"
)

func servicePackages(ctx context.Context) []conns.ServicePackage {
	v := []conns.ServicePackage{
		accessanalyzer.ServicePackage(ctx),
		account.ServicePackage(ctx),
		acm.ServicePackage(ctx),
		acmpca.ServicePackage(ctx),
		amp.ServicePackage(ctx),
		amplify.ServicePackage(ctx),
		apigateway.ServicePackage(ctx),
		apigatewayv2.ServicePackage(ctx),
		appautoscaling.ServicePackage(ctx),
		appconfig.ServicePackage(ctx),
		appfabric.ServicePackage(ctx),
		appflow.ServicePackage(ctx),
		appintegrations.ServicePackage(ctx),
		applicationinsights.ServicePackage(ctx),
		appmesh.ServicePackage(ctx),
		apprunner.ServicePackage(ctx),
		appstream.ServicePackage(ctx),
		appsync.ServicePackage(ctx),
		athena.ServicePackage(ctx),
		auditmanager.ServicePackage(ctx),
		autoscaling.ServicePackage(ctx),
		autoscalingplans.ServicePackage(ctx),
		backup.ServicePackage(ctx),
		batch.ServicePackage(ctx),
		bedrock.ServicePackage(ctx),
		bedrockagent.ServicePackage(ctx),
		budgets.ServicePackage(ctx),
		ce.ServicePackage(ctx),
		chime.ServicePackage(ctx),
		chimesdkmediapipelines.ServicePackage(ctx),
		chimesdkvoice.ServicePackage(ctx),
		cleanrooms.ServicePackage(ctx),
		cloud9.ServicePackage(ctx),
		cloudcontrol.ServicePackage(ctx),
		cloudformation.ServicePackage(ctx),
		cloudfront.ServicePackage(ctx),
		cloudfrontkeyvaluestore.ServicePackage(ctx),
		cloudhsmv2.ServicePackage(ctx),
		cloudsearch.ServicePackage(ctx),
		cloudtrail.ServicePackage(ctx),
		cloudwatch.ServicePackage(ctx),
		codeartifact.ServicePackage(ctx),
		codebuild.ServicePackage(ctx),
		codecatalyst.ServicePackage(ctx),
		codecommit.ServicePackage(ctx),
		codeguruprofiler.ServicePackage(ctx),
		codegurureviewer.ServicePackage(ctx),
		codepipeline.ServicePackage(ctx),
		codestarconnections.ServicePackage(ctx),
		codestarnotifications.ServicePackage(ctx),
		cognitoidentity.ServicePackage(ctx),
		cognitoidp.ServicePackage(ctx),
		comprehend.ServicePackage(ctx),
		computeoptimizer.ServicePackage(ctx),
		configservice.ServicePackage(ctx),
		connect.ServicePackage(ctx),
		connectcases.ServicePackage(ctx),
		controltower.ServicePackage(ctx),
		costoptimizationhub.ServicePackage(ctx),
		cur.ServicePackage(ctx),
		customerprofiles.ServicePackage(ctx),
		dataexchange.ServicePackage(ctx),
		datapipeline.ServicePackage(ctx),
		datasync.ServicePackage(ctx),
		dax.ServicePackage(ctx),
		deploy.ServicePackage(ctx),
		detective.ServicePackage(ctx),
		devicefarm.ServicePackage(ctx),
		devopsguru.ServicePackage(ctx),
		directconnect.ServicePackage(ctx),
		dlm.ServicePackage(ctx),
		dms.ServicePackage(ctx),
		docdb.ServicePackage(ctx),
		docdbelastic.ServicePackage(ctx),
		ds.ServicePackage(ctx),
		dynamodb.ServicePackage(ctx),
		ec2.ServicePackage(ctx),
		ecr.ServicePackage(ctx),
		ecrpublic.ServicePackage(ctx),
		ecs.ServicePackage(ctx),
		efs.ServicePackage(ctx),
		eks.ServicePackage(ctx),
		elasticache.ServicePackage(ctx),
		elasticbeanstalk.ServicePackage(ctx),
		elasticsearch.ServicePackage(ctx),
		elastictranscoder.ServicePackage(ctx),
		elb.ServicePackage(ctx),
		elbv2.ServicePackage(ctx),
		emr.ServicePackage(ctx),
		emrcontainers.ServicePackage(ctx),
		emrserverless.ServicePackage(ctx),
		events.ServicePackage(ctx),
		evidently.ServicePackage(ctx),
		finspace.ServicePackage(ctx),
		firehose.ServicePackage(ctx),
		fis.ServicePackage(ctx),
		fms.ServicePackage(ctx),
		fsx.ServicePackage(ctx),
		gamelift.ServicePackage(ctx),
		glacier.ServicePackage(ctx),
		globalaccelerator.ServicePackage(ctx),
		glue.ServicePackage(ctx),
		grafana.ServicePackage(ctx),
		greengrass.ServicePackage(ctx),
		groundstation.ServicePackage(ctx),
		guardduty.ServicePackage(ctx),
		healthlake.ServicePackage(ctx),
		iam.ServicePackage(ctx),
		identitystore.ServicePackage(ctx),
		imagebuilder.ServicePackage(ctx),
		inspector.ServicePackage(ctx),
		inspector2.ServicePackage(ctx),
		internetmonitor.ServicePackage(ctx),
		iot.ServicePackage(ctx),
		iotanalytics.ServicePackage(ctx),
		iotevents.ServicePackage(ctx),
		ivs.ServicePackage(ctx),
		ivschat.ServicePackage(ctx),
		kafka.ServicePackage(ctx),
		kafkaconnect.ServicePackage(ctx),
		kendra.ServicePackage(ctx),
		keyspaces.ServicePackage(ctx),
		kinesis.ServicePackage(ctx),
		kinesisanalytics.ServicePackage(ctx),
		kinesisanalyticsv2.ServicePackage(ctx),
		kinesisvideo.ServicePackage(ctx),
		kms.ServicePackage(ctx),
		lakeformation.ServicePackage(ctx),
		lambda.ServicePackage(ctx),
		launchwizard.ServicePackage(ctx),
		lexmodels.ServicePackage(ctx),
		lexv2models.ServicePackage(ctx),
		licensemanager.ServicePackage(ctx),
		lightsail.ServicePackage(ctx),
		location.ServicePackage(ctx),
		logs.ServicePackage(ctx),
		lookoutmetrics.ServicePackage(ctx),
		m2.ServicePackage(ctx),
		macie2.ServicePackage(ctx),
		mediaconnect.ServicePackage(ctx),
		mediaconvert.ServicePackage(ctx),
		medialive.ServicePackage(ctx),
		mediapackage.ServicePackage(ctx),
		mediapackagev2.ServicePackage(ctx),
		mediastore.ServicePackage(ctx),
		memorydb.ServicePackage(ctx),
		meta.ServicePackage(ctx),
		mq.ServicePackage(ctx),
		mwaa.ServicePackage(ctx),
		neptune.ServicePackage(ctx),
		networkfirewall.ServicePackage(ctx),
		networkmanager.ServicePackage(ctx),
		oam.ServicePackage(ctx),
		opensearch.ServicePackage(ctx),
		opensearchserverless.ServicePackage(ctx),
		opsworks.ServicePackage(ctx),
		organizations.ServicePackage(ctx),
		osis.ServicePackage(ctx),
		outposts.ServicePackage(ctx),
		pcaconnectorad.ServicePackage(ctx),
		pinpoint.ServicePackage(ctx),
		pipes.ServicePackage(ctx),
		polly.ServicePackage(ctx),
		pricing.ServicePackage(ctx),
		qbusiness.ServicePackage(ctx),
		qldb.ServicePackage(ctx),
		quicksight.ServicePackage(ctx),
		ram.ServicePackage(ctx),
		rbin.ServicePackage(ctx),
		rds.ServicePackage(ctx),
		redshift.ServicePackage(ctx),
		redshiftdata.ServicePackage(ctx),
		redshiftserverless.ServicePackage(ctx),
		rekognition.ServicePackage(ctx),
		resourceexplorer2.ServicePackage(ctx),
		resourcegroups.ServicePackage(ctx),
		resourcegroupstaggingapi.ServicePackage(ctx),
		rolesanywhere.ServicePackage(ctx),
		route53.ServicePackage(ctx),
		route53domains.ServicePackage(ctx),
		route53recoverycontrolconfig.ServicePackage(ctx),
		route53recoveryreadiness.ServicePackage(ctx),
		route53resolver.ServicePackage(ctx),
		rum.ServicePackage(ctx),
		s3.ServicePackage(ctx),
		s3control.ServicePackage(ctx),
		s3outposts.ServicePackage(ctx),
		sagemaker.ServicePackage(ctx),
		scheduler.ServicePackage(ctx),
		schemas.ServicePackage(ctx),
		secretsmanager.ServicePackage(ctx),
		securityhub.ServicePackage(ctx),
		securitylake.ServicePackage(ctx),
		serverlessrepo.ServicePackage(ctx),
		servicecatalog.ServicePackage(ctx),
		servicecatalogappregistry.ServicePackage(ctx),
		servicediscovery.ServicePackage(ctx),
		servicequotas.ServicePackage(ctx),
		ses.ServicePackage(ctx),
		sesv2.ServicePackage(ctx),
		sfn.ServicePackage(ctx),
		shield.ServicePackage(ctx),
		signer.ServicePackage(ctx),
		simpledb.ServicePackage(ctx),
		sns.ServicePackage(ctx),
		sqs.ServicePackage(ctx),
		ssm.ServicePackage(ctx),
		ssmcontacts.ServicePackage(ctx),
		ssmincidents.ServicePackage(ctx),
		ssmsap.ServicePackage(ctx),
		sso.ServicePackage(ctx),
		ssoadmin.ServicePackage(ctx),
		storagegateway.ServicePackage(ctx),
		sts.ServicePackage(ctx),
		swf.ServicePackage(ctx),
		synthetics.ServicePackage(ctx),
		timestreamwrite.ServicePackage(ctx),
		transcribe.ServicePackage(ctx),
		transfer.ServicePackage(ctx),
		verifiedpermissions.ServicePackage(ctx),
		vpclattice.ServicePackage(ctx),
		waf.ServicePackage(ctx),
		wafregional.ServicePackage(ctx),
		wafv2.ServicePackage(ctx),
		wellarchitected.ServicePackage(ctx),
		worklink.ServicePackage(ctx),
		workspaces.ServicePackage(ctx),
		xray.ServicePackage(ctx),
	}

	return slices.Clone(v)
}
