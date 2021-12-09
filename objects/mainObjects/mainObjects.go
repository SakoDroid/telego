package objects

type Update struct {
	Update_id      int     `json:"update_id"`
	Message        Message `json:"message,omitempty"`
	Edited_message Message `json:"edited_message,omitempty"`
}

type Message struct {
}

type User struct {
	Id                      int    `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	Lastname                string `json:"last_name"`
	Username                string `json:"username"`
	LanguageCode            string `json:"language_code,omitempty"`
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
	SuportsInlineQueries    bool   `json:"supports_inline_queries,omitempty"`
}

type Chat struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}
