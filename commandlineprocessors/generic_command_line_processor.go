package commandlineprocessors

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BrunoMCBraga/HayMakerCF/globalstringsproviders"
	"github.com/BrunoMCBraga/HayMakerCF/haymakercfengines"
	"github.com/BrunoMCBraga/HayMakerCF/haymakercfutil"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const defaultKubeconfigPathWithinHome string = ".kube/config"
const defaultZone string = "us-east-1"
const waitTimeBeforeDeletingCFCluster int = 60

func deleteService(kubeConfig *string, deploymentName *string) error {

	loadKubeConfigError := haymakercfengines.KubernetesLoadKubeConfig(kubeConfig)
	if loadKubeConfigError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->deleteService->haymakercfengines.KubernetesLoadKubeConfig:" + loadKubeConfigError.Error() + "|")
	}

	kubernetesDeleteAllDeploymentAndServiceError := haymakercfengines.KubernetesDeleteAllDeploymentAndService(deploymentName)
	if kubernetesDeleteAllDeploymentAndServiceError != nil {
		fmt.Println("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->deleteService->haymakercfengines.KubernetesDeleteAllDeploymentAndService:" + kubernetesDeleteAllDeploymentAndServiceError.Error() + "|")
	}

	return nil

}

func deleteAllServicesAndDeployments(kubeConfig *string) error {

	loadKubeConfigError := haymakercfengines.KubernetesLoadKubeConfig(kubeConfig)
	if loadKubeConfigError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->deleteService->haymakercfengines.KubernetesLoadKubeConfig:" + loadKubeConfigError.Error() + "|")
	}

	kubernetesDeleteAllDeploymentsAndServicesError := haymakercfengines.KubernetesDeleteAllDeploymentsAndServices()

	if kubernetesDeleteAllDeploymentsAndServicesError != nil {
		fmt.Println("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->deleteAllServicesAndDeployments->haymakercfengines.KubernetesDeleteAllDeploymentsAndServices:" + kubernetesDeleteAllDeploymentsAndServicesError.Error() + "|")
	}

	return nil
}

func deleteECR(repoName *string, zone *string) error {

	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(*zone),
	}))

	ecrSession := ecr.New(awsSession, aws.NewConfig().WithRegion(*zone))
	haymakercfengines.InitECREngine(ecrSession)

	eCRDestroyDockerImageError := haymakercfengines.ECRBatchDeleteDockerImagesFromRepository(repoName)
	if eCRDestroyDockerImageError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->deleteECR->haymakercfengines.ECRBatchDeleteDockerImagesFromRepository:" + eCRDestroyDockerImageError.Error() + "|")
	}

	return nil
}

func deployCloudFormationConfig(cloudFormationConfigPath *string, stackName *string, bucketName *string, fileKey *string, zone *string) error {

	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(*zone),
	}))

	s3Session := s3.New(awsSession, aws.NewConfig().WithRegion(*zone))
	s3Uploader := s3manager.NewUploader(awsSession)

	haymakercfengines.InitS3Engine(s3Session, s3Uploader)

	s3CreateBucketError := haymakercfengines.S3CreateBucket(bucketName)
	if s3CreateBucketError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->deployCloudFormationConfig->haymakercfengines.S3CreateBucket:" + s3CreateBucketError.Error() + "|")
	}

	s3UploadFileToBucketResult, s3UploadFileToBucketError := haymakercfengines.S3UploadFileToBucket(cloudFormationConfigPath, bucketName, fileKey)
	if s3UploadFileToBucketError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->deployCloudFormationConfig->haymakercfengines.S3UploadFileToBucket:" + s3UploadFileToBucketError.Error() + "|")
	}

	cloudFormationSession := cloudformation.New(awsSession, aws.NewConfig().WithRegion(*zone))
	haymakercfengines.InitCloudFormationEngine(cloudFormationSession)

	_, cloudFormationCreateStackError := haymakercfengines.CloudFormationCreateStack(stackName, s3UploadFileToBucketResult)
	if cloudFormationCreateStackError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->deployCloudFormationConfig->haymakercfengines.CloudFormationCreateStack:" + cloudFormationCreateStackError.Error() + "|")
	}

	return nil
}

