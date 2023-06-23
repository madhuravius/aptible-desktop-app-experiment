import { BrowserWindow, Menu, Tray, app } from "electron";
import path from "path";

let isQuitting = false;
// const trayIconPath = process.env.VITE_DEV_SERVER_URL ? path.join(__dirname, "../", "favicon.ico") : "dist/favicon.ico"
const iconPath = process.env.VITE_DEV_SERVER_URL
  ? path.join(__dirname, "../", "logo.png")
  : "dist/logo.png";

app.on("before-quit", function () {
  isQuitting = true;
});

app.whenReady().then(() => {
  const tray = new Tray(iconPath);
  const splash = new BrowserWindow({
    width: 330,
    height: 80,
    icon: iconPath,
    transparent: true,
    frame: false,
    alwaysOnTop: true,
  });

  const mainWindow = new BrowserWindow({
    title: "Main window",
    webPreferences: {
      nodeIntegration: true,
      // https://www.electronjs.org/docs/latest/tutorial/security#6-do-not-disable-websecurity
      // TODO - need to return to this and enable later when possible
      webSecurity: false,
    },
    show: false,
    width: 1024,
    height: 768,
  });

  if (process.env.VITE_DEV_SERVER_URL) {
    splash.loadFile("splash.html");
  } else {
    // Load your file
    splash.loadFile("dist/splash.html");
  }

  setTimeout(() => {
    splash.destroy();
    mainWindow.show();

    // You can use `process.env.VITE_DEV_SERVER_URL` when the vite command is called `serve`
    if (process.env.VITE_DEV_SERVER_URL) {
      mainWindow.loadURL(process.env.VITE_DEV_SERVER_URL);
    } else {
      // Load your file
      mainWindow.loadFile("dist/index.html");
    }

    // BEGIN TRAY-RELATED
    // add desktop app-specific code (ex: terminal)
    mainWindow.on("minimize", function (event) {
      event.preventDefault();
      mainWindow.hide();
    });

    mainWindow.on("close", function (event) {
      if (!isQuitting) {
        event.preventDefault();
        mainWindow.hide();
      }
      return false;
    });

    tray.setContextMenu(
      Menu.buildFromTemplate([
        {
          label: "Show Aptible",
          click: function () {
            mainWindow.show();
          },
        },
        {
          label: "Quit",
          click: function () {
            isQuitting = true;
            app.quit();
          },
        },
      ]),
    );
    // END TRAY-RELATED
  }, 1000);
});
