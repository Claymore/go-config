#!/bin/bash

# Author: slowpoke <proxypoke at lavabit dot com>
#
# Copying and distribution of this file, with or without modification,
# are permitted in any medium without royalty provided the copyright
# notice and this notice are preserved.  This file is offered as-is,
# without any warranty.
#
# A pre-commit hook for go projects. In addition to the standard
# checks from the sample hook, it builds the project with go build,
# runs the tests (if any), formats the source code with go fmt, and
# finally go vet to make sure only correct and good code is committed.
#
# Take note that the ASCII filename check of the standard hook was
# removed. Go is unicode, and so should be the filenames. Stop using
# obsolete operating systems without proper Unicode support.

if git rev-parse --verify HEAD >/dev/null 2>&1
then
	against=HEAD
else
	# Initial commit: diff against an empty tree object
	against=4b825dc642cb6eb9a060e54bf8d69288fbee4904
fi

# If there are no go files, it makes no sense to run the other commands
# (and indeed, go build would fail). This is undesirable.
if [ -z "$(find . -name \*.go)" ]
then
	exit 0
fi

go build ./config >/dev/null 2>&1
if [ $? -ne 0 ]
then
	echo "Failed to build project. Please check the output of"
	echo "go build or run commit with --no-verify if you know"
	echo "what you are doing."
	exit 1
fi

go test ./config >/dev/null 2>&1
if [ $? -ne 0 ]
then
	echo "Failed to run tests. Please check the output of"
	echo "go test or run commit with --no-verify if you know"
	echo "what you are doing."
	exit 1
fi

go fmt ./config >/dev/null 2>&1
if [ $? -ne 0 ]
then
	echo "Failed to run go fmt. This shouldn't happen. Please"
	echo "check the output of the command to see what's wrong"
	echo "or run commit with --no-verify if you know what you"
	echo "are doing."

	exit 1
fi

go vet ./config >/dev/null 2>&1
if [ $? -ne 0 ]
then
	echo "go vet has detected potential issues in your project."
	echo "Please check its output or run commit with --no-verify"
	echo "if you know what you are doing."
	exit 1
fi
