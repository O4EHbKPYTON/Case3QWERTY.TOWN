package controllers

import (
	"api/models"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Owner
type OwnerController struct {
	beego.Controller
}

type OwnerResponse struct {
	Err  bool `json:"err"`
	Data any  `json:"data"`
}

func (o *OwnerController) HandlerFunc(rules string) bool {
	switch rules {
	case "GetAll", "Logout":
		auth := AuthController{Controller: o.Controller}
		return auth.CheckAuth()
	default:
		return false
	}
}

// @Title CreateOwner
// @Description create Owners
// @Param	body		body 	models.PostOwnerRequest	true		"body1 for Owner content"
// @Success 200 {object} models.PostOwnerResponse
// @Failure 403 body is empty
// @router / [post]
func (o *OwnerController) Post() {
	var ownerReq models.PostOwnerRequest
	if err := json.Unmarshal(o.Ctx.Input.RequestBody, &ownerReq); err != nil {
		o.Data["json"] = OwnerResponse{Err: true, Data: "Invalid request"}
		o.ServeJSON()
		return
	}

	owner := models.Owner{
		Fullname: ownerReq.FullName,
		Email:    ownerReq.ContactEmail,
		Phone:    ownerReq.ContactPhone,
		Password: ownerReq.Password,
	}

	uid, err := models.AddOwner(owner)
	if err != nil {
		o.Data["json"] = OwnerResponse{Err: true, Data: "Failed to create owner: " + err.Error()}
		o.ServeJSON()
		return
	}

	o.Data["json"] = OwnerResponse{Err: false, Data: models.PostOwnerResponse{Id: uid}}
	o.ServeJSON()
}

// @Title GetAll
// @Description get all Owners
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.Owner
// @router / [get]
func (o *OwnerController) GetAll() {
	owners := models.GetAllOwners()
	o.Data["json"] = OwnerResponse{Err: false, Data: owners}
	o.ServeJSON()
}

//Сначала нужно ввести токен из Login()

// @Title Get
// @Description get Owner by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Owner
// @Failure 403 {string} string "uid is empty"
// @router /:uid [get]
func (o *OwnerController) Get() {
	uid, err := o.GetInt64(":uid")
	if err == nil {
		owner, err := models.GetOwner(uid)
		if err != nil {
			o.Data["json"] = OwnerResponse{Err: true, Data: err.Error()}
		} else {
			o.Data["json"] = OwnerResponse{Err: false, Data: owner}
		}
	}
	o.ServeJSON()
}

// @Title Update
// @Description update the owner
// @Param	body		body 	models.Owner	true		"body for Owner content"
// @Success 200 {object} models.Owner
// @Failure 403 body is empty
// @router / [put]
func (o *OwnerController) Put() {
	var owner models.Owner
	json.Unmarshal(o.Ctx.Input.RequestBody, &owner)
	err := models.UpdateOwner(&owner)
	if err != nil {
		o.Data["json"] = OwnerResponse{Err: true, Data: err.Error()}
	} else {
		o.Data["json"] = OwnerResponse{Err: false, Data: owner}
	}
	o.ServeJSON()
}

// @Title Delete
// @Description delete the Owner
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (o *OwnerController) Delete() {
	uid, err := o.GetInt64(":uid")
	if err != nil {
		o.Data["json"] = OwnerResponse{Err: true, Data: "Invalid ID"}
		o.ServeJSON()
		return
	}

	if err := models.DeleteOwner(uid); err != nil {
		o.Data["json"] = OwnerResponse{Err: true, Data: err.Error()}
	} else {
		o.Data["json"] = OwnerResponse{Err: false, Data: "Owner deleted"}
	}
	o.ServeJSON()
}

// @Title Login
// @Description Авторизация пользователя
// @Param	body		body 	models.OwnerLoginRequest	true	"Данные для входа (логин и пароль)"
// @Success 200 {object} OwnerResponse "Успешный вход, возвращает токен"
// @Failure 400 {object} OwnerResponse "Ошибка в теле запроса"
// @Failure 401 {object} OwnerResponse "Неверные логин или пароль"
// @router /login [post]
func (o *OwnerController) Login() {
	var loginReq models.OwnerLoginRequest

	if err := json.Unmarshal(o.Ctx.Input.RequestBody, &loginReq); err != nil {
		o.Ctx.Output.SetStatus(400)
		o.Data["json"] = OwnerResponse{Err: true, Data: "Invalid request"}
		o.ServeJSON()
		return
	}

	token, err := models.LoginOwner(loginReq)
	if err != nil {
		o.Ctx.Output.SetStatus(401)
		o.Data["json"] = OwnerResponse{Err: true, Data: err.Error()}
		o.ServeJSON()
		return
	}

	o.Data["json"] = OwnerResponse{Err: false, Data: token}
	o.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in Owner session
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} logout success
// @router /logout [get]
func (o *OwnerController) Logout() {
	o.Data["json"] = OwnerResponse{Err: false, Data: "Вышли из сессии"}
	o.ServeJSON()
}
