import {app, BrowserWindow, globalShortcut, ipcMain, Menu, nativeImage, shell, Tray} from "electron";
import {spawn} from "child_process";
import remoteMain from '@electron/remote/main';
import path from "path";

// global garb to prevent gcing and losing
let mainWindow;
let tray; // must be specified globally or will be gc
// end of global garb

// bad code in need of a better store
let activeProcess = null;
let messagesToSendToFrontend = [];
// end bad code

let isQuitting = false;

const sshPath = process.env.VITE_DEV_SERVER_URL ? path.join(__dirname, "../public/ssh") : path.join(__dirname, "../../app.asar.unpacked/ssh");
const sshKeygenPath = process.env.VITE_DEV_SERVER_URL ? path.join(__dirname, "../public/ssh-keygen") : path.join(__dirname, "../../app.asar.unpacked/ssh-keygen");
const cliBinaryPath = process.env.VITE_DEV_SERVER_URL ? path.join(__dirname, "../public/cli") : path.join(__dirname, "../../app.asar.unpacked/cli");
const iconPath = process.env.VITE_DEV_SERVER_URL ? path.join(__dirname, "../build/icon.png") : path.join(__dirname, "../icon.png");
const trayIconPath = process.env.VITE_DEV_SERVER_URL ? path.join(__dirname, "../build/tray-icon.png") : path.join(__dirname, "../tray-icon.png");

app.getPath("home");
app.on("before-quit", function () {
    isQuitting = true;
});

// https://stackoverflow.com/a/69297584
// needed for remote module execution (preload.ts)
remoteMain.initialize();
app.whenReady().then(() => {
    tray = new Tray(nativeImage.createFromPath(trayIconPath));

    const splash = new BrowserWindow({
        width: 800,
        height: 800,
        icon: iconPath,
        transparent: true,
        frame: false,
        alwaysOnTop: true,
    });

    mainWindow = new BrowserWindow({
        title: "Main window",
        webPreferences: {
            preload: path.join(__dirname, 'preload.js'),
            // https://www.electronjs.org/docs/latest/tutorial/security#6-do-not-disable-websecurity
            // TODO - need to return to this and enable later when possible
            webSecurity: false,
        },
        show: false,
        width: 1024,
        height: 768,
    });
    remoteMain.enable(mainWindow.webContents);

    if (process.env.VITE_DEV_SERVER_URL) {
        splash.loadFile("splash.html");
    } else {
        // Load your file
        splash.loadFile("splash.html");
    }

    // disable reloads
    globalShortcut.register('CommandOrControl+R', () => {});
    globalShortcut.register('F5', () => {});

    if (process.env.VITE_DEV_SERVER_URL) {
        mainWindow.loadURL(process.env.VITE_DEV_SERVER_URL);
    } else {
        // Load your file
        mainWindow.loadFile("index.html");
    }

    mainWindow.webContents.on('did-finish-load', () => {
        splash.destroy();
        mainWindow.maximize();
        mainWindow.show();
    })

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

    // taken from: https://stackoverflow.com/a/32415579
    // when external links, send them to browser
    mainWindow.webContents.on('will-navigate', async function (e, url) {
        if (url != mainWindow.webContents.getURL()) {
            e.preventDefault()
            await shell.openExternal(url)
        }
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
});

ipcMain.on("request:cli_command", (_, {cliArgs}) => {
    // do not allow more than one active process at a time
    if (activeProcess) return;

    activeProcess = spawn(
        cliBinaryPath,
        cliArgs,
        { env: { SSH_PATH: sshPath, SSH_KEYGEN_PATH: sshKeygenPath } }
    );
    activeProcess.stdout.setEncoding('utf-8');
    activeProcess.stdout.on('data', (data) => messagesToSendToFrontend.push(data));
    activeProcess.stderr.setEncoding('utf-8');
    activeProcess.stderr.on('data', (data) => messagesToSendToFrontend.push(data));
    activeProcess.on('close', (data) => {
        activeProcess = null;
        mainWindow.webContents.send("received:cli_command", { status: 'success' })
    });
})

// drain the message queue to frontend
setInterval(() => {
    if (messagesToSendToFrontend.length > 0 && mainWindow.webContents) {
        mainWindow.webContents.send("received:term_messages", messagesToSendToFrontend.at(0))
        messagesToSendToFrontend.shift()
    }
}, 50);

ipcMain.on("request:cli_sigint", () => {
    if (!activeProcess) {
        mainWindow.webContents.send("received:cli_sigint")
        return
    }

    activeProcess.on("exit", () => {
        mainWindow.webContents.send("received:cli_sigint")
        activeProcess = null;
    })
    activeProcess.kill("SIGINT")
})