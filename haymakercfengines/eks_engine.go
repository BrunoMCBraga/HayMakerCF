package haymakercfengines

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
)

var eksInstance *eks.EKS

func describeClusterStub(eksClusterNameL *string) (*eks.DescribeClusterOutput, error) {

	describeClusterObject := &eks.DescribeClusterInput{
		Name: aws.String(*eksClusterNameL),
	}

	return eksInstance.DescribeCluster(describeClusterObject)
}

func EKSGetKubernetesClusterAuthParameters(eksClusterName *string) (map[string]interface{}, error) {

	var parameters map[string]interface{} = make(map[string]interface{}, 0)

	describeClusterResult, describeClusterErr := describeClusterStub(eksClusterName)
	if describeClusterErr != nil {
		return parameters, errors.New("|" + "HayMakerCF->haymakercfengines->eks_engine->EKSGetKubernetesClusterAuthParameters->describeClusterStub:" + describeClusterErr.Error() + "|")
	}

	parameters["server"] = describeClusterResult.Cluster.Endpoint
	parameters["cert"] = describeClusterResult.Cluster.CertificateAuthority.Data
	parameters["name"] = describeClusterResult.Cluster.Name

	return parameters, nil
}

func InitEKSEngine(eksIns *eks.EKS) {

	eksInstance = eksIns

}
