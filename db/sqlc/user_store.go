package db

import "context"

func (store *Store) RemoveUser(ctx context.Context, userId int64) error {
	return store.transaction(ctx, func(q *Queries) error {
		err := q.RemoveUserCredential(ctx, userId)
		if err != nil {
			return err
		}

		err = q.RemoveUser(ctx, userId)
		if err != nil {
			return err
		}

		return nil
	})
}
