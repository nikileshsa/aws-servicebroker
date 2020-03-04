package main

import (
	"context"
	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kafka"
	log "github.com/golang/glog"
	"os"
)

var clusterarn = os.Getenv("clusterarn")
var kafkaConnString, zookeeperConnStr string

func mskEndPoints(ctx context.Context, event cfn.Event) (physicalResourceID string, data map[string]interface{}, err error) {

	event.ResourceProperties["PhysicalResourceID"] = lambdacontext.LogStreamName

	data = map[string]interface{}{}

	if event.RequestType == "Create" {

		describeclusterinputarn := &kafka.DescribeClusterInput{
			ClusterArn: aws.String(clusterarn),
		}
		getBootstrapBrokersInput := &kafka.GetBootstrapBrokersInput{
			ClusterArn: aws.String(clusterarn),
		}
		mySession := session.Must(session.NewSession())
		svc := kafka.New(mySession)

		if describeclusteroutput, err := svc.DescribeCluster(describeclusterinputarn); err != nil {
			log.Errorf("Failed to describe cluster: %v", err)
		} else {

			zookeeperConnStr = *describeclusteroutput.ClusterInfo.ZookeeperConnectString
		}
		if getBootstrapBrokersOutput, err := svc.GetBootstrapBrokers(getBootstrapBrokersInput); err != nil {
			log.Errorf("Failed to Bootstrapbroker cluster: %v", err)
		} else {
			kafkaConnString = getBootstrapBrokersOutput.GoString()
		}
		data["BrokerConnectionString"] = kafkaConnString
		data["ZookeeperConnectionString"] =  zookeeperConnStr
	}
	return
}

func main() {
	lambda.Start(cfn.LambdaWrap(mskEndPoints))
}
