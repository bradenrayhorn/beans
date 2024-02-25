package fake

import "github.com/bradenrayhorn/beans/server/beans"

type repository struct{ *database }

func (r repository) txOrNow(tx beans.Tx, do func()) {
	if tx == nil {
		do()
	} else {
		t := tx.(*fTx)
		t.do = append(t.do, do)
	}
}

func (r repository) acquire(do func()) {
	r.database.mu.Lock()
	defer r.database.mu.Unlock()
	do()
}
