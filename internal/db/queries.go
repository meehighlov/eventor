package db

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
)

// idempotent save
// accepts ALL fields of entity and save as is
func (user *User) Save(ctx context.Context) error {
	_, _, _ = user.RefresTimestamps()

	_, err := sqliteConn.ExecContext(
		ctx,
		`INSERT INTO user(id, tgid, name, tgusername, chatid, createdat, updatedat)
        VALUES($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT(tgid) DO UPDATE SET name=$3, tgusername=$4, chatid=$5, updatedat=$6
        RETURNING id;`,
		&user.ID,
		&user.TGId,
		&user.Name,
		&user.TGusername,
		&user.ChatId,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		slog.Error("Error when trying to save user: " + err.Error())
		return err
	}
	slog.Debug("User created/updated")

	return nil
}

func (user *User) Filter(ctx context.Context) ([]User, error) {
	where := []string{}
	if user.TGId != 0 {
		where = append(where, "tgid=$tgid")
	}
	if user.TGusername != "" {
		where = append(where, "tgusername=$tgusername")
	}

	where_ := strings.Join(where, " AND ")
	query := `SELECT id, tgid, name, tgusername, chatid, createdat, updatedat FROM user WHERE ` + where_ + `;`

	rows, err := sqliteConn.QueryContext(
		ctx,
		query,
		sql.Named("tgid", user.TGId),
		sql.Named("tgusername", user.TGusername),
	)
	if err != nil {
		slog.Error("Error when filtering users " + err.Error())
		return nil, err
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		user := User{}
		err := rows.Scan(
			&user.ID,
			&user.TGId,
			&user.Name,
			&user.TGusername,
			&user.ChatId,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			slog.Error("Error fetching users by filter params: " + err.Error())
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

// idempotent save
// accepts ALL fields of entity and save as is
func (e Event) Save(ctx context.Context) error {
	_, _, _ = e.RefresTimestamps()

	_, err := sqliteConn.ExecContext(
		ctx,
		`INSERT INTO event(id, chatid, ownerid, text, notifyat, delta, createdat, updatedat)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT(id) DO UPDATE SET chatid=$2, ownerid=$3, text=$4, notifyat=$5, delta=$6, updatedat=$8
        RETURNING id;`,
		&e.ID,
		&e.ChatId,
		&e.OwnerId,
		&e.Text,
		&e.NotifyAt,
		&e.Delta,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		slog.Error("Error when trying to save event: " + err.Error())
		return err
	}
	slog.Debug("Event created/updated")

	return nil
}

func (event Event) Filter(ctx context.Context) ([]Entity, error) {
	where := []string{}

	if event.OwnerId != 0 {
		where = append(where, "ownerid=$ownerid")
	}
	if event.ID != "" {
		where = append(where, "id=$id")
	}

	where_ := strings.Join(where, " AND ")

	query := `SELECT id, chatid, ownerid, text, notifyat, delta, createdat, updatedat FROM event`

	if len(where) != 0 {
		query += ` WHERE ` + where_
	}

	query += `;`

	rows, err := sqliteConn.QueryContext(
		ctx,
		query,
		sql.Named("ownerid", event.OwnerId),
		sql.Named("id", event.ID),
	)
	if err != nil {
		slog.Error("Error when filtering events " + err.Error())
		return nil, err
	}
	defer rows.Close()

	events := []Entity{}

	for rows.Next() {
		event := Event{}
		err := rows.Scan(
			&event.ID,
			&event.ChatId,
			&event.OwnerId,
			&event.Text,
			&event.NotifyAt,
			&event.Delta,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			slog.Error("Error fetching events by filter params: " + err.Error())
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func (event Event) Delete(ctx context.Context) error {
	_, err := sqliteConn.ExecContext(
		ctx,
		`DELETE FROM event WHERE id = $1;`,
		&event.ID,
	)
	if err != nil {
		slog.Error("Error when trying to delete event row: " + err.Error())
		return err
	}

	slog.Debug("Event row deleted")

	return nil
}
