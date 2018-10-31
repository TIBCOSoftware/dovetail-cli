package contract

//Generator generates smart contract artifacts
type Generator interface {
	Generate() error
}
