const go = new Go();
let bin;
let initialFetchedKeys = false;
let keys; // should be dropped from global state
let term;
let fitAddon;

const terminalElement = document.getElementById("terminal")

// hijack console.log for terminal.js specific things
const consoleLog = console.log;
console.log = (params) => {
    let fromSelf = new Error().stack.split("\n")?.[1].includes("terminal.js")
    if (fromSelf && term) {
        params.split("\n").forEach((line) => term.write(`${line}\n\r`))
    } else {
        consoleLog(params);
    }
}

const waitForTruth = async ({
                                func, interval = 25, callback = () => {
    }
                            }) => {
    return new Promise(r => {
        let timerId = setInterval(checkState, interval);

        function checkState() {
            if (func()) {
                clearInterval(timerId);
                callback()
                r();
            }
        }
    });
}

const getKeys = () => {
    // TODO - this needs to be rewritten with a locker and some checkin awaitable checkout in the CLI flow itself
    ipcRenderer.send("request:keys");
    ipcRenderer.receive("received:keys", () => {
        initialFetchedKeys = true;
    });
};

const runAptibleCLI = async () => {
    showLoadingIndicator();
    const obj = await WebAssembly.instantiate(bin, go.importObject); // reset instance
    await go.run(obj.instance)
    // this must be done because go spawns work that immediately returns a promise in wasm
    // if any async background workers are present (web), we must wait until they are ALL
    // complete. this is not documented very well and is probably subject to change
    await waitForTruth({
        callback: hideLoadingIndicator,
        func: () => go.exited === true,
    })
}

const createTerminal = () => {
    const term = new Terminal({
        convertEol: true,
        cursorBlink: "block",
        fontSize: 12
    });

    fitAddon = new FitAddon.FitAddon();
    term.loadAddon(fitAddon);

    return {fitAddon, term};
}

const runCommandInTerminal = async (command, term) => {
    const splitCommand = command.split(" ");
    const {token: {accessToken}, env: {apiUrl}} = window.reduxStore.getState();
    const cliArgs = ["", "--token", accessToken, "--api-host", apiUrl, ...splitCommand];
    const commandRequiresSshKeys = ["logs", "operation:follow", "ssh"].map((possibleCommand) => {
        return splitCommand.includes(possibleCommand)
    }).some((found) => found);

    if (commandRequiresSshKeys) {
        // This must be done on the nodeJs side, this cannot execute fully client side
        const possibleKeysInHomeDirectory = Object.entries(keys)?.[0];
        if (possibleKeysInHomeDirectory) {
            const [_, {publicKeyData, privateKeyData}] = possibleKeysInHomeDirectory;
            ["--public-key", publicKeyData, "--private-key", privateKeyData]
                .forEach((flagValue) => {
                    cliArgs.splice(cliArgs.length - 2, 0, flagValue)
                });
        }
    } else {
        go.argv = cliArgs;
        return await runAptibleCLI()
    }
}

const waitForReduxStore = async () => {
    await waitForTruth(
        {
            func: () => !!window.reduxStore?.getState,
            callback: () => {
                const {token: {accessToken}} = window.reduxStore.getState();
                if (accessToken) showTerminalButton();
            }
        }
    )
}

