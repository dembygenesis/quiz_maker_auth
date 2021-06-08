package email_notification

import (
	"fmt"
	"github.com/dembygenesis/quiz_maker_auth/src/app/common/config"
	UtilitiesString "github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/string_utils"
	"github.com/dembygenesis/quiz_maker_auth/src/v1/api/database"
	"github.com/dembygenesis/quiz_maker_auth/src/v1/api/models/email"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

type EmailNotificationInterface interface {
	ExecuteCron()
}

type emailNotificationService struct {
}

type EmailDetails struct {
	ContactFirstName string
	Organization     string
	Cases            []struct {
		CaseId    int
		FirstName string
		LastName  string
	}
}

type Organization struct {
	Id int `db:"id"`
}

type CaseStatusInformation struct {
	// Case information
	CaseId                    int    `db:"case_id"`
	CasePatientFirstName      string `db:"case_patient_first_name"`
	CasePatientLastName       string `db:"case_patient_last_name"`
	CaseStatusDaysLastUpdated int    `db:"case_status_days_last_updated"`
	CaseOrganization          string `db:"case_organization"`

	// User information
	CaseContactId        int    `db:"case_contact_id"`
	CaseContactFirstName string `db:"case_contact_first_name"`
	CaseContactLastName  string `db:"case_contact_last_name"`
	CaseContactEmail     string `db:"case_contact_email"`
}

type LienManagerCaseStatusInformation struct {
	// Case information
	CaseId                         int    `db:"case_id"`
	CasePatientFirstName           string `db:"case_patient_first_name"`
	CasePatientLastName            string `db:"case_patient_last_name"`
	CaseStatusDaysLastUpdated      int    `db:"case_status_days_last_updated"`
	TreatmentStatusDaysLastUpdated int    `db:"treatment_status_days_last_updated"`
	CaseOrganization               string `db:"case_organization"`

	// User information
	UserId        int    `db:"user_id"`
	UserFirstName string `db:"user_first_name"`
	UserLastName  string `db:"user_last_name"`
}

func NewEmailNotificationService() EmailNotificationInterface {
	return &emailNotificationService{}
}

func printCronEntries(cronEntries []*cron.Entry) {
	log.Infof("Cron Info: %+v\n", cronEntries)
}

func cronExample() {
	log.Info("Create new cron")
	c := cron.New()
	c.AddFunc("*/1 * * * *", func() { log.Info("[Job 1]Every minute job\n") })

	// Start cron with one scheduled job
	log.Info("Start cron")
	c.Start()
	printCronEntries(c.Entries())
	time.Sleep(2 * time.Minute)

	// Funcs may also be added to a running Cron
	log.Info("Add new job to a running cron")
	_ = c.AddFunc("*/2 * * * *", func() { log.Info("[Job 2]Every two minutes job\n") })
	printCronEntries(c.Entries())
	time.Sleep(5 * time.Minute)

	//Remove Job2 and add new Job2 that run every 1 minute
	log.Info("Remove Job2 and add new Job2 with schedule run every minute")

	c.AddFunc("*/1 * * * *", func() { log.Info("[Job 2]Every one minute job\n") })
	time.Sleep(5 * time.Minute)
}

func (s *emailNotificationService) getOrganizations(t *sqlx.Tx) (*[]Organization, error) {
	var organizations []Organization
	sql := "SELECT id FROM organization"
	err := t.Select(&organizations, sql)
	if err != nil {
		return &organizations, err
	}
	return &organizations, nil
}

func (s *emailNotificationService) sendFollowUpNotifications() {
	t := database.DBInstancePublic.MustBegin()

	var err error
	organizations, err := s.getOrganizations(t)
	if err != nil {
		fmt.Println("has errors fetching the organizations: " + err.Error())
	}

	// Get all active cases under that organization, and send emails to the
	// admins, org members, and lawyers where cases have no updates in the last 30 days
	for _, v := range *organizations {
		s.sendNotifications(t, v.Id, "Law Firm")
		s.sendNotifications(t, v.Id, "Admin")
	}
}

func (s *emailNotificationService) sendNotifications(t *sqlx.Tx, organizationId int, userType string) {
	var sql string
	var err error

	// Fetch notification data
	var caseStatusInformation []CaseStatusInformation
	if userType == "Law Firm" {
		sql = `
			SELECT
			  c.id AS case_id,
			  c.patient_first_name AS case_patient_first_name,
			  c.patient_last_name AS case_patient_last_name,
			  DATEDIFF(
				CURRENT_TIMESTAMP,
				c.case_status_last_updated
			  ) AS case_status_days_last_updated,
			  o.name AS case_organization,
			  cc.case_contact_user_ref_id AS case_contact_id,
			  u.firstname AS case_contact_first_name,
			  u.lastname AS case_contact_last_name,
			  u.email AS case_contact_email
			FROM
			  ` + UtilitiesString.EncloseString("case", "`") + ` c
			  INNER JOIN organization o
				ON 1 = 1
				AND c.organization_ref_id = o.id
			  LEFT JOIN case_contact cc
				ON 1 = 1
				AND c.id = cc.case_ref_id
			  INNER JOIN user u
				ON 1 = 1
				AND cc.case_contact_user_ref_id = u.id
			WHERE 1 = 1
			  AND c.organization_ref_id = ?
			  AND c.is_active = 1
			HAVING case_status_days_last_updated >= 1
		`
	}
	if userType == "Admin" || userType == "Organization Member" {
		sql = `
			SELECT
			  c.id AS case_id,
			  c.patient_first_name AS case_patient_first_name,
			  c.patient_last_name AS case_patient_last_name,
			  DATEDIFF(
				CURRENT_TIMESTAMP,
				c.treatment_status_last_updated
			  ) AS case_status_days_last_updated,
			  o.name AS case_organization,
			  c.created_by AS case_contact_id,
			  u.firstname AS case_contact_first_name,
			  u.lastname AS case_contact_last_name,
			  u.email AS case_contact_email
			FROM
			  ` + UtilitiesString.EncloseString("case", "`") + ` c
			  INNER JOIN organization o
				ON 1 = 1
				AND c.organization_ref_id = o.id
			  INNER JOIN user u
				ON 1 = 1
				AND c.created_by = u.id
			WHERE 1 = 1
			  AND c.organization_ref_id = ?
			  AND c.is_active = 1
			HAVING case_status_days_last_updated >= 1
		`
	}
	err = t.Select(&caseStatusInformation, sql, organizationId)
	if err != nil {
		fmt.Println("error in ubsendNotifications: " + err.Error())
		return
	}

	// Build email struct from notification data
	emailDetails := make(map[int]*EmailDetails)
	for _, v := range caseStatusInformation {
		if emailDetails[v.CaseContactId] == nil {
			emailDetails[v.CaseContactId] = &EmailDetails{
				ContactFirstName: v.CaseContactFirstName,
				Organization:     v.CaseOrganization,
			}
		}
		emailDetails[v.CaseContactId].Cases = append(emailDetails[v.CaseContactId].Cases, struct {
			CaseId    int
			FirstName string
			LastName  string
		}{CaseId: v.CaseId, FirstName: v.CasePatientFirstName, LastName: v.CasePatientLastName})
	}

	// Throttle
	concurrencyLimit := 5
	limiter := make(chan string, concurrencyLimit)

	for _, v := range emailDetails {
		var notificationMsg string
		notificationMsg += fmt.Sprintf("Hi %v,<br/><br/>", v.ContactFirstName)
		notificationMsg += "Please provide an update for the following cases: <br/><br/>"
		for _, y := range v.Cases {
			url := config.BaseUrl + "/law-firm/cases/" + strconv.Itoa(y.CaseId)
			notificationMsg += fmt.Sprintf("<a style=\"#0ad26a\" href=\"%v\">%v</a><br/>", url, y.FirstName+" "+y.LastName)
		}
		notificationMsg += fmt.Sprintf("<br/>Requested by: %v<br/><br/>", v.Organization)
		notificationMsg += "<br/><br/>Best,<br/>The MediLegalRecords team"

		// Throttle email sending
		limiter <- notificationMsg

		go func() {
			defer func() {
				<-limiter
			}()
			err = email.SendMail("dembygenesis@gmail.com", "Case Updates", notificationMsg)
			if err != nil {
				fmt.Println("==================== FAILED TO SEND MESSAGE: " + err.Error())
			} else {
				fmt.Println("Sent message")
			}
		}()
	}
}

func (s *emailNotificationService) sendAdminNotifications(caseId int) {

}

func (s *emailNotificationService) ExecuteCron() {
	log.Info("=============== Initialize CRON ===============")
	c := cron.New()
	if err := c.AddFunc("* * */23 * *", func() {
		s.sendFollowUpNotifications()
	}); err == nil {
		c.Start()
	}
}

