package models

import (
	"api/pkg/logger"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/golang-jwt/jwt/v5"
	"strings"
	"time"
)

func init() {
	orm.RegisterModel(new(Owner))
}

type PostOwnerRequest struct {
	FullName     string `json:"fullname"`
	ContactEmail string `json:"email"`
	ContactPhone string `json:"phone"`
	Password     string `json:"password"`
}

type PostOwnerResponse struct {
	Id int64 `json:"id"`
}

type Owner struct {
	Id       int64  `orm:"auto;column(id)"`
	Fullname string `orm:"column(full_name)"`
	Email    string `orm:"column(contact_email)"`
	Password string `orm:"column(password)"`
	Phone    string `orm:"column(contact_phone)"`
}

type OwnerLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var SecretOwnerKey = []byte("your-secret-key")

func CreateOwnerToken(owner Owner) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    owner.Id,
		"email": owner.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	return claims.SignedString(SecretOwnerKey)
}

func VerifyOwnerToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretOwnerKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}

func AddOwner(o Owner) (int64, error) {
	logFields := map[string]interface{}{
		"fullname": o.Fullname,
		"email":    o.Email,
	}

	logger.InfoAny("Attempting to create Owner", logFields)

	ormer := orm.NewOrmUsingDB("mydatabase")
	id, err := ormer.Insert(&o)
	if err != nil {
		if strings.Contains(err.Error(), "owners_fullname_key") {
			logger.WarnAny("Name already exists", logFields)
			return 0, fmt.Errorf("Name '%s' is already taken", o.Fullname)
		}

		if strings.Contains(err.Error(), "owners_email_key") {
			logger.WarnAny("Email already exists", logFields)
			return 0, fmt.Errorf("email '%s' is already registered", o.Email)
		}

		logFields["error"] = err.Error()
		logger.ErrorAny("Database error when creating owner", logFields)
		return 0, fmt.Errorf("failed to create owner: %v", err)
	}

	logFields["owner_id"] = id
	logger.InfoAny("Owner created successfully", logFields)

	return id, nil
}

func GetOwner(uid int64) (u *Owner, err error) {
	o := orm.NewOrmUsingDB("mydatabase")
	owner := Owner{Id: uid}
	err = o.Read(&owner)
	if err == orm.ErrNoRows {
		return nil, errors.New("owner with this id not found")
	}
	return &owner, nil
}

func GetAllOwners() *[]Owner {
	var owners []Owner
	o := orm.NewOrmUsingDB("mydatabase")
	qb, _ := orm.NewQueryBuilder("postgres")
	qb.Select("id", "full_name", "password").From("owner").Where("id > ?").OrderBy("id").Desc().Limit(10)
	o.Raw(qb.String(), 0).QueryRows(&owners)
	return &owners
}

func UpdateOwner(oo *Owner) error {
	o := orm.NewOrmUsingDB("mydatabase")
	_, err := o.Update(oo, "full_name", "contact_email", "contact_phone", "password")
	return err
}

func DeleteOwner(uid int64) error {
	o := orm.NewOrmUsingDB("mydatabase")
	if _, err := o.Delete(&Owner{Id: uid}); err != nil {
		return fmt.Errorf("failed to delete owner: %v", err)
	}
	return nil
}

func LoginOwner(req OwnerLoginRequest) (string, error) {
	o := orm.NewOrmUsingDB("mydatabase")
	var owner Owner
	err := o.QueryTable("owner").Filter("contact_email", req.Email).One(&owner)
	if err != nil {
		return "", errors.New("invalid email or password")
	}
	if req.Password != owner.Password {
		return "", errors.New("invalid email or password")
	}
	return CreateOwnerToken(owner)
}
