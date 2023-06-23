let bin;

const go = new Go();

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
    const terminalContainer = document.createElement('div');
    terminalContainer.setAttribute("id", "terminal");
    document.body.appendChild(terminalContainer);
    const term = new Terminal({
        cursorBlink: "block"
    });

    const fitAddon = new FitAddon.FitAddon();
    term.loadAddon(fitAddon);

    return { fitAddon, term };
}


const runCommandInTerminal = async (command, term) => {
    go.argv = ["", ...command.split(" ")]
    const lineGroups = await runAptibleCLI()
    lineGroups.forEach((lineGroup) => lineGroup.split("\n").forEach((line) => term.write(`${line}\n\r`)))
}

const startTerminal = () => {
    const { fitAddon, term } = createTerminal();

    let currLine = "";
    const entries = [];

    const newLine = () => term.write("\r\n")
    const userPrompt = () => term.write("")

    term.open(document.getElementById('terminal'));
    fitAddon.fit();

    term.write('Aptible CLI started! \n')
    newLine()
    userPrompt()
    setTimeout(() => {
        runCommandInTerminal("about", term);
        newLine();
    }, 350)

    term.onKey(async ({ key }, ev) => {
        if (key === '\r') {
            // hitting enter
            entries.push(currLine.trim());
            await runCommandInTerminal(currLine.trim(), term)
            userPrompt();
            newLine();
            currLine = "";
        } else if (key === '\x7F') {
            // hitting delete
            if (currLine) {
                currLine = currLine.slice(0, currLine.length - 1)
                term.write("\b \b")
            }
        }
        currLine += key;
        term.write(key);
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
