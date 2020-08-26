package haymakercfengines

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var s3Instance *s3.S3
var s3Uploader *s3manager.Uploader

func createBucketStub(bucketName *string) (*s3.CreateBucketOutput, error) {

	createBucketInputObject := &s3.CreateBucketInput{
		Bucket: aws.String(*bucketName),
		ACL:    aws.String(s3.BucketCannedACLPrivate),
	}

	return s3Instance.CreateBucket(createBucketInputObject)
}

func deleteBucketStub(bucketName *string) (*s3.DeleteBucketOutput, error) {

	deleteBucketInputObject := &s3.DeleteBucketInput{
		Bucket: aws.String(*bucketName),
	}

	return s3Instance.DeleteBucket(deleteBucketInputObject)
}

func deleteObjectStub(bucketName *string, key *string) (*s3.DeleteObjectOutput, error) {

	deleteObjectInputObject := &s3.DeleteObjectInput{
		Bucket: aws.String(*bucketName),
		Key:    aws.String(*key),
	}

	return s3Instance.DeleteObject(deleteObjectInputObject)
}

func uploadStub(filePath *string, bucketName *string, key *string) (*s3manager.UploadOutput, error) {

	fileHandle, openError := os.Open(*filePath)
	if openError != nil {
		return nil, errors.New("|" + "HayMakerCF->haymakercfengines->s3_engine->UploadFile->os.Open:" + openError.Error() + "|")
	}

	// Upload the file to S3.
	uploadResult, uploadError := s3Uploader.Upload(&s3manager.UploadInput{
		ACL:    aws.String(s3.ObjectCannedACLPrivate),
		Bucket: aws.String(*bucketName),
		Key:    aws.String(*key),
		Body:   fileHandle,
	})

	if uploadError != nil {
		return nil, errors.New("|" + "HayMakerCF->haymakercfengines->s3_engine->UploadFile->Upload:" + uploadError.Error() + "|")
	}

	return uploadResult, nil
}

func S3UploadFileToBucket(filePath *string, bucketName *string, fileKey *string) (*string, error) {

	fmt.Println("Uploading File To S3 Bucket")

	uploadStubResult, uploadStubError := uploadStub(filePath, bucketName, fileKey)
	if uploadStubError != nil {
		return nil, errors.New("|" + "HayMakerCF->haymakercfengines->s3_engine->S3UploadFileToBucket->uploadStub:" + uploadStubError.Error() + "|")
	}

	return &uploadStubResult.Location, nil

}

func S3CreateBucket(bucketName *string) error {

	fmt.Println("Creating S3 Bucket")

	_, createBucketStubError := createBucketStub(bucketName)
	if createBucketStubError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->s3_engine->S3CreateBucket->createBucketStub:" + createBucketStubError.Error() + "|")
	}

	return nil
}

func S3DeleteBucket(bucketName *string) error {

	fmt.Println("Deleting S3 Bucket")

	// Setup BatchDeleteIterator to iterate through a list of objects.
	newDeleteListIteratorObject := s3manager.NewDeleteListIterator(s3Instance, &s3.ListObjectsInput{
		Bucket: aws.String(*bucketName),
	})

	// Traverse iterator deleting each object
	deleteError := s3manager.NewBatchDeleteWithClient(s3Instance).Delete(aws.BackgroundContext(), newDeleteListIteratorObject)

	if deleteError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->s3_engine->S3DeleteBucket->Delete:" + deleteError.Error() + "|")
	}

	_, deleteBucketStubError := deleteBucketStub(bucketName)
	if deleteBucketStubError != nil {
		return errors.New("|" + "HayMakerCF->haymakercfengines->s3_engine->S3DeleteBucket->deleteBucketStub:" + deleteBucketStubError.Error() + "|")
	}

	return nil

}

func InitS3Engine(s3Inst *s3.S3, s3Up *s3manager.Uploader) {

	s3Instance = s3Inst
	s3Uploader = s3Up

}
