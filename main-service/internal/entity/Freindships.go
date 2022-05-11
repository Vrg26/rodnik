package entity

import "github.com/google/uuid"

type Freindships struct {
	FriendTo   uuid.UUID `json:"friend_to,omitempty"`
	FriendFrom uuid.UUID `json:"friend_from,omitempty"`
}
