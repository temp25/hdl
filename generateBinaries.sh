#!/bin/bash

# Compiles binary for available os/arch in go version

#Get Gox for Go cross compilation 
go get github.com/mitchellh/gox 	

for osArch in $(gox -osarch-list)
do
 re="[a-z0-9]+\/[a-z0-9]+"
 if [[ $osArch =~ $re ]]; then 
   echo Building binary for $osArch;
   gox -ldflags="-s -w" -osarch="$osArch" -output="artifacts/{{.Dir}}_{{.OS}}_{{.Arch}}"
   echo Binary stored in $(pwd)/artifacts folder
 fi
done

