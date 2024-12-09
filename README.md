# Gopher starter kit for CLI WebAssembly 

Deploy an HTTP server provided by 
a *WebAssembly component*, and then
use Go (in the form of `tinygo`) to
process its incoming HTTP requests.

These instructions use some handy free tools from
[wasmCloud](https://wasmcloud.com).
The same things can be done using only tools from
[Bytecode Alliance](https://bytecodealliance.org),
but it would be more laborious. 

## Prerequisites

- Install `git`, `tinygo`, `rust`, `cargo`,
  (*macOS*) [`brew`](https://brew.sh) 
- [Install `wasmtime`](https://docs.wasmtime.dev/cli-install.html),
  the Bytecode Alliance WebAssembly runtime 
  - *cargo:* `cargo install wasmtime-cli`
  - *macOS:*  `brew install wasmtime`
- [Install `wash`](https://wasmcloud.com/docs/installation),
  the wasmCloud WebAssembly SHell 
  - *cargo:* `cargo install wash-cli`
  - *macOS:*  `brew install wasmcloud/wasmcloud/wash`
- [Install `wadm`](https://wasmcloud.com/docs/deployment/wadm/installing),
  the wasmCloud Application Deployment Manager (from binary or via Docker)
- [Install `wasm-tools`](https://crates.io/crates/wasm-tools)
  - *cargo:* `cargo install --locked wasm-tools`
  - *macOS:*  `brew install wasm-tools`
- [Install the cargo `component`subcommand](https://crates.io/crates/cargo-component)
  - *cargo:* `cargo install cargo-component`

*Note that on macOS,* if at some point you encounter weird
`rust`errors, and you installed `rust` from Homebrew, you
may need to uninstall it and replace it with `rustup`.

## Getting it

- `git clone https://github.com/fbaube/wasm-go-starter.git`
- `cd wasm-go-starter`
- `go mod tidy # (can't hurt eh)`

## Running it

- `wash dev`

The web server at [http://127.0.0.1:8000](http://127.0.0.1:8000)
will say hi and echo back the HTTP request in JSON.

```
~/wasm-go-starter> wash dev 
ℹ️  Resolved wash session ID [orauml]
🚧 Starting a new host...
👀 Found wadm version on the disk: wadm-cli 0.18.0
✅ Using wadm version [v0.18.0]
🔧 Successfully started wasmCloud instance
✅ Successfully started host, logs writing to /Users/kilroy/.wash/dev/orauml/wasmcloud.log
🚧 Building project...
✅ Successfully built project at [/Users/kilroy/wasm-go-starter/build/http_server_s.wasm]
ℹ️  Detected component dependencies: {"http-server"}
🔁 Deployed development manifest for application [dev-orauml-http_server]
✨ HTTP Server: Access your application at http://127.0.0.1:8000
🔁 Reloading component [orauml-http_server]...
👀 Watching for file changes (press Ctrl+c to stop)...
```

You may customise the message by adding a `name`query parameter:

[`http://127.0.0.1:8000?name=kilroy`](http://127.0.0.1:8000?name=kilroy)`

## Tinkering

Modify `main.go` to see what you can make it do.
Every time you run `wash dev`, the `gen/` and `build/`
subdirectories are wiped clean and repopulated.

Modifying any other file may deeply break things.
Be prepared to roll back your modifications. 

<br/> 

*-end-*
