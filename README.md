
 ## S3 compatibility tests

This is a set of integration tests for the S3 Proxy.

It might also be useful for people implementing software that exposes an S3-like API.

The test suite only covers the REST interface and uses [GO amazon SDK](https://aws.amazon.com/sdk-for-go/) and [Golang Environment setup](https://golang.org/doc/install).

### Get the source code

Clone the repository

	git clone https://github.com/huangnauh/go_s3tests

### Edit Configuration

	cd go_s3tests

The config file should look like this:

    DEFAULT :
        host : 127.0.0.1
        port : 5200
        is_secure : false

    fixtures :
        bucket_prefix : test

    s3main :
        access_key : "9d6696bb73ace6af9dfd10e1d50250ee"
        access_secret : "1c63129ae9db9c60c3e8aa94d3e00495"
        bucket : bucket1
        region : us-east-1
        endpoint : 127.0.0.1:5200
        host : 127.0.0.1
        port : 5200
        display_name :
        email : tester@test.com
        is_secure : false
        SSE : aws:kms
        kmskeyid : testkey-1


#### Test dependencies
	cd
	go get -v -d ./...
	go get -v github.com/stretchr/testify

### Run the Tests

To run all tests:

	cd s3tests
	go test -v

To run a specific test e.g. TestSignWithBodyReplaceRequestBody():

	cd s3tests
	go test -v -run TestSuite/TestSignWithBodyReplaceRequestBody
