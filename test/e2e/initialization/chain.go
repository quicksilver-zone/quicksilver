package initialization

const (
	keyringPassphrase = "testpassphrase"
	keyringAppName    = "testnet"
)

// internalChain contains the same info as chain, but with the validator structs instead using the internal validator
// representation, with more derived data
type internalChain struct {
	chainMeta ChainMeta
	nodes     []*internalNode
}

func new(id, dataDir string) (*internalChain, error) {
	chainMeta := ChainMeta{
		Id:      id,
		DataDir: dataDir,
	}
	return &internalChain{
		chainMeta: chainMeta,
	}, nil
}

func (c *internalChain) export() (*Chain, error) {
	exportNodes := make([]*Node, 0, len(c.nodes))
	for _, v := range c.nodes {
		node, err := v.export()
		if err != nil {
			return nil, err
		}
		exportNodes = append(exportNodes, node)
	}

	return &Chain{
		ChainMeta: c.chainMeta,
		Nodes:     exportNodes,
	}, nil
}
