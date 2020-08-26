package haymakercfengines

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
	"golang.org/x/net/context"
)

var dockerHost string = "unix:///var/run/docker.sock"
var dockerVersion string = "1.40"

var tarTempPath string = "/tmp"
var tarName string = "docker.tar"
var defaultDockerTag string = "latest"

func createTarArchiveForDocker(tarPath *string, tarName *string, folderToAdd *string) error {

	tar := new(archivex.TarFile)
	createError := tar.Create(fmt.Sprintf("%s/%s", *tarPath, *tarName))
	if createError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->createTarArchiveForDocker->tar.Create:" + createError.Error() + "|")
	}

	addAllError := tar.AddAll(*folderToAdd, false)
	if addAllError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->createTarArchiveForDocker->tar.AddAll:" + addAllError.Error() + "|")
	}
	defer tar.Close()

	return nil
}

func generateAuthString(base64Token *string, regAddress *string) (*string, error) {

	decodeStringResult, decodeStringError := base64.URLEncoding.DecodeString(*base64Token)
	if decodeStringError != nil {
		return nil, errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->generateAuthString->base64.URLEncoding.DecodeString:" + decodeStringError.Error() + "|")
	}
	usernamAndPassword := strings.Split(string(decodeStringResult), ":")

	authConfig := types.AuthConfig{
		Username:      usernamAndPassword[0],
		Password:      usernamAndPassword[1],
		ServerAddress: fmt.Sprintf("https://%s", *regAddress),
	}

	marshallResult, marshallError := json.Marshal(authConfig)
	if marshallError != nil {
		return nil, errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->generateAuthString->json.Marshal:" + marshallError.Error() + "|")
	}

	authStr := base64.URLEncoding.EncodeToString(marshallResult)

	return &authStr, nil
}

func writeToLog(reader io.ReadCloser) error {
	defer reader.Close()
	rd := bufio.NewReader(reader)
	for {
		n, _, err := rd.ReadLine()
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fmt.Println(string(n))
	}
	return nil
}

func buildImageFromDockerfile(dockHost *string, dockVersion *string, dockFilePath *string, tarPath *string, dockerImageName *string, dockerTag *string, remoteDockerImageName *string) error {

	fmt.Println("Building Docker Image From Dockerfile")

	ctx := context.Background()
	newClientWithOptsResult, newClientWithOptsError := client.NewClient(*dockHost, *dockVersion, nil, nil)
	if newClientWithOptsError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->buildImageFromDockerfile->client.NewClient:" + newClientWithOptsError.Error() + "|")
	}

	openResult, openError := os.Open(*tarPath)
	if openError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->buildImageFromDockerfile->os.Open:" + openError.Error() + "|")

	}
	defer openResult.Close()

	localNameAndTag := fmt.Sprintf("%s:%s", *dockerImageName, *dockerTag)
	tags := []string{localNameAndTag}

	imageBuildResult, imageBuildError := newClientWithOptsResult.ImageBuild(ctx, openResult, types.ImageBuildOptions{
		Tags:        tags,
		ForceRemove: true,
		Remove:      true,
	})
	if imageBuildError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->buildImageFromDockerfile->ImageBuild:" + imageBuildError.Error() + "|")
	}

	//I have to put this here, otherwise, it will fail.
	io.Copy(ioutil.Discard, imageBuildResult.Body)
	defer imageBuildResult.Body.Close()

	imageTagError := newClientWithOptsResult.ImageTag(ctx, localNameAndTag, *remoteDockerImageName)
	if imageTagError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->buildImageFromDockerfile->ImageTag:" + imageTagError.Error() + "|")
	}

	return nil
}

func pushDockerImageToRemoteRepository(dockHost *string, dockVersion *string, authToken *string, regAddress *string, imagName *string) error {

	fmt.Println("Pushing Docker Image To ECR")

	ctx := context.Background()
	newClientWithOptsResult, newClientWithOptsError := client.NewClient(*dockHost, *dockVersion, nil, nil)
	if newClientWithOptsError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->pushDockerImageToRemoteRepository->client.NewClient:" + newClientWithOptsError.Error() + "|")
	}

	generateAuthStringResult, generateAuthStringError := generateAuthString(authToken, regAddress)

	if generateAuthStringError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->pushDockerImageToRemoteRepository->generateAuthString:" + generateAuthStringError.Error() + "|")
	}

	imagePushResult, imagePushError := newClientWithOptsResult.ImagePush(ctx, *imagName, types.ImagePushOptions{
		All:          true,
		RegistryAuth: *generateAuthStringResult,
	})

	if imagePushError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->pushDockerImageToRemoteRepository->ImagePush:" + imagePushError.Error() + "|")
	}

	//I have to put this here, otherwise, it will fail.
	io.Copy(ioutil.Discard, imagePushResult)

	defer imagePushResult.Close()

	return nil
}

