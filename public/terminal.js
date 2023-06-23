const go = new Go();
let bin;

const terminalElement = document.getElementById("terminal")

const runAptibleCLI = async () => {
    let oldLog = console.log;
    let stdOut = [];
    console.log = (line) => {stdOut.push(line);};
    const { instance } = await WebAssembly.instantiate(bin, go.importObject); // reset instance
    await go.run(instance);
    console.log = oldLog;
    return stdOut;
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
    const lineGroups = await runAptibleCLI()
    lineGroups.forEach((lineGroup) => lineGroup.split("\n").forEach((line) => term.write(`${line}\n\r`)))
}

const startTerminal = () => {
    const { fitAddon, term } = createTerminal();

    let currLine = "";
    const entries = [];

    const userPromptText = () => {
        const date = new Date();
        return `\x1b[1;31m${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}\x1b[37m > aptible `
    }
    const newLine = () => term.write("\r\n")
    const userPrompt = () => term.write(userPromptText())

    term.open(terminalElement);
    window.addEventListener("resize", () => fitAddon.fit());
    window.addEventListener("ready", () => fitAddon.fit());

    term.write('Aptible CLI started! \n')
    newLine();
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
    startTerminal();
}).catch((err) => {
    console.error(err);
});


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
        appContainer.classList.add("w-1/2");
        terminalElement.classList.remove("hidden")
        toggleTerminalButton.classList.remove("right-0")
        toggleTerminalButton.classList.add("half-right")
    } else {
        hideTerminal();
    }
})