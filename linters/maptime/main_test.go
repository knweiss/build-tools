// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"go/token"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTimeMapLint(t *testing.T) {
	type lintError struct {
		lineNumber int
		keyStr     string
		valStr     string
	}

	expectedLintErrors := map[string][]lintError{
		"test_file_1.go": []lintError{
			lintError{
				lineNumber: 18,
				keyStr:     "time.Time",
				valStr:     "bool",
			},
			lintError{
				lineNumber: 19,
				keyStr:     "time.Time",
				valStr:     "bool",
			},
		},
		"test_file_2.go": []lintError{
			lintError{
				lineNumber: 7,
				keyStr:     "testdata.timeAlias",
				valStr:     "bool",
			},
			lintError{
				lineNumber: 8,
				keyStr:     "testdata.timeAlias",
				valStr:     "bool",
			},
		},
		"test_file_3.go": []lintError{
			lintError{
				lineNumber: 5,
				keyStr:     "time.Time",
				valStr:     "bool",
			},
		},
		"test_file_4.go": []lintError{
			lintError{
				lineNumber: 4,
				keyStr:     "testdata.timeAlias",
				valStr:     "bool",
			},
		},
		"test_file_5.go": []lintError{
			lintError{
				lineNumber: 9,
				keyStr:     "testdata.structWithInnerTime",
				valStr:     "bool",
			},
			lintError{
				lineNumber: 10,
				keyStr:     "testdata.structWithInnerTime",
				valStr:     "bool",
			},
		},
		"test_file_6.go": []lintError{
			lintError{
				lineNumber: 7,
				keyStr:     "testdata.chanTime",
				valStr:     "bool",
			},
			lintError{
				lineNumber: 8,
				keyStr:     "testdata.chanTime",
				valStr:     "bool",
			},
		},
		"test_file_7.go": []lintError{
			lintError{
				lineNumber: 8,
				keyStr:     "time.Time",
				valStr:     "bool",
			},
			lintError{
				lineNumber: 9,
				keyStr:     "time.Time",
				valStr:     "bool",
			},
		},
	}

	observedLintErrors := map[string][]lintError{}
	testCallback := func(position token.Position, keyStr, valStr string) {
		filePath := position.Filename
		filePathBase := path.Base(filePath)
		_, ok := expectedLintErrors[filePathBase]
		require.True(t, ok, fmt.Sprintf("Failed for file: %s", filePathBase))
		observedLintErrorsForFile, _ := observedLintErrors[filePathBase]
		observedLintErrors[filePathBase] = append(
			observedLintErrorsForFile,
			lintError{
				lineNumber: position.Line,
				keyStr:     keyStr,
				valStr:     valStr,
			},
		)
	}
	handleImportPaths([]string{"./testdata"}, []string{"included"}, testCallback)

	// Make sure all observed errors were expected
	for file, observedErrs := range observedLintErrors {
		expectedErrs, ok := expectedLintErrors[file]
		require.True(t, ok, fmt.Sprintf("Failed for file: %s", file))
		require.Equal(t, expectedErrs, observedErrs, fmt.Sprintf("Failed for file: %s", file))
	}

	// Make sure all expected errors were observed
	for file, expectedErrs := range expectedLintErrors {
		observedErrs, ok := observedLintErrors[file]
		require.True(t, ok, fmt.Sprintf("Failed for file: %s", file))
		require.Equal(t, observedErrs, expectedErrs, fmt.Sprintf("Failed for file: %s", file))
	}
}
