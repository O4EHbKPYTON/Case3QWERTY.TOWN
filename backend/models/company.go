package models

import (
	"api/pkg/geocoder"
	"api/pkg/logger"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"strconv"
	"strings"
	"time"
)

func init() {
	orm.RegisterModel(new(Company))
}

type PostCompanyRequest struct {
	OwnerID          int64  `json:"owner_id"`
	Name             string `json:"name"`
	INN              string `json:"inn"`
	OrganizationType string `json:"organization_type"`
	City             string `json:"city"`
	Address          string `json:"address"`
	BusinessSphere   string `json:"business_sphere"`
	Description      string `json:"description"`
}

type CompanyAddressResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	City        string `json:"city"`
	Address     string `json:"address"`
	FullAddress string `json:"full_address"`
}

type Company struct {
	Id               int64     `orm:"auto;column(id)"`
	Owner            *Owner    `orm:"rel(fk);column(owner_id)"`
	Name             string    `orm:"column(name)"`
	INN              string    `orm:"column(inn)"`
	OrganizationType string    `orm:"column(organization_type)"`
	City             string    `orm:"column(city)"`
	Address          string    `orm:"column(address)"`
	BusinessSphere   string    `orm:"column(business_sphere)"`
	Description      string    `orm:"column(description);null"`
	Lat              float64   `orm:"column(lat);null"`
	Lon              float64   `orm:"column(lon);null"`
	CreatedAt        time.Time `orm:"auto_now_add;type(timestamp);column(created_at)"`
	UpdatedAt        time.Time `orm:"auto_now;type(timestamp);column(updated_at)"`
}

type ManualDescriptionRequest struct {
	Description string `json:"description"`
}

// GeneratedDescriptionResponse представляет ответ с сгенерированным описанием
type GeneratedDescriptionResponse struct {
	Description string `json:"description"`
}

type PostCompanyResponse struct {
	ID int64 `json:"id"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

func GetCompanyAddressById(id int64) (string, error) {
	ormer := orm.NewOrmUsingDB("mydatabase")
	company := Company{Id: id}
	err := ormer.Read(&company)
	if err != nil {
		return "", err
	}
	fullAddress := fmt.Sprintf("%s, %s", company.City, company.Address, company.Lat, company.Lon)
	return fullAddress, nil
}

func UpdateCompanyCoordinates(companyID int64, city, address, name string) {
	logFields := map[string]interface{}{
		"company_id": companyID,
		"city":       city,
		"address":    address,
		"name":       name,
	}

	fullAddress := fmt.Sprintf("%s, %s, %s", city, address, name)
	latStr, lonStr, err := geocoder.GetCoordinates(fullAddress)
	if err != nil {
		logger.ErrorAny("Failed to get coordinates", logFields)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		logFields["error"] = err.Error()
		logFields["lat_str"] = latStr
		logger.ErrorAny("Failed to parse latitude", logFields)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		logFields["error"] = err.Error()
		logFields["lon_str"] = lonStr
		logger.ErrorAny("Failed to parse longitude", logFields)
		return
	}

	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		logFields["lat"] = lat
		logFields["lon"] = lon
		logger.ErrorAny("Invalid coordinates", logFields)
		return
	}

	ormer := orm.NewOrmUsingDB("mydatabase")
	_, err = ormer.QueryTable("company").
		Filter("id", companyID).
		Update(orm.Params{
			"lat": lat,
			"lon": lon,
		})
	if err != nil {
		logFields["error"] = err.Error()
		logger.ErrorAny("Failed to update company coordinates", logFields)
		return
	}

	logFields["lat"] = lat
	logFields["lon"] = lon
	logger.InfoAny("Company coordinates updated successfully", logFields)
}

func AddCompany(c Company) (int64, error) {
	logFields := map[string]interface{}{
		"Id":   c.Id,
		"Name": c.Name,
		"INN":  c.INN,
	}

	logger.InfoAny("Attempting to create Company", logFields)

	ormer := orm.NewOrmUsingDB("mydatabase")
	id, err := ormer.Insert(&c)
	if err != nil {
		errMsg := err.Error()

		switch {
		case strings.Contains(errMsg, "company_inn_key"):
			logger.WarnAny("INN already exists", logFields)
			return 0, fmt.Errorf("company with INN '%s' already exists", c.INN)

		case strings.Contains(errMsg, "owners_fullname_key"):
			logger.WarnAny("Owner name already exists", logFields)
			return 0, fmt.Errorf("owner name '%s' is already taken", c.Name)

		case strings.Contains(errMsg, "violates foreign key constraint"):
			logger.WarnAny("Invalid owner reference", logFields)
			return 0, errors.New("invalid owner ID: referenced owner does not exist")

		default:
			logFields["error"] = errMsg
			logger.ErrorAny("Database error when creating company", logFields)
			return 0, fmt.Errorf("failed to create company: %v", err)
		}
	}

	logFields["company_id"] = id
	logger.InfoAny("Company created successfully", logFields)

	go UpdateCompanyCoordinates(id, c.City, c.Address, c.Name)

	return id, nil
}

func GetCompany(cid int64) (c *Company, err error) {
	o := orm.NewOrmUsingDB("mydatabase")
	company := Company{Id: cid}
	err = o.Read(&company)
	if err == orm.ErrNoRows {
		return nil, errors.New("Company with this id not found")
	}
	return &company, nil
}

func UpdateCompany(c *Company) error {
	o := orm.NewOrmUsingDB("mydatabase")
	_, err := o.Update(c, "name", "city", "address", "organization_type", "business_sphere", "description")
	return err
}

func DeleteCompany(uid int64) error {
	o := orm.NewOrmUsingDB("mydatabase")
	if _, err := o.Delete(&Company{Id: uid}); err != nil {
		return fmt.Errorf("failed to delete company: %v", err)
	}
	return nil
}

func GetOwnerCompanies(ownerId int64) ([]Company, error) {
	o := orm.NewOrmUsingDB("mydatabase")
	var companies []Company

	_, err := o.QueryTable("company").Filter("owner_id", ownerId).All(&companies)
	if err != nil {
		logger.ErrorAny("Error fetching companies for owner", map[string]interface{}{
			"owner_id": ownerId,
			"error":    err.Error(),
		})
		return nil, fmt.Errorf("failed to fetch companies for owner with ID %d: %v", ownerId, err)
	}

	return companies, nil
}

func GetAllCompanies() ([]Company, error) {
	var companies []Company
	o := orm.NewOrmUsingDB("mydatabase")
	qb, _ := orm.NewQueryBuilder("postgres")
	qb.Select("id", "name", "city", "address", "lat", "lon").From("company").Where("id > ?").OrderBy("id").Desc().Limit(10)
	_, err := o.Raw(qb.String(), 0).QueryRows(&companies)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch companies: %v", err)
	}
	return companies, nil
}
