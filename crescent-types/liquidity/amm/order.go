package amm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ Order   = (*BaseOrder)(nil)
	_ Orderer = (*BaseOrderer)(nil)

	DefaultOrderer = BaseOrderer{}
)

// OrderDirection specifies an order direction, either buy or sell.
type OrderDirection int

// OrderDirection enumerations.
const (
	Buy OrderDirection = iota + 1
	Sell
)

func (dir OrderDirection) String() string {
	switch dir {
	case Buy:
		return "Buy"
	case Sell:
		return "Sell"
	default:
		return fmt.Sprintf("OrderDirection(%d)", dir)
	}
}

type Orderer interface {
	Order(dir OrderDirection, price sdk.Dec, amt sdk.Int) Order
}

// BaseOrderer creates new BaseOrder with sufficient offer coin amount
// considering price and amount.
type BaseOrderer struct{}

func (orderer BaseOrderer) Order(dir OrderDirection, price sdk.Dec, amt sdk.Int) Order {
	return NewBaseOrder(dir, price, amt, OfferCoinAmount(dir, price, amt))
}

// Order is the universal interface of an order.
type Order interface {
	GetDirection() OrderDirection
	// GetBatchId returns the batch id where the order was created.
	// Batch id of 0 means the current batch.
	GetBatchId() uint64
	GetPrice() sdk.Dec
	GetAmount() sdk.Int // The original order amount
	GetOfferCoinAmount() sdk.Int
	GetPaidOfferCoinAmount() sdk.Int
	SetPaidOfferCoinAmount(amt sdk.Int)
	GetReceivedDemandCoinAmount() sdk.Int
	SetReceivedDemandCoinAmount(amt sdk.Int)
	GetOpenAmount() sdk.Int
	SetOpenAmount(amt sdk.Int)
	IsMatched() bool
	// HasPriority returns true if the order has higher priority
	// than the other order.
	HasPriority(other Order) bool
	String() string
}

// BaseOrder is the base struct for an Order.
type BaseOrder struct {
	Direction       OrderDirection
	Price           sdk.Dec
	Amount          sdk.Int
	OfferCoinAmount sdk.Int

	// Match info
	OpenAmount               sdk.Int
	PaidOfferCoinAmount      sdk.Int
	ReceivedDemandCoinAmount sdk.Int
}

// NewBaseOrder returns a new BaseOrder.
func NewBaseOrder(dir OrderDirection, price sdk.Dec, amt, offerCoinAmt sdk.Int) *BaseOrder {
	return &BaseOrder{
		Direction:                dir,
		Price:                    price,
		Amount:                   amt,
		OfferCoinAmount:          offerCoinAmt,
		OpenAmount:               amt,
		PaidOfferCoinAmount:      sdk.ZeroInt(),
		ReceivedDemandCoinAmount: sdk.ZeroInt(),
	}
}

// GetDirection returns the order direction.
func (order *BaseOrder) GetDirection() OrderDirection {
	return order.Direction
}

func (order *BaseOrder) GetBatchId() uint64 {
	return 0
}

// GetPrice returns the order price.
func (order *BaseOrder) GetPrice() sdk.Dec {
	return order.Price
}

// GetAmount returns the order amount.
func (order *BaseOrder) GetAmount() sdk.Int {
	return order.Amount
}

func (order *BaseOrder) GetOfferCoinAmount() sdk.Int {
	return order.OfferCoinAmount
}

func (order *BaseOrder) GetPaidOfferCoinAmount() sdk.Int {
	return order.PaidOfferCoinAmount
}

func (order *BaseOrder) SetPaidOfferCoinAmount(amt sdk.Int) {
	order.PaidOfferCoinAmount = amt
}

func (order *BaseOrder) GetReceivedDemandCoinAmount() sdk.Int {
	return order.ReceivedDemandCoinAmount
}

func (order *BaseOrder) SetReceivedDemandCoinAmount(amt sdk.Int) {
	order.ReceivedDemandCoinAmount = amt
}

func (order *BaseOrder) GetOpenAmount() sdk.Int {
	return order.OpenAmount
}

func (order *BaseOrder) SetOpenAmount(amt sdk.Int) {
	order.OpenAmount = amt
}

func (order *BaseOrder) IsMatched() bool {
	return order.OpenAmount.LT(order.Amount)
}

// HasPriority returns whether the order has higher priority than
// the other order.
func (order *BaseOrder) HasPriority(other Order) bool {
	return order.Amount.GT(other.GetAmount())
}

func (order *BaseOrder) String() string {
	return fmt.Sprintf("BaseOrder(%s,%s,%s)", order.Direction, order.Price, order.Amount)
}
