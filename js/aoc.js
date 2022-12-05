const go = new Go();
let wasm;

const setPage = () => {
  const select = document.getElementById("exercise");
  const runBtn = document.getElementById("run-exercise");

  document.getElementById("version").innerText = aocVersion;

  const infos = JSON.parse(getExercises());
  console.log(infos);
  infos.forEach((info) => {
    select.add(new Option(info.name, info.value));
  });
  select.disabled = false;

  select.addEventListener("change", (event) => {
    runBtn.disabled = (event.target.value === "");
  });

  runBtn.addEventListener("click", doRun);
};

const doRun = () => {
  const exercise = document.getElementById("exercise").value;
  const part = document.getElementById("part").value;
  const shell = document.getElementById("console");
  const answerBox = document.getElementById("answer");
  const runBtn = document.getElementById("run-exercise");

  runBtn.disabled = true;

  const typedLine = document.createElement("div");
  typedLine.className = "line";
  typedLine.style.display = "inline";

  const typedText = document.createElement("span");
  typedText.className = "bold green";
  typedText.innerHTML = `run ${exercise} part${part}<br/>`;
  
  typedLine.appendChild(typedText);
  shell.appendChild(typedLine);

  const result = runExercise(exercise, part, shell);

  if (result.error) {
    console.error(result.error);
  }

  answerBox.value = result.answer.toString();
  const cmdPrompt = document.createElement("div");
  cmdPrompt.style.display = "inline";
  cmdPrompt.innerHTML = "$&nbsp;";
  shell.appendChild(cmdPrompt);

  runBtn.disabled = false;
}

if (!WebAssembly.instantiateStreaming) {
  // polyfill
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

WebAssembly.instantiateStreaming(fetch("js/aoc.wasm"), go.importObject).then(result => {
  wasm = result.instance;
  go.run(wasm);
  setPage();
});