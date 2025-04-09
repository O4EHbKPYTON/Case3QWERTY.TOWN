package controllers

import (
	"api/models"
	"api/pkg/logger"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

// Operations about YandexMaps
// controllers/geocoder.go
// @APIVersion 1.0.0
// @Title Geo Controller
// @Description Геокодирование адресов компаний
type GeoController struct {
	web.Controller
}
type GetCoordinatesRequest struct {
	CompanyID int64
}

type GetCoordinatesResponse struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Address   string `json:"address"`
}

// @Title Get Coordinates by Company ID
// @Description Возвращает координаты (широту и долготу) и адрес компании по её ID
// @Param id path int true "ID компании"
// @Success 200 {object} controllers.GetCoordinatesResponse
// @Success 202 {object} controllers.GetCoordinatesResponse "Координаты еще не получены, начат процесс геокодирования"
// @Failure 400 {string} string "Invalid company ID"
// @Failure 404 {string} string "Company not found"
// @Failure 500 {string} string "Failed to get coordinates"
// @router /geo/company/:id [get]
func (c *GeoController) GetCoordinatesByCompanyId() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.CustomAbort(400, "Invalid company ID")
		return
	}

	// Получаем компанию из базы
	company, err := models.GetCompany(id)
	if err != nil {
		if err == orm.ErrNoRows {
			c.CustomAbort(404, "Company not found")
		} else {
			c.CustomAbort(500, "Failed to get company: "+err.Error())
		}
		return
	}

	// Если координаты уже есть - возвращаем их
	if company.Lat != 0 && company.Lon != 0 {
		fullAddress := fmt.Sprintf("%s, %s", company.City, company.Address)
		res := GetCoordinatesResponse{
			Latitude:  strconv.FormatFloat(company.Lat, 'f', 6, 64),
			Longitude: strconv.FormatFloat(company.Lon, 'f', 6, 64),
			Address:   fullAddress,
		}
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	go func() {
		models.UpdateCompanyCoordinates(company.Id, company.City, company.Address, company.Name)
	}()

	fullAddress := fmt.Sprintf("%s, %s", company.City, company.Address)
	res := GetCoordinatesResponse{
		Address: fullAddress,
	}

	c.Data["json"] = res
	c.ServeJSON()
}

// GetAllCompaniesCoordinatesResponse представляет структуру ответа для метода GetAllCompaniesCoordinates
// swagger:model
type GetAllCompaniesCoordinatesResponse struct {
	CompanyAddressResponse
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// CompanyAddressResponse представляет базовую информацию о компании и её адресе
// swagger:model
type CompanyAddressResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	City    string `json:"city"`
	Address string `json:"address"`
}

// @Title Get Coordinates for All Companies
// @Description Возвращает координаты (широту и долготу) и адреса для всех компаний (первые 10)
// @Success 200 {array} GeoController.GetAllCompaniesCoordinatesResponse
// @Success 206 {array} GeoController.GetAllCompaniesCoordinatesResponse "Часть компаний не имеет координат (начат процесс геокодирования)"
// @Failure 404 {string} string "Companies not found"
// @Failure 500 {string} string "Failed to get companies"
// @router /geo/companies [get]
func (c *GeoController) GetAllCompaniesCoordinates() {
	// Получаем список компаний из модели
	companies, err := models.GetAllCompanies()
	if err != nil {
		c.CustomAbort(500, "Failed to get companies: "+err.Error())
		return
	}

	if len(companies) == 0 {
		c.CustomAbort(404, "Companies not found")
		return
	}

	// Подготавливаем ответ
	var response []GetAllCompaniesCoordinatesResponse
	var hasMissingCoordinates bool

	for _, company := range companies {

		item := GetAllCompaniesCoordinatesResponse{
			CompanyAddressResponse: CompanyAddressResponse{
				ID:      company.Id,
				City:    company.City,
				Address: company.Address,
				Name:    company.Name,
			},
		}

		if company.Lat != 0 && company.Lon != 0 {
			item.Latitude = strconv.FormatFloat(company.Lat, 'f', 6, 64)
			item.Longitude = strconv.FormatFloat(company.Lon, 'f', 6, 64)
		} else {
			logFields := map[string]interface{}{
				"company_id": company.Id,
				"db_lat":     company.Lat,
				"db_lon":     company.Lon,
			}
			//go models.UpdateCompanyCoordinates(company.Id, company.City, company.Address, company.Name)
			logFields["lat"] = company.Lat
			logFields["lon"] = company.Lon
			logger.InfoAny("Company coordinates updated successfully", logFields)
			hasMissingCoordinates = true
		}

		response = append(response, item)
	}

	if hasMissingCoordinates {
		c.Ctx.ResponseWriter.WriteHeader(206)
	}

	c.Data["json"] = response
	c.ServeJSON()
}
