package types

import (
    "time"
)

const (
    TypeLand     = "land"
    TypeItem     = "item"
    TypeBadge    = "badge"
    TypeAttachment = "attachment"
)

// NFT represents a non-fungible token in the game
type NFT struct {
    ID          string    `json:"id"`
    Type        string    `json:"type"`
    Owner       string    `json:"owner"`
    Metadata    Metadata  `json:"metadata"`
    Created     time.Time `json:"created"`
    Transferable bool     `json:"transferable"`
}


// Metadata contains NFT-specific attributes
type Metadata struct {
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Image       string            `json:"image"`
    Properties  map[string]string `json:"properties"`
}

// LandMetadata contains land-specific properties
type LandMetadata struct {
    Metadata
    Location    Location `json:"location"`
    Size        uint32   `json:"size"`
    BuildRights bool     `json:"build_rights"`
}

// Location represents coordinates in the game world
type Location struct {
    X int32 `json:"x"`
    Y int32 `json:"y"`
}

// ItemMetadata contains item-specific properties
type ItemMetadata struct {
    Metadata
    ItemType    string   `json:"item_type"`
    Rarity      string   `json:"rarity"`
    Season      uint32   `json:"season"`
    Attachments []string `json:"attachments"`
}

// BadgeMetadata contains badge-specific properties
type BadgeMetadata struct {
    Metadata
    Achievement string    `json:"achievement"`
    DateEarned  time.Time `json:"date_earned"`
    Permanent   bool      `json:"permanent"`
}
