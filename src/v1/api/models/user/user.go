package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/response_builder"
	UtilitiesString "github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/string_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/v1/api/database"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
	"time"
)

func (u *User) GetUsersByLawFirm(l string) (*[]User, error) {
	var container []User

	sql := `
		SELECT 
			u.id, 
		    u.firstname, 
		    u.lastname,
		    CONCAT(u.lastname, ', ', u.firstname) AS name
		FROM user u
		INNER JOIN user_type ut 
			ON 1 = 1
		 		AND u.user_type_id = ut.id
				AND ut.name = "Law Firm"
		WHERE 1 = 1
			AND u.law_firm_ref_id = ?
	`

	err := database.DBInstancePublic.Select(&container, sql, l)

	return &container, err
}

func (u *User) GetDetailsById() (*ResponseLoginUserInfo, error) {
	var responseLoginUserInfo ResponseLoginUserInfo

	sql := `
		SELECT 
			u.id,
			u.firstname,
			u.lastname,
			ut.name AS role
		FROM
		user u
		INNER JOIN user_type ut 
			ON 1 = 1
				AND u.user_type_id = ut.id 
		WHERE 1 = 1
			AND u.id = ?
	`

	err := database.DBInstancePublic.Get(&responseLoginUserInfo, sql, u.ID)

	return &responseLoginUserInfo, err
}

func (u *User) GetOne() ([]ResponseUserSingleDisplay, error) {
	var responseUserSingleDisplay []ResponseUserSingleDisplay

	sql := `
		SELECT 
			u.id,
			IF(u.firstname IS NULL, "", u.firstname) AS firstname,
			IF(u.lastname IS NULL, "", u.lastname) AS lastname,
			IF(u.email IS NULL, "", u.email) AS email,
			IF(u.mobile_number IS NULL, "", u.mobile_number) AS mobile_number,
			IF(ut.name IS NULL, "", ut.name) AS role,
			IF(bt.name IS NULL, "", bt.name) AS bank_type,
			IF(u.bank_no IS NULL, "", u.bank_no) AS bank_no,
			IF(u.address IS NULL, "", u.address) AS address,
			IF(u.birthday IS NULL, "", u.birthday) AS birthday,
			IF(u.gender IS NULL, "", u.gender) AS gender,
			IF(u.user_type_id IS NULL, 0, u.user_type_id) AS user_type_id,
			IF(u.bank_type_id IS NULL, 0, u.bank_type_id) AS bank_type_id,
			IF(u.m88_account IS NULL, "", u.m88_account) AS m88_account,
			IF(u.region_id IS NULL, 0, u.region_id) AS region_id
		FROM user u 
		INNER JOIN user_type ut 
			ON 1 = 1
				AND u.user_type_id = ut.id
		INNER JOIN bank_type bt 
			ON 1 = 1
				AND u.bank_type_id = bt.id 
		WHERE 1 = 1      
			and u.is_active = 1
			and u.id = ?
	`

	err := database.DBInstancePublic.Select(&responseUserSingleDisplay, sql, u.ID)

	return responseUserSingleDisplay, err
}

func (u *User) EmailNotTaken() (bool, error) {
	hasEmail := false
	sql := `
		SELECT 
			IF(COUNT(*) > 0, false, true) AS has_email 
		FROM user 
		WHERE 1 = 1
			AND email = ? 
	`

	err := database.DBInstancePublic.Get(&hasEmail, sql, u.Email)

	return hasEmail, err
}

func (u *User) Create(paramsInsert ParamsInsert) (sql.Result, error) {
	fmt.Println("I will print a new user", paramsInsert)

	if paramsInsert.UserTypeId == 15 {
		return nil, errors.New("you can't add a new admin account")
	}

	sql := `
		INSERT INTO user (
			firstname, 
			lastname, 
			email, 
			mobile_number, 
			password, 
			user_type_id, 
			address,
			birthday,
			gender,
			organization_ref_id
		)
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)
	`

	encryptedPassword, _ := UtilitiesString.Encrypt(paramsInsert.Password)

	sqlResult, err := database.DBInstancePublic.Exec(
		sql,
		paramsInsert.FirstName,
		paramsInsert.LastName,
		paramsInsert.Email,
		paramsInsert.MobileNumber,
		encryptedPassword,
		paramsInsert.UserTypeId,
		paramsInsert.Address,
		paramsInsert.Birthday,
		paramsInsert.Gender,
		paramsInsert.OrganizationID,
	)

	return sqlResult, err
}

