package installer

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// renderSchemaExampleRegion turns a canonical schemas/<name>.schema.yaml into the
// FULL inline marker-region BODY for its producing step .md — the bytes that sit
// strictly between the begin and end markers. It is the schema-sot generator
// counterpart of renderInlineStepsRegion: a pure, deterministic function with no
// I/O, so committed == rendered can be asserted byte-for-byte.
//
// The output mirrors renderInlineStepsRegion's leading/trailing newline padding
// so the begin marker, the fence label, the fenced block, and the end marker
// each occupy their own line after splicing:
//
//	"\n" + fenceLabel + "\n" + "```yaml\n" + <payload mapping> + "```\n"
//
// The <payload mapping> is a `payload:` mapping whose value mirrors the schema's
// top-level `properties:` as a single, complete example (required AND optional
// keys, one element per array). Determinism is guaranteed by driving key order
// strictly from the schema `properties:` MappingNode Content order (never a Go
// map) and by setting every scalar's Style explicitly.
func renderSchemaExampleRegion(schemaBytes []byte, fenceLabel string) (string, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(schemaBytes, &root); err != nil {
		return "", fmt.Errorf("parse schema: %w", err)
	}

	doc := &root
	if doc.Kind == yaml.DocumentNode {
		if len(doc.Content) == 0 {
			return "", fmt.Errorf("schema is empty")
		}
		doc = doc.Content[0]
	}
	if doc.Kind != yaml.MappingNode {
		return "", fmt.Errorf("schema root is not a mapping")
	}

	props := schemaMappingValue(doc, "properties")
	if props == nil || props.Kind != yaml.MappingNode {
		return "", fmt.Errorf("schema has no top-level properties mapping")
	}

	exampleProps, err := schemaNodeToExample(props)
	if err != nil {
		return "", err
	}

	// Wrap the rendered properties under a `payload:` key.
	payload := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
		Content: []*yaml.Node{
			keyScalar("payload"),
			exampleProps,
		},
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(payload); err != nil {
		return "", fmt.Errorf("encode example: %w", err)
	}
	if err := enc.Close(); err != nil {
		return "", fmt.Errorf("close encoder: %w", err)
	}

	rendered := buf.String()
	// Strip any trailing blank lines beyond the single newline the fence needs.
	rendered = strings.TrimRight(rendered, "\n") + "\n"

	region := "\n" + fenceLabel + "\n" + "```yaml\n" + rendered + "```\n"
	return region, nil
}

// schemaNodeToExample builds an example mapping `*yaml.Node` from a schema
// `properties:` MappingNode. Key order follows the schema property order exactly
// (Content is walked in pairs, never a Go map). Each property value is rendered
// once, as a complete example.
func schemaNodeToExample(props *yaml.Node) (*yaml.Node, error) {
	out := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	for i := 0; i+1 < len(props.Content); i += 2 {
		keyNode := props.Content[i]
		propSchema := props.Content[i+1]

		valNode, err := schemaPropertyToExample(propSchema)
		if err != nil {
			return nil, fmt.Errorf("property %q: %w", keyNode.Value, err)
		}

		out.Content = append(out.Content, keyScalar(keyNode.Value), valNode)
	}
	return out, nil
}

// schemaPropertyToExample renders ONE schema property definition to an example
// value node, attaching a trailing line-comment derived from `description` and
// any `enum` hint when present.
func schemaPropertyToExample(propSchema *yaml.Node) (*yaml.Node, error) {
	if propSchema.Kind != yaml.MappingNode {
		return scalar(""), nil
	}

	typ := scalarValue(propSchema, "type")
	desc := scalarValue(propSchema, "description")
	enum := schemaMappingValue(propSchema, "enum")

	var node *yaml.Node

	switch typ {
	case "object":
		nested := schemaMappingValue(propSchema, "properties")
		if nested != nil && nested.Kind == yaml.MappingNode {
			n, err := schemaNodeToExample(nested)
			if err != nil {
				return nil, err
			}
			node = n
		} else {
			node = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		}
	case "array":
		items := schemaMappingValue(propSchema, "items")
		if items != nil && items.Kind == yaml.MappingNode && scalarValue(items, "type") == "object" {
			nested := schemaMappingValue(items, "properties")
			elem := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
			if nested != nil && nested.Kind == yaml.MappingNode {
				n, err := schemaNodeToExample(nested)
				if err != nil {
					return nil, err
				}
				elem = n
			}
			node = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: []*yaml.Node{elem}}
		} else {
			// items is scalar (or unspecified) → empty sequence rendered as [].
			node = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Style: yaml.FlowStyle}
		}
	case "integer", "number":
		node = explicitScalar("0", "!!int")
		if typ == "number" {
			node.Tag = "!!float"
		}
	case "boolean":
		node = explicitScalar("false", "!!bool")
	case "string":
		if enum != nil && enum.Kind == yaml.SequenceNode && len(enum.Content) > 0 {
			node = scalar(enum.Content[0].Value)
		} else {
			node = scalar("")
		}
	default:
		// Unknown/absent type: fall back to an empty string scalar.
		node = scalar("")
	}

	if c := buildComment(desc, enum); c != "" {
		node.LineComment = c
	}

	return node, nil
}

// buildComment composes the trailing line-comment for a property from its
// description and enum hint. Deterministic: "<description> | one of: A|B|C",
// either part omitted when absent.
func buildComment(desc string, enum *yaml.Node) string {
	var parts []string
	if d := strings.TrimSpace(collapseWS(desc)); d != "" {
		parts = append(parts, d)
	}
	if enum != nil && enum.Kind == yaml.SequenceNode && len(enum.Content) > 0 {
		members := make([]string, len(enum.Content))
		for i, m := range enum.Content {
			members[i] = m.Value
		}
		parts = append(parts, "one of: "+strings.Join(members, "|"))
	}
	if len(parts) == 0 {
		return ""
	}
	return "# " + strings.Join(parts, " | ")
}

// collapseWS folds any run of whitespace (including newlines from folded YAML
// `description: >` scalars) into single spaces so comments stay one line.
func collapseWS(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// scalar returns a double-quoted string scalar node with an explicit style so
// the encoder output is deterministic.
func scalar(v string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: v,
		Style: yaml.DoubleQuotedStyle,
	}
}

// keyScalar returns a plain (bare, unquoted) string scalar node for use as a
// mapping key, matching the hand-authored inline-block style where keys are
// unquoted and only values carry quotes.
func keyScalar(v string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: v,
		Style: 0,
	}
}

// explicitScalar returns a non-string scalar (int/float/bool) with no quoting
// style, so it renders as a bare token (0, false) and stays valid YAML.
func explicitScalar(v, tag string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   tag,
		Value: v,
		Style: 0,
	}
}

// schemaMappingValue returns the value node for key in a MappingNode, or nil.
// Walks Content in pairs — never a Go map — to preserve determinism, and
// nil-guards the node (the package-level mappingValue does not).
func schemaMappingValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

// scalarValue returns the string value for key when its value is a scalar.
func scalarValue(node *yaml.Node, key string) string {
	v := schemaMappingValue(node, key)
	if v == nil || v.Kind != yaml.ScalarNode {
		return ""
	}
	return v.Value
}
