// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version v0.3.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/discord-gophers/goapi-gen/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// PostJSONBody defines parameters for Post.
type PostJSONBody struct {
	Email string `json:"email"`

	// Full name of the user
	Name *string `json:"name,omitempty"`

	// List of role IDs to assign to the user
	RoleIds []string `json:"role_ids,omitempty"`

	// Unique username for the user
	Username string `json:"username"`
}

// PutIDJSONBody defines parameters for PutID.
type PutIDJSONBody struct {
	// Full name of the user
	Name *string `json:"name,omitempty"`

	// List of role IDs to assign to the user
	RoleIds []string `json:"role_ids,omitempty"`

	// Unique username for the user
	Username *string `json:"username,omitempty"`
}

// PostJSONRequestBody defines body for Post for application/json ContentType.
type PostJSONRequestBody PostJSONBody

// Bind implements render.Binder.
func (PostJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PutIDJSONRequestBody defines body for PutID for application/json ContentType.
type PutIDJSONRequestBody PutIDJSONBody

// Bind implements render.Binder.
func (PutIDJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Response is a common response struct for all the API calls.
// A Response object may be instantiated via functions for specific operation responses.
// It may also be instantiated directly, for the purpose of responding with a single status code.
type Response struct {
	body        interface{}
	Code        int
	contentType string
}

// Render implements the render.Renderer interface. It sets the Content-Type header
// and status code based on the response definition.
func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", resp.contentType)
	render.Status(r, resp.Code)
	return nil
}

// Status is a builder method to override the default status code for a response.
func (resp *Response) Status(code int) *Response {
	resp.Code = code
	return resp
}

// ContentType is a builder method to override the default content type for a response.
func (resp *Response) ContentType(contentType string) *Response {
	resp.contentType = contentType
	return resp
}

// MarshalJSON implements the json.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.body)
}

// MarshalXML implements the xml.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(resp.body)
}

