#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

package_split=(${package//\// })
package_name=${package_split[-1]}

# add the arch for which we want to build the image
platforms=("linux/amd64" "linux/arm64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=build/_output/bin/chaos-operator$GOARCH

    # The script executes for the argument passed (in package variable)
    # here the arg will be "github.com/litmuschaos/chaos-operator/cmd/manager" for creating binary
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package

    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
