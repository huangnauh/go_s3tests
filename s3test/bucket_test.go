package s3test

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/huangnauh/go_s3tests/helpers"
)

func (suite *S3Suite) TestBucketCreateReadDelete() {

	/*
		Resource : bucket, method: create/delete
		Scenario : create and delete bucket.
		Assertion: bucket exists after create and is gone after delete.
	*/

	assert := suite
	bucket := helpers.GetBucketName()

	err := helpers.CreateBucket(svc, bucket)
	assert.Nil(err)

	bkts, err := helpers.ListBuckets(svc)
	assert.Equal(true, helpers.Contains(bkts, bucket))

	err = helpers.DeleteBucket(svc, bucket)

	//ensure it doesnt exist
	err = helpers.DeleteBucket(svc, bucket)
	assert.NotNil(err)

	awsErr, ok := err.(awserr.Error)
	assert.True(ok)
	assert.Equal(awsErr.Code(), "NoSuchBucket")
}

func (suite *S3Suite) TestBucketDeleteNotExist() {

	/*
		Resource : bucket, method: delete
		Scenario : non existant bucket
		Assertion: fails NoSuchBucket.
	*/

	assert := suite
	bucket := helpers.GetBucketName()

	err := helpers.DeleteBucket(svc, bucket)
	assert.NotNil(err)

	awsErr, ok := err.(awserr.Error)
	assert.True(ok)
	assert.Equal(awsErr.Code(), "NoSuchBucket")
}

func (suite *S3Suite) TestBucketDeleteNotEmpty() {

	/*
		Resource : bucket, method: delete
		Scenario : bucket not empty
		Assertion: fails BucketNotEmpty.
	*/

	assert := suite
	bucket := helpers.GetBucketName()
	objects := map[string]string{"key1": "echo"}

	err := helpers.CreateBucket(svc, bucket)
	assert.Nil(err)

	err = helpers.CreateObjects(svc, bucket, objects)

	err = helpers.DeleteBucket(svc, bucket)
	assert.NotNil(err)

	awsErr, ok := err.(awserr.Error)
	assert.True(ok)
	assert.Equal(awsErr.Code(), "BucketNotEmpty")
}

func (suite *S3Suite) TestBucketListEmpty() {

	/*
		Resource : object, method: list
		Scenario : bucket not empty
		Assertion: empty buckets return no contents.
	*/

	assert := suite
	bucket := helpers.GetBucketName()
	var empty_list []*s3.Object

	err := helpers.CreateBucket(svc, bucket)
	assert.Nil(err)

	resp, err := helpers.GetObjects(svc, bucket)
	assert.Nil(err)
	assert.Equal(empty_list, resp.Contents)
}

func (suite *S3Suite) TestBucketListDistinct() {

	/*
		Resource : object, method: list
		Scenario : bucket not empty
		Assertion: distinct buckets have different contents.
	*/

	assert := suite
	bucket1 := helpers.GetBucketName()
	bucket2 := helpers.GetBucketName()
	objects1 := map[string]string{"key1": "Hello"}
	objects2 := map[string]string{"key2": "Manze"}

	err := helpers.CreateBucket(svc, bucket1)
	err = helpers.CreateBucket(svc, bucket2)
	assert.Nil(err)

	err = helpers.CreateObjects(svc, bucket1, objects1)
	err = helpers.CreateObjects(svc, bucket2, objects2)

	obj1, _ := helpers.GetObject(svc, bucket1, "key1")
	obj2, _ := helpers.GetObject(svc, bucket2, "key2")

	assert.Equal(obj1, "Hello")
	assert.Equal(obj2, "Manze")

}

func (suite *S3Suite) TestObjectAclCreateContentlengthNone() {

	/*
		Resource : bucket, method: acls
		Scenario :set w/no content length.
		Assertion: suceeds
	*/

	assert := suite
	conLength := map[string]string{"Content-Length": ""}
	// acl := map[string]string{"ACL": "public-read"}
	content := "bar"

	bucket := helpers.GetBucketName()
	key := "key1"
	err := helpers.CreateBucket(svc, bucket)

	err = helpers.SetupObjectWithHeader(svc, bucket, key, content, conLength)
	_, err = helpers.SetACL(svc, bucket, "public-read")
	assert.Nil(err)
}

func (suite *S3Suite) TestBucketPutCanned_acl() {

	/*
		Resource : bucket, method: put
		Scenario :set w/invalid permission.
		Assertion: fails
	*/

	assert := suite

	bucket := helpers.GetBucketName()
	err := helpers.CreateBucket(svc, bucket)

	_, err = helpers.SetACL(svc, bucket, "public-ready")
	_, err = helpers.SetACL(svc, bucket, "public-read")
	assert.Nil(err)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			assert.Equal(awsErr.Code(), "AccessDenied")
			assert.Equal(awsErr.Message(), "")
		}
	}
}

