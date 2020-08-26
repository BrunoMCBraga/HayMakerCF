package commandlinegenerators

import (
	"flag"
	"strings"

	"github.com/BrunoMCBraga/HayMakerCF/globalstringsproviders"
)

var option *string
var cloudFormationTemplate *string
var cloudFormationStackName *string
var s3BucketName *string
var s3CloudFormationTemplatefileKey *string
var sessionZone *string
var kubeconfigFile *string

var repositoryName *string
var registryId *string
var dockerFileFolder *string
var deleteLocalImagesAfterPush *bool

var clusterName *string

var kubernetesServiceport *int
var deploymentName *string
var imageName *string
var kubernetesProtocol *string
var kubernetesReplicas *int

func PrepareCommandLineProcessing() {

	optionHelp := globalstringsproviders.GetOptionsMenu()

	option = flag.String("cm", "", strings.TrimLeft(optionHelp, "\n"))
	cloudFormationTemplate = flag.String("t", "", "CloudFormation Template.")
	cloudFormationStackName = flag.String("sn", "", "CloudFormation Stack Name.")

	s3BucketName = flag.String("bn", "", "Bucket Name.")
	s3CloudFormationTemplatefileKey = flag.String("fk", "", "File key.")

	sessionZone = flag.String("sz", "", "Session zone.")

	//Repository stuff
	repositoryName = flag.String("rn", "", "Repository Name (e.g. haymaker-docker-repo/haymaker-docker). ")
	registryId = flag.String("ri", "", "ECR registry ID")

	dockerFileFolder = flag.String("df", "", "Dockerfile folder (i.e. folder where you would run \"docker build .\")")
	deleteLocalImagesAfterPush = flag.Bool("di", false, "Delete local Docker images after build and push to ECR.")

	//Kubernetes cluster
	clusterName = flag.String("cn", "", "Kubernetes cluster name.")
	kubeconfigFile = flag.String("kf", "", "Path used to save the Kubeconfig file. Used with gks option. If not provided, the default ~/.kube/config will be used.")
	kubernetesServiceport = flag.Int("kp", 0, "Port on which the service will run.")
	deploymentName = flag.String("dn", "", "Kubernetes deployment name.")
	imageName = flag.String("in", "", "Docker image name (e.g. 965440066241.dkr.ecr.us-east-1.amazonaws.com/haymaker-docker-repo/haymaker-docker:latest).")
	kubernetesProtocol = flag.String("pr", "", "Kubernetes service protocol (i.e. TCP, UDP).")
	kubernetesReplicas = flag.Int("kr", 0, "Number of Kubernetes service replicas.")
}

func ParseCommandLine() {
	flag.Parse()
}

func GetParametersDict() map[string]interface{} {

	parameters := make(map[string]interface{}, 0)
	parameters["option"] = *option
	parameters["cf_template"] = *cloudFormationTemplate
	parameters["cf_stack_name"] = *cloudFormationStackName
	parameters["s3_bucket_name"] = *s3BucketName
	parameters["s3_cf_template_file_key"] = *s3CloudFormationTemplatefileKey
	parameters["session_zone"] = *sessionZone

	parameters["repo_name"] = *repositoryName
	parameters["registry_id"] = *registryId
	parameters["docker_file_folder"] = *dockerFileFolder
	parameters["delete_local_images_after_push"] = *deleteLocalImagesAfterPush

	parameters["cluster_name"] = *clusterName

	parameters["kubeconfig_file"] = *kubeconfigFile

	parameters["kubernetes_port"] = *kubernetesServiceport
	parameters["kubernetes_deployment_name"] = *deploymentName
	parameters["kubernetes_image_name"] = *imageName
	parameters["kubernetes_protocol"] = *kubernetesProtocol
	parameters["kubernetes_replicas"] = *kubernetesReplicas

	return parameters
}
