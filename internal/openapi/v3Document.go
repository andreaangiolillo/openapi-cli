package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	model "github.com/pb33f/libopenapi/datamodel/high/base"
	v3model "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/index"
	"github.com/tufin/oasdiff/load"
	"gopkg.in/yaml.v3"
	"os"
)

type V3Document struct {
	// Version is the version of OpenAPI being used, extracted from the 'openapi: x.x.x' definition.
	// This is not a standard property of the OpenAPI model, it's a convenience mechanism only.
	Version string `json:"openapi,omitempty" yaml:"openapi,omitempty"`

	// Info represents a specification Info definitions
	// Provides metadata about the API. The metadata MAY be used by tooling as required.
	// - https://spec.openapis.org/oas/v3.1.0#info-object
	Info *Info `json:"info,omitempty" yaml:"info,omitempty"`

	// Servers is a slice of Server instances which provide connectivity information to a target server. If the servers
	// property is not provided, or is an empty array, the default value would be a Server Object with an url value of /.
	// - https://spec.openapis.org/oas/v3.1.0#server-object
	Servers []*v3model.Server `json:"servers,omitempty" yaml:"servers,omitempty"`

	// Tags is a slice of base.Tag instances defined by the specification
	// A list of tags used by the document with additional metadata. The order of the tags can be used to reflect on
	// their order by the parsing tools. Not all tags that are used by the Operation Object must be declared.
	// The tags that are not declared MAY be organized randomly or based on the tools’ logic.
	// Each tag name in the list MUST be unique.
	// - https://spec.openapis.org/oas/v3.1.0#tag-object
	Tags []*Tag `json:"tags,omitempty" yaml:"tags,omitempty"`

	// Paths contains all the PathItem definitions for the specification.
	// The available paths and operations for the API, The most important part of ths spec.
	// - https://spec.openapis.org/oas/v3.1.0#paths-object
	Paths *v3model.Paths `json:"paths,omitempty" yaml:"paths,omitempty"`

	// Components is an element to hold various schemas for the document.
	// - https://spec.openapis.org/oas/v3.1.0#components-object
	Components *Components `json:"components,omitempty" yaml:"components,omitempty"`
}

// Info represents a high-level Info object as defined by both OpenAPI 2 and OpenAPI 3.
type Info struct {
	Summary        string         `json:"summary,omitempty" yaml:"summary,omitempty"`
	Title          string         `json:"title,omitempty" yaml:"title,omitempty"`
	Description    string         `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string         `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *model.Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *model.License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string         `json:"version,omitempty" yaml:"version,omitempty"`
}

// Components represents a high-level OpenAPI 3+ Components Object
type Components struct {
	Schemas         map[string]*model.Schema           `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]*v3model.Response       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]*v3model.Parameter      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	SecuritySchemes map[string]*v3model.SecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
}

// Tag represents a high-level Tag instance that is backed by a low-level one.
type Tag struct {
	Name         string             `json:"name,omitempty" yaml:"name,omitempty"`
	Description  string             `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *model.ExternalDoc `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// NewV3Document returns and openAPI V3 document from a path
func NewV3Document(path string) (*V3Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Read the first few bytes to check for BOM
	bom := make([]byte, 3)
	_, err = file.Read(bom)
	if err != nil {
		return nil, err
	}

	// Check if BOM exists
	if bom[0] != 0xEF || bom[1] != 0xBB || bom[2] != 0xBF {
		// No BOM, reset the read pointer to the start of the file
		_, err = file.Seek(0, 0)
		if err != nil {
			return nil, err
		}
	}

	var doc *V3Document

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func NewSpecInfo(path string) (*load.SpecInfo, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	openapi3.CircularReferenceCounter = 10

	spec, err := load.LoadSpecInfo(loader, load.NewSource(path))
	if err != nil {
		return nil, err
	}

	return spec, nil
}

// NewDocument returns and openAPI V3 document from a path
func NewDocument(path string) (*V3Document, error) {

	// load an OpenAPI 3 specification with circular refs from bytes
	circularBytes, _ := os.ReadFile(path)

	// create a new root node *yaml.Node reference
	var rootNode yaml.Node

	// unmarshal the specification bytes into the root node
	_ = yaml.Unmarshal(circularBytes, &rootNode)

	// create and index from the specification root node
	idx := index.NewSpecIndex(&rootNode)

	// create a new resolver from the index
	resolverRef := index.NewResolver(idx)

	// resolve the openapi specifications
	resolvingErrors := resolverRef.Resolve()

	// any errors found during resolution? Print them out.
	for _, err := range resolvingErrors {
		fmt.Printf("Error: %s\n", err.Error())
	}

	// create a new document from specification bytes
	document, _ := libopenapi.NewDocument(circularBytes)

	document.SetConfiguration(&datamodel.DocumentConfiguration{
		BasePath:                            "test/data/external_ref",
		AllowFileReferences:                 true,
		AllowRemoteReferences:               false,
		AvoidIndexBuild:                     false,
		BypassDocumentCheck:                 false,
		IgnorePolymorphicCircularReferences: false,
		IgnoreArrayCircularReferences:       false,
		SkipCircularReferenceCheck:          true,
	})
	// build a v3 model from the document.
	v3Model, errs := document.BuildV3Model()
	if errs != nil {
		for _, e := range errs {
			print(e.Error())
		}
	}

	v3Model.Index.GetAllCombinedReferences()

	components, err := newComponents(v3Model.Model.Components)
	if err != nil {
		return nil, err
	}
	return &V3Document{
		Version: v3Model.Model.Version,
		Info: &Info{
			Summary:        v3Model.Model.Info.Summary,
			Title:          v3Model.Model.Info.Title,
			Description:    v3Model.Model.Info.Description,
			TermsOfService: v3Model.Model.Info.TermsOfService,
			Contact:        v3Model.Model.Info.Contact,
			License:        v3Model.Model.Info.License,
			Version:        v3Model.Model.Info.Version,
		},
		Servers:    v3Model.Model.Servers,
		Tags:       newTag(v3Model.Model.Tags),
		Paths:      v3Model.Model.Paths,
		Components: components,
	}, nil
}

func newTag(tags []*model.Tag) []*Tag {
	out := make([]*Tag, len(tags))
	for i, v := range tags {
		out[i] = &Tag{
			Name:         v.Name,
			Description:  v.Description,
			ExternalDocs: v.ExternalDocs,
		}
	}

	return out
}

//func newPaths(paths *v3model.Paths) *v3model.Paths {
//	//for _, v := range paths.PathItems {
//	//		//res2B, _ := json.MarshalIndent(v.GoLow().Description, "", "	")
//	//		//fmt.Println(string(res2B))
//	//		print("\n\n \n\n")
//	//}
//	print("\n\n ANDREA PATH LOW: \n\n")
//	b, _ := paths.Render()
//	print(string(b))
//
//	return paths
//}

func newComponents(components *v3model.Components) (*Components, error) {
	out := &Components{
		SecuritySchemes: components.SecuritySchemes,
		Responses:       components.Responses,
	}
	schemas := map[string]*model.Schema{}
	for k, s := range components.Schemas {
		schemas[k] = s.Schema()
		if err := s.GetBuildError(); err != nil {
			return nil, err
		}
	}
	out.Schemas = schemas

	for _, v := range components.Parameters {
		v.Schema.Schema()
		if err := v.Schema.GetBuildError(); err != nil {
			return nil, err
		}
	}

	out.Parameters = components.Parameters
	return out, nil
}
