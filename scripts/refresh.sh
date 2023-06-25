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
    --exclude 'index.html' \
    ./tmp/ .

rm -rf tmp

echo "Setting up dependencies"
yarn
yarn add vite-plugin-electron crypto text-encoding-polyfill
yarn add -D electron-builder electron "@electron/remote"@latest

# electron does not support browser router, so we'll need to switch to memory router for that
# sed on mac does not allow overwriting in place easily, so just copy it to tmp and overwrite
sed -i '' -e 's/createBrowserRouter/createMemoryRouter/g' ./src/app/router.tsx

# relative image paths to root are not allowed in electron either, so this has to be fixed
for pathToEdit in database-types language-types resource-types; do
    find ./src -type f -exec sed -i '' -e "s|/$pathToEdit|$pathToEdit|g" {} +
done

echo "Setting up package.json properly for use"
npm pkg set 'main'="dist-electron/main.js"
npm pkg set 'type'="commonjs"
npm pkg set 'name'="aptible"
npm pkg set 'version'="0.0.1"
