package types

import (
	"fmt"
	"github.com/cosmos/gogoproto/proto"
)

// WebConfig represents the configuration for the web interface
type WebConfig struct {
	// enabled indicates if the web interface is enabled
	Enabled bool `json:"enabled"`
	
	// port is the port number for the web interface
	Port uint32 `json:"port"`
	
	// host is the host address for the web interface
	Host string `json:"host"`
	
	// api_endpoint is the API endpoint URL
	ApiEndpoint string `json:"api_endpoint"`
	
	// websocket_endpoint is the WebSocket endpoint URL
	WebsocketEndpoint string `json:"websocket_endpoint"`
	
	// theme is the UI theme configuration
	Theme string `json:"theme"`
	
	// features are the enabled features
	Features []string `json:"features"`
	
	// fee_distribution configures the transaction fee distribution
	FeeDistribution FeeDistribution `json:"fee_distribution"`
}

// ProtoMessage implements proto.Message interface
func (wc *WebConfig) ProtoMessage() {}

// Reset implements proto.Message interface
func (wc *WebConfig) Reset() {
	*wc = WebConfig{}
}

// String implements proto.Message interface
func (wc *WebConfig) String() string {
	return fmt.Sprintf("WebConfig{Enabled: %t, Port: %d, Host: %s, ApiEndpoint: %s, WebsocketEndpoint: %s, Theme: %s, Features: %v}",
		wc.Enabled, wc.Port, wc.Host, wc.ApiEndpoint, wc.WebsocketEndpoint, wc.Theme, wc.Features)
}

// Marshal implements ProtoMarshaler interface
func (wc *WebConfig) Marshal() ([]byte, error) {
	return proto.Marshal(wc)
}

// Unmarshal implements ProtoMarshaler interface
func (wc *WebConfig) Unmarshal(data []byte) error {
	return proto.Unmarshal(data, wc)
}

// MarshalTo implements ProtoMarshaler interface
func (wc *WebConfig) MarshalTo(data []byte) (int, error) {
	marshaled, err := wc.Marshal()
	if err != nil {
		return 0, err
	}
	copy(data, marshaled)
	return len(marshaled), nil
}

// Size implements ProtoMarshaler interface
func (wc *WebConfig) Size() int {
	marshaled, _ := wc.Marshal()
	return len(marshaled)
}

// MarshalToSizedBuffer implements ProtoMarshaler interface
func (wc *WebConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	marshaled, err := wc.Marshal()
	if err != nil {
		return 0, err
	}
	copy(dAtA[len(dAtA)-len(marshaled):], marshaled)
	return len(marshaled), nil
}

// DefaultWebConfig returns the default web configuration
func DefaultWebConfig() WebConfig {
	return WebConfig{
		Enabled:           true,
		Port:              8090,
		Host:              "0.0.0.0",
		ApiEndpoint:       "http://localhost:1317",
		WebsocketEndpoint: "ws://localhost:26657/websocket",
		Theme:             "default",
		Features: []string{
			"tokenfactory",
			"nft",
			"marketplace",
			"governance",
			"staking",
		},
		FeeDistribution: DefaultFeeDistribution(),
	}
}
