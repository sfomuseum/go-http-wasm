package server

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/aaronland/go-http-bootstrap"
	aa_server "github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-http-wasm/v2"
)

func ServeWithFS(ctx context.Context, server_uri string, server_fs fs.FS) error {

	s, err := aa_server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %v", err)
	}

	mux := http.NewServeMux()

	bootstrap_opts := bootstrap.DefaultBootstrapOptions()

	err = bootstrap.AppendAssetHandlers(mux, bootstrap_opts)

	if err != nil {
		return fmt.Errorf("Failed to append Bootstrap asset handlers, %v", err)
	}

	wasm_opts := wasm.DefaultWASMOptions()

	err = wasm.AppendAssetHandlers(mux, wasm_opts)

	if err != nil {
		return fmt.Errorf("Failed to append wasm assets handler, %v", err)
	}

	http_fs := http.FS(server_fs)
	fs_handler := http.FileServer(http_fs)

	fs_handler = bootstrap.AppendResourcesHandler(fs_handler, bootstrap_opts)

	fs_handler = wasm.AppendResourcesHandler(fs_handler, wasm_opts)

	mux.Handle("/", fs_handler)

	log.Printf("Listening on %s", s.Address())
	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to start server, %v", err)
	}

	return nil
}
