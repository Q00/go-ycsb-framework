package tikv

import (
	"fmt"

	"github.com/magiconair/properties"
	"github.com/q00/golang-mongo/pkg/ycsb"
	"github.com/tikv/client-go/config"
)

const (
	tikvPD = "tikv.pd"
	// raw, txn, or coprocessor
	tikvType      = "tikv.type"
	tikvConnCount = "tikv.conncount"
	tikvBatchSize = "tikv.batchsize"
)

type tikvCreator struct {
}

func (c tikvCreator) Create(p *properties.Properties) (ycsb.DB, error) {
	conf := config.Default()
	conf.RPC.MaxConnectionCount = p.GetUint(tikvConnCount, 128)
	conf.RPC.Batch.MaxBatchSize = p.GetUint(tikvBatchSize, 128)

	tp := p.GetString(tikvType, "raw")
	switch tp {
	case "raw":
		return createRawDB(p, conf)
	case "txn":
		return createTxnDB(p, conf)
	default:
		return nil, fmt.Errorf("unsupported type %s", tp)
	}
}

func init() {
	ycsb.RegisterDBCreator("tikv", tikvCreator{})
}
