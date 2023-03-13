#!/bin/bash

source $(dirname $0)/common.sh

# check_if_exists returns 1 if an "osmosis" image exists, 0 otherwise.
check_if_exists() {
    if [[ "$(docker images -q q quicksilver > /dev/null)" != "" ]]; then
        return 1
    fi
    return 0
}

# check_if_exists returns 1 if an "quicksilver" image is built from the same commit SHA
# as the current commit, 0 otherwise.
# It assumes that the "quicksilver" image was specifically tagged with Git SHA at build
# time. Please see "build-docker" Makefile step for details.
check_if_up_to_date() {
    sha_from_image=$LIST_DOCKER_IMAGE_HASHES
    local_git_sha=$(git rev-parse HEAD)
    echo "Local Git Commit SHA: $local_git_sha"
    for cur_image_sha in $sha_from_image; do
        echo "Found Docker Tag Git SHA  : $cur_image_sha"
        if [[ "$cur_image_sha" == "$local_git_sha" ]]; then
            return 1
        fi
    done
    return 0
}

check_if_exists
exists=$?

if [[ "$exists" -eq 1 ]]; then
    echo "osmosis:debug image found"

    check_if_up_to_date
    up_to_date=$?

    if [[ "$up_to_date" -eq 1 ]]; then
        echo "quicksilver image is up to date; nothing is done"
        exit 0
    else
        echo "quicksilver image is not up to date; rebuilding"
    fi
else
    echo "quicksilver image not found; building"
fi

# Rebuild the image
make build-docker

check_if_up_to_date