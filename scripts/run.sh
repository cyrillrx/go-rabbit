#!/usr/bin/env bash

BASH_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT_DIR=$BASH_DIR/..

# Build
$BASH_DIR/build.sh

# Environment variables to define
# AMQP_URL    : URL to the AMQP server.
# SEGMENT_KEY : The Segment.io key

if [ -f $BASH_DIR/env.sh ]; then
if [ ! -x $BASH_DIR/env.sh ]; then
chmod u+x $BASH_DIR/env.sh
fi
echo "Set environment variables..."
. $BASH_DIR/env.sh
fi

# Run
echo "Running..."
$GOPATH/bin/go-rabbit -amqp-url="$AMQP_URL"