func (p *ParamsUpdate) Update(userType string) (sql.Result, error) {
	var sqlResult sql.Result
	var sql = ""
	var encryptedPassword = ""
	var error error

	if p.UserTypeId == 15 {
		return sqlResult, errors.New("you can't update a user to become an Admin")
	}

	if userType == "Admin" {
		if p.Password != "" {
			sql = `
				UPDATE
				  user
				SET
				  firstname = ?,
				  lastname = ?,
				  email = ?,
				  mobile_number = ?,
				  user_type_id = ?,
				  last_updated = NOW(),
				  bank_type_id = ?,
				  bank_no = ?,
				  address = ?,
				  birthday = ?,
				  gender = ?,
				  m88_account = ?,
				  password = ?
				WHERE id = ?;
			`

			encryptedPassword, error = UtilitiesString.Encrypt(p.Password)

			if error != nil {
				return sqlResult, error
			}

			sqlResult, error = database.DBInstancePublic.Exec(
				sql,
				p.FirstName,
				p.LastName,
				p.Email,
				p.MobileNumber,
				p.UserTypeId,
				p.BankTypeId,
				p.BankNo,
				p.Address,
				p.Birthday,
				p.Gender,
				p.M88Account,
				encryptedPassword,
				p.ID,
			)
		} else {
			sql = `
				UPDATE
				  user
				SET
				  firstname = ?,
				  lastname = ?,
				  email = ?,
				  mobile_number = ?,
				  user_type_id = ?,
				  last_updated = NOW(),
				  bank_type_id = ?,
				  bank_no = ?,
				  address = ?,
				  birthday = ?,
				  gender = ?,
				  m88_account = ?
				WHERE id = ?;
			`

			sqlResult, error = database.DBInstancePublic.Exec(
				sql,
				p.FirstName,
				p.LastName,
				p.Email,
				p.MobileNumber,
				p.UserTypeId,
				p.BankTypeId,
				p.BankNo,
				p.Address,
				p.Birthday,
				p.Gender,
				p.M88Account,
				p.ID,
			)
		}
	}

	/**
		These fields cannot be edited by Dropshippers and sellers
		1. Firstname
		2. Lastname
		3. UserType
		4. Birthday
	 */

	if userType == "Dropshipper" || userType == "Seller" || userType == "Rider" {
		fmt.Println("Knocked him out")
		if p.Password != "" {
			sql = `
				UPDATE
				  user
				SET
				  mobile_number = ?,
				  last_updated = NOW(),
				  bank_type_id = ?,
				  bank_no = ?,
				  address = ?,
				  gender = ?,
				  m88_account = ?,
				  password = ?
				WHERE id = ?;
			`

			encryptedPassword, error = UtilitiesString.Encrypt(p.Password)

			if error != nil {
				return sqlResult, error
			}

			sqlResult, error = database.DBInstancePublic.Exec(
				sql,
				p.MobileNumber,
				p.BankTypeId,
				p.BankNo,
				p.Address,
				p.Gender,
				p.M88Account,
				encryptedPassword,
				p.ID,
			)
		} else {
			sql = `
				UPDATE
				  user
				SET
				  mobile_number = ?,
				  last_updated = NOW(),
				  bank_type_id = ?,
				  bank_no = ?,
				  address = ?,
				  gender = ?,
				  m88_account = ?
				WHERE id = ?;
			`

			sqlResult, error = database.DBInstancePublic.Exec(
				sql,
				p.MobileNumber,
				p.BankTypeId,
				p.BankNo,
				p.Address,
				p.Gender,
				p.M88Account,
				p.ID,
			)
		}
	}


	return sqlResult, error
}

