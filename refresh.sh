#!/usr/bin/env bash

# How does this script work?
#
# It ensures that the remote version of app-ui is in harmony with this electron app.
# The vite.config.mjs DOES NOT get carried over, so be warned!
#
# It will pull from a BRANCH variable you specify or main

BRANCH=${BRANCH:-main}
echo "Rehydrating against remote repository on \"$BRANCH\" branch."

echo "Starting to clone remote \"$BRANCH\" to tmp"
git clone -b $BRANCH git@github.com:aptible/app-ui.git tmp
rsync -r \
    --exclude '.git' \
    --exclude '.gitignore' \
    --exclude 'vite.config.mjs' \
    ./tmp .

rm -rf tmp
yarn add electron vite-plugin-electron
yarn
