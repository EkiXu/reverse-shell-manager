package typing

type HubInterface interface {
	Construct()
	Run()
	BroadcastWSData(WSData) error
}
