package sanityref

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/samber/lo"
	"github.com/wearemojo/mojo-public-go/lib/gjson"
	"github.com/wearemojo/mojo-public-go/lib/merr"
	"github.com/wearemojo/mojo-public-go/lib/slicefn"
)

// roughly equivalent to `resolveSanityDocuments` in the TypeScript codebase

const (
	ErrDocumentIDDuplicated      = merr.Code("document_id_duplicated")
	ErrDocumentIDMissing         = merr.Code("document_id_missing")
	ErrInvalidJSONType           = merr.Code("invalid_json_type")
	ErrReferenceFieldInvalid     = merr.Code("reference_field_invalid")
	ErrReferencedDocumentMissing = merr.Code("referenced_document_missing")
)

// Document is a Document from a Sanity dataset
//
// It must have a string `_id` field uniquely identifying the Document
//
// Any nested fields wishing to reference another Document must be a map with a
// `_ref` key, whose value matches the `_id` of the Document to reference
//
// Unmarshal directly into this type to use this package
type Document = map[string]any

type StrictOptions struct {
	RejectMissingImages bool
}

// ResolveReferencesStrict is a strict version of ResolveReferences, which will
// return ErrReferencedDocumentMissing if any referenced Document is missing
func ResolveReferencesStrict(ctx context.Context, documents []Document, opts *StrictOptions) (map[string]Document, error) {
	if opts == nil {
		opts = &StrictOptions{}
	}

	documentMap, missingDocumentIDs, err := ResolveReferences(ctx, documents)
	if err != nil {
		return nil, err
	}

	missingIDs := missingDocumentIDs.ToSlice()
	if !opts.RejectMissingImages {
		missingIDs = slicefn.Filter(missingIDs, func(id string) bool { return !strings.HasPrefix(id, "image-") })
	}
	if len(missingIDs) > 0 {
		return nil, merr.New(ctx, ErrReferencedDocumentMissing, merr.M{
			"missing_ids": missingIDs,
		})
	}

	return documentMap, nil
}

// ResolveReferences recursively resolves all references in the provided
// documents, returning a map of all documents by their ID, and a set of IDs of
// any referenced documents that are missing
func ResolveReferences(ctx context.Context, documents []Document) (
	documentMap map[string]Document,
	missingDocumentIDs mapset.Set[string],
	err error,
) {
	// prevents mutation of the original input
	documents = inefficientlyDeepCopy(documents)

	// beyond this point we should only copy pointers, not the actual data, to
	// ensure we're consistently referencing the original data in memory

	// allows us to efficiently look up documents by their ID
	documentMap, err = slicefn.ReduceE(documents, func(acc map[string]Document, doc Document) (map[string]Document, error) {
		id, ok := doc["_id"].(string)
		if !ok {
			return nil, merr.New(ctx, ErrDocumentIDMissing, merr.M{"document": doc})
		}

		if _, ok := acc[id]; ok {
			return nil, merr.New(ctx, ErrDocumentIDDuplicated, merr.M{"id": id})
		}

		acc[id] = doc
		return acc, nil
	}, map[string]Document{})
	if err != nil {
		return nil, nil, err
	}

	documentsAny := lo.Map(documents, func(doc Document, _ int) any { return doc })
	missingDocumentIDs = mapset.NewThreadUnsafeSet[string]()
	if _, err := recursivelyResolve(ctx, documentsAny, documentMap, missingDocumentIDs); err != nil {
		return nil, nil, err
	}

	return documentMap, missingDocumentIDs, nil
}

// Only works with types that `json` will unmarshal to when targeting `any`
// (https://pkg.go.dev/encoding/json#Unmarshal) - like `[]any`,
// `map[string]any`, and some primitives
//
// Providing e.g. `[]map[string]any` will not work - map to `[]any` first
func recursivelyResolve(
	ctx context.Context,
	data any,
	documentMap map[string]Document,
	missingDocumentIDs mapset.Set[string],
) (replacementValue any, err error) {
	// this function must always mutate, never return new slices/maps, as we're
	// making use of the existing pointers to ensure everything gets updated,
	// even with potentially-infinite recursion

	switch data := data.(type) {
	case map[string]any:
		// if the map has a `_ref` key, this is what we're looking for!
		//
		// so now try to replace it with the actual document
		if refField, ok := data["_ref"]; ok {
			return handleRefField(ctx, data, refField, documentMap, missingDocumentIDs)
		}

		// not a reference, so just keep walking the map
		for key, value := range data {
			if replacementValue, err := recursivelyResolve(ctx, value, documentMap, missingDocumentIDs); err != nil {
				return nil, err
			} else {
				data[key] = replacementValue
			}
		}

		return data, nil

	case []any:
		// slices can't be references themselves, so just walk the slice
		for idx, value := range data {
			if replacementValue, err := recursivelyResolve(ctx, value, documentMap, missingDocumentIDs); err != nil {
				return nil, err
			} else {
				data[idx] = replacementValue
			}
		}

		return data, nil

	case bool, float64, string, nil, json.Number:
		// expected primitive JSON types - also can't be refs, so leave as-is
		return data, nil

	default:
		// we should now have covered all types documented at:
		// https://pkg.go.dev/encoding/json#Unmarshal
		//
		// so if we get here, that indicates the json package has changed (not
		// likely!), or someone's unmarshaled into specific non-`any` types,
		// preventing its basic behavior from applying
		//
		// we'll then try to provide handy tips to help the caller understand what
		// to do, esp if they've unmarshaled into `[]map[string]any` or similar
		//
		// we also only do reflection at this point, to keep the happy path fast
		return nil, handleUnexpectedType(ctx, data)
	}
}

func handleRefField(
	ctx context.Context,
	data map[string]any,
	refField any,
	documentMap map[string]Document,
	missingDocumentIDs mapset.Set[string],
) (replacementValue any, err error) {
	ref, ok := refField.(string)
	if !ok {
		return nil, merr.New(ctx, ErrReferenceFieldInvalid, merr.M{
			"data": data,
			"ref":  refField,
		})
	}

	if otherDoc, ok := documentMap[ref]; ok {
		return otherDoc, nil
	}

	missingDocumentIDs.Add(ref)

	// in theory we could keep resolving within the map, but we never expect
	// to have a `_ref` that isn't an actual reference, so we don't want to
	// enable potential misuse
	return data, nil
}

func handleUnexpectedType(ctx context.Context, data any) error {
	typ := reflect.TypeOf(data)
	//nolint:exhaustive // we only need to validate these types
	switch typ.Kind() {
	case reflect.Slice:
		if _, ok := data.([]any); ok {
			panic("should be caught by the `[]any` case above")
		}
		return merr.New(ctx, ErrInvalidJSONType, merr.M{
			"expected": "[]any",
			"actual":   typ.String(),
			"tip":      "map to []any first",
		})

	case reflect.Map:
		if _, ok := data.(map[string]any); ok {
			panic("should be caught by the `map[string]any` case above")
		}
		return merr.New(ctx, ErrInvalidJSONType, merr.M{
			"expected": "map[string]any",
			"actual":   typ.String(),
			"tip":      "map to map[string]any first",
		})

	default:
		return merr.New(ctx, ErrInvalidJSONType, merr.M{
			"expected": "[]any, map[string]any, bool, float64, string, nil, json.Number",
			"actual":   typ.String(),
			"tip":      "ensure you're unmarshaling into `any`-oriented types, not custom structs",
		})
	}
}

func inefficientlyDeepCopy[T any](v T) T {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return gjson.MustUnmarshal[T](data)
}
