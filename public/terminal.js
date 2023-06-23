const go = new Go();
let bin;
let term;

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

const runAptibleCLI = async () => {
    const obj = await WebAssembly.instantiate(bin, go.importObject); // reset instance
    await go.run(obj.instance)
    //
    return new Promise(r => {
      let timerId = setInterval(checkState, 25);
      function checkState () {
        if (go.exited == true) {
          clearInterval(timerId);
          r();
        }
      }
    });
}

const createTerminal = () => {
    const term = new Terminal({
        convertEol: true,
        cursorBlink: "block",
        fontSize: 12
    });

    const fitAddon = new FitAddon.FitAddon();
    term.loadAddon(fitAddon);

    return { fitAddon, term };
}

const runCommandInTerminal = async (command, term) => {
    const { token: { accessToken }, env: { apiUrl }} = window.reduxStore.getState();
    go.argv = ["", "--token", accessToken, "--api-host", apiUrl, ...command.split(" ")]
    return await runAptibleCLI()
}

const waitForReduxStore = async () => {
    while (true) {
        if (window.reduxStore?.getState) {
            const { token: { accessToken }} = window.reduxStore.getState();
            if (accessToken) {
                showTerminalButton();
                break;
            }
        }
        
        console.log("No store found, waiting...")
        await new Promise(r => setTimeout(r, 3000));
    }
}

const startTerminal = async () => {
    const { fitAddon, term: terminalToSet } = createTerminal();
    term = terminalToSet;

    let currLine = "";
    const entries = [];

    const userPromptText = () => {
        const date = new Date();
        const hours = ('0'+date.getHours()).slice(-2);
        const minutes = ('0'+date.getMinutes()).slice(-2)
        const seconds = ('0'+date.getSeconds()).slice(-2)
        return `\x1b[1;31m${hours}:${minutes}:${seconds}\x1b[37m > aptible `
    }
    const newLine = () => term.write("\r\n")
    const userPrompt = () => term.write(userPromptText())

    term.open(terminalElement);
    window.addEventListener("resize", () => fitAddon.fit());
    window.addEventListener("ready", () => fitAddon.fit());

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
        const { key } = char;
        if (["\u0038","\u0040"].includes(key)) {
            // ignore up/down arrows
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
        }
        else if (key === "\u0003") { // ctrl + c
            term.write('^C');
            newLine();
            userPrompt();
            currLine = "";
        } else if (key === '\r') {
            // hitting enter
            newLine();
            entries.push(currLine.trim());
            await runCommandInTerminal(currLine.trim(), term)
            newLine();
            userPrompt();
            currLine = "";
        } else if (key === '\u007F') {
            // hitting delete
            if (term._core.buffer.x > 2 && currLine) {
                term.write("\b \b")
                currLine = currLine.slice(0, currLine.length - 1)
            } else {
                return;
            }
        } else {
            currLine += key;
            term.write(key);
        }
    })
}

console.log("Loading WASM binary for use")
fetch('/cli.wasm').then(response => response.arrayBuffer()).then((binData) => {
    bin = binData;
    console.log("Loaded WASM bin, starting terminal")
}).catch((err) => {
    console.error(err);
}).finally(() => {
    startTerminal();
});

const showTerminalButton = () => {
    toggleTerminalButton.classList.remove("hidden");
}
const hideTerminalButton = () => {
    toggleTerminalButton.classList.add("hidden");
}

setInterval(async () => {
    if (window.reduxStore?.getState) {
        const { token: { accessToken }} = window.reduxStore.getState();
        if (!accessToken) hideTerminalButton();
    }
}, 1000)


const appContainer = document.getElementById("electron-app-container");
const toggleTerminalButton = document.getElementById("show-hide-terminal");
const hideTerminal = () => {
    appContainer.classList.add("w-full");
    terminalElement.classList.add("hidden")
    toggleTerminalButton.classList.remove("half-right")
    toggleTerminalButton.classList.add("right-0")
}
toggleTerminalButton.addEventListener("click", () => {
    if (terminalElement.classList.contains("hidden")) {
        appContainer.classList.remove("w-full")
        appContainer.classList.add("w-1/2");
        terminalElement.classList.remove("hidden")
        toggleTerminalButton.classList.remove("right-0")
        toggleTerminalButton.classList.add("half-right")
    } else {
        hideTerminal();
    }
})