package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	f := awslambdago.NewGoFunction(stack, jsii.String("sqlite-function"), &awslambdago.GoFunctionProps{
		Architecture: awslambda.Architecture_ARM_64(),
		LogRetention: awslogs.RetentionDays_ONE_WEEK,
		MemorySize:   jsii.Number(1024),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(15)),
		Entry:        jsii.String("function"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		Bundling: &awslambdago.BundlingOptions{
			CgoEnabled:           jsii.Bool(true),
			DockerImage:          awscdk.DockerImage_FromBuild(jsii.String("function"), nil),
			ForcedDockerBundling: jsii.Bool(true),
			Environment: &map[string]*string{
				"GOMODCACHE": jsii.String("/tmp/"),
				"GOCACHE":    jsii.String("/tmp/"),
			},
		},
	})
	u := f.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})

	awscdk.NewCfnOutput(stack, jsii.String("sqlite-function-url"), &awscdk.CfnOutputProps{
		Value: u.Url(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkStack(app, "go-sqlite-test", nil)

	app.Synth(nil)
}
