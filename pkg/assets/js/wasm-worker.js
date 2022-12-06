// this code is designed to be run in a web worker

importScripts("./wasm_exec.js");

// relies on wasm_exec.js
const go = new self.Go();

// polyfill 
if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

// receive messages from the main thread here
self.onmessage = async (e) => {
  const { id, type, name, arguments } = e.data;
  switch (type) {
    case "wasmInit":
      await wasmInit();
      return;
    case "wasmFunction":
      callWasmFunction(id, name, arguments);
      return;
    case "wasmVariable":
      getWasmVariable(id, name);
      return;
  }
}

const wasmInit = async () => {
  const result = await WebAssembly.instantiateStreaming(fetch("./aoc.wasm"), go.importObject);
  console.log(result);

  go.run(result.instance);
  self.postMessage({ type: "wasm-ready" });
}

const callWasmFunction = (id, name, args) => {
  const func = self[name];
  if (!func) {
    throw new Error(`function ${name} does not exist in worker scope`);
  }

  const result = func.apply(self, args);
  self.postMessage({ type: "results", functionName: name, results: result.answer, error: result.error, id: id });
}

const getWasmVariable = (id, name) => {
  self.postMessage({ type: "var", name: name, value: self[name], id: id });
}

