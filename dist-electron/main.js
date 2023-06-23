"use strict";
const electron = require("electron");
electron.app.whenReady().then(() => {
  const splash = new electron.BrowserWindow({
    width: 330,
    height: 80,
    transparent: true,
    frame: false,
    alwaysOnTop: true
  });
  const mainWindow = new electron.BrowserWindow({
    title: "Main window",
    webPreferences: {
      // https://www.electronjs.org/docs/latest/tutorial/security#6-do-not-disable-websecurity
      // TODO - need to return to this and enable later when possible
      webSecurity: false
    },
    show: false
  });
  if (process.env.VITE_DEV_SERVER_URL) {
    splash.loadFile("splash.html");
  } else {
    splash.loadFile("dist/splash.html");
  }
  setTimeout(() => {
    splash.destroy();
    mainWindow.show();
    if (process.env.VITE_DEV_SERVER_URL) {
      mainWindow.loadURL(process.env.VITE_DEV_SERVER_URL);
    } else {
      mainWindow.loadFile("dist/index.html");
    }
  }, 1e3);
});
