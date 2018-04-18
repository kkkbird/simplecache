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

func (suite *SimpleCacheTestSuite) testSetAndGetAndDel(key string,
	value1 interface{}, value2 interface{}) {
	suite.cache.Set(key, value1)
	rlt, err := suite.cache.Get(key)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), value1, rlt)

	suite.cache.Set(key, value2)
	rlt, err = suite.cache.Get(key)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), value2, rlt)

	func() {
		defer func() {
			err := recover()
			assert.NotNil(suite.T(), err)
		}()
		suite.cache.HSet(key, "1", "One")
	}()

	suite.cache.Del(key)
	rlt, err = String(suite.cache.Get(key))
	assert.NotNil(suite.T(), err)
}

func (suite *SimpleCacheTestSuite) TestSetAndGetAndDel() {
	suite.testSetAndGetAndDel("KeyA", "AAA", "BBB")
	suite.testSetAndGetAndDel("KeyA", 1, 2)

	a := struct{ A, B int }{1, 2}
	b := struct{ A, B int }{3, 4}
	suite.testSetAndGetAndDel("KeyA", a, b)
}

func (suite *SimpleCacheTestSuite) TestSetAndGetMapAndDel() {
	suite.cache.HSet("TestA", "1", "One")
	rlt, err := String(suite.cache.HGet("TestA", "1"))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "One", rlt)

	suite.cache.HSet("TestA", "2", "Two")
	rlts, err := Strings(suite.cache.HMGet("TestA", "1", "2"))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "One", rlts[0])
	assert.Equal(suite.T(), "Two", rlts[1])

	suite.cache.HMSet("TestA", "1", "NewOne", "2", "Two", "4", "Four")
	rlts, err = Strings(suite.cache.HMGet("TestA", "1", "2", "3", "4"))
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "NewOne", rlts[0])
	assert.Equal(suite.T(), "Two", rlts[1])
	assert.Equal(suite.T(), "", rlts[2])
	assert.Equal(suite.T(), "Four", rlts[3])

	func() {
		defer func() {
			err := recover()
			assert.NotNil(suite.T(), err)
		}()
		suite.cache.Set("TestA", "BBB")
	}()

	suite.cache.Del("TestA")
	rlts, err = Strings(suite.cache.HMGet("TestA", "1"))
	assert.NotNil(suite.T(), err)

	func() {
		defer func() {
			err := recover()
			assert.Nil(suite.T(), err)

			suite.cache.Del("TestA")
		}()
		suite.cache.Set("TestA", "BBB")
	}()
}

func TestSimpleCacheTestSuite(t *testing.T) {
	suite.Run(t, new(SimpleCacheTestSuite))
}
