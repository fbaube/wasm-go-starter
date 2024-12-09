//go:generate go run github.com/bytecodealliance/wasm-tools-go/cmd/wit-bindgen-go generate --world example --out gen ./wit

// The above compiler directive is all that is needed to
// re-generate the low-level Go code that interfaces to
// wasm. It is re-generated every time the project is built.

package main

import (
	"fmt"
        S "strings"
        "io"
	"net/http"
	"log/slog" 
	"encoding/json"

	// Equivalent:
	// https://pkg.go.dev/go.bytecodealliance.org
	// https://github.com/bytecodealliance/go-modules
	baCM "go.bytecodealliance.org/cm"
	
	// WasmCloud uses obnoxious vanity URLs.
	// "go.wasmcloud.dev/component/" is used in go.mod and at
	// pkg.go.dev, but note that the corresponding source code 
	// is found at https://github.com/wasmCloud/component-sdk-go/
	wcLog   "go.wasmcloud.dev/component/log/wasilog"
	wcHttp  "go.wasmcloud.dev/component/net/wasihttp"
	wcEnvmt "go.wasmcloud.dev/component/gen/wasi/cli/environment"
//	wcFSpre "go.wasmcloud.dev/component/gen/wasi/filesystem/preopens"
//	wcFStps "go.wasmcloud.dev/component/gen/wasi/filesystem/types"
	
	// These other wasmCloud imports are described at
	// https://wasmcloud.com/docs/tour/add-features?lang=tinygo
//	wcAtomics "github.com/wasmcloud/wasmcloud/examples/golang/components/http-hello-world/gen/wasi/keyvalue/atomics"
//	wcStore "github.com/wasmcloud/wasmcloud/examples/golang/components/http-hello-world/gen/wasi/keyvalue/store"
)

/* wasi:filesystem
type DescriptorType uint8
const (
        DescriptorTypeUnknown DescriptorType = iota
        DescriptorTypeBlockDevice
        DescriptorTypeCharacterDevice
        DescriptorTypeDirectory
        DescriptorTypeFIFO
        DescriptorTypeSymbolicLink
        DescriptorTypeRegularFile
        DescriptorTypeSocket
) */

/*
interface preopens {
  use types.{descriptor};
  get-directories: func() -> list<tuple<descriptor, string>>;
} */

var logger *slog.Logger
var hdr = "<!doctype html> \n  <html> \n    <body> \n"
var ftr = "  </body> \n</html>\n"
var addressee string
var execEnvmt ExecutionEnvironment

type ExecutionEnvironment struct {
     Argmts baCM.List[string]
     Envars baCM.List[[2]string]
     CWD string
}
     
// PartsOfRequest contains request info, sorted into fields. 
// It is sent back as a response JSON from the server.
type PartsOfRequest struct {
        Method      string `json:"method"`
        Path        string `json:"path"`
        QueryString string `json:"query_string,omitempty"`
        Body        string `json:"body,omitempty"`
}

// init handles all setup; main does zilch. 
func init() {
     	// [ContextLogger] returns a [DefaultLogger] implementation 
	// that has an additional "wasi-context" [slog.Attr] attached
	// to it. DefaultLogger is the WASI default logger implementation 
	// that adapts the wasi:logging interface to a slog.Handler.
	logger = wcLog.ContextLogger("hdlr")

	// Open the KV store
	/*
	kvStore := wcStore.Open("default")
	logger.Info("==> OPEN KV STORE <==", "type", fmt.Sprintf("%T", kvStore))
	if err := kvStore.Err(); err != nil {
	   logger.Error("Error: ", err.String())
	   return
	   }
	   */
	execEnvmt.Argmts = wcEnvmt.GetArguments()
	execEnvmt.Envars = wcEnvmt.GetEnvironment()
	logger.Info("==> ARGMTs <==", "all", 
		fmt.Sprintf("%#v", execEnvmt.Argmts))
	logger.Info("==> ENVARs <==", "all", 
		fmt.Sprintf("%#v", execEnvmt.Envars))
	/*
	logger.Info("==> FSYS <==", "all",
		fmt.Sprintf("DIRR<%d> SYML<%d> FILE<%d>",
		wcFStps.DescriptorTypeDirectory,
		wcFStps.DescriptorTypeSymbolicLink,
		wcFStps.DescriptorTypeRegularFile))
	*/	
	wcHttp.HandleFunc(handler)
}

// handler implements export `wasi-http:incoming-handler`.
// 
func handler(w http.ResponseWriter, r *http.Request) {
     	 var e error
	
	// Get the PartsOfRequest 
	var porp *PartsOfRequest
	porp, e = newPartsOfHttpRequest(r)
	if e != nil {
                http.Error(w, e.Error(), http.StatusInternalServerError)
                return
	}
	// println("GOT A REQUEST") // println does nothing! :-/
	// Some of these log fields might need to be output in quotes
	logger.Info("==> incmg HTTP REQ <==", // "host", r.Host,
		"path", r.URL.Path, // "agent", r.Header.Get("User-Agent"))
		"queryString", porp.QueryString, "body", porp.Body)
	
	// Content-Type is important!
        // w.Header().Add("Content-Type", "application/json")
        w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	addressee = "wasm-ey world"
	// FIXME 
	if S.HasPrefix(porp.QueryString, "name=") {
	   addressee = porp.QueryString[5:]
	   }
	/*
	kvStore := wcStore.Open("default")
	value := wcAtomics.Increment(*kvStore.OK(), addressee, 1)
	if err := value.Err(); err != nil {
	   logger.Error("Error: ", err.String())
	   return
	   }
	*/ 
	// Write HTML 
	wrt(w, hdr)
	wrt(w, "<h1> Hello " + addressee + "! </h1><hr/>\n")
/*	wrt(w, fmt.Sprintf("<p> For name <%s>, call number <%d> </p>",
	       addressee, value)) */
	wrt(w, "    <pre> \n")

	// Marshal the PartsOfRequest to JSON
        enc := json.NewEncoder(w)
	enc.SetIndent("", "  ") // no prefix
        e = enc.Encode(porp)
        if e != nil {
                http.Error(w, e.Error(), http.StatusInternalServerError)
                return
        }
	wrt(w, "    </pre> \n")
	wrt(w, ftr)
}

func wrt(w http.ResponseWriter, s string) {
     	_, e := w.Write([]byte(s)) // nBytes is not of interest 
	if e != nil {
                logger.Error("handler: cannot write to " +
			"http response body", "error", e)
		}
	}    

// main would have a body if we were (also?) running this as a CLI app.
func main() {}

func newPartsOfHttpRequest(r *http.Request) (*PartsOfRequest, error) {
     	p := new(PartsOfRequest)
	p.Method = r.Method
	// Split the path to retrieve the query element,
	// building the EchoRespons object 
        splitPathQuery := S.Split(r.URL.RequestURI(), "?")
        p.Path = splitPathQuery[0]
        if len(splitPathQuery) > 1 {
                p.QueryString = splitPathQuery[1]
        }
	// Consume the request body
        body, err := io.ReadAll(r.Body)
        if err != nil {
                // TODO: Log an error 
                return nil, err 
        }
        p.Body = string(body)
	return p, nil
}
