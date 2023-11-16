package openapi

import (
	"andreaangiolillo/openapi-cli/internal/openapi/errors"
	"fmt"
	model "github.com/pb33f/libopenapi/datamodel/high/base"
	v3model "github.com/pb33f/libopenapi/datamodel/high/v3"
	"os"
	"reflect"
)

type V3Merge struct {
	base *V3Document
}

func NewV3Merge(base *V3Document) *V3Merge {
	return &V3Merge{
		base: base,
	}
}

func (o V3Merge) mergeSpecIntoBase(spec *V3Document) (*V3Document, error) {
	federatedSpec := &V3Document{
		Version: o.base.Version,
		Info:    o.base.Info,
		Servers: o.base.Servers,
	}

	paths, err := mergePaths(o.base.Paths, spec.Paths.PathItems)
	if err != nil {
		return nil, err
	}

	components, err := mergeComponents(o.base.Components, spec.Components)
	if err != nil {
		return nil, err
	}

	tags, err := mergeTags(o.base.Tags, spec.Tags)
	if err != nil {
		return nil, err
	}

	federatedSpec.Paths = paths
	federatedSpec.Components = components
	federatedSpec.Tags = tags
	return federatedSpec, nil
}

func mergePaths(basePaths *v3model.Paths, pathsToMerge map[string]*v3model.PathItem) (*v3model.Paths, error) {
	outPathItems := map[string]*v3model.PathItem{}

	// Copy base path to the federated spec paths
	for k, v := range basePaths.PathItems {
		outPathItems[k] = v
	}

	for k, v := range pathsToMerge {
		if _, ok := outPathItems[k]; !ok {
			outPathItems[k] = v
		} else {
			return nil, errors.PathConflictError{
				Entry: k,
			}
		}
	}

	return &v3model.Paths{
		PathItems: outPathItems,
	}, nil
}

func mergeTags(baseTags []*Tag, tagsToMerge []*Tag) ([]*Tag, error) {
	var out []*Tag
	tagsSet := make(map[string]bool)

	// Copy base tags to the federated spec tags
	for _, v := range baseTags {
		out = append(out, v)
		tagsSet[v.Name] = true
	}

	for _, v := range tagsToMerge {
		if _, ok := tagsSet[v.Name]; !ok {
			out = append(out, v)
		} else {
			return nil, errors.TagConflictError{
				Entry:       v.Name,
				Description: v.Description,
			}
		}
	}

	return out, nil
}

func mergeComponents(baseCps *Components, cpsToMerge *Components) (*Components, error) {
	outComponents := &Components{
		SecuritySchemes: baseCps.SecuritySchemes,
		Parameters:      baseCps.Parameters,
		Responses:       baseCps.Responses,
		Schemas:         baseCps.Schemas,
	}

	if err := mergeParameters(outComponents, cpsToMerge.Parameters); err != nil {
		return nil, err
	}

	if err := mergeResponses(outComponents, cpsToMerge.Responses); err != nil {
		return nil, err
	}

	if err := mergeSchemas(outComponents, cpsToMerge.Schemas); err != nil {
		return nil, err
	}

	return outComponents, nil
}

func mergeParameters(baseCps *Components, params map[string]*v3model.Parameter) error {
	for k, v := range params {
		if _, ok := baseCps.Parameters[k]; !ok {
			baseCps.Parameters[k] = v
		} else {
			if reflect.DeepEqual(baseCps.Schemas[k], v) {
				// if the params are the same, we skip
				_, _ = fmt.Fprintf(os.Stderr, "We silently resolved the conflict with the params '%s' because the definition was identical.", k)
				continue
			}

			// The params have the same name but different definitions
			return errors.ParamConflictError{
				Entry: k,
			}
		}
	}

	return nil
}

func mergeResponses(baseCps *Components, responses map[string]*v3model.Response) error {
	for k, v := range responses {
		if _, ok := baseCps.Responses[k]; !ok {
			baseCps.Responses[k] = v
		} else {
			if reflect.DeepEqual(baseCps.Schemas[k], v) {
				// if the schemas are the same, we skip
				_, _ = fmt.Fprintf(os.Stderr, "We silently resolved the conflict with the response '%s' because the definition was identical.", k)
				continue
			}

			// The responses have the same name but different definitions
			return errors.ResponseConflictError{
				Entry: k,
			}
		}
	}

	return nil
}

func mergeSchemas(baseCps *Components, schemas map[string]*model.Schema) error {
	for k, schemaToMerge := range schemas {
		if _, ok := baseCps.Schemas[k]; !ok {
			baseCps.Schemas[k] = schemaToMerge
		} else {

			if reflect.DeepEqual(baseCps.Schemas[k], schemaToMerge) {
				// if the schemas are the same, we skip
				_, _ = fmt.Fprintf(os.Stderr, "We silently resolved the conflict with the schema '%s' because the definition was identical.", k)
				continue
			}

			// The schemas have the same name but different definitions
			return errors.SchemaConflictError{
				Entry: k,
			}
		}
	}

	return nil
}
