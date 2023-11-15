package openapi

type Merger interface {
	Merge([]string) (*V3Document, error)
}

func (o V3Merge) Merge(paths []string) (*V3Document, error) {
	var federatedSpec *V3Document
	for _, p := range paths {
		spec, err := NewV3Document(p)
		if err != nil {
			return nil, err
		}

		federatedSpec, err = o.mergeSpecIntoBase(spec)
		if err != nil {
			return nil, err
		}
	}
	return federatedSpec, nil
}
