package sanityref

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/wearemojo/mojo-public-go/lib/gerrors"
	"github.com/wearemojo/mojo-public-go/lib/gjson"
	"github.com/wearemojo/mojo-public-go/lib/merr"
)

// Queried with:
//
//	*[
//		(_type == 'FlowGroup' && condition != null)
//		|| _type match 'Expression*'
//	] {
//		_id,
//		_type,
//		_type == 'FlowGroup' => {
//			condition,
//		},
//		_type match 'Expression*' => {
//			...,
//		},
//		_type == 'ExpressionPollAnswerKeys' => {
//			poll-> {
//				id,
//			},
//		},
//	}
//
//go:embed sanityref_test_input.json
var testInputRaw []byte
var testInput = gjson.MustUnmarshal[[]Document](testInputRaw)

//go:embed sanityref_test_output.json
var testOutputRaw []byte
var testOutput = gjson.MustUnmarshal[map[string]Document](testOutputRaw)

func TestResolveReferencesNormal(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)

	copiedTestInput := inefficientlyDeepCopy(testInput)

	documentMap, missingDocumentIDs, err := ResolveReferences(ctx, copiedTestInput)
	is.NoErr(err)
	is.Equal(documentMap, testOutput)
	is.True(missingDocumentIDs != nil)
	is.Equal(missingDocumentIDs.ToSlice(), []string{})

	// ensures no infinite recursion
	_, err = json.Marshal(documentMap)
	is.NoErr(err)

	// ensures no mutation
	is.Equal(copiedTestInput, testInput)
}

func TestResolveReferencesInfiniteRecursion(t *testing.T) {
	// this is not exactly recommended, but it's technically valid and supported
	// by Sanity, and therefore also supported by this package

	// this test also helps prove we're mutating rather than copying anything, as
	// any cycles would break or cause a stack overflow

	ctx := context.Background()
	is := is.New(t)

	documentMap, missingDocumentIDs, err := ResolveReferences(ctx, []Document{
		{
			"_id": "id1",
			"ref": map[string]any{
				"_ref": "id2",
			},
		},
		{
			"_id": "id2",
			"ref": map[string]any{
				"_ref": "id1",
			},
		},
		{
			"_id": "id3",
			"ref": map[string]any{
				"_ref": "id3",
			},
		},
	})
	is.NoErr(err)
	is.Equal(len(documentMap), 3)
	is.True(documentMap["id1"] != nil)
	is.Equal(documentMap["id1"]["_id"], "id1")
	is.Equal(documentMap["id1"]["ref"], documentMap["id2"])
	is.True(documentMap["id2"] != nil)
	is.Equal(documentMap["id2"]["_id"], "id2")
	is.Equal(documentMap["id2"]["ref"], documentMap["id1"])
	is.True(documentMap["id3"] != nil)
	is.Equal(documentMap["id3"]["_id"], "id3")
	is.Equal(documentMap["id3"]["ref"], documentMap["id3"])
	is.True(missingDocumentIDs != nil)
	is.Equal(missingDocumentIDs.ToSlice(), []string{})

	_, err = json.Marshal(documentMap)
	is.Equal(err.Error(), "json: unsupported value: encountered a cycle via map[string]interface {}")
}

func TestResolveReferencesMissingDocuments(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)

	documentMap, missingDocumentIDs, err := ResolveReferences(ctx, []Document{
		{
			"_id": "id1",
			"ref": map[string]any{
				"_ref": "id2",
			},
		},
	})
	is.NoErr(err)

	is.Equal(documentMap, map[string]Document{
		"id1": {
			"_id": "id1",
			"ref": map[string]any{
				"_ref": "id2",
			},
		},
	})

	is.True(missingDocumentIDs != nil)
	is.Equal(missingDocumentIDs.ToSlice(), []string{"id2"})
}

func TestResolveReferencesStrictMissingDocuments(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)

	documentMap, err := ResolveReferencesStrict(ctx, []Document{
		{
			"_id": "id1",
			"ref": map[string]any{
				"_ref": "id2",
			},
		},
	}, nil)
	is.Equal(documentMap, nil)

	err2, ok := gerrors.As[merr.E](err)
	is.True(ok)
	is.Equal(err2.Code, ErrReferencedDocumentMissing)
	is.Equal(err2.Meta, merr.M{"missing_ids": []string{"id2"}})
}

func TestResolveReferencesMissingImages(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)

	documentMap, missingDocumentIDs, err := ResolveReferences(ctx, []Document{
		{
			"_id": "id1",
			"image": map[string]any{
				"_ref": "image-id2",
			},
			"ref": map[string]any{
				"_ref": "id3",
			},
		},
		{
			"_id": "id3",
			"image": map[string]any{
				"_ref": "image-id4",
			},
		},
	})
	is.NoErr(err)

	is.Equal(documentMap, map[string]Document{
		"id1": {
			"_id": "id1",
			"image": map[string]any{
				"_ref": "image-id2",
			},
			"ref": map[string]any{
				"_id": "id3",
				"image": map[string]any{
					"_ref": "image-id4",
				},
			},
		},
		"id3": {
			"_id": "id3",
			"image": map[string]any{
				"_ref": "image-id4",
			},
		},
	})

	is.Equal(missingDocumentIDs.ToSlice(), []string{"image-id2", "image-id4"})

	// ensures no infinite recursion
	_, err = json.Marshal(documentMap)
	is.NoErr(err)
}

func TestResolveReferencesStrictMissingImagesDefault(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)

	documentMap, err := ResolveReferencesStrict(ctx, []Document{
		{
			"_id": "id1",
			"image": map[string]any{
				"_ref": "image-id2",
			},
			"ref": map[string]any{
				"_ref": "id3",
			},
		},
		{
			"_id": "id3",
			"image": map[string]any{
				"_ref": "image-id4",
			},
		},
	}, nil)
	is.NoErr(err)

	is.Equal(documentMap, map[string]Document{
		"id1": {
			"_id": "id1",
			"image": map[string]any{
				"_ref": "image-id2",
			},
			"ref": map[string]any{
				"_id": "id3",
				"image": map[string]any{
					"_ref": "image-id4",
				},
			},
		},
		"id3": {
			"_id": "id3",
			"image": map[string]any{
				"_ref": "image-id4",
			},
		},
	})

	// ensures no infinite recursion
	_, err = json.Marshal(documentMap)
	is.NoErr(err)
}

func TestResolveReferencesStrictMissingImagesReject(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)

	documentMap, err := ResolveReferencesStrict(ctx, []Document{
		{
			"_id": "id1",
			"image": map[string]any{
				"_ref": "image-id2",
			},
			"ref": map[string]any{
				"_ref": "id3",
			},
		},
		{
			"_id": "id3",
			"image": map[string]any{
				"_ref": "image-id4",
			},
		},
	}, &StrictOptions{
		RejectMissingImages: true,
	})
	is.Equal(documentMap, nil)

	err2, ok := gerrors.As[merr.E](err)
	is.True(ok)
	is.Equal(err2.Code, ErrReferencedDocumentMissing)
	is.Equal(err2.Meta, merr.M{"missing_ids": []string{"image-id2", "image-id4"}})
}
