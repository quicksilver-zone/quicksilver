package types

type ProtocolDataI interface {
	ValidateBasic() error
}
