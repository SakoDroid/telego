package objects

/*This object represents a forum topic.*/
type ForumTopic struct {
	//Unique identifier of the forum topic
	MessageThreadId int `json:"message_thread_id"`
	//Name of the topic
	Name string `json:"name"`
	//Color of the topic icon in RGB format
	IconColor int `json:"icon_color"`
	//Optional. Unique identifier of the custom emoji shown as the topic icon
	IconCustomEmojiId string `json:"icon_custom_emoji_id,omitempty"`
}

/*This object represents a service message about a new forum topic created in the chat.*/
type ForumTopicCreated struct {
	/*Name of the topic*/
	Name string `json:"name"`
	/*Color of the topic icon in RGB format*/
	IconColor int `json:"icon_color"`
	/*Optional. Unique identifier of the custom emoji shown as the topic icon*/
	IconCustomEmojiId string `json:"icon_custom_emoji_id"`
}

/*This object represents a service message about a forum topic closed in the chat. Currently holds no information.*/
type ForumTopicClosed struct{}

/*This object represents a service message about an edited forum topic.*/
type ForumTopicEdited struct {
	/*Optional. New name of the topic, if it was edited*/
	Name string `json:"name"`
	/*Optional. New identifier of the custom emoji shown as the topic icon, if it was edited; an empty string if the icon was removed*/
	IconCustomEmojiId string `json:"icon_custom_emoji_id"`
}

/*This object represents a service message about a forum topic reopened in the chat. Currently holds no information.*/
type ForumTopicReopened struct{}

/*This object represents a service message about General forum topic hidden in the chat. Currently holds no information.*/
type GeneralForumTopicHidden struct{}

/*This object represents a service message about General forum topic unhidden in the chat. Currently holds no information.*/
type GeneralForumTopicUnhidden struct{}