func (suite *S3Suite) TestBucketCreateBadExpectMismatch() {

	/*
		Resource : bucket, method: put
		Scenario :create w/expect 200.
		Assertion: garbage, but S3 succeeds!
	*/

	assert := suite
	acl := map[string]string{"Expect": "200"}

	bucket := helpers.GetBucketName()
	err := helpers.CreateBucketWithHeader(svc, bucket, acl)
	assert.NotNil(err)
	awsErr, ok := err.(awserr.Error)
	assert.True(ok)
	assert.Equal(awsErr.Code(), "ExpectationFailed")
}

func (suite *S3Suite) TestBucketCreateBadExpectEmpty() {

	/*
		Resource : bucket, method: put
		Scenario :create w/expect empty.
		Assertion: garbage, but S3 succeeds!
	*/

	assert := suite
	acl := map[string]string{"Expect": " "}

	bucket := helpers.GetBucketName()

	err := helpers.CreateBucketWithHeader(svc, bucket, acl)
	assert.Nil(err)
}

func (suite *S3Suite) TestBucketCreateBadExpectUnreadable() {

	/*
		Resource : bucket, method: put
		Scenario :create w/expect nongraphic.
		Assertion: fails with "invalid header field value"
	*/

	assert := suite
	acl := map[string]string{"Expect": "\x07"}

	bucket := helpers.GetBucketName()

	err := helpers.CreateBucketWithHeader(svc, bucket, acl)

	assert.NotNil(err)
}

func (suite *S3Suite) TestBucketCreateBadContentLengthEmpty() {

	/*
		Resource : bucket, method: put
		Scenario :create w/empty content length.
		Assertion: fails
	*/

	assert := suite
	acl := map[string]string{"Content-Length": " "}

	bucket := helpers.GetBucketName()

	err := helpers.CreateBucketWithHeader(svc, bucket, acl)
	assert.NotNil(err)
	awsErr, ok := err.(awserr.Error)
	assert.True(ok)
	assert.Equal(awsErr.Code(), "MissingContentLength")
}

func (suite *S3Suite) TestBucketCreateBadContentlengthNegative() {

	/*
		Resource : bucket, method: put
		Scenario :create w/negative content length.
		Assertion: fails
	*/

	assert := suite
	acl := map[string]string{"Content-Length": "-1"}

	bucket := helpers.GetBucketName()

	err := helpers.CreateBucketWithHeader(svc, bucket, acl)
	assert.NotNil(err)
	awsErr, ok := err.(awserr.Error)
	assert.True(ok)
	assert.Equal(awsErr.Code(), "MissingContentLength")
}

func (suite *S3Suite) TestBucketCreateBadContentlengthNone() {

	/*
		Resource : bucket, method: put
		Scenario :create w/no content length.
		Assertion: suceeds
	*/

	assert := suite
	acl := map[string]string{"Content-Length": ""}

	bucket := helpers.GetBucketName()

	err := helpers.CreateBucketWithHeader(svc, bucket, acl)
	assert.Nil(err)
}

func (suite *S3Suite) TestBucket_CreateBadContentlengthUnreadable() {

	/*
		Resource : bucket, method: put
		Scenario :create w/unreadable content length.
		Assertion: fails
	*/

	assert := suite
	acl := map[string]string{"Content-Length": "\x07"}

	bucket := helpers.GetBucketName()

	err := helpers.CreateBucketWithHeader(svc, bucket, acl)
	assert.NotNil(err)
	awsErr, ok := err.(awserr.Error)
	assert.True(ok)
	assert.Equal(awsErr.Code(), "MissingContentLength")
}

//TODO:
// func (suite *S3Suite) TestBucketCreateBadAuthorizationUnreadable() {

// 	/*
// 		Resource : bucket, method: put
// 		Scenario :create w/non-graphic authorization.
// 		Assertion: expected to fail..but suceeded
// 	*/

// 	assert := suite
// 	acl := map[string]string{"Authorization": "\x07"}

// 	bucket := helpers.GetBucketName()

// 	err := helpers.CreateBucketWithHeader(svc, bucket, acl)
// 	assert.Nil(err)
// 	if err != nil {
// 		if awsErr, ok := err.(awserr.Error); ok {

// 			assert.Equal(awsErr.Code(), "AccessDenied")
// 			assert.Equal(awsErr.Message(), "")
// 		}
// 	}
// }

//TODO:
// func (suite *S3Suite) TestBucketCreateBadAuthorizationEmpty() {

// 	/*
// 		Resource : bucket, method: put
// 		Scenario :create w/empty authorization.
// 		Assertion: expected to fail..but suceeded
// 	*/

// 	assert := suite

// 	acl := map[string]string{"Authorization": " "}

// 	bucket := helpers.GetBucketName()
// 	err := helpers.CreateBucket(svc, bucket)

