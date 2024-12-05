package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExpenseType string

const (
	Type001 ExpenseType = "001"
	Type002 ExpenseType = "002"
)

func (et ExpenseType) IsValid() error {
	if et != Type001 && et != Type002 {
		return errors.New("invalid type: must be '001' or '002'")
	}
	return nil
}

type LoginDetails struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AboutMe struct {
	Name        string `json:"name" bson:"name"`
	Role        string `json:"role" bson:"role"`
	Description string `json:"description" bson:"description"`
	Github      string `json:"github" bson:"github"`
	LinkedIn    string `json:"linkedIn" bson:"linkedIn"`
	Facebook    string `json:"facebook" bson:"facebook"`
	Telegram    string `json:"telegram" bson:"telegram"`
	Image       string `json:"image" bson:"image"`
	T1          string `json:"t1" bson:"t1"`
	T2          string `json:"t2" bson:"t2"`
}

type ServiceInfo struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

type ProjectInfo struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

type Blog struct {
	Title       string `json:"title" bson:"title"`
	SubTitle    string `json:"sub_title" bson:"sub_title"`
	Description string `json:"description" bson:"description"`
	Image       string `json:"image" bson:"image"`
	Link        string `json:"link" bson:"link"`
}

type Service struct {
	Service_ID primitive.ObjectID `json:"_id" bson:"_id"`
	Title      *string            `json:"title" bson:"title"`
	Content    *string            `json:"content" bson:"content"`
	Image      *string            `json:"image" bson:"image"`
	T1         *string            `json:"t1" bson:"t1"`
	T2         *string            `json:"t2" bson:"t2"`
	Created_At time.Time          `json:"created_at" bson:"created_at"`
	Updated_At time.Time          `json:"updated_at" bson:"updated_at"`
}

type Project struct {
	Project_ID  primitive.ObjectID `json:"_id" bson:"_id"`
	Title       *string            `json:"title" bson:"title"`
	Description *string            `json:"description" bson:"description"`
	Role        *string            `json:"role" bson:"role"`
	DemoLink    *string            `json:"demo_link" bson:"demo_link"`
	CodeLink    *string            `json:"code_link" bson:"code_link"`
	Tag         *string            `json:"tag" bson:"tag"`
	Image       *string            `json:"image" bson:"image"`
	T1          *string            `json:"t1" bson:"t1"`
	T2          *string            `json:"t2" bson:"t2"`
	Created_At  time.Time          `json:"created_at" bson:"created_at"`
	Updated_At  time.Time          `json:"updated_at" bson:"updated_at"`
}

type Certificate struct {
	Certificate_ID primitive.ObjectID `json:"_id" bson:"_id"`
	Title          *string            `json:"title" bson:"title"`
	Content        *string            `json:"content" bson:"content"`
	Image          *string            `json:"image" bson:"image"`
	DemoLink       *string            `json:"demo_link" bson:"demo_link"`
	T1             *string            `json:"t1" bson:"t1"`
	T2             *string            `json:"t2" bson:"t2"`
	Created_At     time.Time          `json:"created_at" bson:"created_at"`
	Updated_At     time.Time          `json:"updated_at" bson:"updated_at"`
}

type Message struct {
	Message_ID  primitive.ObjectID `json:"_id" bson:"_id"`
	Name        *string            `json:"name" bson:"name"`
	Email       *string            `json:"email" bson:"email"      validate:"email,required"`
	Phone       *string            `json:"phone" bson:"phone"`
	CompanyName *string            `json:"company_name" bson:"company_name"`
	Message     *string            `json:"message" bson:"message"`
	T1          *string            `json:"t1" bson:"t1"`
	T2          *string            `json:"t2" bson:"t2"`
	Created_At  time.Time          `json:"created_at" bson:"created_at"`
	Updated_At  time.Time          `json:"updated_at" bson:"updated_at"`
}

type User struct {
	User_ID    primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"password" bson:"password"`
	Avatar     string             `json:"avatar" bson:"avatar"`
	Role       int                `json:"role" bson:"role"`
	T1         string             `json:"t1" bson:"t1"`
	T2         string             `json:"t2" bson:"t2"`
	Created_At time.Time          `json:"created_at" bson:"created_at"`
	Updated_At time.Time          `json:"updated_at" bson:"updated_at"`
}

type ExpenseCategory struct {
	Category_ID primitive.ObjectID `json:"_id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Type        ExpenseType        `json:"type" bson:"type"`
	User_ID     string             `json:"user_id" bson:"user_id"`
	T1          string             `json:"t1" bson:"t1"`
	T2          string             `json:"t2" bson:"t2"`
	Created_At  time.Time          `json:"created_at" bson:"created_at"`
	Updated_At  time.Time          `json:"updated_at" bson:"updated_at"`
}

type ExpenseItem struct {
	Item_ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Category_ID string             `json:"category_id" bson:"category_id"`
	User_ID     string             `json:"user_id" bson:"user_id"`
	Type        ExpenseType        `json:"type" bson:"type"`
	Title       string             `json:"title" bson:"title"`
	Remark      string             `json:"remark" bson:"remark"`
	Amount      float64            `json:"amount" bson:"amount"`
	T1          string             `json:"t1" bson:"t1"`
	T2          string             `json:"t2" bson:"t2"`
	Created_At  time.Time          `json:"created_at" bson:"created_at"`
	Updated_At  time.Time          `json:"updated_at" bson:"updated_at"`
}
