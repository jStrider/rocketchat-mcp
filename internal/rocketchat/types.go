package rocketchat

// Channel represents a Rocket.Chat channel or group.
type Channel struct {
	ID           string   `json:"_id"`
	Name         string   `json:"name"`
	Fname        string   `json:"fname"`
	Type         string   `json:"t"`
	MembersCount int      `json:"usersCount"`
	MsgCount     int      `json:"msgs"`
	Topic        string   `json:"topic"`
	Description  string   `json:"description"`
	ReadOnly     bool     `json:"ro"`
	Archived     bool     `json:"archived"`
	UpdatedAt    string   `json:"_updatedAt"`
	LastMessage  *Message `json:"lastMessage,omitempty"`
}

// DisplayName returns the best available name for the channel.
func (c Channel) DisplayName() string {
	if c.Fname != "" {
		return c.Fname
	}
	return c.Name
}

// MessageUser is the user sub-object inside a message.
type MessageUser struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// Attachment represents a file or media attachment on a message.
type Attachment struct {
	Title     string `json:"title"`
	Type      string `json:"type"`
	TitleLink string `json:"title_link"`
}

// FileInfo represents a file reference on a message.
type FileInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Message represents a Rocket.Chat message.
type Message struct {
	ID          string            `json:"_id"`
	RoomID      string            `json:"rid"`
	Text        string            `json:"msg"`
	User        MessageUser       `json:"u"`
	Timestamp   string            `json:"ts"`
	EditedAt    string            `json:"editedAt,omitempty"`
	Reactions   map[string]any    `json:"reactions,omitempty"`
	ThreadID    string            `json:"tmid,omitempty"`
	ThreadCount int               `json:"tcount,omitempty"`
	Replies     []string          `json:"replies,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
	File        *FileInfo         `json:"file,omitempty"`
}

// Email represents an email entry on a user profile.
type Email struct {
	Address  string `json:"address"`
	Verified bool   `json:"verified"`
}

// User represents a Rocket.Chat user.
type User struct {
	ID       string   `json:"_id"`
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Status   string   `json:"status"`
	Emails   []Email  `json:"emails,omitempty"`
	Roles    []string `json:"roles"`
	Active   bool     `json:"active"`
}

// RoomFile represents a file uploaded to a room.
type RoomFile struct {
	ID       string      `json:"_id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Size     int64       `json:"size"`
	User     MessageUser `json:"user"`
	URL      string      `json:"url"`
	UploadAt string      `json:"uploadedAt"`
}

// ListOptions holds pagination and filter parameters.
type ListOptions struct {
	Count  int
	Offset int
	Query  string
}

// Pagination holds pagination metadata from list responses.
type Pagination struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// ChannelListResponse wraps the channels.list API response.
type ChannelListResponse struct {
	Channels []Channel `json:"channels"`
	Pagination
}

// MessageListResponse wraps message list API responses.
type MessageListResponse struct {
	Messages []Message `json:"messages"`
	Pagination
}

// UserListResponse wraps the users.list API response.
type UserListResponse struct {
	Users []User `json:"users"`
	Pagination
}

// FileListResponse wraps the channels.files API response.
type FileListResponse struct {
	Files []RoomFile `json:"files"`
	Pagination
}

// Room represents a DM or group room (used in im.create response).
type Room struct {
	ID string `json:"_id"`
}
