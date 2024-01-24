package storge

import (
	"database/sql"
	"handlar_tes/model"
)

func connect() (*sql.DB, error) {
	dsn := "user=asadbek password=1234 dbname=handlar sslmode=disable"
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateUser(user *model.User) (*model.User, error) {
	db, err := connect()
	defer db.Close()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	var respUser model.User

	query := `INSERT INTO users (id, first_name, last_name) VALUES ($1, $2, $3) RETURNING id , first_name, last_name`

	err = db.QueryRow(query, user.ID, user.FirstName, user.LastName).Scan(&respUser.ID, &respUser.FirstName, &respUser.LastName)

	if err != nil {
		return nil, err
	}

	return &respUser, nil
}

func Get(id string) (*model.User, error) {
	db, err := connect()
	defer db.Close()

	if err != nil {
		panic(err)
	}

	var respUser model.User

	query := `SELECT id, first_name, last_name FROM users WHERE id = $1`

	err = db.QueryRow(query, id).Scan(&respUser.ID, &respUser.FirstName, &respUser.LastName)

	if err != nil {
		return nil, err
	}

	return &respUser, nil
}

func GetAll(page, limit int) (users []*model.User, err error) {
	db, err := connect()
	defer db.Close()

	if err != nil {
		panic(err)
	}

	offset := limit * (page - 1)

	rows, err := db.Query(`SELECT id, first_name, last_name FROM users LIMIT $1 offset $2`, limit, offset)
	if err != nil {
		return nil, err
	}

	// var users []*model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)

	}

	return users, nil
}

func DeleteUser(id string) error {
	db, err := connect()
	defer db.Close()

	if err != nil {
		return err
	}

	_, err = db.Exec(`DELETE FROM users WHERE id = $1`, id)

	if err != nil {
		return err
	}

	return nil
}

func UptadeUser(id string, name *model.User) (*model.User, error) {
	db, err := connect()
	defer db.Close()

	if err != nil {
		return nil, err
	}

	var respoUser model.User

	query := `
	UPDATE 	
		users 
	SET 
		first_name = $1, 
		last_name = $2
	WHERE 
		id = $3 
	RETURNING
		id, first_name, last_name`

	if err = db.QueryRow(query, name.FirstName, name.LastName, id).Scan(&respoUser.ID, &respoUser.FirstName, &respoUser.LastName); err != nil {
		panic(err)
	}

	return &respoUser, err

}
