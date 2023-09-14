package gamm

type LiquidityChangeType int

const (
	AddLiquidity LiquidityChangeType = iota
	RemoveLiquidity
)

// LiquidityChangeMsg defines a simple interface for determining if an LP msg
// is removing or adding liquidity.
type LiquidityChangeMsg interface {
	LiquidityChangeType() LiquidityChangeType
}

var (
	_ LiquidityChangeMsg = MsgExitPool{}
	_ LiquidityChangeMsg = MsgExitSwapShareAmountIn{}
	_ LiquidityChangeMsg = MsgExitSwapExternAmountOut{}
)

var (
	_ LiquidityChangeMsg = MsgJoinPool{}
	_ LiquidityChangeMsg = MsgJoinSwapExternAmountIn{}
	_ LiquidityChangeMsg = MsgJoinSwapShareAmountOut{}
)

func (MsgExitPool) LiquidityChangeType() LiquidityChangeType {
	return RemoveLiquidity
}

func (MsgExitSwapShareAmountIn) LiquidityChangeType() LiquidityChangeType {
	return RemoveLiquidity
}

func (MsgExitSwapExternAmountOut) LiquidityChangeType() LiquidityChangeType {
	return RemoveLiquidity
}

func (MsgJoinPool) LiquidityChangeType() LiquidityChangeType {
	return AddLiquidity
}

func (MsgJoinSwapExternAmountIn) LiquidityChangeType() LiquidityChangeType {
	return AddLiquidity
}

func (MsgJoinSwapShareAmountOut) LiquidityChangeType() LiquidityChangeType {
	return AddLiquidity
}
