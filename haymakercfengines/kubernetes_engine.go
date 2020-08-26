package haymakercfengines

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const defaultKubeconfigPathWithinHome string = ".kube/config"

var clientConfig *kubernetes.Clientset = nil

func KubernetesLoadKubeConfig(kubeconfig *string) error {

	var kubeconfigTemp string

	if kubeconfig == nil {

		userHomeDirResult, userHomeDirError := os.UserHomeDir()
		if userHomeDirError != nil {
			return errors.New("|" + "HayMakerCF->haymakercfengines->kubernetes_engine->KubernetesLoadKubeConfig->os.UserHomeDir:" + userHomeDirError.Error() + "|")
		}
		kubeconfigTemp = fmt.Sprintf("%s/%s", userHomeDirResult, defaultKubeconfigPathWithinHome)
	} else {
		kubeconfigTemp = *kubeconfig
	}

	buildConfigFromFlagsResult, buildConfigFromFlagsError := clientcmd.BuildConfigFromFlags("", kubeconfigTemp)
	if buildConfigFromFlagsError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->kubernetes_engine->KubernetesLoadKubeConfig->BuildConfigFromFlags:" + buildConfigFromFlagsError.Error() + "|")
	}

	newForConfigResult, newForConfigError := kubernetes.NewForConfig(buildConfigFromFlagsResult)
	if newForConfigError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->kubernetes_engine->KubernetesLoadKubeConfig->NewForConfig:" + newForConfigError.Error() + "|")
	}

	clientConfig = newForConfigResult
	return nil
}

func kubernetesDeleteDeployment(deploymentName *string) error {

	fmt.Println("Deleting Kubernetes Deployment")

	deploymentsClient := clientConfig.AppsV1().Deployments(apiv1.NamespaceDefault)

	// Create Deployment
	deploymentDeleteError := deploymentsClient.Delete(*deploymentName, &metav1.DeleteOptions{})

	if deploymentDeleteError != nil {
		if !strings.Contains(deploymentDeleteError.Error(), "not found") {
			return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->kubernetesDeleteDeployment->Deployments.Delete:" + deploymentDeleteError.Error() + "|")
		}
	}

	return nil
}

func kubernetesDeleteService(deploymentName *string) error {

	fmt.Println("Deleting Service")

	deleteError := clientConfig.CoreV1().Services(apiv1.NamespaceDefault).Delete(*deploymentName, &metav1.DeleteOptions{})

	if deleteError != nil {
		if !strings.Contains(deleteError.Error(), "not found") {
			return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->kubernetesDeleteService->Services.Create:" + deleteError.Error() + "|")
		}
	}

	return nil

}

func KubernetesDeleteAllDeploymentAndService(deploymentName *string) error {

	if clientConfig == nil {
		loadKubeConfigError := KubernetesLoadKubeConfig(nil)
		if loadKubeConfigError != nil {
			return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesCopyFileFromRemotePods->KubernetesLoadKubeConfig:" + loadKubeConfigError.Error() + "|")
		}

	}

	kubernetesDeleteServiceError := kubernetesDeleteService(deploymentName)
	if kubernetesDeleteServiceError != nil {
		return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesDeleteAllDeploymentAndService->KubernetesDeleteService:" + kubernetesDeleteServiceError.Error() + "|")
	}

	kubernetesDeleteDeploymentError := kubernetesDeleteDeployment(deploymentName)
	if kubernetesDeleteDeploymentError != nil {
		return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesDeleteAllDeploymentAndService->KubernetesDeleteDeployment:" + kubernetesDeleteDeploymentError.Error() + "|")
	}

	return nil
}

