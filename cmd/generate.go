package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/parser"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate supergraph and node-resolver graph",
	Run: func(cmd *cobra.Command, args []string) {
		generate(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

func generate(ctx context.Context) {
	sources := []*ast.Source{}
	interfaces := map[string]ast.DefinitionList{}

	for _, f := range []string{
		"../metadata-api/schema.graphql",
		"../tenant-api/schema.graphql",
		"../location-api/schema.graphql",
		"../resource-provider-api/schema.graphql",
		"../load-balancer-api/schema.graphql",
		"../ipam-api/schema.graphql",
	} {
		b, err := os.ReadFile(f)
		if err != nil {
			logger.Fatalw("error reading file", "file", f, "error", err)
		}

		sources = append(sources, &ast.Source{
			Name:  f,
			Input: string(b),
		})

	}

	schemaDoc, err := parser.ParseSchemas(sources...)
	if err != nil {
		logger.Fatalw("error parsing schemas", "error", err)
	}

	resolverSchema := ast.Schema{}
	resolverObjects := map[string]*ast.Definition{}

	for _, s := range schemaDoc.Definitions {
		for _, i := range s.Interfaces {
			interfaces[i] = append(interfaces[i], s)
		}
	}

	for i, defs := range interfaces {
		intDef := &ast.Definition{
			Kind:       ast.Interface,
			Name:       i,
			Directives: ast.DirectiveList{&keyIDDirective},
			Fields: ast.FieldList{
				{
					Name: "id",
					Type: ast.NonNullNamedType("ID", nil),
				},
			},
		}

		resolverSchema.AddTypes(intDef)

		for _, d := range defs {
			obj, exists := resolverObjects[d.Name]
			if !exists {
				obj = &ast.Definition{
					Kind:       ast.Object,
					Name:       d.Name,
					Directives: ast.DirectiveList{&keyIDDirective},
					Fields: ast.FieldList{
						{
							Name: "id",
							Type: ast.NonNullNamedType("ID", nil),
						},
					},
				}

				pd := d.Directives.ForName("prefixedID")
				if pd != nil {
					obj.Directives = append(obj.Directives, pd)
				} else {
					fmt.Printf("WARNING: %s doesn't have a @prefixedID directive, lookups will fail\n", d.Name)
				}

				resolverObjects[d.Name] = obj
				resolverSchema.AddTypes(obj)
			}
			obj.Interfaces = append(obj.Interfaces, i)
		}
	}

	// for _, t := range resolverObjects {
	// 	resolverSchema.AddTypes(t)
	// }

	f, err := os.Create("schema.graphql")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fmtr := formatter.NewFormatter(f)
	fmtr.FormatSchema(&resolverSchema)
}

var keyIDDirective = ast.Directive{
	Name: "key",
	Arguments: []*ast.Argument{
		{
			Name: "fields",
			Value: &ast.Value{
				Raw:  "id",
				Kind: ast.StringValue,
			},
		},
	},
}

var externalDirective = ast.Directive{
	Name: "external",
}
