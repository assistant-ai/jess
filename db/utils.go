package db

import (
	"database/sql"
	"log"
	"path/filepath"

	"github.com/assistant-ai/jess/model"
	"github.com/google/uuid"

	"github.com/b0noi/go-utils/v2/fs"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var err error
	folderPath, err := fs.MaybeCreateProgramFolder("assistent.ai")
	if err != nil {
		log.Fatal(err)
	}
	dbFilePath := filepath.Join(folderPath, "messages.db")
	db, err = sql.Open("sqlite3", dbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		dialog_id TEXT,
		timestamp DATETIME,
		role TEXT,
		content TEXT
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func RemoveMessagesByDialogId(dialogId string) error {
	statement := `DELETE FROM messages WHERE dialog_id = ?`

	result, err := db.Exec(statement, dialogId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("Removed %d messages with dialog_id %s", rowsAffected, dialogId)
	return nil
}

func StoreMessage(m model.Message) (string, error) {
	id := uuid.New().String()
	stmt, err := db.Prepare("INSERT INTO messages(id, dialog_id, timestamp, role, content) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return "", err
	}

	_, err = stmt.Exec(id, m.DialogId, m.Timestamp, m.Role, m.Content)
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetDialogIDs() ([]string, error) {
	rows, err := db.Query("SELECT DISTINCT dialog_id FROM messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dialogIDs []string
	for rows.Next() {
		var dialogID string
		err := rows.Scan(&dialogID)
		if err != nil {
			return nil, err
		}
		dialogIDs = append(dialogIDs, dialogID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dialogIDs, nil
}

func GetMessageByID(id string) (model.Message, error) {
	var m model.Message
	err := db.QueryRow("SELECT id, dialog_id, timestamp, role, content FROM messages WHERE id=?", id).Scan(&m.ID, &m.DialogId, &m.Timestamp, &m.Role, &m.Content)
	if err != nil {
		return model.Message{}, err
	}

	return m, nil
}

func DeleteMessageByID(id string) error {
	stmt, err := db.Prepare("DELETE FROM messages WHERE id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func GetMessagesByDialogID(dialogID string) ([]model.Message, error) {
	rows, err := db.Query("SELECT id, dialog_id, timestamp, role, content FROM messages WHERE dialog_id=? ORDER BY timestamp ASC", dialogID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []model.Message{}
	for rows.Next() {
		var m model.Message
		err := rows.Scan(&m.ID, &m.DialogId, &m.Timestamp, &m.Role, &m.Content)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}
