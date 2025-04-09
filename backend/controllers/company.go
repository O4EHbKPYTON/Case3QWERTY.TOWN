package controllers

import (
	"api/models"
	"api/pkg/generator"
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"log"
	"regexp"
	"strconv"
)

// Operations about companies
type CompanyController struct {
	web.Controller
}

type DescriptionGenerator interface {
	GenerateCompanyDescription(name, businessSphere string) (string, error)
}

func isValidINN(inn string) bool {
	re := regexp.MustCompile(`^\d{10}$|^\d{12}$`)
	return re.MatchString(inn)
}

// @Title GenerateCompanyDescription
// @Description Генерация описания компании по ID компании с автоматическим сохранением
// @Param id path int true "ID компании"
// @Success 200 {object} models.Company
// @Failure 400 {string} string "Invalid company ID"
// @Failure 404 {string} string "Company not found"
// @Failure 500 {string} string "Generation failed"
// @router /:id/generate-description [get]
func (c *CompanyController) GenerateDescription() {
	// Получаем ID компании из URL
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.CustomAbort(400, "Invalid company ID")
		return
	}

	company, err := models.GetCompany(id)
	if err != nil {
		if err.Error() == "Company with this id not found" {
			c.CustomAbort(404, "Company not found")
		} else {
			c.CustomAbort(500, "Failed to get company: "+err.Error())
		}
		return
	}

	if company.Name == "" || company.BusinessSphere == "" {
		c.CustomAbort(400, "Company name or business sphere is empty")
		return
	}

	generator := generator.NewOpenRouterClient("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")

	description, err := generator.GenerateCompanyDescription(company.Name, company.BusinessSphere)
	if err != nil {
		log.Printf("Ошибка генерации: %v", err)
		c.CustomAbort(500, "Description generation failed: "+err.Error())
		return
	}

	company.Description = description
	if err := models.UpdateCompany(company); err != nil {
		c.CustomAbort(500, "Failed to save generated description: "+err.Error())
		return
	}

	c.Data["json"] = company
	c.ServeJSON()
}

// @Title UpdateCompanyDescription
// @Description Ручное обновление описания компании по ID
// @Param id path int true "ID компании"
// @Param description body models.ManualDescriptionRequest true "Описание компании"
// @Success 200 {object} models.Company
// @Failure 400 {string} string "Invalid company ID or description"
// @Failure 404 {string} string "Company not found"
// @Failure 500 {string} string "Failed to update company description"
// @router /:id/update-description [put]
func (c *CompanyController) UpdateDescription() {
	// Получаем ID компании из URL
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.CustomAbort(400, "Invalid company ID")
		return
	}

	// Получаем данные запроса
	var request models.ManualDescriptionRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &request); err != nil {
		c.CustomAbort(400, "Invalid request body: "+err.Error())
		return
	}

	// Проверяем, что описание не пустое
	if request.Description == "" {
		c.CustomAbort(400, "Description cannot be empty")
		return
	}

	// Получаем компанию из базы данных
	company, err := models.GetCompany(id)
	if err != nil {
		if err.Error() == "Company with this id not found" {
			c.CustomAbort(404, "Company not found")
		} else {
			c.CustomAbort(500, "Failed to get company: "+err.Error())
		}
		return
	}

	company.Description = request.Description
	if err := models.UpdateCompany(company); err != nil {
		c.CustomAbort(500, "Failed to update company description: "+err.Error())
		return
	}

	c.Data["json"] = company
	c.ServeJSON()
}

