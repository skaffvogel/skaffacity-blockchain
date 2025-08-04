package types

type BankKeeper interface{}

type bankKeeperStub struct{}

// StoreKeyType is a stub for sdk.StoreKey
type StoreKeyType struct{}

// Stub methods for BankKeeper
func (b *bankKeeperStub) SendCoinsFromAccountToModule() {}
func (b *bankKeeperStub) SendCoinsFromModuleToAccount() {}
