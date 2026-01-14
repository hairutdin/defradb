// Copyright 2026 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package one_to_many

import (
	"testing"

	testUtils "github.com/sourcenetwork/defradb/tests/integration"
)

// TestQueryOneToMany_WithShowDeletedOnNestedList tests that showDeleted can be specified
// on nested list relations to control whether deleted documents are returned for that
// specific relation, independent of the parent query's showDeleted setting.
func TestQueryOneToMany_WithShowDeletedOnNestedList(t *testing.T) {
	test := testUtils.TestCase{
		Actions: []any{
			testUtils.CreateDoc{
				CollectionID: 1, // Author
				Doc: `{
					"name": "John Grisham",
					"age": 65
				}`,
			},
			testUtils.CreateDoc{
				CollectionID: 0, // Book
				DocMap: map[string]any{
					"name":      "Painted House",
					"rating":    4.9,
					"author_id": testUtils.NewDocIndex(1, 0),
				},
			},
			testUtils.CreateDoc{
				CollectionID: 0, // Book
				DocMap: map[string]any{
					"name":      "A Time for Mercy",
					"rating":    4.5,
					"author_id": testUtils.NewDocIndex(1, 0),
				},
			},
			// Delete the first book
			testUtils.DeleteDoc{
				CollectionID: 0,
				DocID:        0,
			},
			// Query author without showDeleted - should not show deleted books by default
			testUtils.Request{
				Request: `query {
					Author {
						name
						age
						published {
							_deleted
							name
							rating
						}
					}
				}`,
				Results: map[string]any{
					"Author": []map[string]any{
						{
							"name": "John Grisham",
							"age":  int64(65),
							"published": []map[string]any{
								{
									"_deleted": false,
									"name":     "A Time for Mercy",
									"rating":   4.5,
								},
							},
						},
					},
				},
			},
			// Query author with showDeleted: true on nested relation - should show deleted books
			testUtils.Request{
				Request: `query {
					Author {
						name
						age
						published(showDeleted: true) {
							_deleted
							name
							rating
						}
					}
				}`,
				Results: map[string]any{
					"Author": []map[string]any{
						{
							"name": "John Grisham",
							"age":  int64(65),
							"published": []map[string]any{
								{
									"_deleted": true,
									"name":     "Painted House",
									"rating":   4.9,
								},
								{
									"_deleted": false,
									"name":     "A Time for Mercy",
									"rating":   4.5,
								},
							},
						},
					},
				},
				NonOrderedResults: true,
			},
			// Query author with showDeleted: false explicitly on nested relation
			testUtils.Request{
				Request: `query {
					Author {
						name
						age
						published(showDeleted: false) {
							_deleted
							name
							rating
						}
					}
				}`,
				Results: map[string]any{
					"Author": []map[string]any{
						{
							"name": "John Grisham",
							"age":  int64(65),
							"published": []map[string]any{
								{
									"_deleted": false,
									"name":     "A Time for Mercy",
									"rating":   4.5,
								},
							},
						},
					},
				},
			},
		},
	}

	executeTestCase(t, test)
}

// TestQueryOneToMany_WithShowDeletedOnNestedSingle tests that showDeleted can be specified
// on nested single object relations.
func TestQueryOneToMany_WithShowDeletedOnNestedSingle(t *testing.T) {
	test := testUtils.TestCase{
		Actions: []any{
			testUtils.CreateDoc{
				CollectionID: 1, // Author
				Doc: `{
					"name": "John Grisham",
					"age": 65
				}`,
			},
			testUtils.CreateDoc{
				CollectionID: 0, // Book
				DocMap: map[string]any{
					"name":      "Painted House",
					"rating":    4.9,
					"author_id": testUtils.NewDocIndex(1, 0),
				},
			},
			testUtils.CreateDoc{
				CollectionID: 0, // Book
				DocMap: map[string]any{
					"name":      "A Time for Mercy",
					"rating":    4.5,
					"author_id": testUtils.NewDocIndex(1, 0),
				},
			},
			// Delete the author
			testUtils.DeleteDoc{
				CollectionID: 1,
				DocID:        0,
			},
			// Query books without showDeleted on author relation - should return null author
			testUtils.Request{
				Request: `query {
					Book {
						name
						rating
						author {
							_deleted
							name
							age
						}
					}
				}`,
				Results: map[string]any{
					"Book": []map[string]any{
						{
							"name":   "Painted House",
							"rating": 4.9,
							"author": nil,
						},
						{
							"name":   "A Time for Mercy",
							"rating": 4.5,
							"author": nil,
						},
					},
				},
				NonOrderedResults: true,
			},
			// Query books with showDeleted: true on author relation - should show deleted author
			testUtils.Request{
				Request: `query {
					Book {
						name
						rating
						author(showDeleted: true) {
							_deleted
							name
							age
						}
					}
				}`,
				Results: map[string]any{
					"Book": []map[string]any{
						{
							"name":   "Painted House",
							"rating": 4.9,
							"author": map[string]any{
								"_deleted": true,
								"name":     "John Grisham",
								"age":      int64(65),
							},
						},
						{
							"name":   "A Time for Mercy",
							"rating": 4.5,
							"author": map[string]any{
								"_deleted": true,
								"name":     "John Grisham",
								"age":      int64(65),
							},
						},
					},
				},
				NonOrderedResults: true,
			},
		},
	}

	executeTestCase(t, test)
}

