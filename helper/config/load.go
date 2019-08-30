package config

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
)

var configSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "artifact", LabelNames: []string{"type", "name"}},
	},
}

func Load(filename string) (Artifacts, hcl.Diagnostics) {
	var err error
	filename, err = translateBuilder(filename)
	if err != nil {
		var diags hcl.Diagnostics
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "translateBuilder error",
			Detail:   fmt.Sprintf("Err %v", err),
		})
		return nil, diags
	}

	// We'll return as complete a config as we can manage, even if there are
	// errors, since a partial result can be useful for careful analysis by
	// development tools such as text editor extensions.
	artifacts := Artifacts(make(map[ArtifactRef]*Artifact))

	f, diags := parser.ParseFile(filename)
	if diags.HasErrors() {
		return artifacts, diags
	}

	content, moreDiags := f.Body.Content(configSchema)
	diags = append(diags, moreDiags...)

	for _, block := range content.Blocks {
		switch block.Type {
		case "artifact":
			artifact, moreDiags := decodeArtifactConfig(block)
			diags = append(diags, moreDiags...)
			if artifact != nil {
				ref := artifact.Ref()
				if existing := artifacts[ref]; existing != nil {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Duplicate artifact block",
						Detail:   fmt.Sprintf("This artifact block has the same builder type and name as the previous block declared at %s. Each artifact must have a unique name per builder type.", existing.DeclRange),
						Subject:  &artifact.DeclRange,
					})
					continue
				}
				artifacts[ref] = artifact
			}
		default:
			// Only "artifact" is in our schema, so we can never get here
			panic(fmt.Sprintf("unexpected block type %q", block.Type))
		}
	}

	return artifacts, diags
}
