package haymakercfengines

import (
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

var cloudFormationInstance *cloudformation.CloudFormation

const cloudFormationDeployWait int = 10
const cloudFormationDeleteWait int = 10

func createStackStub(stackName *string, templateURL *string) (*cloudformation.CreateStackOutput, error) {

	createStackInputObject := &cloudformation.CreateStackInput{
		OnFailure:   aws.String(cloudformation.OnFailureDelete),
		StackName:   aws.String(*stackName),
		TemplateURL: aws.String(*templateURL),
	}

	return cloudFormationInstance.CreateStack(createStackInputObject)
}

func deleteStackStub(stackName *string) (*cloudformation.DeleteStackOutput, error) {

	deleteStackInputObject := &cloudformation.DeleteStackInput{
		StackName: aws.String(*stackName),
	}

	return cloudFormationInstance.DeleteStack(deleteStackInputObject)
}

func describeStacksStub(stackName *string) (*cloudformation.DescribeStacksOutput, error) {

	describeStacksInputObject := &cloudformation.DescribeStacksInput{
		StackName: aws.String(*stackName),
	}

	return cloudFormationInstance.DescribeStacks(describeStacksInputObject)
}

func CloudFormationCreateStack(stackName *string, templateURL *string) (*cloudformation.CreateStackOutput, error) {

	fmt.Println("Creating CloudFormation Stack")

	createStackStubResult, createStackStubError := createStackStub(stackName, templateURL)
	if createStackStubError != nil {
		return nil, errors.New("|" + "HayMakerCF->haymakercfengines->cloudformation_engine->CloudFormationCreateStack->createStackStub:" + createStackStubError.Error() + "|")
	}

	for true {
		describeStacksStubResult, describeStacksStubError := describeStacksStub(stackName)
		if describeStacksStubError != nil {
			return nil, errors.New("|" + "HayMakerCF->haymakercfengines->cloudformation_engine->CloudFormationCreateStack->describeStacksStub:" + describeStacksStubError.Error() + "|")
		}

		if *describeStacksStubResult.Stacks[0].StackStatus == cloudformation.StackStatusCreateComplete {
			break
		}

		time.Sleep(time.Duration(cloudFormationDeployWait) * time.Second)
	}

	return createStackStubResult, nil

}

func CloudFormationDeleteStack(stackName *string) error {

	fmt.Println("Deleting CloudFormation Stack")

	_, deleteStackStubError := deleteStackStub(stackName)
	if deleteStackStubError != nil {
		fmt.Println("LOOOOOL3")
		return errors.New("|" + "HayMakerCF->haymakercfengines->cloudformation_engine->CloudFormationDeleteStack->deleteStackStub:" + deleteStackStubError.Error() + "|")
	}

	for true {
		fmt.Println("LOOOOOL")
		describeStacksStubResult, describeStacksStubError := describeStacksStub(stackName)
		if describeStacksStubError != nil {
			fmt.Println("LOOOOOL1")
			fmt.Println(*describeStacksStubResult.Stacks[0].StackStatus)
			fmt.Println(*describeStacksStubResult.Stacks[0].StackStatus)
			fmt.Println(*describeStacksStubResult.Stacks[0].StackStatus)
			return errors.New("|" + "HayMakerCF->haymakercfengines->cloudformation_engine->CloudFormationDeleteStack->describeStacksStub:" + describeStacksStubError.Error() + "|")
		}

		if *describeStacksStubResult.Stacks[0].StackStatus == cloudformation.StackStatusDeleteComplete {
			break
		} else if *describeStacksStubResult.Stacks[0].StackStatus == cloudformation.StackStatusDeleteFailed {
			fmt.Println("LOOOOOL2")
			fmt.Println(*describeStacksStubResult.Stacks[0].StackStatus)
			return errors.New("|HayMakerCF->haymakercfengines->cloudformation_engine->CloudFormationDeleteStack->describeStacksStub: Failed to delete stack. Consider deleting dependencies manually and retrying.|")
		}

		time.Sleep(time.Duration(cloudFormationDeleteWait) * time.Second)
	}

	return nil

}

func InitCloudFormationEngine(cfInst *cloudformation.CloudFormation) {

	cloudFormationInstance = cfInst
}
