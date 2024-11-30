package TablesDB

import (
	"context"
	postgresql "goydamess/pkg/data_base"
)

type Repository struct {
	Client postgresql.Client
}

func (r Repository) CreateUsersTable(ctx context.Context) error {
	q := `
CREATE TABLE IF NOT EXISTS public."Users"
(
    "Login" character varying COLLATE pg_catalog."default" NOT NULL,
    "Password" character varying COLLATE pg_catalog."default" NOT NULL,
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    CONSTRAINT "Users_pkey" PRIMARY KEY (id)
)

TABLESPACE pg_default;

`

	_, err := r.Client.Exec(ctx, q)
	if err != nil {
		return err
	}
	return nil
}
