package bench

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

// MultiDBSetBits sets bits with increasing profile id and bitmap id.
type MultiDBSetBits struct {
	HasClient
	BaseBitmapID  int `json:"base-bitmap-id"`
	BaseProfileID int `json:"base-profile-id"`
	Iterations    int `json:"iterations"`
}

func (b *MultiDBSetBits) Usage() string {
	return `
multi-db-set-bits sets bits with increasing profile id and bitmap id using a different DB for each agent.

Usage: multi-db-set-bits [arguments]

The following arguments are available:

	-base-bitmap-id int
		bits being set will all be greater than base-bitmap-id

	-base-profile-id int
		profile id num to start from

	-iterations int
		number of bits to set

	-client-type string
		Can be 'single' (all agents hitting one host) or 'round_robin'

`[1:]
}

func (b *MultiDBSetBits) ConsumeFlags(args []string) ([]string, error) {
	fs := flag.NewFlagSet("MultiDBSetBits", flag.ContinueOnError)
	fs.SetOutput(ioutil.Discard)
	fs.IntVar(&b.BaseBitmapID, "base-bitmap-id", 0, "")
	fs.IntVar(&b.BaseProfileID, "base-profile-id", 0, "")
	fs.IntVar(&b.Iterations, "iterations", 100, "")
	fs.StringVar(&b.ClientType, "client-type", "single", "")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return fs.Args(), nil
}

// Run runs the MultiDBSetBits benchmark
func (b *MultiDBSetBits) Run(ctx context.Context, agentNum int) map[string]interface{} {
	results := make(map[string]interface{})
	if b.cli == nil {
		results["error"] = fmt.Errorf("No client set for MultiDBSetBits agent: %v", agentNum)
		return results
	}
	s := NewStats()
	var start time.Time
	for n := 0; n < b.Iterations; n++ {
		query := fmt.Sprintf("SetBit(%d, 'frame.n', %d)", b.BaseBitmapID+n, b.BaseProfileID+n)
		start = time.Now()
		_, err := b.cli.ExecuteQuery(ctx, "multidb"+strconv.Itoa(agentNum), query, true)
		if err != nil {
			results["error"] = err
			return results
		}
		s.Add(time.Now().Sub(start))
	}
	AddToResults(s, results)
	return results
}