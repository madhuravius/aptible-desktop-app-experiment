let initialFetchedKeys = false;
let keys; // should be dropped from global state
let term;
let termActivity = false;
let fitAddon;

const terminalElement = document.getElementById("terminal")

const waitFor = async ({
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
    ipcRenderer.send("request:keys");
    ipcRenderer.receive("received:keys", () => {
        initialFetchedKeys = true;
    });
};

// the only purpose of this is to drain the queue with a subscriber, it does nothing else but log
// to term
(() => {
    ipcRenderer.receive("received:term_messages", (event) => {
        if (term) {
            event.toString().split("\n").forEach((line) => term.write(`${line}\n\r`))
        }

        reconcileLoadingIndicatorWithScroll();
    })
})()

const createTerminal = () => {
    const term = new Terminal({
        convertEol: true,
        cursorBlink: "block",
        fontSize: 12,
    });

    fitAddon = new FitAddon.FitAddon();
    term.loadAddon(fitAddon);

    return {fitAddon, term};
}


const userPromptText = () => {
    const date = new Date();
    const hours = ('0' + date.getHours()).slice(-2);
    const minutes = ('0' + date.getMinutes()).slice(-2)
    const seconds = ('0' + date.getSeconds()).slice(-2)
    return `\x1b[1;31m${hours}:${minutes}:${seconds}\x1b[37m > aptible `
}
const newLine = () => term.write("\r\n")
const userPrompt = () => term.write(userPromptText())

const runCommandInTerminal = async (command) => {
    const splitCommand = command.split(" ");
    const {token: {accessToken}, env: {apiUrl}} = window.reduxStore.getState();
    const cliArgs = ["--token", accessToken, "--api-host", apiUrl, ...splitCommand];

    termActivity = true;
    showLoadingIndicator();
    // This must be done on the nodeJs side, this cannot execute fully client side
    ipcRenderer.send("request:cli_command", ({cliArgs}));
    let completedRemoteTask = false;
    ipcRenderer.receive("received:cli_command", ({ status }) => {
        completedRemoteTask = true;
    })
    await waitFor({
        callback: () => {
            hideLoadingIndicator();
            setTimeout(() => {
                newLine();
                userPrompt();
            }, 250); // if this is done too fast, things end badly
        },
        func: () => completedRemoteTask
    })
    termActivity = false;
}

const sendInterruptToTerminal = async () => {
    termActivity = true;
    showLoadingIndicator();
    // This must be done on the nodeJs side, this cannot execute fully client side
    ipcRenderer.send("request:cli_sigint");
    let completedRemoteTask = false;
    ipcRenderer.receive("received:cli_sigint", () => {
        completedRemoteTask = true;
    })
    await waitFor({
        callback: () => {
            hideLoadingIndicator();
        },
        func: () => completedRemoteTask
    })
    termActivity = false;
}

const waitForReduxStore = async () => {
    await waitFor(
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

    term.open(terminalElement);
    window.addEventListener("resize", () => fitAddon.fit());

    term.write('Aptible CLI started! \n')
    newLine();
    await waitForReduxStore();

    setTimeout(async () => {
        await runCommandInTerminal("about");
    }, 350)

    term.clear();
    // todo - https://github.com/EDDYMENS/interactive-terminal/blob/main/frontend.js#L21
    // main loop
    term.onKey(async (char, ev) => {
        reconcileLoadingIndicatorWithScroll();

        const {key} = char;
        // overrides: interrupts and anything that ignores something has been sent to the terminal
         if (key === "\u0003") { // ctrl + c
            await sendInterruptToTerminal();
            term.write('^C');
            newLine();
            userPrompt();
            currLine = "";
        }

        // do not allow other actions while an activity is ongoing
        if (termActivity) return
        // ignore left/right arrows for now
        if (["\x1B[D", "\x1B[C"].includes(key)) return

        // normal flows
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
        } else if (key === '\r') { // hitting enter
            newLine();
            entries.push(currLine.trim());
            await runCommandInTerminal(currLine.trim())
            currLine = "";
            lastPositionInHistory = entries.length - 1;
        } else if (key === '\u007F') { // hitting delete
            if (term._core.buffer.x > 3 && currLine) {
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

const shiftLoadingIndicatorForScroll = () => {
    loadingIndicator.classList.add("mr-2");
}

const unshiftLoadingIndicatorForScroll = () => {
    loadingIndicator.classList.remove("mr-2");
}

const reconcileLoadingIndicatorWithScroll = () => {
    // if scroll height exceeds a certain value, adjust margin of loader
    if (!terminalElement.classList.contains("hidden")) {
        if (term._core.viewport._activeBuffer.lines.length > term._core.viewport._activeBuffer._rows) {
            shiftLoadingIndicatorForScroll()
        } else {
            unshiftLoadingIndicatorForScroll()
        }
    }
}

const hideTerminal = () => {
    appContainer.classList.add("w-full");
    terminalElement.classList.add("hidden")
    unshiftLoadingIndicatorForScroll()
    loadingIndicator.classList.remove("overlaid");
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
    reconcileLoadingIndicatorWithScroll()
    loadingIndicator.classList.add("overlaid");
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
getKeys();
showLoadingIndicator();
(async () => {
    await waitFor({
        func: () => initialFetchedKeys === true,
    });
    startTerminal();
    hideLoadingIndicator();
})();