package TablesDB

import (
	"context"
	"fmt"
)

func (r Repository) CreateMessegesTable(ctx context.Context, name string) error {
	q := `
CREATE TABLE IF NOT EXISTS public.$1
(
    "Author" character varying COLLATE pg_catalog."default" NOT NULL,
    "Text" character varying COLLATE pg_catalog."default" NOT NULL,
    "TimeOfSending" character varying COLLATE pg_catalog."default" NOT NULL,
    "Viewed" text[],
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    CONSTRAINT $2 PRIMARY KEY (id)
)

TABLESPACE pg_default;

`

	_, err := r.Client.Exec(ctx, q, name, fmt.Sprintf("%s_pkay", name))
	if err != nil {
		return err
	}
	return nil
}
