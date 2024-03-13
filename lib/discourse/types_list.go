package discourse

// list types wrap arrays

// named as `FooList`, containing a `[]Foo` field called `foos`
// may contain additional fields for pagination or other metadata

type CategoryList struct {
	Categories []Category `json:"categories"`
}

type PostList struct {
	Posts []Post `json:"posts"`
}
