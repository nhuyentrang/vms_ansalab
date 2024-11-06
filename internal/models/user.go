package models

// Model for users
/*************************************************************** Users ***************************************************************/
type User struct {
	ID   int64 `json:"id"`
	Role int64 `json:"role"`
}

type DTO_User struct {
	ID   int64 `json:"id" validate:"required"`
	Role int64 `json:"role" validate:"required"`
}

type UserHome struct {
	ID         int64 `json:"id"`
	UserID     int64 `json:"user_id"`
	HomeID     int64 `json:"home_id"`
	Permission int64 `json:"permission"`
}

type DTO_UserHome struct {
	ID         int64 `json:"id" validate:"required"`
	UserID     int64 `json:"user_id" validate:"required"`
	HomeID     int64 `json:"home_id" validate:"required"`
	Permission int64 `json:"permission" validate:"required"`
}

type Login struct {
	ID     int64  `json:"id"`
	Token  string `json:"token"`
	UserID int64  `json:"user_id"`
}

type DTO_Login struct {
	ID     int64  `json:"id" validate:"required"`
	Token  string `json:"token" validate:"required"`
	UserID int64  `json:"user_id" validate:"required"`
}
