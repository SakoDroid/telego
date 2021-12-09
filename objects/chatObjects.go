package objects

type Chat struct {
	Id                   int             `json:"id"`
	Type                 string          `json:"type"`
	Title                string          `json:"title,omitempty"`
	Username             string          `json:"username,omitempty"`
	FirstName            string          `json:"first_name,omitempty"`
	LastName             string          `json:"last_name,omitempty"`
	Photo                ChatPhoto       `json:"photo,omitempty"`
	Bio                  string          `json:"bio,omitempty"`
	HasPrivateForwards   bool            `json:"has_private_forwards,omitempty"`
	Description          string          `json:"description,omitempty"`
	InviteLink           string          `json:"invite_link,omitempty"`
	PinnedMessage        Message         `json:"pinned_message,omitempty"`
	Permissions          ChatPermissions `json:"permissions,omitempty"`
	SlowModeDelay        int             `json:"slow_mode_delay,omitempty"`
	MessageAutoDeletTime int             `json:"message_auto_delete_time,omitempty"`
	HasProtectedContent  bool            `json:"has_protected_content,omitempty"`
	StickerSetName       string          `json:"sticker_set_name,omitempty"`
	CanSetStickerSet     bool            `json:"can_set_sticker_set,omitempty"`
	LinkedChatId         int             `json:"linked_chat_id,omitempty"`
	Location             ChatLocation    `json:"location,omitempty"`
}

type ChatPhoto struct {
	SmallFileId       string `json:"small_file_id"`
	SmallFileUniqueId string `json:"small_file_unique_id"`
	BigFileId         string `json:"big_file_id"`
	BigFileUniqueId   string `json:"big_file_unique_id"`
}

type ChatInviteLink struct {
	InviteLink              string `json:"invite_link"`
	Creator                 User   `json:"user"`
	CreatesJoinRequest      bool   `json:"creates_join_request"`
	IsPrimary               bool   `json:"is_primary"`
	IsRevoked               bool   `json:"is_revoked"`
	Name                    string `json:"name,omitempty"`
	ExpireDate              int    `json:"expire_date,omitempty"`
	MemberLimit             int    `json:"member_limit,omitempty"`
	PendingJoinRequestCount int    `json:"pending_join_request_count,omitempty"`
}

type ChatLocation struct {
}

type ChatPermissions struct {
}