// GetJSON200Response is a constructor method for a Get response.
// A *Response is returned with the configured status code and content type from the spec.
func GetJSON200Response(body []struct {
	Email *string `json:"email,omitempty"`
	ID    *string `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`

	// List of assigned roles with their details
	Roles []struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	} `json:"roles,omitempty"`
	Username *string `json:"username,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostJSON200Response is a constructor method for a Post response.
// A *Response is returned with the configured status code and content type from the spec.
func PostJSON200Response(body struct {
	Email *string `json:"email,omitempty"`
	ID    *string `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`

	// List of assigned roles with their details
	Roles []struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	} `json:"roles,omitempty"`
	Username *string `json:"username,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetIDJSON200Response is a constructor method for a GetID response.
// A *Response is returned with the configured status code and content type from the spec.
func GetIDJSON200Response(body struct {
	Email *string `json:"email,omitempty"`
	ID    *string `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`

	// List of assigned roles with their details
	Roles []struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	} `json:"roles,omitempty"`
	Username *string `json:"username,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PutIDJSON200Response is a constructor method for a PutID response.
// A *Response is returned with the configured status code and content type from the spec.
func PutIDJSON200Response(body struct {
	Email *string `json:"email,omitempty"`
	ID    *string `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`

	// List of assigned roles with their details
	Roles []struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	} `json:"roles,omitempty"`
	Username *string `json:"username,omitempty"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get a list of users
	// (GET /)
	Get(w http.ResponseWriter, r *http.Request) *Response
	// Create a new user
	// (POST /)
	Post(w http.ResponseWriter, r *http.Request) *Response
	// Delete user by UUID
	// (DELETE /{id})
	DeleteID(w http.ResponseWriter, r *http.Request, id string) *Response
	// Get user details by UUID
	// (GET /{id})
	GetID(w http.ResponseWriter, r *http.Request, id string) *Response
	// Update user details by UUID
	// (PUT /{id})
	PutID(w http.ResponseWriter, r *http.Request, id string) *Response
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// Get operation middleware
func (siw *ServerInterfaceWrapper) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.Get(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// Post operation middleware
func (siw *ServerInterfaceWrapper) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.Post(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// DeleteID operation middleware
func (siw *ServerInterfaceWrapper) DeleteID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.DeleteID(w, r, id)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetID operation middleware
func (siw *ServerInterfaceWrapper) GetID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetID(w, r, id)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PutID operation middleware
func (siw *ServerInterfaceWrapper) PutID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PutID(w, r, id)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter %s: %v", err.paramName, err.err)
}

func (err UnescapedCookieParamError) Unwrap() error { return err.err }

type UnmarshalingParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnmarshalingParamError) Error() string {
	return fmt.Sprintf("error unmarshaling parameter %s as JSON: %v", err.paramName, err.err)
}

func (err UnmarshalingParamError) Unwrap() error { return err.err }

type RequiredParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err RequiredParamError) Error() string {
	if err.err == nil {
		return fmt.Sprintf("query parameter %s is required, but not found", err.paramName)
	} else {
		return fmt.Sprintf("query parameter %s is required, but errored: %s", err.paramName, err.err)
	}
}

func (err RequiredParamError) Unwrap() error { return err.err }

type RequiredHeaderError struct {
	paramName string
}

// Error implements error.
func (err RequiredHeaderError) Error() string {
	return fmt.Sprintf("header parameter %s is required, but not found", err.paramName)
}

type InvalidParamFormatError struct {
	err       error
	paramName string
}

// Error implements error.
func (err InvalidParamFormatError) Error() string {
	return fmt.Sprintf("invalid format for parameter %s: %v", err.paramName, err.err)
}

func (err InvalidParamFormatError) Unwrap() error { return err.err }

type TooManyValuesForParamError struct {
	NumValues int
	paramName string
}

// Error implements error.
func (err TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("expected one value for %s, got %d", err.paramName, err.NumValues)
}

// ParameterName is an interface that is implemented by error types that are
// relevant to a specific parameter.
type ParameterError interface {
	error
	// ParamName is the name of the parameter that the error is referring to.
	ParamName() string
}

func (err UnescapedCookieParamError) ParamName() string  { return err.paramName }
func (err UnmarshalingParamError) ParamName() string     { return err.paramName }
func (err RequiredParamError) ParamName() string         { return err.paramName }
func (err RequiredHeaderError) ParamName() string        { return err.paramName }
func (err InvalidParamFormatError) ParamName() string    { return err.paramName }
func (err TooManyValuesForParamError) ParamName() string { return err.paramName }

type ServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

type ServerOption func(*ServerOptions)

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface, opts ...ServerOption) http.Handler {
	options := &ServerOptions{
		BaseURL:    "/",
		BaseRouter: chi.NewRouter(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	for _, f := range opts {
		f(options)
	}

	r := options.BaseRouter
	wrapper := ServerInterfaceWrapper{
		Handler:          si,
		ErrorHandlerFunc: options.ErrorHandlerFunc,
	}

	r.Route(options.BaseURL, func(r chi.Router) {
		r.Get("/", wrapper.Get)
		r.Post("/", wrapper.Post)
		r.Delete("/{id}", wrapper.DeleteID)
		r.Get("/{id}", wrapper.GetID)
		r.Put("/{id}", wrapper.PutID)
	})
	return r
}

func WithRouter(r chi.Router) ServerOption {
	return func(s *ServerOptions) {
		s.BaseRouter = r
	}
}

func WithServerBaseURL(url string) ServerOption {
	return func(s *ServerOptions) {
		s.BaseURL = url
	}
}

func WithErrorHandler(handler func(w http.ResponseWriter, r *http.Request, err error)) ServerOption {
	return func(s *ServerOptions) {
		s.ErrorHandlerFunc = handler
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xVTW/bMAz9KwLPRpNtN9+2BRsCbEAvORXFoFpMwsKWVIlqYRj+7wPltEsWJ9nWYluB",
	"XGJH5tcj3xM7qFzjnUXLEcquL4Ds0kHZARPXCCUsIgb1VVu9wgYtq/eXcyjgHkMkZ6GENxfTiyn0BTiP",
	"VnuCEt7lowK85rVEhYn8rJDl4TwGzeTs3EAJn5GhgIDROxsxG7+dTuVROctos4v2vqYqO01uo2TtIFZr",
	"bLS8EWOTHX2Q2ExDGGw01RlI6wVH5EB2JYWSkeOlC41mKCElMlDsm1nd4Kh/cPWQwmCsAnkeGvGFIiu3",
	"VDpGWlk0KtupB+K14jVSUAZZUx2hOFTz8yrrn0zdzS1WDD8OdAi6lf8pYvhjdzEZxyxhYw4RU9Po0A6j",
	"VVrVuxYFeBdHeHApp0KEu4SRPzjT/hYHfnX0j9B3UXxKda3kkxTKa8zFjjVeJvqNzJHhi4Waz6JityGC",
	"vG3FfJr8ySkfG91u8oWluzRkyCiWLhyB0Q9tpoAGyqtNs7bCX+/xYNeFQ8L+mZo9S/WFpbovzXxvVwE1",
	"o1ExVRXGuEx13f4k04/ZRGll8WFgjHyfdGT6oXE1Mu4LdpbP57N8zQfdIIu+y6sOSLLL1Q+P8CH3ZpdB",
	"xRYbTnSxvx5n2wjeodqjeIfCM1J106rFYj6Tbh/aTv8U4VlP/5meNjBUQA6E98eZJgswbXtt0c2nsSWY",
	"/ibdXmLTnvfpIfKc9+Vr1nfy5tTeXGSTAwLv+/57AAAA//+zrhbOWg0AAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
