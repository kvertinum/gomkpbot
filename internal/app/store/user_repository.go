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