func KubernetesDeleteAllDeploymentsAndServices() error {

	fmt.Println("Deleting All Deployments And Services")

	if clientConfig == nil {
		loadKubeConfigError := KubernetesLoadKubeConfig(nil)
		if loadKubeConfigError != nil {
			return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesCopyFileFromRemotePods->KubernetesLoadKubeConfig:" + loadKubeConfigError.Error() + "|")
		}

	}

	servicesClient := clientConfig.CoreV1().Services(apiv1.NamespaceDefault)
	servicesListResult, servicesListError := servicesClient.List(metav1.ListOptions{})

	var serviceName string
	var deploymentName string

	if servicesListError != nil {
		return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesDeleteAllDeploymentsAndServices->services.List:" + servicesListError.Error() + "|")
	}

	for _, service := range servicesListResult.Items {
		serviceName = service.Name
		kubernetesDeleteServiceError := kubernetesDeleteService(&serviceName)
		if kubernetesDeleteServiceError != nil {
			return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesDeleteAllDeploymentsAndServices->services.List:" + kubernetesDeleteServiceError.Error() + "|")
		}

	}

	deploymentsClient := clientConfig.AppsV1().Deployments(apiv1.NamespaceDefault)
	deploymentListResult, deploymentListError := deploymentsClient.List(metav1.ListOptions{})

	if deploymentListError != nil {
		return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesDeleteAllDeploymentsAndServices->deployments.List:" + deploymentListError.Error() + "|")
	}

	for _, deployment := range deploymentListResult.Items {

		deploymentName = deployment.Name
		kubernetesDeleteDeploymentError := kubernetesDeleteDeployment(&deploymentName)
		if kubernetesDeleteDeploymentError != nil {
			return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesDeleteAllDeploymentsAndServices->deployments.List:" + kubernetesDeleteDeploymentError.Error() + "|")
		}

	}

	return nil

}

func kubernetesCreateDeployment(kubernetesPort int, deploymentName *string, imageName *string, kubernetesProtocol *string, replicas int) error {

	fmt.Println("Creating Kubernetes Deployment")

	replicasLabelsTemp := make(map[string]string, 0)
	replicasLabelsTemp["app"] = *deploymentName

	podsLabelsTemp := make(map[string]string, 0)
	podsLabelsTemp["app"] = *deploymentName

	deploymentsClient := clientConfig.AppsV1().Deployments(apiv1.NamespaceDefault)

	var replicasInt32 int32 = int32(replicas)

	var kubernetesPortStructure []apiv1.ContainerPort

	if kubernetesPort > 0 {
		kubernetesPortStructure = []apiv1.ContainerPort{
			{
				Name:          fmt.Sprintf("%s-%d", "port", kubernetesPort),
				Protocol:      apiv1.Protocol(*kubernetesProtocol),
				ContainerPort: int32(kubernetesPort),
			},
		}

	} else {
		kubernetesPortStructure = []apiv1.ContainerPort{}
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: *deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicasInt32,
			Selector: &metav1.LabelSelector{
				MatchLabels: podsLabelsTemp,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: podsLabelsTemp,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  fmt.Sprintf("%s-%s", *deploymentName, "container"),
							Image: *imageName,
							Ports: kubernetesPortStructure,
						},
					},
				},
			},
		},
	}

	_, deploymentsCreateError := deploymentsClient.Create(deployment)
	if deploymentsCreateError != nil {
		return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->kubernetesCreateDeployment->Deployments.Create:" + deploymentsCreateError.Error() + "|")
	}

	return nil
}

func kubernetesCreateService(kubernetesPort int, deploymentName *string, kubernetesProtocol *string) error {

	fmt.Println("Creating Kubernetes Service")

	podsLabelsTemp := make(map[string]string, 0)
	podsLabelsTemp["app"] = *deploymentName

	var kubernetesPortStructure []v1.ServicePort

	if kubernetesPort > 0 {
		kubernetesPortStructure = []v1.ServicePort{{
			Protocol: apiv1.Protocol(*kubernetesProtocol),
			Port:     int32(kubernetesPort),
		}}

		_, createError := clientConfig.CoreV1().Services(apiv1.NamespaceDefault).Create(&apiv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", *deploymentName, "service"),
				Namespace: apiv1.NamespaceDefault,
			},
			Spec: apiv1.ServiceSpec{
				Ports:    kubernetesPortStructure,
				Selector: podsLabelsTemp,
				Type:     apiv1.ServiceTypeLoadBalancer,
			},
		})

		if createError != nil {
			return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->kubernetesCreateService->Services.Create:" + createError.Error() + "|")
		}
	}

	return nil

}

