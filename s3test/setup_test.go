package s3test

import (
	"testing"

	"github.com/huangnauh/go_s3tests/helpers"
	"github.com/stretchr/testify/suite"
)

var svc = helpers.GetConn()

type S3Suite struct {
	suite.Suite
}

func (suite *S3Suite) SetupTest() {

}

type HeadSuite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {

	suite.Run(t, new(HeadSuite))
	suite.Run(t, new(S3Suite))
}

func (suite *S3Suite) TearDownTest() {

	helpers.DeletePrefixedBuckets(svc)
}

func (suite *HeadSuite) TearDownTest() {

	helpers.DeletePrefixedBuckets(svc)
}
