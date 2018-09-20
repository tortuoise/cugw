#!/bin/sh

# go insists on absolute path.
export GOBIN=`pwd`/dist
export DISTDIR=`pwd`/dist
export GOPATH=`pwd`
echo "GOPATH=$GOPATH"

# space seperated packages
PACKAGES=`cd src/golang.gurusys.co.uk && ls -1 |grep -v vendor`

# gitlab nonsense
[ -z "${BUILD_NUMBER}" ] && export BUILD_NUMBER=${CI_PIPELINE_ID}
[ -z "${PROJECT_NAME}" ] && export PROJECT_NAME=${CI_PROJECT_NAME}
[ -z "${COMMIT_ID}" ] && export COMMIT_ID=${CI_COMMIT_SHA}
[ -z "${GIT_BRANCH}" ] && export GIT_BRANCH=${CI_COMMIT_REF_NAME}

# use this to tar stuff
finishhook() {
	echo finished.
}

fatal_error() {
    echo $@
    exit 10
}

buildall() {
    echo building ${GOOS}/$GOARCH
    GOBIN=${DISTDIR}/${GOOS}/${GOARCH}
    mkdir -p $GOBIN
    for pkg in ${PACKAGES}; do
	MYSRC=src/golang.gurusys.co.uk/${pkg}
	( cd ${MYSRC} && make all ) || fatal_error failed to build $MYSRC
    done
}

echo
echo
echo "autobuilder.sh for ${PROJECT_NAME}, build #${BUILD_NUMBER}"
echo

if [ -d dist ]; then
    rm -rf dist
fi
mkdir dist

# we only build for amd64 atm
export GOARCH=amd64

# this allows local builds on -dev machines
# to quickly build only a single arch
# intent is for devs to set DEVOS=[localos] permanently
# on their machine and
# the autobuild.sh will do 'The Right Thing'
if [ ! -z "${DEVOS}" ]; then
    GOOS=${DEVOS}
    buildall
    finishhook
    exit 0
fi

#========= build linux
export GOOS=linux ; buildall

#========= build mac
export GOOS=darwin ; buildall

#========= build windows
# export GOOS=windows ; buildall

tar -cf dist/echoservice.tar dist/

tar -jcf dist/skel-src.tar.bz2 src/

finishhook

# we're not on a build server..
if [ -z "${BUILD_NUMBER}" ]; then
    echo no build number - not submitting
    exit 0
fi

# we are on a build server, so submit it:
build-repo-client -branch=${GIT_BRANCH} -build=${BUILD_NUMBER} -commitid=${COMMIT_ID} -commitmsg="commit msg unknown" -repository=${PROJECT_NAME} -server_addr=buildrepo:5004 

