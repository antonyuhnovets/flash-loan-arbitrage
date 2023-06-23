package contract

type _WHEI interface {
	deposit()

	withdraw(uint)

	transfer(to string, value uint) bool
}

type event func(...interface{})

type Transfer event
type Approval event

type IERC20 interface {
	totalSupply() uint
	balanceOf(account string) uint
	transfer(recipient string, amount uint) bool
	allowance(owner string, spender string) uint
	approve(spender string, amount uint) bool
	transferFrom(sender string, recipient string, amount uint) bool

	Transfer(from string, to string, value uint)
	Approval(owner string, spender string, value uint)
}