func imageRemoveStub(dockHost *string, dockVersion *string, imageName *string) error {

	fmt.Println("Removing Local Docker Image")

	ctx := context.Background()
	newClientWithOptsResult, newClientWithOptsError := client.NewClient(*dockHost, *dockVersion, nil, nil)
	if newClientWithOptsError != nil {
		return errors.New("|" + "HayMakerCF->haymakerengines->docker_engine->imageRemoveStub->NewClient:" + newClientWithOptsError.Error() + "|")
	}

	_, imageRemoveError := newClientWithOptsResult.ImageRemove(ctx, *imageName, types.ImageRemoveOptions{
		Force: true,
	})

	if imageRemoveError != nil {
		return errors.New("|" + "HayMakerCF->haymakerengines->docker_engine->imageRemoveStub->ImageRemove:" + imageRemoveError.Error())
	}

	return nil
}

func DockerDeleteStagingImages(dockerImageName *string, registryAddress *string) error {

	localDockerImageNameWithTagAndRemoteURL := fmt.Sprintf("%s/%s:%s", *registryAddress, *dockerImageName, defaultDockerTag)

	imageRemoveStubError := imageRemoveStub(&dockerHost, &dockerVersion, &localDockerImageNameWithTagAndRemoteURL)
	if imageRemoveStubError != nil {
		return errors.New("|" + "HayMakerCF->haymakerengines->docker_engine->DockerDeleteStagingImages->imageRemoveStub(Local-Local):" + imageRemoveStubError.Error() + "|")
	}

	localNameAndTag := fmt.Sprintf("%s:%s", *dockerImageName, defaultDockerTag)

	imageRemoveStubError = imageRemoveStub(&dockerHost, &dockerVersion, &localNameAndTag)
	if imageRemoveStubError != nil {
		return errors.New("|" + "HayMakerCF->haymakerengines->docker_engine->DockerDeleteStagingImages->imageRemoveStub(Local-Remote):" + imageRemoveStubError.Error() + "|")
	}

	return nil
}

func DockerBuildImageAndPush(dockerFilePath *string, dockerImageName *string, authorizationToken *string, registryAddress *string) error {

	localDockerImageNameWithTagAndRemoteURL := fmt.Sprintf("%s/%s:%s", *registryAddress, *dockerImageName, defaultDockerTag)

	createTarArchiveForDockerError := createTarArchiveForDocker(&tarTempPath, &tarName, dockerFilePath)
	if createTarArchiveForDockerError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->DockerBuildImageAndPush->createTarArchiveForDocker:" + createTarArchiveForDockerError.Error() + "|")
	}

	dockerTarPath := fmt.Sprintf("%s/%s", tarTempPath, tarName)

	buildContainerFromDockerfileError := buildImageFromDockerfile(&dockerHost, &dockerVersion, dockerFilePath, &dockerTarPath, dockerImageName, &defaultDockerTag, &localDockerImageNameWithTagAndRemoteURL)
	if buildContainerFromDockerfileError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->DockerBuildImageAndPush->buildImageFromDockerfile:" + buildContainerFromDockerfileError.Error() + "|")
	}

	pushDockerImageToRemoteRepositoryError := pushDockerImageToRemoteRepository(&dockerHost, &dockerVersion, authorizationToken, registryAddress, &localDockerImageNameWithTagAndRemoteURL)
	if pushDockerImageToRemoteRepositoryError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->DockerBuildImageAndPush->pushDockerImageToRemoteRepository:" + pushDockerImageToRemoteRepositoryError.Error() + "|")
	}

	removeError := os.Remove(dockerTarPath)
	if removeError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->docker_engine->DockerBuildImageAndPush->os.Remove:" + removeError.Error() + "|")
	}

	return nil
}
