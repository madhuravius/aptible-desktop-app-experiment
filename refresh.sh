#!/usr/bin/env bash

# How does this script work?
#
# It ensures that the remote version of app-ui is in harmony with this electron app.
# The vite.config.mjs DOES NOT get carried over, so be warned!
#
# It will pull from a BRANCH variable you specify or main
#
# Standalone binary can be built by simply running `npx electron-builder`

BRANCH=${BRANCH:-main}
echo "Rehydrating against remote repository on \"$BRANCH\" branch."

echo "Starting to clone remote \"$BRANCH\" to tmp"
git clone -b $BRANCH git@github.com:aptible/app-ui.git tmp
rsync -a \
    --exclude '.git' \
    --exclude '.gitignore' \
    --exclude 'vite.config.mjs' \
    --exclude 'README.md' \
    ./tmp/ .

rm -rf tmp

echo "Setting up dependencies"
yarn
yarn add vite-plugin-electron
yarn add -D electron-builder electron

echo "Setting up package.json properly for use"
npm pkg set 'main'="dist-electron/main.js"
npm pkg set 'type'="commonjs"
npm pkg set 'name'="aptible"
npm pkg set 'version'="0.0.1"