func teardownCloudFormationConfig(stackName *string, bucketName *string, repoName *string, zone *string, kubeConfig *string) error {

	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(*zone),
	}))

	s3Session := s3.New(awsSession, aws.NewConfig().WithRegion(*zone))
	s3Uploader := s3manager.NewUploader(awsSession)

	if *bucketName != "" {
		haymakercfengines.InitS3Engine(s3Session, s3Uploader)

		s3DeleteBucketError := haymakercfengines.S3DeleteBucket(bucketName)
		if s3DeleteBucketError != nil {
			fmt.Println("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->teardownCloudFormationConfig->haymakercfengines.DeleteS3Bucket:" + s3DeleteBucketError.Error() + "|")
		}
	}

	deleteAllServicesAndDeploymentsError := deleteAllServicesAndDeployments(kubeConfig)
	if deleteAllServicesAndDeploymentsError != nil {
		fmt.Println("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->teardownCloudFormationConfig->deleteAllServicesAndDeployments:" + deleteAllServicesAndDeploymentsError.Error() + "|")
	}

	time.Sleep(time.Duration(waitTimeBeforeDeletingCFCluster) * time.Second)

	//We don't stop just becayse ecr delete failed..
	deleteECRError := deleteECR(repoName, zone)
	if deleteECRError != nil {
		fmt.Println("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->teardownCloudFormationConfig->deleteECR:" + deleteECRError.Error() + "|")
	}

	cloudFormationSession := cloudformation.New(awsSession, aws.NewConfig().WithRegion(*zone))
	haymakercfengines.InitCloudFormationEngine(cloudFormationSession)

	cloudFormationDeleteStackError := haymakercfengines.CloudFormationDeleteStack(stackName)
	if cloudFormationDeleteStackError != nil {
		fmt.Println("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->teardownCloudFormationConfig->haymakercfengines.DeleteCloudFormationStack:" + cloudFormationDeleteStackError.Error() + "|")
	}

	return nil

}

func buildDockerImageAndPushToECR(repoName *string, zone *string, dockerFilePath *string, deleteLocalImagesAfterPush bool, registryId string) error {

	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(*zone),
	}))

	ecrSession := ecr.New(awsSession, aws.NewConfig().WithRegion(*zone))

	haymakercfengines.InitECREngine(ecrSession)

	eCRGetAuthorizationTokenResult, eCRGetAuthorizationTokenError := haymakercfengines.ECRGetAuthorizationToken(registryId)
	if eCRGetAuthorizationTokenError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->buildDockerImageAndPushToECR->haymakercfengines.GetAuthorizationToken:" + eCRGetAuthorizationTokenError.Error() + "|")
	} else if eCRGetAuthorizationTokenResult == nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->buildDockerImageAndPushToECR->haymakercfengines.BuildImageAndPush: Unable to retrieve authorization token for ECR. |")
	}

	dockerBuildImageAndPushError := haymakercfengines.DockerBuildImageAndPush(dockerFilePath, repoName, eCRGetAuthorizationTokenResult["token"], eCRGetAuthorizationTokenResult["endpoint"])
	if dockerBuildImageAndPushError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->buildDockerImageAndPushToECR->haymakercfengines.BuildImageAndPush:" + dockerBuildImageAndPushError.Error() + "|")
	}

	if deleteLocalImagesAfterPush {
		dockerDeleteStagingImagesError := haymakercfengines.DockerDeleteStagingImages(repoName, eCRGetAuthorizationTokenResult["endpoint"])
		if dockerDeleteStagingImagesError != nil {
			return errors.New("|" + "HayMaker->commandlineprocessors->generic_command_line_processor->deleteImageFromLocalDocker->haymakercfengines.DockerDeleteStagingImages:" + dockerDeleteStagingImagesError.Error() + "|")
		}
	}

	return nil
}

func writeFile(stringToWrite *string, filePath *string) error {

	fileHandle, createError := os.Create(*filePath)
	if createError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->writeFile->os.Create:" + createError.Error() + "|")
	}

	defer fileHandle.Close()

	_, writeError := fileHandle.Write([]byte(*stringToWrite))
	if writeError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->writeFile->file.Write:" + writeError.Error() + "|")
	}

	return nil
}

