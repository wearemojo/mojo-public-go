package discourse

import (
	"time"
)

// these are the actual types

// name as `FooNew` if it's a type that's used to create a new `Foo`

type User struct {
	ID                 int                     `json:"id"`
	Username           string                  `json:"username"`
	AvatarTemplate     *string                 `json:"avatar_template"`
	SingleSignOnRecord *UserSingleSignOnRecord `json:"single_sign_on_record"`
}

type UserSingleSignOnRecord struct {
	ExternalID string `json:"external_id"`
}

type Category struct {
	ID              int        `json:"id"`
	Name            string     `json:"name"`
	Slug            string     `json:"slug"`
	ReadRestricted  bool       `json:"read_restricted"`
	SubcategoryIDs  []int      `json:"subcategory_ids"`
	SubcategoryList []Category `json:"subcategory_list"`
}

type Topic struct {
	ID                int          `json:"id"`
	Title             string       `json:"title"`
	FancyTitle        string       `json:"fancy_title"`
	Slug              string       `json:"slug"`
	PostsCount        int          `json:"posts_count"`
	ReplyCount        int          `json:"reply_count"`
	HighestPostNumber int          `json:"highest_post_number"`
	ImageURL          *string      `json:"image_url"`
	CreatedAt         time.Time    `json:"created_at"`
	LastPostedAt      *time.Time   `json:"last_posted_at"`
	Archetype         Archetype    `json:"archetype"`
	Pinned            bool         `json:"pinned"`
	Unpinned          *time.Time   `json:"unpinned"`
	Visible           bool         `json:"visible"`
	Closed            bool         `json:"closed"`
	Archived          bool         `json:"archived"`
	Tags              []string     `json:"tags"`
	CategoryID        int          `json:"category_id"`
	PinnedGlobally    bool         `json:"pinned_globally"`
	FeaturedLink      *string      `json:"featured_link"`
	Details           TopicDetails `json:"details"`
}

type TopicDetails struct {
	CreatedBy TopicCreatedBy `json:"created_by"`
}

type TopicCreatedBy struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	AvatarTemplate string `json:"avatar_template"`
}

type PostNew struct {
	Title             *string `json:"title"`
	Raw               string  `json:"raw"`
	TopicID           *int    `json:"topic_id"`
	Category          *int    `json:"category"`
	ReplyToPostNumber *int    `json:"reply_to_post_number"`
}

type Post struct {
	ID                int        `json:"id"`
	Username          string     `json:"username"`
	AvatarTemplate    string     `json:"avatar_template"`
	CreatedAt         time.Time  `json:"created_at"`
	Raw               string     `json:"raw"`
	Cooked            string     `json:"cooked"`
	PostNumber        int        `json:"post_number"`
	PostType          PostType   `json:"post_type"`
	UpdatedAt         time.Time  `json:"updated_at"`
	ReplyToPostNumber *int       `json:"reply_to_post_number"`
	TopicID           int        `json:"topic_id"`
	Version           int        `json:"version"`
	UserID            int        `json:"user_id"`
	Hidden            bool       `json:"hidden"`
	DeletedAt         *time.Time `json:"deleted_at"`
	Wiki              bool       `json:"wiki"`

	// Discourse Reactions plugin
	Reactions           []PluginDCReactionsPostReaction           `json:"reactions"`
	CurrentUserReaction *PluginDCReactionsPostReactionCurrentUser `json:"current_user_reaction"`
}
