# vue-electron-go-wasm

## Project setup
```
npm install
sh install_wasm_exec.sh
```

### Edit go/main.go file and Make
```
make
```

### Compiles and hot-reloads for development
```
npm run serve
npm run electron:serve
```
Note: Force reloading browser window is recommended (Shift-reload).

### Compiles and minifies for production
```
npm run build
npm run electron:build
```

Note: An executable built with electron:build did not work properly
because of different practices of file loading (use of Fetch API)
in server-side and Electron.

The following code is added to `background.js` to enable Fetch API.
However, there may be security risk discussions when allowing Fetch API.

```
protocol.registerSchemesAsPrivileged([
  {
    scheme: 'app',
    privileges: {
      secure: true, standard: true, supportFetchAPI: true
    },
  },
])
```

### Run your unit tests
```
npm run test:unit
```

### Lints and fixes files
```
npm run lint
```

### What you can learn from this sample
- How to load `wasm_exec.js` and Go binary for use in index.html
  - Look at `index.html`:
    ```
      <script src="/wasm/wasm_exec.js"></script>
      <script>
        const go = new Go()
        WebAssembly.instantiateStreaming(
          fetch("/wasm/main.wasm"), go.importObject
        ).then((result) => {
          go.run(result.instance)
        })
      </script>
    ```
    then find that a bare `echo()` can be called within the same html.
    ```
      <button onClick="console.log('index.html: ' + echo('Hello, World!'))">
        index.html: say hello to console
      </button>
    ```

- How to load `wasm_exec.js` and Go binary for use in Vue.
  - Note that `wasm_exec.js` is alredy loaded in `index.html`
  - Find that the Go-wasm loader is stored in variable `global`.
  - Use `global.Go()`
    ```
      const go = new global.Go() // Loaded in public/index.html
      const WASM_URL = '/wasm/main.wasm'

      WebAssembly.instantiateStreaming(
        fetch(WASM_URL), go.importObject
      ).then(
        function (obj) {
          const wasm = obj.instance
          go.run(wasm)
        }
      )
    ```
  - Declare a global function using Vue.mixin
    ```
      Vue.mixin({
        methods: {
          run: function (program, data) {
            return global.run(program, data)
          },
        },
      })
    ```
  - Consume the functions in vue files (HTML part).
    ```
      <input v-model="message" type="textarea" />
      <button @click="result = run('echo', message)">
        Say Hello
      </button>
      <button @click="result = run('uppercase', message)">
        SAY UPPER
      </button>
      <button @click="result = run('lowercase', message)">
        say lower
      </button>
      <div>Result: {{ result }}</div>
    ```
    The functions can be bare.
    The entry point is `run()` and it takes two args:
    (1) program and (2) data.

    Currently three simple programs are supported:
    `echo`, `uppercase`, and `lowercase`.
    There is no options, no spacing, strictly.

  - Consume the functions in vue files (script part).
    ```
    methods: {
      update(program, data) {
        this.liveResult = run(program, data) // eslint-disable-line no-undef
      },
    ```
    The functions can be bare.
    Add `eslint-disable-line no-undef` to avoid a lint error,
    as it seems unable to find the declarations in vue mixin.

- How to write Go files
  - Declare a bare function. Functions can be in lowercase.
    ```
      func run(this js.Value, inputs []js.Value) interface{} {
    ```
    `this` is the context from JavaScript.
    The function args are stored in `inputs`.
  - The args must be cast using String() or an appropriate function.
    ```
      var program = inputs[0].String()
      var data = inputs[1].String()
    ```
  - Return a plain string value or an appropriate numeric value.
    ```
      return result
    ```
  - Don't forget to register the function to expose to JavaScript.
    ```
      func main(){
        js.Global().Set("run", js.FuncOf(run))
    ```
  - Use channel for keeping the process alive.
    ```
      ch := make(chan struct{}, 0)
      <-ch
    ```

- How to include an external go package
  - Import the external package

    ```
      package main

      import (
        "syscall/js"
        "local.packages/gossi"
      )
    ```

    This `gossi` is a sample project that provides a core functionality:
    `echo`, `uppercase`, and `lowercase`.

  - Put go.mod at the project home directory,
    such as `~/github/vue-electron-go-wasm/go.mod`.

    This go.mod file controls the location of the external project.
    The content will be like below, pointing to another local go package.
    ```
      module github.com/<USER>/vue-electron-go-wasm

      go 1.15

      replace local.packages/gossi => ./go_external/gossi

      require local.packages/gossi v0.0.0-00010101000000-000000000000 // indirect
    ```
    `replace` is redirecting to the current project's `./go_external/gossi` directory.

  - The package can be located outside of the main project.
    Suppose you want to maintain these as separate projects, like below.
    ```
      ~/github/
        vue-electron-go-wasm/
          go.mod
          go/
            main.go
        gossi/
          go.mod
          main.go
    ```
    In this case, the `replace` path should look like:
    ```
      replace local.packages/gossi => ../gossi
    ```
    (up one directory, which is `~/github/`, then down to `gossi`)

  - The content of the external `go.mod` should simply look like:
  ```
    module github.com/<USER>/gossi
    go 1.15
  ```

### FAQ
- Q: I get `Uncaught SyntaxError: Unexpected token '<' [wasm_exec.js:1]`
  - A: Missing `wasm_exec.js` at `./public/wasm`.
    Use the script `install_wasm_exec.sh` to install it.