// @Title CreateCompany
// @Description Создание компании
// @Param body body models.PostCompanyRequest true "Данные компании"
// @Success 200 {object} models.PostCompanyResponse "ID созданной компании"
// @Failure 400 {string} string "Invalid input"
// @router / [post]
func (c *CompanyController) Post() {
	var req models.PostCompanyRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.CustomAbort(400, "Invalid input JSON")
		return
	}

	if req.INN == "" || req.Name == "" {
		c.CustomAbort(400, "Missing required fields: name, inn")
		return
	}

	if req.INN != "string" {
		if !isValidINN(req.INN) {
			c.CustomAbort(400, "Invalid INN: must be 10 or 12 digits")
			return
		}
	}

	company := models.Company{
		Owner:            &models.Owner{Id: req.OwnerID},
		Name:             req.Name,
		INN:              req.INN,
		OrganizationType: req.OrganizationType,
		City:             req.City,
		Address:          req.Address,
		Description:      req.Description,
		BusinessSphere:   req.BusinessSphere,
	}

	id, err := models.AddCompany(company)
	if err != nil {
		c.CustomAbort(400, err.Error())
		return
	}

	c.Data["json"] = models.PostCompanyResponse{ID: id}
	c.ServeJSON()
}

// @Title GetCompany
// @Description Получить компанию по ID
// @Param id path int true "ID компании"
// @Success 200 {object} models.Company
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Company not found"
// @router /:id [get]
func (c *CompanyController) Get() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.CustomAbort(400, "Invalid ID")
		return
	}

	company, err := models.GetCompany(id)
	if err != nil {
		c.CustomAbort(404, err.Error())
		return
	}

	c.Data["json"] = company
	c.ServeJSON()
}

// @Title UpdateCompany
// @Description Обновить данные компании
// @Param id path int true "ID компании"
// @Param body body models.PostCompanyRequest true "Данные для обновления"
// @Success 200 {object} models.StatusResponse "status: updated"
// @Failure 400 {string} string "Invalid input or ID"
// @Failure 404 {string} string "Company not found"
// @Failure 500 {string} string "Update failed"
// @router /:id [put]
func (c *CompanyController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.CustomAbort(400, "Invalid ID")
		return
	}

	var req models.PostCompanyRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.CustomAbort(400, "Invalid input")
		return
	}

	company, err := models.GetCompany(id)
	if err != nil {
		c.CustomAbort(404, err.Error())
		return
	}

	company.Name = req.Name
	company.City = req.City
	company.Address = req.Address
	company.OrganizationType = req.OrganizationType
	company.BusinessSphere = req.BusinessSphere

	if err := models.UpdateCompany(company); err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	c.Data["json"] = models.StatusResponse{Status: "updated"}
	c.ServeJSON()
}

// @Title DeleteCompany
// @Description Удалить компанию по ID
// @Param id path int true "ID компании"
// @Success 200 {object} models.StatusResponse "status: deleted"
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Delete failed"
// @router /:id [delete]
func (c *CompanyController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.CustomAbort(400, "Invalid ID")
		return
	}

	if err := models.DeleteCompany(id); err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	c.Data["json"] = models.StatusResponse{Status: "deleted"}
	c.ServeJSON()
}

// @Title GetAllCompanies
// @Description Получить список всех компаний
// @Success 200 {array} models.Company
// @Failure 500 {string} string "Internal server error"
// @router / [get]
func (c *CompanyController) GetAll() {
	companies, err := models.GetAllCompanies()
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	c.Data["json"] = companies
	c.ServeJSON()
}

// @Title GetOwnerCompanies
// @Description Получить все компании по ID владельца
// @Param owner_id path int true "ID владельца"
// @Success 200 {array} models.Company
// @Failure 400 {string} string "Invalid owner ID"
// @Failure 500 {string} string "Internal server error"
// @router / [get]
func (c *CompanyController) GetOwnerCompanies() {
	ownerIdStr := c.Ctx.Input.Param(":owner_id")
	ownerId, err := strconv.ParseInt(ownerIdStr, 10, 64)
	if err != nil {
		c.CustomAbort(400, "Invalid owner ID")
		return
	}

	companies, err := models.GetOwnerCompanies(ownerId)
	if err != nil {
		c.CustomAbort(500, err.Error())
		return
	}

	c.Data["json"] = companies
	c.ServeJSON()
}
