#!/bin/bash
gittag= $(git describe --tags  --abbrev=0)
now=$(date)
gitcommit=$(git rev-parse --verify HEAD)
user=$USER

echo $(go version)
GOOS=darwin 
GOARCH=arm64
go build -v -ldflags "-X 'main.buildGitTag=$gittag' -X 'main.buildGitCommit=$gitcommit' -X 'main.buildTime=$now' -X 'main.buildUser=$user'" -o ./bin/studyrun

echo "buildVersion $gittag"
echo "buildGitCommit $gitcommit"
echo "buildTime $now"
echo "buildUser $user"