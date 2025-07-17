package dto

type FriendRequestInput struct {
    FromUserID int `json:"from_user_id"`
    ToUserID   int `json:"to_user_id"`
}

type FriendActionInput struct {
	Action      string `json:"action"`
	UserID      int    `json:"user_id"`
	FromUserID  int    `json:"from_user_id"`
}

type BlockFriendInput struct {
	UserID   int `json:"user_id"`    
	FriendID int `json:"friend_id"`  
}