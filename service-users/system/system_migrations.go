package system

func Migrations() error {
	var err error
	// Create the users table
	_, err = Db.Exec(`
        create table if not exists users (
            id text primary key not null,
            created datetime not null default current_timestamp,
            updated datetime not null default current_timestamp,
            deleted datetime not null default '2400-01-01 00:00:00',
            email text unique not null,
            role int not null,
            sub text unique not null,
            avatar text not null default '',
            subscription_id text not null default '',
            subscription_end datetime not null default '2000-01-01 00:00:00',
            subscription_check datetime not null default '2000-01-01 00:00:00'
        )`)
	if err != nil {
		return err
	}

	// Create tokens table
	_, err = Db.Exec(`
        create table if not exists tokens (
            id text primary key not null,
            created datetime not null default current_timestamp,
            updated datetime not null default current_timestamp,
            deleted datetime not null default '2400-01-01 00:00:00',
            user_id text not null,
            provider text not null,
            access_token text not null,
            refresh_token text not null,
            token_type text not null,
            expires timestamp not null
        )`)
	if err != nil {
		return err
	}
	return nil
}
