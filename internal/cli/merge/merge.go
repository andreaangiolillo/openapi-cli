package merge

import (
	"andreaangiolillo/openapi-cli/internal/openapi"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type Opts struct {
	Base       *openapi.V3Document
	Merger     openapi.Merger
	outputPath string
}

func (o *Opts) Run(args []string) error {
	federated, err := o.Merger.Merge(args[1:])
	if err != nil {
		return err
	}

	return o.SaveFile(federated)
}

func (o *Opts) SaveFile(federated *openapi.V3Document) error {
	data, err := json.MarshalIndent(federated, "", "    ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(o.outputPath, data, 0644); err != nil {
		return err
	}

	_, _ = fmt.Printf("Federated Spec was saved in '%s'\n", o.outputPath)
	return nil
}

func (o *Opts) PreRunE(args []string) error {
	d, err := openapi.NewV3Document(args[0])
	if err != nil {
		return err
	}
	o.Base = d
	o.Merger = openapi.NewV3Merge(d)
	return nil
}

func Builder() *cobra.Command {
	opts := &Opts{}

	cmd := &cobra.Command{
		Use:   "merge [base-spec] [spec-1] [spec-2] [spec-3] ... [spec-n]",
		Short: "Merge Open API specifications into a base spec.",
		Args:  cobra.MinimumNArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.PreRunE(args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.Run(args)
		},
	}

	cmd.Flags().StringVarP(&opts.outputPath, "output", "o", "federated.json", "File name of the merged spec")
	return cmd
}
