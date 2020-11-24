package msgraph

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type FilterCriteriaTestSuite struct {
	suite.Suite
}

func (suite *FilterCriteriaTestSuite) SetupSuite() {
	fmt.Println("SetupSuite")
}

func (suite *FilterCriteriaTestSuite) TearDownSuite() {
	fmt.Println("TearDownSuite")
}

func (suite *FilterCriteriaTestSuite) SetupTest() {
	fmt.Println("SetupTest")
}

func (suite *FilterCriteriaTestSuite) TearDownTest() {
	fmt.Println("TearDownTest")
}

func (suite *FilterCriteriaTestSuite) TestMyFunc() {
	fmt.Println("TestMyFunc")

	filter := new(FilterCriteria)
	criteria := filter.LogicOr(filter.StartWith("field1", "startwith1"), filter.StartWith("field2", "startwith2"))
	assert.Equal(
		suite.T(),
		"startwith(field1,'startwith1') OR startwith(field2,'startwith2')",
		(*criteria).String(),
	)
}
func (suite *FilterCriteriaTestSuite) TestMyFunc2() {
	fmt.Println("TestMyFunc2")
}

func TestMyTestSuite(t *testing.T) {
	tests := new(FilterCriteriaTestSuite)
	suite.Run(t, tests)
}
