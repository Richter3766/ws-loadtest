package wsclient

// MessageRequest : 서버에 publish할 메시지 구조체
type MessageRequest struct {
    MessageId string `json:"messageId"`
    SentAt int64 `json:"sentAt"`
    RoomId  int    `json:"roomId"`
    Content string `json:"content"`
}

type User struct {
    Email string `json:"email"`
    JWT   string `json:"jwt"`
    RefreshToken string `json:"refreshToken"`
}

type MessageResponseDto struct {
    RoomId           int64  `json:"roomId"`
    MessageId        int64  `json:"messageId"`
    ClientMessageId  string `json:"clientMessageId"`
    ClientSentAt     int64  `json:"clientSentAt"`

    ProfileImage     string `json:"profileImage"`
    Writer           string `json:"writer"`
    Position         string `json:"position"`
    Content          string `json:"content"`

    LikesNum         int64  `json:"likesNum"`
    ThreadNum        int64  `json:"threadNum"`

    CreatedTime      string `json:"createdTime"` // "yyyy-MM-dd HH:mm:ss" 형식의 문자열로 전달됨
}
