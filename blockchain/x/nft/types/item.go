package types

// ProtoMessage implements the proto.Message interface for Item.
func (i *Item) ProtoMessage() {}

// Reset implements the proto.Message interface for Item.
func (i *Item) Reset() {}

// String implements the fmt.Stringer interface for Item.
func (i *Item) String() string { return "Item stub" }