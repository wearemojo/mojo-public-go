package discourse

type Archetype string

const (
	// https://github.com/discourse/discourse/blob/8cf2f909f5053b13d8f6a79703aaf7fbdb3f6423/lib/archetype.rb

	ArchetypeRegular        Archetype = "regular"
	ArchetypePrivateMessage Archetype = "private_message"
	ArchetypeBanner         Archetype = "banner"
)

type PostType int

const (
	// https://github.com/discourse/discourse/blob/ce3f592275295f201f1332c0c5069897341d5f47/app/models/post.rb#L172

	PostTypeRegular         PostType = 1
	PostTypeModeratorAction PostType = 2
	PostTypeSmallAction     PostType = 3
	PostTypeWhisper         PostType = 4
)
