# aptible

Catch-all aptible application that contains:

* CLI
* Electron app that will pull the latest UI

## Requirements

* go (tested on >=1.20.x)
* nodejs (>=18.x)
* patience

## Running it locally

To run with shell scripts (also see Makefile)

```sh
./scripts/refresh.sh
./scripts/install-cli.sh

# local and test in dev server
yarn start

# dist and look in releases directory for dist
rm -rf release && yarn build && npx electron-builder
```

## Misc

To generate images from pngs:

```sh
cd build
sips -s format icns icon.png --out icon.icns
```