func (u *User) ValidateToken(role string) (int, error) {
	var userId int

	// Attempt to parse
	decrypted, err := UtilitiesString.ParseJWT(u.Token)

	if err != nil {
		return 0, errors.New("error parsing the JWT token")
	}

	// Ensure timestamp is valid for (X) hours depending on your configuration
	err = godotenv.Load()

	if err != nil {
		return 0, errors.New("error opening .env file")
	}

	var tokenExpiry = os.Getenv("TOKEN_EXPIRY")
	tokenExpiryInt, err := strconv.Atoi(tokenExpiry)

	if err != nil {
		return 0, errors.New("error checking token expiry settings")
	}

	timestamp := fmt.Sprintf("%v", decrypted["timestamp"])
	timestampFormat, err := time.Parse("2006-01-02 15:04:05", timestamp)

	if err != nil {
		return 0, errors.New("error comparing times")
	}

	hourDifference := timestampFormat.Sub(time.Now()).Hours() // time.Now().Sub(timestampFormat).Hours()

	if hourDifference > float64(tokenExpiryInt) {
		return 0, errors.New("token expired")
	}

	// Check if valid user and the roles match
	userId, err = strconv.Atoi(fmt.Sprintf("%v", decrypted["userId"]))

	if err != nil {
		return 0, errors.New("something went wrong when trying to convert the userId to string")
	}

	sql := `
		SELECT 
			COUNT(*) AS user_count	
		FROM user u
		INNER JOIN user_type ut
			ON 1 = 1
				AND u.user_type_id = ut.id
		WHERE 1 = 1
			AND u.id = ?	
			AND ut.name = ?
			AND u.is_active = 1
	`

	var (
		Count int
	)

	err = database.DBInstancePublic.Get(&Count, sql, userId, role)

	if err != nil {
		return 0, errors.New("something went wrong when trying to check the user roles")
	}

	if Count == 0 {
		return 0, errors.New("====== user does not have the correct privileges and/or is disabled")
	}

	fmt.Println("userId", userId)

	return userId, nil
}

func (u *User) ValidateTokenV2(roles []string) (int, string, error) {

	// These user variables will be populated then returned at
	// the end of this function
	var userId int

	// Attempt to parse JWT
	decrypted, err := UtilitiesString.ParseJWT(u.Token)

	if err != nil {
		return 0, "", errors.New("error parsing the JWT token")
	}

	// Ensure timestamp is valid for (X) hours depending on your configuration
	err = godotenv.Load()

	if err != nil {
		return 0, "", errors.New("error opening .env file")
	}

	var tokenExpiry = os.Getenv("TOKEN_EXPIRY")
	tokenExpiryInt, err := strconv.Atoi(tokenExpiry)

	if err != nil {
		return 0, "", errors.New("error checking token expiry settings")
	}

	timestamp := fmt.Sprintf("%v", decrypted["timestamp"])
	timestampFormat, err := time.Parse("2006-01-02 15:04:05", timestamp)

	if err != nil {
		return 0, "", errors.New("error comparing times")
	}

	hourDifference := timestampFormat.Sub(time.Now()).Hours() // time.Now().Sub(timestampFormat).Hours()

	if hourDifference > float64(tokenExpiryInt) {
		return 0, "", errors.New("token expired")
	}

	// Check if valid user and the roles match
	userId, err = strconv.Atoi(fmt.Sprintf("%v", decrypted["userId"]))

	if err != nil {
		return 0, "", errors.New("something went wrong when trying to convert the userId to string")
	}

	sql := `
		SELECT 
			COUNT(*) AS user_count,
			IF(ut.name IS NULL, "", ut.name) AS user_type
		FROM `+ "`user`" + ` u
		INNER JOIN user_type ut
			ON 1 = 1
				AND u.user_type_id = ut.id
		WHERE 1 = 1
			AND u.id = ?	
			AND ut.name IN (USER_TYPES)
			AND u.is_active = 1
	`

	sql = strings.Replace(sql, "USER_TYPES", UtilitiesString.GetJoinedStringForWhereIn(roles), -1)

	var userMiddlewareDetails UserMiddlewareDetails

	err = database.DBInstancePublic.Get(
		&userMiddlewareDetails,
		sql,
		userId,
	)

	if err != nil {
		fmt.Println("HERE")
		return 0, "", err
	}

	if userMiddlewareDetails.UserCount == 0 {
		return 0, "", errors.New("========= user does not have the correct privileges and/or is disabled")
	}

	return userId, userMiddlewareDetails.UserType, nil
}

