// this code is designed to be run in a web worker

// relies on wasm_exec.js
const go = new Go();

// polyfill 
if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

// receive messages from the main thread here
self.onmessage = async (e) => {
  const { type, message } = e.data;
  switch (type) {
    case "start":
      await wasmInit();
      return;
    case "call":
      callWasmFunction(message.name, message.arguments);
      return;
    case "globalVar":
      getWasmVariable(message.name);
      return;
  }
}

const wasmInit = async () => {
  const result = await WebAssembly.instantiateStreaming(fetch("js/aoc.wasm"), go.importObject);
  console.log(result);

  go.run(result.instance);
  self.postMessage({ type: "wasm-ready" });
}

const callWasmFunction = (name, arguments) => {
  const func = self[name];
  if (!func) {
    throw new Error(`function ${name} does not exist in worker scope`);
  }

  const result = func(...arguments);
  self.postMessage({ type: "results", functionName: name, results: result });
}

