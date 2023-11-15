package openapi

import (
	"fmt"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	model "github.com/pb33f/libopenapi/datamodel/high/base"
	v3model "github.com/pb33f/libopenapi/datamodel/high/v3"
	"os"
)

type V3Document struct {
	// Version is the version of OpenAPI being used, extracted from the 'openapi: x.x.x' definition.
	// This is not a standard property of the OpenAPI model, it's a convenience mechanism only.
	Version string `json:"openapi,omitempty" yaml:"openapi,omitempty"`

	// Info represents a specification Info definitions
	// Provides metadata about the API. The metadata MAY be used by tooling as required.
	// - https://spec.openapis.org/oas/v3.1.0#info-object
	Info *base.Info `json:"info,omitempty" yaml:"info,omitempty"`

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
	Tags []*model.Tag `json:"tags,omitempty" yaml:"tags,omitempty"`

	// Paths contains all the PathItem definitions for the specification.
	// The available paths and operations for the API, The most important part of ths spec.
	// - https://spec.openapis.org/oas/v3.1.0#paths-object
	Paths *v3model.Paths `json:"paths,omitempty" yaml:"paths,omitempty"`

	// Components is an element to hold various schemas for the document.
	// - https://spec.openapis.org/oas/v3.1.0#components-object
	Components *v3model.Components `json:"components,omitempty" yaml:"components,omitempty"`
}

// NewV3Document returns and openAPI V3 document from a path
func NewV3Document(path string) (*V3Document, error) {
	// load the base OpenAPI 3 specification from bytes
	baseSpecBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// create a new document from specification bytes
	d, err := libopenapi.NewDocument(baseSpecBytes)
	if err != nil {
		return nil, err
	}

	// because we know this is a v3 spec, we can build a ready to go v3model from it.
	v3Model, errors := d.BuildV3Model()
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %e\n", errors[i])
		}
		panic(fmt.Sprintf("cannot create v3 v3model from document: %d errors reported",
			len(errors)))
	}

	return &V3Document{
		Version:    v3Model.Model.Version,
		Info:       v3Model.Model.Info,
		Servers:    v3Model.Model.Servers,
		Tags:       v3Model.Model.Tags,
		Paths:      v3Model.Model.Paths,
		Components: v3Model.Model.Components,
	}, nil
}
