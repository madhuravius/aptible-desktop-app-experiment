# aptible

Catch-all aptible application that contains:

* CLI
* Electron app that will pull the latest UI
* Bundled WASM embed on the electron app

## Requirements

* go (tested on >=1.20.x)
* nodejs (>=18.x)
* patience

## Running it locally

To run

```sh
./refresh.sh

yarn 

# to start and run local
yarn start

# to build
yarn dist
npx electron-builder
```

## Misc

To generate images from pngs:

```sh
cd build
sips -s format icns icon.png --out icon.icns
```

To generate `wasm_exec.js` for your go version:

```sh
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./public/wasm_exec.js
```
