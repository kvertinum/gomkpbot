package store

import "github.com/Kvertinum01/gomkpbot/internal/app/models"

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *models.User) error {
	// Init the user in database
	return r.store.db.QueryRow(
		"INSERT INTO users (user_id, user_name, peer_id, wins, loses) VALUES ($1, $2, $3, $4, $5)",
		u.UserID, u.UserName, u.PeerID, u.Wins, u.Loses,
	).Err()
}

func (r *UserRepository) FindByID(userID int) (*models.User, error) {
	// Find user by id
	u := &models.User{}
	if err := r.store.db.QueryRow(
		"SELECT * FROM users WHERE user_id = $1", userID,
	).Scan(
		&u.UserID, &u.UserName, &u.PeerID, &u.Wins, &u.Loses,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) WinByID(winUserID int, loseUserId int) error {
	// Update wins and loses by id
	if _, err := r.store.db.Exec(
		"UPDATE users SET wins = wins + 1 WHERE user_id = $1", winUserID,
	); err != nil {
		return err
	}
	if _, err := r.store.db.Exec(
		"UPDATE users SET loses = loses + 1 WHERE user_id = $1", loseUserId,
	); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) NameByID(userID int, name string) error {
	// Update name by id
	if _, err := r.store.db.Exec(
		"UPDATE users SET user_name = $1 WHERE user_id = $2",
		name, userID,
	); err != nil {
		return err
	}
	return nil
}
