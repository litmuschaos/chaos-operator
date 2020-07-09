#!/usr/bin/env bash
#
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")" echo SOURCE; done
DIR="$( cd -P "$( dirname "$SOURCE" )/../" && pwd )"

# Change into that directory
cd "$DIR"

# Get the version details
VERSION="$(cat $GOPATH/src/github.com/litmuschaos/admission-controllers/VERSION)"
VERSION_META="$(cat $GOPATH/src/github.com/litmuschaos/admission-controllers/BUILDMETA)"

# Determine the arch/os combos we're building for
UNAME=$(uname)
ARCH=$(uname -m)
if [ "$UNAME" != "Linux" -a "$UNAME" != "Darwin" ] ; then
    echo "Sorry, this OS is not supported yet."
    exit 1
fi

if [ -z "${PNAME}" ];
then
    echo "Project name not defined"
    exit 1
fi

if [ -z "${CTLNAME}" ];
then
    echo "CTLNAME not defined"
    exit 1
fi

# Delete the old dir
echo "==> Removing old directory..."
rm -rf bin/${PNAME}/*
mkdir -p bin/${PNAME}/

# Build!
echo "==> Building ${CTLNAME} using $(go version)... "

output_name="bin/${PNAME}/"$CTLNAME

env GOOS=$GOOS GOARCH=$GOARCH go build ${BUILD_TAG} -ldflags \
    "-X main.CtlName='${CTLNAME}' \
    -X github.com/litmuschaos/admission-controllers/pkg/version.Version=${VERSION} \
    -X github.com/litmuschaos/admission-controllers/pkg/version.VersionMeta=${VERSION_META}"\
    -o $output_name\
    ./cmd/${CTLNAME}

echo ""
# Done!
echo "==> Results:"
ls -hl bin/${PNAME}/
