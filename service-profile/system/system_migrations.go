package system

func (s Storage) Migrations() error {
	var err error
	// Create profile table
	_, err = s.Conn.Exec(`
        create table if not exists profiles (
            id text primary key not null,
            created datetime not null default current_timestamp,
            updated datetime not null default current_timestamp,
            user_id text unique not null,
            active boolean not null,
            username text not null,
            about text not null,
            first_name text not null,
            last_name text not null,
            email text not null,
            country text not null,
            street_address text not null,
            city text not null,
            state text not null,
            zip text not null,
            email_notifications text not null,
            push_notification text not null,
            resume text not null,
            cover text not null,
            position text not null,
            skills text not null
        )`)
	if err != nil {
		return err
	}
	// Index user_id
	_, err = s.Conn.Exec(`create index if not exists profile_user_id on profile (user_id)`)
	if err != nil {
		return err
	}
	return nil
}