// 	err = helpers.CreateBucketWithHeader(svc, bucket, acl)
// 	assert.Nil(err)
// 	if err != nil {
// 		if awsErr, ok := err.(awserr.Error); ok {

// 			assert.Equal(awsErr.Code(), "AccessDenied")
// 			assert.Equal(awsErr.Message(), "")
// 		}
// 	}
// }

//TODO:
// func (suite *S3Suite) TestBucketCreateBadAuthorizationNone() {

// 	/*
// 		Resource : bucket, method: put
// 		Scenario :create w/no authorization.
// 		Assertion: expected to fail..but suceeded
// 	*/

// 	assert := suite
// 	acl := map[string]string{"Authorization": ""}

// 	bucket := helpers.GetBucketName()
// 	err := helpers.CreateBucket(svc, bucket)

// 	err = helpers.CreateBucketWithHeader(svc, bucket, acl)
// 	assert.Nil(err)
// 	if err != nil {
// 		if awsErr, ok := err.(awserr.Error); ok {

// 			assert.Equal(awsErr.Code(), "AccessDenied")
// 			assert.Equal(awsErr.Message(), "")
// 		}
// 	}
// }

//TODO:
// func (suite *S3Suite) TestLifecycleGetNoLifecycle() {

// 	/*
// 		Resource : bucket, method: get
// 		Scenario : get lifecycle config that has not been set.
// 		Assertion: fails
// 	*/

// 	assert := suite
// 	//acl := map[string]string{"Authorization": ""}

// 	bucket := helpers.GetBucketName()
// 	err := helpers.CreateBucket(svc, bucket)

// 	_, err = helpers.GetLifecycle(svc, bucket)
// 	assert.NotNil(err)
// 	if err != nil {
// 		if awsErr, ok := err.(awserr.Error); ok {

// 			assert.Equal(awsErr.Code(), "NoSuchLifecycleConfiguration")
// 			assert.Equal(awsErr.Message(), "")
// 		}
// 	}
// }

//TODO:
// func (suite *S3Suite) TestLifecycleInvalidMD5() {

// 	/*
// 		Resource : bucket, method: get
// 		Scenario : set lifecycle config with invalid md5.
// 		Assertion: fails
// 	*/

// 	assert := suite

// 	bucket := helpers.GetBucketName()
// 	err := helpers.CreateBucket(svc, bucket)

// 	content := strings.NewReader("Enabled")
// 	h := md5.New()
// 	content.WriteTo(h)
// 	sum := h.Sum(nil)
// 	b := make([]byte, base64.StdEncoding.EncodedLen(len(sum)))
// 	base64.StdEncoding.Encode(b, sum)

// 	md5 := string(b)

// 	_, err = helpers.SetLifecycle(svc, bucket, "", "Enabled", md5)
// 	assert.NotNil(err)
// 	if err != nil {
// 		if awsErr, ok := err.(awserr.Error); ok {
// 			assert.Equal(awsErr.Code(), "MalformedXML")
// 			assert.Equal(awsErr.Message(), "")
// 		}
// 	}
// }

//TODO:
// func (suite *S3Suite) TestLifecycleInvalidStatus() {

// 	/*
// 		Resource : bucket, method: get
// 		Scenario : invalid status in lifecycle rule.
// 		Assertion: fails
// 	*/

// 	assert := suite

// 	bucket := helpers.GetBucketName()
// 	err := helpers.CreateBucket(svc, bucket)

// 	content := strings.NewReader("Enabled")
// 	h := md5.New()
// 	content.WriteTo(h)
// 	sum := h.Sum(nil)
// 	b := make([]byte, base64.StdEncoding.EncodedLen(len(sum)))
// 	base64.StdEncoding.Encode(b, sum)

// 	md5 := string(b)

// 	_, err = helpers.SetLifecycle(svc, bucket, "rule1", "enabled", md5)
// 	assert.NotNil(err)
// 	if err != nil {
// 		if awsErr, ok := err.(awserr.Error); ok {

// 			assert.Equal(awsErr.Code(), "MalformedXML")
// 			assert.Equal(awsErr.Message(), "")
// 		}
// 	}

// 	_, err = helpers.SetLifecycle(svc, bucket, "rule1", "disabled", md5)
// 	assert.NotNil(err)
// 	if err != nil {
// 		if awsErr, ok := err.(awserr.Error); ok {

// 			assert.Equal(awsErr.Code(), "MalformedXML")
// 			assert.Equal(awsErr.Message(), "")
// 		}
// 	}

// 	_, err = helpers.SetLifecycle(svc, bucket, "rule1", "invalid", md5)
// 	assert.NotNil(err)
// 	if err != nil {
// 		if awsErr, ok := err.(awserr.Error); ok {

// 			assert.Equal(awsErr.Code(), "MalformedXML")
// 			assert.Equal(awsErr.Message(), "")
// 		}
// 	}
// }
