#!/usr/bin/env bash
docker login -u $BINTRAY_USERNAME -p $BINTRAY_APIKEY -e $BINTRAY_EMAIL aacebedo-docker-cfdsnupdater.bintray.io 
docker build -f  environments/run/Dockerfile.amd64 -t aacebedo/cfdnsupdater-amd64 .
#docker tag $ID aacebedo/cfdnsupdater-amd64:$TRAVIS_BRANCH
#docker push aacebedo/cfdnsupdater-amd64:$TRAVIS_BRANCH
