package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CopySuite struct {
	suite.Suite
	inputFilePath  string
	tmpDir         string
	tmpOutFilePath string
}

func (c *CopySuite) SetupTest() {
	c.inputFilePath = "testdata/input.txt"
	tmpDir, err := ioutil.TempDir("", "copy")
	if err != nil {
		c.T().Fatal("can't create temp dir:", err)
	}
	c.tmpDir = tmpDir

	var builder strings.Builder
	builder.WriteString(c.tmpDir)
	builder.WriteString("out.txt")
	c.tmpOutFilePath = builder.String()
}

func (c *CopySuite) ReferenceOutput(offset, limit int64) string {
	var pathBuilder strings.Builder
	pathBuilder.WriteString("testdata/")
	pathBuilder.WriteString("out_offset")
	pathBuilder.WriteString(strconv.Itoa(int(offset)))
	pathBuilder.WriteString("_limit")
	pathBuilder.WriteString(strconv.Itoa(int(limit)))
	pathBuilder.WriteString(".txt")
	out, err := ioutil.ReadFile(pathBuilder.String())
	if err != nil {
		c.T().Fatal(err)
	}

	var outBuilder strings.Builder
	outBuilder.Write(out)

	return outBuilder.String()
}

func (c *CopySuite) CopyOutput() string {
	out, err := ioutil.ReadFile(c.tmpOutFilePath)
	if err != nil {
		c.T().Fatal(err)
	}

	var builder strings.Builder
	builder.Write(out)

	return builder.String()
}

func (c *CopySuite) TeardownTest() {
	os.RemoveAll(c.tmpDir)
}

func (c *CopySuite) TestCopySuccess() {
	// offset 0, limit 0
	err := Copy(c.inputFilePath, c.tmpOutFilePath, 0, 0)

	c.Require().NoError(err)
	c.Require().Equal(c.ReferenceOutput(0, 0), c.CopyOutput())

	// offset 0, limit 10
	err = Copy(c.inputFilePath, c.tmpOutFilePath, 0, 10)

	c.Require().NoError(err)
	c.Require().Equal(c.ReferenceOutput(0, 10), c.CopyOutput())

	// offset 0, limit 1000
	err = Copy(c.inputFilePath, c.tmpOutFilePath, 0, 1000)

	c.Require().NoError(err)
	c.Require().Equal(c.ReferenceOutput(0, 1000), c.CopyOutput())

	// offset 100, limit 1000
	err = Copy(c.inputFilePath, c.tmpOutFilePath, 100, 1000)

	c.Require().NoError(err)
	c.Require().Equal(c.ReferenceOutput(100, 1000), c.CopyOutput())

	// offset 6000, limit 1000
	err = Copy(c.inputFilePath, c.tmpOutFilePath, 6000, 1000)

	c.Require().NoError(err)
	c.Require().Equal(c.ReferenceOutput(6000, 1000), c.CopyOutput())
}

func (c *CopySuite) TestFileDoesNotExist() {
	err := Copy("testdata/NaN.txt", c.tmpOutFilePath, 0, 0)

	c.Require().Error(err)
}

func (c *CopySuite) TestInputFileIsEmpty() {
	err := Copy(c.tmpOutFilePath, "", 0, 0)

	c.Require().Error(err)
}

func (c *CopySuite) TestOffsetBiggerThanFileSize() {
	err := Copy(c.inputFilePath, c.tmpOutFilePath, 999999, 0)

	c.Require().Error(err)
}

func (c *CopySuite) TestLimitBiggerThanFileSize() {
	err := Copy(c.inputFilePath, c.tmpOutFilePath, 0, 99999999999999)

	c.Require().NoError(err)
}

func TestCopy(t *testing.T) {
	suite.Run(t, new(CopySuite))
}
