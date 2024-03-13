package discourse

type PostType int

const (
	// https://github.com/discourse/discourse/blob/ce3f592275295f201f1332c0c5069897341d5f47/app/models/post.rb#L172

	PostTypeRegular         PostType = 1
	PostTypeModeratorAction PostType = 2
	PostTypeSmallAction     PostType = 3
	PostTypeWhisper         PostType = 4
)
