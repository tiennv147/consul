// +build !consulent

package state

import (
	"fmt"
	"strings"

	memdb "github.com/hashicorp/go-memdb"

	"github.com/hashicorp/consul/agent/structs"
)

func indexFromConfigEntryKindName(arg interface{}) ([]byte, error) {
	n, ok := arg.(ConfigEntryKindName)
	if !ok {
		return nil, fmt.Errorf("invalid type for ConfigEntryKindName query: %T", arg)
	}

	var b indexBuilder
	b.String(strings.ToLower(n.Kind))
	b.String(strings.ToLower(n.Name))
	return b.Bytes(), nil
}

func indexFromConfigEntry(raw interface{}) ([]byte, error) {
	c, ok := raw.(structs.ConfigEntry)
	if !ok {
		return nil, fmt.Errorf("type must be structs.ConfigEntry: %T", raw)
	}

	if c.GetName() == "" || c.GetKind() == "" {
		return nil, errMissingValueForIndex
	}

	var b indexBuilder
	b.String(strings.ToLower(c.GetKind()))
	b.String(strings.ToLower(c.GetName()))
	return b.Bytes(), nil
}

func validateConfigEntryEnterprise(_ ReadTxn, _ structs.ConfigEntry) error {
	return nil
}

func getAllConfigEntriesWithTxn(tx ReadTxn, _ *structs.EnterpriseMeta) (memdb.ResultIterator, error) {
	return tx.Get(tableConfigEntries, "id")
}

func getConfigEntryKindsWithTxn(tx ReadTxn, kind string, _ *structs.EnterpriseMeta) (memdb.ResultIterator, error) {
	return tx.Get(tableConfigEntries, "kind", kind)
}

func configIntentionsConvertToList(iter memdb.ResultIterator, _ *structs.EnterpriseMeta) structs.Intentions {
	var results structs.Intentions
	for v := iter.Next(); v != nil; v = iter.Next() {
		entry := v.(*structs.ServiceIntentionsConfigEntry)
		for _, src := range entry.Sources {
			results = append(results, entry.ToIntention(src))
		}
	}
	return results
}