func generateKubeconfigFile(kubeconfigPath *string, zone *string, eksClusterName *string) error {

	awsSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(*zone),
	}))

	eksSession := eks.New(awsSession, aws.NewConfig().WithRegion(*zone))
	haymakercfengines.InitEKSEngine(eksSession)

	eKSGetKubernetesClusterAuthParametersResult, eKSGetKubernetesClusterAuthParametersError := haymakercfengines.EKSGetKubernetesClusterAuthParameters(eksClusterName)
	if eKSGetKubernetesClusterAuthParametersError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->generateKubeconfigFile-> haymakercfengines.EKSGetKubernetesClusterAuthParameters:" + eKSGetKubernetesClusterAuthParametersError.Error() + "|")
	}

	fileWriteError := writeFile(haymakercfutil.GenerateKubectlConfFileFromStruct(eKSGetKubernetesClusterAuthParametersResult), kubeconfigPath)
	if fileWriteError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->generateKubeconfigFile->writeFile:" + fileWriteError.Error() + "|")
	}

	return nil

}

func createKubernetesService(kubeConfig *string, kubernetesPort int, deploymentName *string, imageName *string, kubernetesProtocol *string, kubernetesReplicas int) error {

	loadKubeConfigError := haymakercfengines.KubernetesLoadKubeConfig(kubeConfig)
	if loadKubeConfigError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->createKubernetesService->haymakercfengines.KubernetesLoadKubeConfig:" + loadKubeConfigError.Error() + "|")
	}

	kubernetesDeployContainerError := haymakercfengines.KubernetesDeployContainer(kubernetesPort, deploymentName, imageName, kubernetesProtocol, kubernetesReplicas)
	if kubernetesDeployContainerError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->createKubernetesService->haymakercfengines.KubernetesDeployContainer:" + kubernetesDeployContainerError.Error() + "|")
	}

	return nil
}

