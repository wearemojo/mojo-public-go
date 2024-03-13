package discourse

// result types are the root response types for most endpoints

// name as `FooResult` if solely wrapping 1 `Foo` field called `foo`
// more complex result types may need to be endpoint-specific

type UserResult struct {
	User User `json:"user"`
}

type CategoryListResult struct {
	CategoryList CategoryList `json:"category_list"`
}

type PostStreamResult struct {
	PostStream PostList `json:"post_stream"`
}

type PostIDsResult struct {
	PostIDs []int `json:"post_ids"`
}
