#!/bin/bash

#if [[ $TRAVIS_TEST_RESULT == 0 ]]; then 
  go get github.com/mitchellh/gox
  go get github.com/tcnksm/ghr
  gox -ldflags="-s -w" -output="artifacts/{{.Dir}}_{{.OS}}_{{.Arch}}"
  echo $TRAVIS_TAG
  goReleaseCmd="ghr \""$TRAVIS_TAG"\" \"artifacts/\""
  echo $goReleaseCmd
  eval $goReleaseCmd
  #ls -lah
#else
#  travis_terminate 1
#fi