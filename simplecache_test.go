package simplecache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type SimpleCacheTestSuite struct {
	suite.Suite

	cache SimpleCache
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *SimpleCacheTestSuite) SetupTest() {

}

func (suite *SimpleCacheTestSuite) TestSetAndGetAndDel() {
	suite.cache.Set("TestA", "AAAAA")
	rlt, err := String(suite.cache.Get("TestA"))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "AAAAA", rlt)

	func() {
		defer func() {
			err := recover()
			assert.NotNil(suite.T(), err)
		}()
		suite.cache.HSet("TestA", "1", "One")
	}()

	suite.cache.Del("TestA")
	rlt, err = String(suite.cache.Get("TestA"))
	assert.NotNil(suite.T(), err)
}

func (suite *SimpleCacheTestSuite) TestSetAndGetMapAndDel() {
	suite.cache.HMSet("TestA", "1", "One", "2", "Two")
	rlt, err := Strings(suite.cache.HMGet("TestA", "1", "2"))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "One", rlt[0])
	assert.Equal(suite.T(), "Two", rlt[1])

	func() {
		defer func() {
			err := recover()
			assert.NotNil(suite.T(), err)
		}()
		suite.cache.Set("TestA", "BBB")
	}()

	suite.cache.Del("TestA")
	rlt, err = Strings(suite.cache.HMGet("TestA", "1"))
	assert.NotNil(suite.T(), err)
}

func TestSimpleCacheTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleCacheTestSuite))
}
