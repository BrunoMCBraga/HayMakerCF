package haymakercfengines

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ecr"
)

var ecrInstance *ecr.ECR
var docketTag string = "latest"

func describeRepositoriesStub(repositoryNames []*string) (*ecr.DescribeRepositoriesOutput, error) {

	describeRepositoryInputObject := &ecr.DescribeRepositoriesInput{
		RepositoryNames: repositoryNames,
	}

	return ecrInstance.DescribeRepositories(describeRepositoryInputObject)
}

func listImagesStub(repositoryName *string) (*ecr.ListImagesOutput, error) {

	listImagesInputObject := &ecr.ListImagesInput{

		RepositoryName: repositoryName,
	}

	return ecrInstance.ListImages(listImagesInputObject)
}

func getAuthorizationTokenStub(registryIds []*string) (*ecr.GetAuthorizationTokenOutput, error) {

	var describeRepositoryInputObject *ecr.GetAuthorizationTokenInput

	if len(registryIds) > 0 {
		describeRepositoryInputObject = &ecr.GetAuthorizationTokenInput{
			RegistryIds: registryIds,
		}
	} else {
		describeRepositoryInputObject = &ecr.GetAuthorizationTokenInput{}
	}

	return ecrInstance.GetAuthorizationToken(describeRepositoryInputObject)
}

func batchDeleteImageStub(repositoryName *string, imageIDs []*ecr.ImageIdentifier) (*ecr.BatchDeleteImageOutput, error) {

	var batchDeleteImageInputObject *ecr.BatchDeleteImageInput

	if len(imageIDs) > 0 {
		batchDeleteImageInputObject = &ecr.BatchDeleteImageInput{
			ImageIds:       imageIDs,
			RepositoryName: repositoryName,
		}
	} else {
		batchDeleteImageInputObject = &ecr.BatchDeleteImageInput{
			RepositoryName: repositoryName,
		}
	}

	return ecrInstance.BatchDeleteImage(batchDeleteImageInputObject)
}

func ECRBatchDeleteDockerImagesFromRepository(repoName *string) error {

	fmt.Println("Destroying Docker Image On ECR")

	listImagesStubResult, listImagesStubError := listImagesStub(repoName)
	if listImagesStubError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->->ecr_engine->ECRBatchDeleteDockerImagesFromRepository->listImagesStub:" + listImagesStubError.Error() + "|")
	}

	if len(listImagesStubResult.ImageIds) > 0 {
		_, batchDeleteImageError := batchDeleteImageStub(repoName, listImagesStubResult.ImageIds)

		if batchDeleteImageError != nil {
			return errors.New("|" + "HayMakerCF->haymakercfengines->->ecr_engine->ECRBatchDeleteDockerImagesFromRepository->batchDeleteImageStub:" + batchDeleteImageError.Error() + "|")
		}
	}
	/*
		imageIDs := []*ecr.ImageIdentifier{}
		imageIDs = append(imageIDs, &ecr.ImageIdentifier{
			ImageTag: &docketTag,
		})
		imageIDs = append(imageIDs, &ecr.ImageIdentifier{
			ImageTag: &docketTag,
		})

		_, batchDeleteImageError := batchDeleteImageStub(repoName, imageIDs)

		if batchDeleteImageError != nil {
			return errors.New("|" + "HayMakerCF->haymakercfengines->->ecr_engine->ECRBatchDeleteDockerImagesFromRepository->batchDeleteImageStub:" + batchDeleteImageError.Error() + "|")
		}*/

	return nil
}

func ECRGetAuthorizationToken(registryId string) (map[string]*string, error) {

	fmt.Println("Obtaining ECR Authorization Token")

	registryIdsList := make([]*string, 0)
	if registryId != "" {
		registryIdsList = append(registryIdsList, &registryId)
	}

	getAuthorizationTokenStubResult, getAuthorizationTokenStubError := getAuthorizationTokenStub(registryIdsList)

	if getAuthorizationTokenStubError != nil {
		return nil, errors.New("|" + "HayMakerCF->haymakercfengines->->ecr_engine->ECRGetAuthorizationToken->getAuthorizationTokenStub:" + getAuthorizationTokenStubError.Error() + "|")
	}

	authorizationTokenStruct := make(map[string]*string, 0)

	if len(getAuthorizationTokenStubResult.AuthorizationData) > 0 {
		authorizationTokenStruct["token"] = getAuthorizationTokenStubResult.AuthorizationData[0].AuthorizationToken
		protocolStrippedEndpoing := strings.TrimLeft(*getAuthorizationTokenStubResult.AuthorizationData[0].ProxyEndpoint, "https://")
		authorizationTokenStruct["endpoint"] = &protocolStrippedEndpoing
		return authorizationTokenStruct, nil
	} else {
		return nil, errors.New("|" + "HayMakerCF->haymakercfengines->->ecr_engine->ECRGetAuthorizationToken: no repository found.")
	}

}

//Returns the first one one the list. The repository name should be specific enough.
func ECRGetRepositoryURI(repoName string, registryId string) (string, error) {

	repoNames := []*string{&repoName}
	var repositoryURI string = ""

	describeRepositoriesStubResult, describeRepositoriesStubError := describeRepositoriesStub(repoNames)

	if describeRepositoriesStubError != nil {
		return repositoryURI, errors.New("|" + "HayMakerCF->haymakercfengines->->ecr_engine->ECRGetRepositoryURI->describeRepositoriesStub:" + describeRepositoriesStubError.Error() + "|")
	}

	if len(describeRepositoriesStubResult.Repositories) > 0 {
		if registryId == "" {
			repositoryURI = *describeRepositoriesStubResult.Repositories[0].RepositoryUri
		} else {
			for _, repository := range describeRepositoriesStubResult.Repositories {
				if *repository.RegistryId == registryId {
					repositoryURI = *describeRepositoriesStubResult.Repositories[0].RepositoryUri
					break
				}
			}
		}

	} else {
		return repositoryURI, errors.New("|" + "HayMakerCF->haymakercfengines->->ecr_engine->ECRGetRepositoryURI: no repository found.")
	}

	return repositoryURI, nil

}

func ECRGetRepositoryURIWithDefaultRegistry(repoName string) (string, error) {

	eCRGetRepositoryURIResult, eCRGetRepositoryURIError := ECRGetRepositoryURI(repoName, "")

	if eCRGetRepositoryURIError != nil {
		return "", errors.New("|" + "HayMakerCF->haymakercfengines->->ecr_engine->ECRGetRepositoryURI->describeRepositoriesStub:" + eCRGetRepositoryURIError.Error() + "|")
	} else {
		return eCRGetRepositoryURIResult, nil
	}

}

func ECRGetRepositoriesRegistriesAndURIs(repoName string) (map[string]string, error) {

	repoNames := []*string{&repoName}
	registryAndURI := make(map[string]string, 0)

	describeRepositoriesStubResult, describeRepositoriesStubError := describeRepositoriesStub(repoNames)

	if describeRepositoriesStubError != nil {
		return nil, errors.New("|" + "HayMakerCF->haymakercfengines->->ecr_engine->ECRGetRepositoriesRegistriesAndURIs->describeRepositoriesStub:" + describeRepositoriesStubError.Error() + "|")
	}

	for _, repository := range describeRepositoriesStubResult.Repositories {
		registryAndURI[*repository.RegistryId] = *repository.RepositoryUri
	}

	return registryAndURI, nil

}

func InitECREngine(ecrInst *ecr.ECR) {
	ecrInstance = ecrInst

}
