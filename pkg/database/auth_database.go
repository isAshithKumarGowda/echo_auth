package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/isAshithKumarGowda/Echo_Auth/internals/models"
)

func (q *Query) Register(user models.AuthModel, userType string) error {
	tx, err := q.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	switch userType {
	case "admin":
		_, err := tx.Exec("INSERT INTO admin(admin_id, admin_name, admin_email, admin_hash_password) VALUES($1, $2, $3, $4)", user.ID, user.Name, user.Email, user.Password)
		if err != nil {
			return err
		}
	case "user":
		_, err := tx.Exec("INSERT INTO users(user_id, user_name, user_email, user_hash_password) VALUES($1, $2, $3, $4)", user.ID, user.Name, user.Email, user.Password)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid user type")
	}
	return nil
}

func (q *Query) Login(user models.AuthLoginModel, userType string) (models.AuthModel, time.Time, error) {
	var User models.AuthModel

	var exists bool
	err := q.db.QueryRow("SELECT EXISTS(SELECT 1 FROM verified_emails WHERE email = $1)", user.Email).Scan(&exists)

	if err != nil {
		return User, time.Now(), fmt.Errorf("email is not verified: %v", err)
	}

	switch userType {
	case "admin":
		err = q.db.QueryRow("SELECT admin_id, admin_name, admin_email, admin_hash_password FROM admin WHERE admin_email = $1", user.Email).Scan(&User.ID, &User.Name, &User.Email, &User.Password)
		if err != nil {
			return User, time.Now(), err
		}
		_, err = q.db.Exec("INSERT INTO admin_login_history(email) VALUES ($1)", user.Email)
		if err != nil {
			return User, time.Now(), err
		}
	case "user":
		err = q.db.QueryRow("SELECT user_id, user_name, user_email, user_hash_password FROM users WHERE user_email = $1", user.Email).Scan(&User.ID, &User.Name, &User.Email, &User.Password)
		if err != nil {
			return User, time.Now(), err
		}
		_, err = q.db.Exec("INSERT INTO user_login_history(email) VALUES ($1)", user.Email)
		if err != nil {
			return User, time.Now(), err
		}
	default:
		return User, time.Now(), fmt.Errorf("invalid UserType")
	}

	return User, time.Now(), nil
}

// func (q *Query) StoreOtp(email string, otp string) error {
// 	query := `INSERT INTO otps(email, otp) VALUES ($1, $2) ON CONFLICT (email) DO UPDATE SET otp = EXCLUDED.otp`

// 	_, err := q.db.Exec(query, email, otp)

// 	return err
// }

// func (q *Query) ClearOtp(email string, otp string) error {
// 	query := `DELETE from otps WHERE email = $1 AND otp = $2`
// 	_, err := q.db.Exec(query, email, otp)
// 	return err
// }

// func (q *Query) IsOtpValid(email string, otp string) (bool, error) {
// 	query1 := `SELECT EXISTS(
// 				SELECT 1 FROM otps WHERE email = $1 AND otp = $2
// 			)`

// 	query2 := `DELETE FROM otps WHERE email = $1 AND otp = $2`

// 	query3 := `INSERT INTO verified_emails(email) VALUES ($1)`

// 	tx, err := q.db.Begin()

// 	if err != nil {
// 		return false, fmt.Errorf("error in q.db.Begin %v", err)
// 	}

// 	if _, err = tx.Exec(query1, email, otp); err != nil {
// 		tx.Rollback()
// 		return false, fmt.Errorf("error while executing query1 %v", err)
// 	}

// 	if _, err = tx.Exec(query2, email, otp); err != nil {
// 		tx.Rollback()
// 		return false, fmt.Errorf("error while executing query2 %v", err)
// 	}

// 	if _, err = tx.Exec(query3, email); err != nil {
// 		tx.Rollback()
// 		return false, fmt.Errorf("error while executing query3 %v", err)
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return false, fmt.Errorf("error while commiting the querys %v", err)
// 	}

// 	return true, nil

// }

func (q *Query) GetUserDetails(email string, userType string) (models.AuthModel, error) {
	var User models.AuthModel

	switch userType {
	case "admin":
		err := q.db.QueryRow("SELECT admin_id, admin_name FROM admin WHERE admin_email = $1", email).
			Scan(&User.ID, &User.Name)
		if err != nil {
			return User, fmt.Errorf("error while executing admin query %v", err)
		}
	case "user":
		err := q.db.QueryRow("SELECT user_id, user_name FROM users WHERE user_email = $1", email).
			Scan(&User.ID, &User.Name)
		if err != nil {
			return User, fmt.Errorf("error while executing user query %v", err)
		}
	default:
		return User, fmt.Errorf("invalid user type")
	}
	return User, nil
}

func (q *Query) IfUserExists(email string) error {
	_, err := q.db.Query("SELECT * FROM verified_emails WHERE email = $1", email)
	if err != nil {
		return err
	}
	return nil
}

func (q *Query) SetVerifiedEmail(email string) error {
	_, err := q.db.Exec("INSERT INTO verified_emails(email) VALUES ($1)", email)
	return err
}

func (q *Query) GetLoginHistory(limit int, offset int, userType string) ([]models.AuthLoginHistoryModel, error) {
	var users []models.AuthLoginHistoryModel
	var rows *sql.Rows
	var err error

	switch userType {
	case "admin":
		rows, err = q.db.Query("SELECT email, login_time, logout_time FROM admin_login_history LIMIT $1 OFFSET $2", limit, offset)
		if err != nil {
			return nil, err
		}
	case "user":
		rows, err = q.db.Query("SELECT email, login_time, logout_time FROM user_login_history LIMIT $1 OFFSET $2", limit, offset)
		if err != nil {
			return nil, err
		}
	}

	for rows.Next() {
		var user models.AuthLoginHistoryModel
		rows.Scan(&user.Email, &user.Login, &user.Logout)
		users = append(users, user)
	}

	return users, nil
}

func (q *Query) GetTotalCount(userType string) (int, error) {
	var total int

	switch userType {
	case "admin":
		err := q.db.QueryRow("SELECT COUNT(*) FROM admin_login_history").Scan(&total)
		if err != nil {
			return -1, err
		}
	case "user":
		err := q.db.QueryRow("SELECT COUNT(*) FROM user_login_history").Scan(&total)
		if err != nil {
			return -1, err
		}
	}

	return total, nil
}

func (q *Query) Logout(userType, email string, login time.Time) error {
	switch userType {
	case "admin":
		query := "UPDATE admin_login_history SET logout_time = NOW() WHERE email = $1 AND login = $2"
		_, err := q.db.Exec(query, email, login)
		if err != nil {
			return err
		}
	case "user":
		query := "UPDATE user_login_history SET logout_time = NOW() WHERE email = $1 AND login = $2"
		_, err := q.db.Exec(query, email, login)
		if err != nil {
			return err
		}
	}

	return nil
}