const startTerminal = async () => {
    const {fitAddon: fitAddonToSet, term: terminalToSet} = createTerminal();
    term = terminalToSet;
    fitAddon = fitAddonToSet;

    let currLine = "";
    let lastPositionInHistory = 0;
    const entries = [];

    const userPromptText = () => {
        const date = new Date();
        const hours = ('0' + date.getHours()).slice(-2);
        const minutes = ('0' + date.getMinutes()).slice(-2)
        const seconds = ('0' + date.getSeconds()).slice(-2)
        return `\x1b[1;31m${hours}:${minutes}:${seconds}\x1b[37m > aptible `
    }
    const newLine = () => term.write("\r\n")
    const userPrompt = () => term.write(userPromptText())

    term.open(terminalElement);
    window.addEventListener("resize", () => fitAddon.fit());

    term.write('Aptible CLI started! \n')
    newLine();
    await waitForReduxStore();

    setTimeout(async () => {
        await runCommandInTerminal("about", term);
        newLine();
        userPrompt();
    }, 350)

    term.clear();

    // todo - https://github.com/EDDYMENS/interactive-terminal/blob/main/frontend.js#L21
    // main loop
    term.onKey(async (char, ev) => {
        const {key} = char;
        if (["\x1B[D", "\x1B[C"].includes(key)) { // ignore left/right arrows for now
            return
        }

        if (key === "\u0004") {
            term.write('^D');
            newLine();
            term.write('User requested to close Terminal, hiding!');
            newLine();
            userPrompt();
            currLine = "";
            hideTerminal();
        } else if (key === "\x1B[A") { // up arrow
            currLine = "";
            if (entries.length > 0) {
                if (lastPositionInHistory === entries.length) {
                    currLine = entries.at(-1);
                    lastPositionInHistory--;
                } else if (lastPositionInHistory >= 0) {
                    currLine = entries.at(lastPositionInHistory)
                    lastPositionInHistory--;
                }
                term.write('\x1b[2K\r'); // clear CURRENT line
                userPrompt();
                term.write(currLine);
            }
        } else if (key === "\x1B[B") { // down arrow
            if (lastPositionInHistory < entries.length - 1) {
                currLine = entries.at(lastPositionInHistory)
                lastPositionInHistory++;
                term.write('\x1b[2K\r'); // clear CURRENT line
                userPrompt();
                term.write(currLine);
            }
        } else if (key === "\f") { // ctrl + l / clear
            term.clear();
            newLine();
            userPrompt();
            currLine = "";
        } else if (key === "\u0003") { // ctrl + c
            term.write('^C');
            newLine();
            userPrompt();
            currLine = "";
        } else if (key === '\r') { // hitting enter
            newLine();
            entries.push(currLine.trim());
            await runCommandInTerminal(currLine.trim(), term)
            userPrompt();
            currLine = "";
            lastPositionInHistory = entries.length - 1;
        } else if (key === '\u007F') { // hitting delete
            if (term._core.buffer.x > 2 && currLine) {
                term.write("\b \b")
                currLine = currLine.slice(0, currLine.length - 1)
            } else {

            }
        } else {
            currLine += key;
            term.write(key);
        }
    })
}

setInterval(async () => {
    if (window.reduxStore?.getState) {
        const {token: {accessToken}} = window.reduxStore.getState();
        if (!accessToken) hideTerminalButton();
    }
}, 1000)


const appContainer = document.getElementById("electron-app-container");
const toggleTerminalButton = document.getElementById("show-hide-terminal");
const loadingIndicator = document.getElementById("loading-indicator");

const showTerminalButton = () => {
    toggleTerminalButton.classList.remove("hidden");
}
const hideTerminalButton = () => {
    toggleTerminalButton.classList.add("hidden");
}

const showLoadingIndicator = () => {
    loadingIndicator.classList.remove("hidden");
}

const hideLoadingIndicator = () => {
    loadingIndicator.classList.add("hidden");
}

const hideTerminal = () => {
    appContainer.classList.add("w-full");
    terminalElement.classList.add("hidden")
    toggleTerminalButton.classList.remove("half-right")
    toggleTerminalButton.classList.add("right-0")
    toggleTerminalButton.innerHTML = `        <div class="flex">
  <span class="leading-4">
    ‹ Terminal <br />
    <span class="text-xs">Ctrl + Shift + T</span>
  </span>
  <img class="inline-block ml-2 h-8" src="resource-types/logo-service.png" />
</div>`;
}
const showTerminal = () => {
    appContainer.classList.remove("w-full")
    appContainer.classList.add("w-1/2");
    terminalElement.classList.remove("hidden")
    toggleTerminalButton.classList.remove("right-0")
    toggleTerminalButton.classList.add("half-right")
    toggleTerminalButton.innerHTML = `<div class="flex">
      <span class="leading-4">
        › Terminal <br />
        <span class="text-xs">Ctrl + Shift + T</span>
      </span>
      <img class="inline-block ml-2 h-8" src="resource-types/logo-service.png" />
    </div>`;
    setTimeout(() => {
        fitAddon.fit();
        term.focus();
    }, 100);
}

toggleTerminalButton.addEventListener("click", () => {
    if (terminalElement.classList.contains("hidden")) {
        showTerminal();
    } else {
        hideTerminal();
    }
})

document.addEventListener('keydown', (event) => {
    if (event.ctrlKey && event.shiftKey && event.key === "T") {
        if (terminalElement.classList.contains("hidden")) {
            showTerminal();
        } else {
            hideTerminal();
        }
    }
});

// MAIN LOOP
console.log("Loading WASM binary for use")
getKeys();
showLoadingIndicator();
fetch('cli.wasm')
    .then(response => response.arrayBuffer()).then((binData) => {
    bin = binData;
    console.log("Loaded WASM bin, starting terminal")
}).catch((err) => {
    console.error(err);
}).finally(async () => {
    await waitForTruth({
        func: () => initialFetchedKeys === true,
    });
    startTerminal();
    hideLoadingIndicator();
});