import {app, BrowserWindow, ipcMain, Menu, nativeImage, Tray} from "electron";
import {exec, spawn} from "child_process";
import remoteMain from '@electron/remote/main';
import path from "path";
import {readFileSync} from "fs";

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
        width: 330,
        height: 80,
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

    setTimeout(() => {
        splash.destroy();
        mainWindow.show();

        if (process.env.VITE_DEV_SERVER_URL) {
            mainWindow.loadURL(process.env.VITE_DEV_SERVER_URL);
        } else {
            // Load your file
            mainWindow.loadFile("index.html");
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

// THIS WILL NEED TO BE REFACTORED TO ONLY QUERY ON DEMAND
let keyData: {
    [key: string]: {
        privateKeyFilename: string;
        publicKeyFilename: string;
        privateKeyPath: string;
        publicKeyPath: string;
        privateKeyData: string;
        publicKeyData: string;
    }
} = {};

// key mounting adapted heavily from this example app:
// https://github.com/bluprince13/ssh-key-manager/blob/master/public/electron.js
ipcMain.on("request:keys", (_) => {
    const sshdir = `${process.env.HOME}/.ssh`;
    exec("ls -a " + sshdir, (err, stdout, _) => {
        if (err) {
            console.error(`exec error when getting ssh keys: ${err}`);
            return;
        }

        const filenames = stdout.split(/\r?\n/).slice(2);
        filenames.pop();

        filenames.forEach((filename) => {
            if (filename.endsWith(".pub")) {
                const publicKeyFilename = filename;
                const privateKeyFilename = filename.split(".")[0];
                const privateKeyPath = sshdir + "/" + privateKeyFilename;
                const publicKeyPath = sshdir + "/" + publicKeyFilename;
                keyData[publicKeyFilename] = {
                    privateKeyFilename,
                    publicKeyFilename,
                    privateKeyPath,
                    publicKeyPath,
                    privateKeyData: readFileSync(privateKeyPath, 'utf-8'),
                    publicKeyData: readFileSync(publicKeyPath, 'utf-8'),
                };
            }
        })
        mainWindow.webContents.send("received:keys", true);
    });
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