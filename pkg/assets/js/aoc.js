const workerPromises = {};
let worker;

const receiveWorkerMessage = (event) => {
  const msg = event.data;
  const id = msg.id;
  switch (msg.type) {
    case "wasm-ready":
      initializePage();
      return;
    case "output":
      const shell = document.getElementById("console");
      const output = msg.output;

      const line = document.createElement("div");
      line.innerHTML = output;

      shell.appendChild(line);
      return;
    case "results":
      const { results, error } = msg;
      const fp = workerPromises[id];
      if (fp) {
        delete workerPromises[id];
        if (error) {
          fp.reject(error);
          return;
        }

        fp.resolve(results);
      }
      return;
    case "var":
      const { value } = msg;
      const vp = workerPromises[id];
      if (vp) {
        delete workerPromises[id];
        vp.resolve(value);
      }
      return;
  }
}

const getVariableFromWasm = async (varName) => {
  const id = `gv-${Date.now()}`;
  return new Promise((resolve, reject) => {
    workerPromises[id] = { resolve, reject };
    worker.postMessage({
      id: id,
      type: "wasmVariable",
      name: varName
    });
  });
}

const callWasmFunction = async (functionName, ...args) => {
  const id = `fc-${Date.now()}`;
  return new Promise((resolve, reject) => {
    workerPromises[id] = { resolve, reject };
    worker.postMessage({
      id: id,
      type: "wasmFunction",
      name: functionName,
      arguments: args,
    });
  });
}

const initializePage = () => {
  const select = document.getElementById("exercise");
  const runBtn = document.getElementById("run-exercise");

  getVariableFromWasm("aocVersion").then((result) => {
    document.getElementById("version").innerText = result;
  }).catch((err) => {
    console.error(err);
  })

  callWasmFunction("getExercises").then((exercisesJSON) => {
    const exercises = JSON.parse(exercisesJSON);
    exercises.forEach((exercise) => {
      select.add(new Option(exercise.name, exercise.value));
    });

    select.disabled = false;
    select.addEventListener("change", (event) => {
      runBtn.disabled = (event.target.value === "");
    });

    runBtn.addEventListener("click", runExercise);
  })
};

const runExercise = (event) => {
  event.target.disabled = true;

  const exercise = document.getElementById("exercise").value;
  const part = document.getElementById("part").value;
  const shell = document.getElementById("console");
  const answerBox = document.getElementById("answer");

  const typedLine = document.createElement("div");
  typedLine.style.display = "inline";

  const typedText = document.createElement("span");
  typedText.className = "bold green";
  typedText.innerHTML = `run ${exercise} part${part}<br/>`;

  typedLine.appendChild(typedText);
  shell.appendChild(typedLine);

  callWasmFunction("runExercise", exercise, part).then((answer) => {
    answerBox.value = answer;
  }).catch((err) => {
    console.error(err);
  }).finally(() => {
    const cmdPrompt = document.createElement("div");
    cmdPrompt.style.display = "inline";
    cmdPrompt.innerHTML = "$&nbsp;";
    shell.appendChild(cmdPrompt);

    event.target.disabled = false;
  })
}

window.addEventListener("DOMContentLoaded", () => {
  worker = new Worker("./js/wasm-worker.js");
  worker.addEventListener("message", receiveWorkerMessage);
  worker.addEventListener("error", (e) => {
    console.error(e);
  });

  worker.postMessage({ type: "wasmInit" });
});