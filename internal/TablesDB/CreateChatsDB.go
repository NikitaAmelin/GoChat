package TablesDB

import "context"

func (r Repository) CreateChatsTable(ctx context.Context) error {
	q := `
CREATE TABLE IF NOT EXISTS public."Chats"
(
    "Name" character varying COLLATE pg_catalog."default" NOT NULL,
    "Messeges" text[],
    "NameMessegesDB" character varying COLLATE pg_catalog."default" NOT NULL,
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    CONSTRAINT "Chats_pkey" PRIMARY KEY (id)
)

TABLESPACE pg_default;

`

	_, err := r.Client.Exec(ctx, q)
	if err != nil {
		return err
	}
	return nil
}
