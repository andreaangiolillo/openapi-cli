package openapi

import (
	"encoding/json"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	model "github.com/pb33f/libopenapi/datamodel/high/base"
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
	Servers *any `json:"servers,omitempty" yaml:"servers,omitempty"`

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
	Paths map[string]any `json:"paths,omitempty" yaml:"paths,omitempty"`

	// Components is an element to hold various schemas for the document.
	// - https://spec.openapis.org/oas/v3.1.0#components-object
	Components *Components `json:"components,omitempty" yaml:"components,omitempty"`
}

// Info represents a high-level Info object as defined by both OpenAPI 2 and OpenAPI 3.
type Info struct {
	Summary        string        `json:"summary,omitempty" yaml:"summary,omitempty"`
	Title          string        `json:"title,omitempty" yaml:"title,omitempty"`
	Description    string        `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string        `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *base.Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *base.License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string        `json:"version,omitempty" yaml:"version,omitempty"`
}

// Components represents a high-level OpenAPI 3+ Components Object
type Components struct {
	Schemas         map[string]*any `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]*any `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]*any `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	SecuritySchemes map[string]any  `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
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
