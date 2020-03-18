#!/usr/bin/env bash

CMD=$1
NEW_CHECKSUM=$LETS_CHECKSUM
LABEL=$LETS_CMD
LETS_DIR=.lets

# ensure lets dir exists
if [[ ! -d "$LETS_DIR" ]]; then
    mkdir $LETS_DIR
fi;
# ensure checksum file exists
if [[ ! -f "$LETS_DIR/$LABEL" ]]; then
    touch $LETS_DIR/$LABEL
fi;

# read
PREV_CHECKSUM=$(head -n 1 $LETS_DIR/$LABEL)

if [[ "$NEW_CHECKSUM" != "$PREV_CHECKSUM" ]]; then
    # update checksum
    echo $NEW_CHECKSUM > $LETS_DIR/$LABEL
    eval $CMD
fi;