func KubernetesDeployContainer(kubernetesPort int, deploymentName *string, imageName *string, kubernetesProtocol *string, replicas int) error {

	fmt.Println("Creating Kubernetes Deployment And Service")

	if clientConfig == nil {
		loadKubeConfigError := KubernetesLoadKubeConfig(nil)
		if loadKubeConfigError != nil {
			return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesCopyFileFromRemotePods->KubernetesLoadKubeConfig:" + loadKubeConfigError.Error() + "|")
		}

	}

	kubernetesCreateDeploymentError := kubernetesCreateDeployment(kubernetesPort, deploymentName, imageName, kubernetesProtocol, replicas)
	if kubernetesCreateDeploymentError != nil {
		return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesDeployContainer->kubernetesCreateDeployment:" + kubernetesCreateDeploymentError.Error() + "|")
	}

	kubernetesCreateServiceError := kubernetesCreateService(kubernetesPort, deploymentName, kubernetesProtocol)
	if kubernetesCreateServiceError != nil {
		return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesDeployContainer->kubernetesCreateService:" + kubernetesCreateServiceError.Error() + "|")
	}

	return nil

}

/*
func NewIOStreams() (genericclioptions.IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	return genericclioptions.IOStreams{
		In:     in,
		Out:    out,
		ErrOut: errOut,
	}, in, out, errOut
}

func TestCopyToPod(deploymentName *string, remotePodFile *string, localFile *string) error {

	tf := cmdutil.NewFactory(&genericclioptions.ConfigFlags{})
	ioStreams, _, _, _ := NewIOStreams()

	cmd := cp.NewCmdCp(tf, ioStreams)

	podsClient := clientConfig.CoreV1().Pods(apiv1.NamespaceDefault)
	podsListResult, podsListError := podsClient.List(context.Background(), metav1.ListOptions{})

	if podsListError != nil {
		return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->TestCopyToPod->services.List:" + podsListError.Error() + "|")
	}

	opts := cp.NewCopyOptions(ioStreams)

	for _, pod := range podsListResult.Items {
		for podLabelKey, podLabelValue := range pod.Labels {
			if podLabelKey == "app" && podLabelValue == *deploymentName {

				sufixedLocalPath := fmt.Sprintf("%s-%s", *localFile, pod.Name)

				src := map[string]string{
					"PodNamespace": apiv1.NamespaceDefault,
					"PodName":      pod.Name,
					"File":         *remotePodFile,
				}

				dest := map[string]string{
					"File": sufixedLocalPath,
				}

				opts.Complete(tf, cmd)
				options := &kexec.ExecOptions{}
				copyFromPodError := opts.copyFromPod(src, dest, options)
				if copyFromPodError != nil {
					return errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->TestCopyToPod->copyFromPod:" + copyFromPodError.Error() + "|")
				}

			}
		}

	}
	return nil
}
*/

func KubernetesGetPodsForDeployment(deploymentName string) (map[string]bool, error) {

	fmt.Println("Retrieving List Of Pods For Deployment")

	var podsListForDeployment map[string]bool = make(map[string]bool, 0)

	if clientConfig == nil {
		loadKubeConfigError := KubernetesLoadKubeConfig(nil)
		if loadKubeConfigError != nil {
			return nil, errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesGetPodsForDeployment->KubernetesLoadKubeConfig:" + loadKubeConfigError.Error() + "|")
		}

	}

	podsClient := clientConfig.CoreV1().Pods(apiv1.NamespaceDefault)
	podsListResult, podsListError := podsClient.List(metav1.ListOptions{})

	if podsListError != nil {
		return nil, errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesGetPodsForDeployment->services.List:" + podsListError.Error() + "|")
	}

	for _, pod := range podsListResult.Items {
		for podLabelKey, podLabelValue := range pod.Labels {
			if podLabelKey == "app" && podLabelValue == deploymentName {
				podsListForDeployment[pod.Name] = true
			}
		}
	}

	return podsListForDeployment, nil
}

/*
toPod = true

/loca/file->/remote/file


*/
func KubernetesCopyFileFromRemotePods(podsMap map[string]bool, remotePodFile string, localFolder string) (map[string]string, error) {

	fmt.Println("Copying Files From Pods")

	var podNameAndLocalFileMap map[string]string = make(map[string]string, 0)

	if clientConfig == nil {
		loadKubeConfigError := KubernetesLoadKubeConfig(nil)
		if loadKubeConfigError != nil {
			return nil, errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesCopyFileFromOrToRemotePods->KubernetesLoadKubeConfig:" + loadKubeConfigError.Error() + "|")
		}
	}

	var kubeCtlCommand string

	var kubeCtlFormat string = "kubectl cp %s:%s %s"

	podsTempList := make([]string, 0)
	for podName, _ := range podsMap {
		podsTempList = append(podsTempList, podName)
	}

	for _, podName := range podsTempList {

		fileFullLocalPath := fmt.Sprintf("%s/%s-%s", localFolder, podName, path.Base(remotePodFile)) //May have issues with Windows paths

		kubeCtlCommand = fmt.Sprintf(kubeCtlFormat, podName, remotePodFile, fileFullLocalPath)

		kubeCtlCommandSplitted := strings.Split(kubeCtlCommand, " ")

		commandResult := exec.Command(kubeCtlCommandSplitted[0], kubeCtlCommandSplitted[1], kubeCtlCommandSplitted[2], kubeCtlCommandSplitted[3])

		runError := commandResult.Run()
		if runError != nil {
			fmt.Println("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesCopyFileFromOrToRemotePods->Run:" + runError.Error() + "|")
		} else {
			podNameAndLocalFileMap[podName] = fileFullLocalPath
			//delete(podsMap, podName)
		}
	}
	return podNameAndLocalFileMap, nil
}

