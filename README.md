# go-http-wasm

Go HTTP middleware package for bundling, serving and appending pointers to `wasm_exec.js` and related JavaScript assets.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-http-wasm.svg)](https://pkg.go.dev/github.com/sfomuseum/go-http-wasm)

## Deprecation Notice

This package has been deprecated and is no longer maintained. You should just use/load [sfomuseum/js-sfomuseum-golang-wasm](https://github.com/sfomuseum/js-sfomuseum-golang-wasm) in your code as an embeddable file instead.

## Motivation

This is a simple Go HTTP middleware package for bundling, serving and appending pointers to `wasm_exec.js` which is required by web applications using WASM binaries derived from Go applications.

Given that all this package does is manage two files (the `wasm_exec.js` file that is bundled with the Go programming language and the `sfomuseum.wasm.js` helper library) it borders on the absurd. It would arguably be easier for those packages that consume this one to simple manage bundling and serving that same file themselves. That is a perfectly good way to do the same thing.

## Example

_Note that error handling omitted for the sake of brevity._
 
First, import the `net/http` and `go-http-wasm` packages.

```
import (
       "net/http"

	"github.com/sfomuseum/go-http-wasm/v2"
)
```

Next, create standard `net/http` mux and handler instances. The details of these instances are left to individual applications.

```
example_mux := http.NewServeMux()
example_handler := SomeHTTPHandler()

wasm_opts := wasm.DefaultWASMOptions()
```

Next, append relevant handlers for serving WASM-related assets (in this case `wasm_exec.js`) to the mux instance.

```
wasm.AppendAssetHandlers(mux, wasm_opts)
```

Finally, wrap the handler instance with `wasm.AppendResourcesHandler`. This will append the relevant JavaScript directives to HTML output reference the WASM-related assets that were added using the `wasm.AppendAssetHandlers` method.

```
example_handler = wasm.AppendResourcesHandler(example_handler, wasm_opts)
	
example_mux.Handle("/", example_handler)
```

Notably this will rewrite HTML output to append the following JavaScript directives:

```
<script type="text/javascript" src="/javascript/wasm_exec.js"></script>
<script type="text/javascript" src="/javascript/sfomuseum.wasm.js"></script></head>
```

* `wasm_exec.js` is the JavaScript library for interacting with WASM binaries bundled with the Go programming language.
* `sfomuseum.wasm.js` is a convenience library for fetching and invoking WASM. It exposes a single `fetch` method that take the URL of a WASM binary and returns a JavaScript Promise. For example:

```
sfomuseum.wasm.fetch("wasm/whosonfirst_placetypes.wasm").then(rsp => {
	...
}).catch(err => {
	...
});
``

Here's a concrete example taken from the [go-whosonfirst-placetypes-wasm](https://github.com/whosonfirst/go-whosonfirst-placetypes-wasm/tree/main/cmd/example) package.
 
```
package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"

	placetypes_wasm "github.com/whosonfirst/go-whosonfirst-placetypes-wasm/http"
	"github.com/sfomuseum/go-http-wasm/v2"
)

//go:embed index.html example.*
var FS embed.FS

func main() {

	host := flag.String("host", "localhost", "The host name to listen for requests on")
	port := flag.Int("port", 8080, "The host port to listen for requests on")

	flag.Parse()

	mux := http.NewServeMux()

	http_fs := http.FS(FS)
	example_handler := http.FileServer(http_fs)

	wasm_opts := wasm.DefaultWASMOptions()
	
	wasm.AppendAssetHandlers(mux, wasm_opts)
	placetypes_wasm.AppendAssetHandlers(mux)

	example_handler = wasm.AppendResourcesHandler(example_handler, wasm_opts)
	
	mux.Handle("/", example_handler)

	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Printf("Listening for requests on %s\n", addr)

	err = http.ListenAndServe(addr, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
```

## See also

* https://github.com/aaronland/go-http-static
* https://github.com/golang/go/wiki/WebAssembly
