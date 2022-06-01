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

func (r *UserRepository) FindByID(peerID int, userID int) (*models.User, error) {
	// Find user by id
	u := &models.User{}
	if err := r.store.db.QueryRow(
		"SELECT * FROM users WHERE user_id = $1 AND peer_id = $2",
		userID, peerID,
	).Scan(
		&u.UserID, &u.UserName, &u.PeerID, &u.Wins, &u.Loses,
	); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) WinByID(peerID int, winUserID int, loseUserId int) error {
	// Update wins and loses by id
	if _, err := r.store.db.Exec(
		"UPDATE users SET wins = wins + 1 WHERE user_id = $1 AND peer_id = $2",
		winUserID, peerID,
	); err != nil {
		return err
	}
	if _, err := r.store.db.Exec(
		"UPDATE users SET loses = loses + 1 WHERE user_id = $1 AND peer_id = $2",
		loseUserId, peerID,
	); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) NameByID(peerID int, userID int, name string) error {
	// Update name by id
	if _, err := r.store.db.Exec(
		"UPDATE users SET user_name = $1 WHERE user_id = $2 AND peer_id = $3",
		name, userID, peerID,
	); err != nil {
		return err
	}
	return nil
}