// TestQueryOneToMany_WithShowDeletedOnParentAndChild tests interaction between
// showDeleted on parent query and nested relation.
func TestQueryOneToMany_WithShowDeletedOnParentAndChild(t *testing.T) {
	test := testUtils.TestCase{
		Actions: []any{
			testUtils.CreateDoc{
				CollectionID: 1, // Author
				Doc: `{
					"name": "John Grisham",
					"age": 65
				}`,
			},
			testUtils.CreateDoc{
				CollectionID: 1, // Author
				Doc: `{
					"name": "Cornelia Funke",
					"age": 55
				}`,
			},
			testUtils.CreateDoc{
				CollectionID: 0, // Book
				DocMap: map[string]any{
					"name":      "Painted House",
					"rating":    4.9,
					"author_id": testUtils.NewDocIndex(1, 0),
				},
			},
			testUtils.CreateDoc{
				CollectionID: 0, // Book
				DocMap: map[string]any{
					"name":      "A Time for Mercy",
					"rating":    4.5,
					"author_id": testUtils.NewDocIndex(1, 0),
				},
			},
			testUtils.CreateDoc{
				CollectionID: 0, // Book
				DocMap: map[string]any{
					"name":      "Theif Lord",
					"rating":    4.8,
					"author_id": testUtils.NewDocIndex(1, 1),
				},
			},
			// Delete the first author and first book
			testUtils.DeleteDoc{
				CollectionID: 1,
				DocID:        0,
			},
			testUtils.DeleteDoc{
				CollectionID: 0,
				DocID:        0,
			},
			// Query with showDeleted: true on parent, showDeleted: false on nested
			// Should show deleted author but not deleted books
			testUtils.Request{
				Request: `query {
					Author(showDeleted: true) {
						_deleted
						name
						age
						published(showDeleted: false) {
							_deleted
							name
							rating
						}
					}
				}`,
				Results: map[string]any{
					"Author": []map[string]any{
						{
							"_deleted": true,
							"name":     "John Grisham",
							"age":      int64(65),
							"published": []map[string]any{
								{
									"_deleted": false,
									"name":     "A Time for Mercy",
									"rating":   4.5,
								},
							},
						},
						{
							"_deleted": false,
							"name":     "Cornelia Funke",
							"age":      int64(55),
							"published": []map[string]any{
								{
									"_deleted": false,
									"name":     "Theif Lord",
									"rating":   4.8,
								},
							},
						},
					},
				},
				NonOrderedResults: true,
			},
			// Query with showDeleted: true on both parent and nested
			// Should show all deleted documents
			testUtils.Request{
				Request: `query {
					Author(showDeleted: true) {
						_deleted
						name
						age
						published(showDeleted: true) {
							_deleted
							name
							rating
						}
					}
				}`,
				Results: map[string]any{
					"Author": []map[string]any{
						{
							"_deleted": true,
							"name":     "John Grisham",
							"age":      int64(65),
							"published": []map[string]any{
								{
									"_deleted": true,
									"name":     "Painted House",
									"rating":   4.9,
								},
								{
									"_deleted": false,
									"name":     "A Time for Mercy",
									"rating":   4.5,
								},
							},
						},
						{
							"_deleted": false,
							"name":     "Cornelia Funke",
							"age":      int64(55),
							"published": []map[string]any{
								{
									"_deleted": false,
									"name":     "Theif Lord",
									"rating":   4.8,
								},
							},
						},
					},
				},
				NonOrderedResults: true,
			},
		},
	}

	executeTestCase(t, test)
}