func ProcessCommandLine(commandLineMap map[string]interface{}) error {

	userHomeDirResult, userHomeDirError := os.UserHomeDir()
	if userHomeDirError != nil {
		return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine->os.UserHomeDir:" + userHomeDirError.Error() + "|")
	}
	defaultKubeconfigFile := fmt.Sprintf("%s/%s", userHomeDirResult, defaultKubeconfigPathWithinHome)

	if opt, ok := commandLineMap["option"].(string); ok {
		switch opt {
		case "td":
			//go run ./main.go -cm td -t /Users/brubraga/go/src/github.com/haymakercf/CloudFormationFiles/cloudformation_cluster.json -sn haymakerstack -fk something -bn haymakerbucket -cn haymaker-eks
			fmt.Println("HaymakerCF CloudFormation Template Deployment")

			var cloudFormationTemplateTemp string
			if cloudFormationTemplate, ok := commandLineMap["cf_template"].(string); ok && cloudFormationTemplate != "" {
				cloudFormationTemplateTemp = cloudFormationTemplate
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -t (CloudFormation template file) |")
			}

			var cloudFormationStackNameTemp string
			if cloudFormationStackName, ok := commandLineMap["cf_stack_name"].(string); ok && cloudFormationStackName != "" {
				cloudFormationStackNameTemp = cloudFormationStackName
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -sn (CloudFormation stack name |")
			}

			var bucketNameTemp string
			if bucketName, ok := commandLineMap["s3_bucket_name"].(string); ok && bucketName != "" {
				bucketNameTemp = bucketName
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -bn (S3 bucket name to store CloudFormation template file) |")
			}

			var fileKeyTemp string
			if fileKey, ok := commandLineMap["s3_cf_template_file_key"].(string); ok && fileKey != "" {
				fileKeyTemp = fileKey
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -fk (S3 file key for CloudFormation template file) |")
			}

			var zoneTemp string
			if zone, ok := commandLineMap["session_zone"].(string); ok && zone != "" {
				zoneTemp = zone
			} else {
				zoneTemp = defaultZone
			}

			deployCloudFormationConfigError := deployCloudFormationConfig(&cloudFormationTemplateTemp, &cloudFormationStackNameTemp, &bucketNameTemp, &fileKeyTemp, &zoneTemp)
			if deployCloudFormationConfigError != nil {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine->deployCloudFormationConfig:" + deployCloudFormationConfigError.Error() + "|")
			}

			var clusterNameTemp string
			if clusterName, ok := commandLineMap["cluster_name"].(string); ok && clusterName != "" {
				clusterNameTemp = clusterName
			} else {
				fmt.Println("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -cn (Kubernetes cluster name). |")
			}

			var kubeConfigTemp string
			if kubeConfig, ok := commandLineMap["kubeconfig_file"].(string); ok && kubeConfig != "" {
				kubeConfigTemp = kubeConfig
			} else {
				kubeConfigTemp = defaultKubeconfigFile
			}

			generateKubectlFileError := generateKubeconfigFile(&kubeConfigTemp, &zoneTemp, &clusterNameTemp)
			if generateKubectlFileError != nil {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine->generateKubectlFile:" + generateKubectlFileError.Error() + "|")
			}

		case "tt": //go run ./main.go -cm tt -sn haymakerstack -bn haymakercfbucket -rn haymaker-docker-repo
			fmt.Println("HaymakerCF CloudFormation Teardown")

			var kubeConfigTemp string
			if kubeConfig, ok := commandLineMap["kubeconfig_file"].(string); ok && kubeConfig != "" {
				kubeConfigTemp = kubeConfig
			} else {
				kubeConfigTemp = defaultKubeconfigFile
			}

			var cloudFormationStackNameTemp string
			if cloudFormationStackName, ok := commandLineMap["cf_stack_name"].(string); ok && cloudFormationStackName != "" {
				cloudFormationStackNameTemp = cloudFormationStackName
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -sn (CloudFormation stack name). |")
			}

			var bucketNameTemp string
			if bucketName, ok := commandLineMap["s3_bucket_name"].(string); ok {
				bucketNameTemp = bucketName
			}

			var repoNameTemp string
			if repoName, ok := commandLineMap["repo_name"].(string); ok && repoName != "" {
				repoNameTemp = repoName
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -rn (ECR repo name). |")
			}

			var zoneTemp string
			if zone, ok := commandLineMap["session_zone"].(string); ok && zone != "" {
				zoneTemp = zone
			} else {
				zoneTemp = defaultZone
			}

			teardownCloudFormationConfig := teardownCloudFormationConfig(&cloudFormationStackNameTemp, &bucketNameTemp, &repoNameTemp, &zoneTemp, &kubeConfigTemp)
			if teardownCloudFormationConfig != nil {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine->teardownCloudFormationConfig:" + teardownCloudFormationConfig.Error() + "|")
			}
		case "pi":
			//go run ./main.go -cm pi -rn haymaker-docker-repo/haymaker-docker -df /Users/brubraga/go/src/github.com/haymakercf/Docker -di -ri [registry_id]
			fmt.Println("HayMakerCF Docker Image Build And Push To ECR")

			var repoNameTemp string
			if repoName, ok := commandLineMap["repo_name"].(string); ok && repoName != "" {
				repoNameTemp = repoName
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -rn (ECR repo name). |")
			}

			var dockerFileTemp string
			if dockerFile, ok := commandLineMap["docker_file_folder"].(string); ok && dockerFile != "" {
				dockerFileTemp = dockerFile
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -df (Dockerfile folder). |")
			}

			var zoneTemp string
			if zone, ok := commandLineMap["session_zone"].(string); ok && zone != "" {
				zoneTemp = zone
			} else {
				zoneTemp = defaultZone
			}

			var deleteLocalImagesAfterPushTemp bool
			if deleteLocalImagesAfterPush, ok := commandLineMap["delete_local_images_after_push"].(bool); ok {
				deleteLocalImagesAfterPushTemp = deleteLocalImagesAfterPush
			} else {
				deleteLocalImagesAfterPush = false
			}

			var registryIDTemp string
			if registryID, ok := commandLineMap["registry_id"].(string); ok {
				registryIDTemp = registryID
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -ri (ECR registry ID). |")
			}

			buildDockerImageAndPushToECRError := buildDockerImageAndPushToECR(&repoNameTemp, &zoneTemp, &dockerFileTemp, deleteLocalImagesAfterPushTemp, registryIDTemp)
			if buildDockerImageAndPushToECRError != nil {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine->buildDockerImageAndPushToECR:" + buildDockerImageAndPushToECRError.Error() + "|")
			}
		case "gk":
			fmt.Println("HayMakerCF Generate Kubeconfig File")
			//go run ./main.go -cm gk -cn haymaker-eks

			var zoneTemp string
			if zone, ok := commandLineMap["session_zone"].(string); ok && zone != "" {
				zoneTemp = zone
			} else {
				zoneTemp = defaultZone
			}

			var clusterNameTemp string
			if clusterName, ok := commandLineMap["cluster_name"].(string); ok && clusterName != "" {
				clusterNameTemp = clusterName
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -cn (Kubernetes cluster name).|")
			}

			generateKubectlFileError := generateKubeconfigFile(&defaultKubeconfigFile, &zoneTemp, &clusterNameTemp)
			if generateKubectlFileError != nil {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine->generateKubectlFile:" + generateKubectlFileError.Error() + "|")
			}
		case "sc":
			//go run ./main.go -cm sc -kp 80 -dn haymaker -in 965440066241.dkr.ecr.us-east-1.amazonaws.com/haymaker-docker-repo/haymaker-docker:latest -pr TCP -kr 2

			fmt.Println("HayMakerCF Deploy Container And Create Service")

			var kubeConfigTemp string
			if kubeConfig, ok := commandLineMap["kubeconfig_file"].(string); ok && kubeConfig != "" {
				kubeConfigTemp = kubeConfig
			} else {
				kubeConfigTemp = defaultKubeconfigFile
			}

			var kubernetesPortTemp int
			if kubernetesPort, ok := commandLineMap["kubernetes_port"].(int); ok && kubernetesPort != 0 {
				kubernetesPortTemp = kubernetesPort
			} else {
				kubernetesPortTemp = 0
				//return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -kp (Kubernetes container port). |")
			}

			var kubernetesDeploymentNameTemp string
			if kubernetesDeploymentName, ok := commandLineMap["kubernetes_deployment_name"].(string); ok && kubernetesDeploymentName != "" {
				kubernetesDeploymentNameTemp = kubernetesDeploymentName
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -dn (Kubernetes deployment name). |")
			}

			var imageNameTemp string
			if imageName, ok := commandLineMap["kubernetes_image_name"].(string); ok && imageName != "" {
				imageNameTemp = imageName
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -in (ECR image name, i.e., URI). |")
			}

			var kubernetesProtocolTemp string
			if kubernetesProtocol, ok := commandLineMap["kubernetes_protocol"].(string); ok && (kubernetesProtocol != "" && kubernetesPortTemp != 0 || kubernetesPortTemp == 0) {
				kubernetesProtocolTemp = kubernetesProtocol
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -pr (Kubernetes container port protocol). |")
			}

			var kubernetesReplicasTemp int
			if kubernetesReplicas, ok := commandLineMap["kubernetes_replicas"].(int); ok && kubernetesReplicas > 0 {
				kubernetesReplicasTemp = kubernetesReplicas
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -kr (Kubernetes number of replicas for container). |")
			}

			spinUpServiceError := createKubernetesService(&kubeConfigTemp, kubernetesPortTemp, &kubernetesDeploymentNameTemp, &imageNameTemp, &kubernetesProtocolTemp, kubernetesReplicasTemp)
			if spinUpServiceError != nil {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine->spinUpService:" + spinUpServiceError.Error() + "|")
			}
		case "sd":
			fmt.Println("HayMakerCF Delete Service And Associated Deployment")

			//go run ./main.go -cm ds -dn haymaker

			var kubeConfigTemp string
			if kubeConfig, ok := commandLineMap["kubeconfig_file"].(string); ok && kubeConfig != "" {
				kubeConfigTemp = kubeConfig
			} else {
				kubeConfigTemp = defaultKubeconfigFile
			}

			var kubernetesDeploymentNameTemp string
			if kubernetesDeploymentName, ok := commandLineMap["kubernetes_deployment_name"].(string); ok && kubernetesDeploymentName != "" {
				kubernetesDeploymentNameTemp = kubernetesDeploymentName
			} else {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine: missing -dn (Kubernetes deployment name). |")
			}

			deleteServiceError := deleteService(&kubeConfigTemp, &kubernetesDeploymentNameTemp)
			if deleteServiceError != nil {
				return errors.New("|" + "HayMakerCF->commandlineprocessors->generic_command_line_processor->ProcessCommandLine->deleteService:" + deleteServiceError.Error() + "|")
			}

		default:
			fmt.Println("Invalid option")
			fmt.Println(globalstringsproviders.GetMenuPictureStringWithOptions())
		}
	}

	return nil

}
