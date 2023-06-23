import { BrowserWindow, app } from "electron";

app.whenReady().then(() => {
  const splash = new BrowserWindow({
    width: 330,
    height: 80,
    transparent: true,
    frame: false,
    alwaysOnTop: true,
  });

  const mainWindow = new BrowserWindow({
    title: "Main window",
    webPreferences: {
      // https://www.electronjs.org/docs/latest/tutorial/security#6-do-not-disable-websecurity
      // TODO - need to return to this and enable later when possible
      webSecurity: false,
    },
    show: false,
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
  }, 1000);
});
