package sqldb

import (
	"context"
)

// PingContext function
func (db *DB) PingContext(ctx context.Context) error {
	errCh := make(chan error, 2)

	go func() {
		errCh <- db.master.PingContext(ctx)
	}()

	go func() {
		errCh <- db.follower.PingContext(ctx)
	}()

	for i := 0; i < 2; i++ {
		err := <-errCh
		if err != nil {
			return err
		}
	}

	return nil
}
