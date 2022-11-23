/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/nullc4t/og/pkg/editor"
	"github.com/nullc4t/og/pkg/extract"
	"github.com/nullc4t/og/pkg/generator"
	"github.com/nullc4t/og/pkg/names"
	"github.com/nullc4t/og/pkg/templates"
	"github.com/nullc4t/og/pkg/writer"
	"github.com/spf13/viper"
	"go/parser"
	"go/token"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

// protocolCmd represents the protocol command
var protocolCmd = &cobra.Command{
	Use:   "protocol",
	Short: "Create request & response types",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("protocol called")
		logger.Println("files:", viper.GetStringSlice("files"))

		tmpl := template.Must(template.New("struct").Parse(templates.StructTemplate))
		tmpl = template.Must(tmpl.New("").Parse(templates.ProtocolTemplate))

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, args[0], nil, parser.ParseComments)
		if err != nil {
			logger.Fatal(err)
		}

		ifaces := extract.Interfaces(file)

		type StructDot struct {
			StructName string
			Fields     extract.Args
		}

		var sds []StructDot
		for _, iface := range ifaces {
			for _, method := range iface.Methods {
				//logger.Println(fmt.Sprintf("%sRequest", method.Name))

				sd := StructDot{StructName: fmt.Sprintf("%sRequest", method.Name)}
				for _, arg := range method.Args {
					if !ArgIsContext(arg) {
						sd.Fields = append(sd.Fields, arg)
					}
					//logger.Println(arg)
				}
				sds = append(sds, sd)

				//logger.Println(fmt.Sprintf("%sResponse", method.Name))

				sd = StructDot{StructName: fmt.Sprintf("%sResponse", method.Name)}
				for _, arg := range method.Results.Args {
					sd.Fields = append(sd.Fields, arg)
					//logger.Println(arg)
				}
				sds = append(sds, sd)
				//logger.Println()
			}

			for _, i := range iface.Imports {
				logger.Println(iface.Name, i.Name, i.Path)
			}

			unit := generator.NewUnit(nil, tmpl, map[string]any{
				"Package": file.Name.Name,
				"Structs": sds,
			}, []editor.CodeEditor{
				//editor.AddNamedImportsFactory(iface.Imports...),
			},
				[]editor.ASTEditor{
					editor.ASTImportsFactory(iface.Imports...),
				}, filepath.Join(
					filepath.Dir(args[0]),
					names.FileNameWithSuffix(iface.Name, "protocol"),
				), writer.File)
			err = unit.Generate()
			if err != nil {
				logger.Fatal("generate protocol error:", err)
			}
		}
	},
}

func ArgIsContext(arg extract.Arg) bool {
	return arg.Type.String() == "context.Context"
}

func ArgIsError(arg extract.Arg) bool {
	return arg.Type.String() == "error"
}

func init() {
	rootCmd.AddCommand(protocolCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// protocolCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// protocolCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