// Pagination alpha - to be used in all of it , these are nice to haves.
func (u *User) GetAll2(page int, rows int) (*[]UserListDisplay, response_builder.Pagination, error) {

	var container []UserListDisplay
	var pagination response_builder.Pagination

	sql := `
		SELECT 
			u.id,
			IF(u.firstname IS NULL, "", u.firstname) AS firstname,
			IF(u.lastname IS NULL, "", u.lastname) AS lastname,
			IF(u.email IS NULL, "", u.email) AS email,
			IF(u.mobile_number IS NULL, "", u.mobile_number) AS mobile_number,
			IF(ut.name IS NULL, "", ut.name) AS role,
			IF(bt.name IS NULL, "", bt.name) AS bank_type,
			IF(u.bank_no IS NULL, "", u.bank_no) AS bank_no,
			IF(u.address IS NULL, "", u.address) AS address,
			IF(u.birthday IS NULL, "", u.birthday) AS birthday,
			IF(u.gender IS NULL, "", u.gender) AS gender,
			IF(u.m88_account IS NULL, "", u.m88_account) AS m88_account
		FROM user u 
		INNER JOIN user_type ut 
			ON 1 = 1
				AND u.user_type_id = ut.id
		INNER JOIN bank_type bt 
			ON 1 = 1
				AND u.bank_type_id = bt.id 
		WHERE 1 = 1     
			and u.is_active = 1
	`

	paginate := func () (*[]UserListDisplay, response_builder.Pagination, error) {
		// Execute count
		// This is an example if you want to pass args
		// count, err := database.GetQueryCount(sql, 26)

		// Get count first
		count, err := database.GetQueryCount(sql)

		// Fail error if error
		if err != nil {
			return &container, pagination, err
		}

		// Just return blank entries if there's no count (save operations)
		if count == 0 {
			pagination.Pages = make([]int, 0)
			return &container, pagination, nil
		}

		rows = 1000

		sql, pages, rowsPerPage, offset, page, totalCount, resultCount := database.GetPaginationDetails(
			sql,
			count,
			page,
			rows,
			1000,
		)

		pagination.SetData(rowsPerPage, offset, pages, rows, page, totalCount, resultCount)

		// Attempt to perform
		err = database.DBInstancePublic.Select(&container, sql)

		if err != nil {
			return &container, pagination, err
		}

		return &container, pagination, nil
	}

	res, pagination, err := paginate()

	return res, pagination, err
}

func (u *User) GetAll() (*[]UserListDisplay, error) {
	var userListDisplay []UserListDisplay

	sql := `
		SELECT 
			u.id,
			IF(u.firstname IS NULL, "", u.firstname) AS firstname,
			IF(u.lastname IS NULL, "", u.lastname) AS lastname,
			IF(u.email IS NULL, "", u.email) AS email,
			IF(u.mobile_number IS NULL, "", u.mobile_number) AS mobile_number,
			IF(ut.name IS NULL, "", ut.name) AS role,
			IF(bt.name IS NULL, "", bt.name) AS bank_type,
			IF(u.bank_no IS NULL, "", u.bank_no) AS bank_no,
			IF(u.address IS NULL, "", u.address) AS address,
			IF(u.birthday IS NULL, "", u.birthday) AS birthday,
			IF(u.gender IS NULL, "", u.gender) AS gender,
			IF(u.m88_account IS NULL, "", u.m88_account) AS m88_account
		FROM user u 
		INNER JOIN user_type ut 
			ON 1 = 1
				AND u.user_type_id = ut.id
		INNER JOIN bank_type bt 
			ON 1 = 1
				AND u.bank_type_id = bt.id 
		WHERE 1 = 1     
			and u.is_active = 1
	`

	err := database.DBInstancePublic.Select(&userListDisplay, sql)

	if err != nil {
		return &userListDisplay, err
	}

	return &userListDisplay, err
}

