#!/bin/bash

#if [[ $TRAVIS_TEST_RESULT == 0 ]]; then 
  go get github.com/mitchellh/gox
  go get github.com/tcnksm/ghr
  gox -ldflags="-s -w" -output="artifacts/{{.Dir}}_{{.OS}}_{{.Arch}}"
  eval "$(ghr $TRAVIS_TAG artifacts/*)"
  ls -lah
#else
#  travis_terminate 1
#fi