func KubernetesCopyFilesToRemotePods(podsMap map[string]bool, remotePodFile string, localFilesNamesMap map[string]bool) (map[string]string, error) {

	fmt.Println("Copying Files To Pods")

	var localFileAndRemotePodMap map[string]string = make(map[string]string, 0)

	if clientConfig == nil {
		loadKubeConfigError := KubernetesLoadKubeConfig(nil)
		if loadKubeConfigError != nil {
			return nil, errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesCopyFilesToRemotePods->KubernetesLoadKubeConfig:" + loadKubeConfigError.Error() + "|")
		}
	}

	var kubeCtlCommand string

	var kubeCtlFormat string = "kubectl cp %s %s:%s"

	podsTempList := make([]string, 0)
	for podName, _ := range podsMap {
		podsTempList = append(podsTempList, podName)
	}

	localFilesNamesTempList := make([]string, 0)
	for localFileName, _ := range localFilesNamesMap {
		localFilesNamesTempList = append(localFilesNamesTempList, localFileName)
	}

	for localFileNameIndex, localFileName := range localFilesNamesTempList {

		podName := podsTempList[localFileNameIndex]
		kubeCtlCommand = fmt.Sprintf(kubeCtlFormat, localFileName, podName, remotePodFile)
		fmt.Println(kubeCtlCommand)
		kubeCtlCommandSplitted := strings.Split(kubeCtlCommand, " ")
		commandResult := exec.Command(kubeCtlCommandSplitted[0], kubeCtlCommandSplitted[1], kubeCtlCommandSplitted[2], kubeCtlCommandSplitted[3])

		runError := commandResult.Run()
		if runError != nil {
			fmt.Println("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesCopyFilesToRemotePods->Run:" + runError.Error() + "|")
		} else {
			localFileAndRemotePodMap[localFileName] = podName
			//delete(podsMap, podName)
		}

	}
	return localFileAndRemotePodMap, nil
}

//Useless because pod.Status.Phase == apiv1.PodRunning verifies when in crash loop....
func KubernetesGetListOfPodsInReadyState(deploymentName string) (map[string]bool, error) {

	fmt.Println("Retrieving List Of Ready Pods For Deployment")

	var runningPodsMap map[string]bool = make(map[string]bool, 0)

	if clientConfig == nil {
		loadKubeConfigError := KubernetesLoadKubeConfig(nil)
		if loadKubeConfigError != nil {
			return nil, errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesGetListOfPodsInReadyState->KubernetesLoadKubeConfig:" + loadKubeConfigError.Error() + "|")
		}

	}

	podsClient := clientConfig.CoreV1().Pods(apiv1.NamespaceDefault)
	podsListResult, podsListError := podsClient.List(metav1.ListOptions{})

	if podsListError != nil {
		return nil, errors.New("|" + "HaymakerCF->haymakercfengines->kubernetes_engine->KubernetesGetListOfPodsInReadyState->Pods.List:" + podsListError.Error() + "|")
	}

	for _, pod := range podsListResult.Items {
		for podLabelKey, podLabelValue := range pod.Labels {
			failedConditionDetected := false
			for _, podCondition := range pod.Status.Conditions {
				if podCondition.Status == apiv1.ConditionFalse {
					failedConditionDetected = true
					break
				}
			}
			if podLabelKey == "app" && podLabelValue == deploymentName && !failedConditionDetected {
				runningPodsMap[pod.Name] = true
			}
		}
	}

	return runningPodsMap, nil

}