func (u *User) GetAllByUserType(userType string) (*[]UserListDisplay, error) {
	var userListDisplay []UserListDisplay

	sql := `
		SELECT 
			u.id,
			IF(u.firstname IS NULL, "", u.firstname) AS firstname,
			IF(u.lastname IS NULL, "", u.lastname) AS lastname,
			IF(u.email IS NULL, "", u.email) AS email,
			IF(u.mobile_number IS NULL, "", u.mobile_number) AS mobile_number,
			IF(ut.name IS NULL, "", ut.name) AS role,
			IF(bt.name IS NULL, "", bt.name) AS bank_type,
			IF(u.bank_no IS NULL, "", u.bank_no) AS bank_no,
			IF(u.address IS NULL, "", u.address) AS address,
			IF(u.birthday IS NULL, "", u.birthday) AS birthday,
			IF(u.gender IS NULL, "", u.gender) AS gender,
			IF(u.m88_account IS NULL, "", u.m88_account) AS m88_account
		FROM user u 
		INNER JOIN user_type ut 
			ON 1 = 1
				AND u.user_type_id = ut.id
		INNER JOIN bank_type bt 
			ON 1 = 1
				AND u.bank_type_id = bt.id 
		WHERE 1 = 1     
			and u.is_active = 1
			and ut.name = ?
		ORDER BY u.lastname ASC, u.firstname ASC
	`

	err := database.DBInstancePublic.Select(&userListDisplay, sql, userType)

	if err != nil {
		return &userListDisplay, err
	}

	return &userListDisplay, err
}

func (u *User) GetLoginDetails() (ResponseLoginUserInfo, error) {
	var responseLoginUserInfo ResponseLoginUserInfo

	sql := `
		SELECT 
			u.id,
			u.firstname,
			u.lastname,
			IF(u.law_firm_ref_id IS NULL, 0, u.law_firm_ref_id) AS law_firm_id,
			ut.name AS role
		FROM
		user u
		INNER JOIN user_type ut 
			ON 1 = 1
				AND u.user_type_id = ut.id 
		WHERE 1 = 1
			AND email = ?
	`

	err := database.DBInstancePublic.Get(&responseLoginUserInfo, sql, u.Email)

	return responseLoginUserInfo, err
}

func (u *User) GetPasswordAndIdByEmail() (string, int, error) {
	var user User

	sql := `
		SELECT 
			id, 
			password
		FROM user 
		WHERE 1 = 1
			AND email = ?
	`

	fmt.Println("u.Email", u.Email)

	err := database.DBInstancePublic.Get(&user, sql, u.Email)

	return user.Password, user.ID, err
}

func (u *UserLogin) IsUniqueEmail() (bool, error) {
	var isValidEmail bool

	exists, err := database.ValidEntry(u.Email, "email", "user")

	if err != nil {
		return isValidEmail, err
	}

	isValidEmail = exists

	return isValidEmail, nil
}

func (u *UserLogin) ValidEmail() (bool, error) {
	var isValidEmail bool

	exists, err := database.ValidEntry(u.Email, "email", "user")

	if err != nil {
		return isValidEmail, err
	}

	isValidEmail = exists

	return isValidEmail, nil
}

func (u *User) ValidId() (bool, error) {
	hasId := false
	sql := `
		SELECT 
			IF(COUNT(*) > 0, true, false) AS has_id 
		FROM user 
		WHERE 1 = 1
			AND id = ? 
			AND is_active = 1
	`

	err := database.DBInstancePublic.Get(&hasId, sql, u.ID)

	return hasId, err
}

func (u *User) Delete() (sql.Result, error) {
	sql := `
		UPDATE user  
		SET is_active = 0 
		WHERE id = ?
	`

	res, err := database.DBInstancePublic.Exec(sql, u.ID)

	return res, err
}

