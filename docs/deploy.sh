#!/bin/bash

cp ../install.sh ./static/install.sh

GIT_USER=lets-cli CURRENT_BRANCH=master USE_SSH=true npm run deploy

rm ./static/install.sh