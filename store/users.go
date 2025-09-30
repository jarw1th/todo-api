package store

import (
	models "ToDoProject/models"
	"ToDoProject/safety"
	"errors"

	_ "github.com/lib/pq"
)

func (s *TodoStore) CreateUser(username string, password string) (models.User, error) {
	users, errGet := s.getAllUsers()
	if errGet != nil {
		return models.User{}, errGet
	}
	for _, user := range users {
		if user.Username == username {
			return models.User{}, errors.New("username already exists")
		}
	}

	hashedPassword, errPassword := s.hashPassword(password)
	if errPassword != nil {
		return models.User{}, errPassword
	}

	var u models.User
	err := s.DB.QueryRow(
		"INSERT INTO users(username, password) VALUES($1, $2) RETURNING id, username, password",
		username, hashedPassword,
	).Scan(&u.ID, &u.Username, &u.Password)
	return u, err
}

func (s *TodoStore) CheckUserCredentials(username string, password string) (int, error) {
	users, errGet := s.getAllUsers()
	if errGet != nil {
		return 0, errGet
	}
	for _, user := range users {
		if user.Username == username {
			err := safety.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return 0, errors.New("invalid password")
			}
			return user.ID, nil
		}
	}
	return 0, errors.New("user not found")
}

func (s *TodoStore) DeleteUser(id int, password string) (models.User, error) {
	user, errGer := s.getUserBy(id)
	if errGer != nil {
		return models.User{}, errGer
	}
	if user.Password != password {
		return models.User{}, errors.New("invalid password")
	}
	var u models.User
	err := s.DB.QueryRow(
		"DELETE FROM users WHERE id=$1 RETURNING id, username, password",
		id,
	).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (s *TodoStore) getUserBy(id int) (models.User, error) {
	var u models.User
	err := s.DB.QueryRow(
		"SELECT id, username, password FROM users WHERE id=$1",
		id,
	).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (s *TodoStore) getAllUsers() ([]models.User, error) {
	rows, err := s.DB.Query("SELECT id, username, password FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Username, &u.Password)
		users = append(users, u)
	}
	return users, nil
}

func (s *TodoStore) hashPassword(password string) (string, error) {
	hashed, err := safety.GenerateFromPassword([]byte(password))
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