func (p *ParamsUpdate) NoEmptyFields() []string {

	var emptyFields []string

	// Fuck that reflection, I'm doing this manually
	if p.ID == 0 {
		emptyFields = append(emptyFields, "id is empty")
	}

	if p.FirstName == "" {
		emptyFields = append(emptyFields, "firstname is empty")
	}

	if p.LastName == "" {
		emptyFields = append(emptyFields, "lastname is empty")
	}

	if p.MobileNumber == "" {
		emptyFields = append(emptyFields, "mobile_number is empty")
	}

	if p.Password == "" {
		// emptyFields = append(emptyFields, "password is empty")
	}

	if p.UserTypeId == 0 {
		emptyFields = append(emptyFields, "user_type_id is empty")
	}

	if p.BankTypeId == 0 {
		emptyFields = append(emptyFields, "bank_type_id is empty")
	}

	if p.BankNo == "" {
		emptyFields = append(emptyFields, "bank_no is empty")
	}

	if p.Address == "" {
		emptyFields = append(emptyFields, "address is empty")
	}

	if p.Birthday == "" {
		emptyFields = append(emptyFields, "birthday is empty")
	}

	if p.M88Account == "" {
		emptyFields = append(emptyFields, "m88_account is empty")
	}

	if p.Gender == "" {
		emptyFields = append(emptyFields, "gender is empty")
	}

	return emptyFields
}

func (p *ParamsInsert) NoEmptyFields() []string {

	var emptyFields []string

	if p.FirstName == "" {
		emptyFields = append(emptyFields, "firstname is empty")
	}

	if p.LastName == "" {
		emptyFields = append(emptyFields, "lastname is empty")
	}

	if p.Email == "" {
		emptyFields = append(emptyFields, "lastname is empty")
	}

	if p.MobileNumber == "" {
		emptyFields = append(emptyFields, "mobile_number is empty")
	}

	if p.Password == "" {
		emptyFields = append(emptyFields, "password is empty")
	}

	if p.UserTypeId == 0 {
		emptyFields = append(emptyFields, "user_type_id is empty")
	}

	if p.BankNo == "" {
		emptyFields = append(emptyFields, "bank_no is empty")
	}

	if p.Address == "" {
		emptyFields = append(emptyFields, "address is empty")
	}

	if p.Birthday == "" {
		emptyFields = append(emptyFields, "birthday is empty")
	}

	if p.Gender == "" {
		emptyFields = append(emptyFields, "gender is empty")
	}

	return emptyFields
}

func InsertToSAmpleTable() {

}

// Checks whether law firm contact can be added because it's law firm is
// representing the case.
func (u *User) ValidateLawFirmContactId(e bool) error {
	var count int
	var sql string

	sql = `
		SELECT 
			COUNT(*) AS count 
		FROM ` + UtilitiesString.EncloseString("case", "`") + ` c
		WHERE 1 = 1
			AND c.id = ?
			AND c.law_firm_ref_id = (
				SELECT law_firm_ref_id
				FROM ` + UtilitiesString.EncloseString("user", "`") + ` u
				WHERE 1 = 1
					AND u.id = ?
			)
	`

	err := database.DBInstancePublic.Get(&count, sql, u.CaseId, u.ID)

	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("no user id found for law firm")
	}

	return nil
}

func (u *User) GetDetailsByEmail(e string) (*ResponseLoginUserInfo, error) {
	var responseLoginUserInfo ResponseLoginUserInfo

	sql := `
		SELECT 
			u.id,
			u.firstname,
			u.lastname,
			u.email,
			ut.name AS role,
			IF(u.organization_ref_id IS NULL, 0, u.organization_ref_id) AS organization_id,
			IF(u.law_firm_ref_id IS NULL, "", u.law_firm_ref_id) AS law_firm_id
		FROM
		user u
		INNER JOIN user_type ut 
			ON 1 = 1
				AND u.user_type_id = ut.id 
		WHERE 1 = 1
			AND u.email = ?
	`

	fmt.Println("email ==================", e)

	err := database.DBInstancePublic.Get(&responseLoginUserInfo, sql, e)

	return &responseLoginUserInfo, err
}