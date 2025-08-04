package types

// ProtoMessage implements the proto.Message interface for NFT.
func (n *NFT) ProtoMessage() {}

// Reset implements the proto.Message interface for NFT.
func (n *NFT) Reset() {}

// String implements the fmt.Stringer interface for NFT.
func (n *NFT) String() string { return "NFT stub" }

// MarshalTo implements codec.ProtoMarshaler for NFT (stub).
func (n *NFT) MarshalTo(data []byte) (int, error) { return 0, nil }

// Marshal implements codec.ProtoMarshaler for NFT (stub).
func (n *NFT) Marshal() ([]byte, error) { return nil, nil }

// Unmarshal implements codec.ProtoMarshaler for NFT (stub).
func (n *NFT) Unmarshal([]byte) error { return nil }

// Size implements codec.ProtoMarshaler for NFT (stub).
func (n *NFT) Size() int { return 0 }

// MarshalToSizedBuffer implements codec.ProtoMarshaler for NFT (stub).
func (n *NFT) MarshalToSizedBuffer(data []byte) (int, error) { return 0, nil }