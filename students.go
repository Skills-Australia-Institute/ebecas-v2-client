package ebecasv2client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Student represents an eBECAS student.
type Student struct {
	ID                      int64             `json:"id"`
	StudentNumber           string            `json:"student_number"`
	Name                    string            `json:"name"`
	FirstName               string            `json:"first_name"`
	MiddleName              string            `json:"middle_name"`
	LastName                string            `json:"last_name"`
	DOB                     string            `json:"dob"`
	LocalEmail              string            `json:"local_email"`
	LocalPhone              string            `json:"local_phone"`
	LocalMobile             string            `json:"local_mobile"`
	OverseasEmail           string            `json:"os_email"`
	OverseasPhone           string            `json:"os_phone"`
	OverseasMobile          string            `json:"os_mobile"`
	OverseasAddress         Address           `json:"os_address"`
	LocalAddress            Address           `json:"local_address"`
	PostalAddress           Address           `json:"postal_address"`
	Age                     int               `json:"age"`
	Country                 Country           `json:"country"`
	Gender                  Gender            `json:"gender"`
	Title                   Title             `json:"title"`
	OverseasStudent         bool              `json:"os_student"`
	OnlineStudent           bool              `json:"online_student"`
	Language                Language          `json:"language"`
	CitizenCountry          Country           `json:"citizen_country"`
	BirthCountry            Country           `json:"birth_country"`
	PassportCode            string            `json:"passport_code"`
	PassportTypeCode        string            `json:"passport_type_code"`
	PassportIssuer          string            `json:"passport_issuer"`
	PassportIssueDate       string            `json:"passport_issue_date"`
	PassportExpiryDate      string            `json:"passport_expiry_date"`
	PassportCountry         Country           `json:"passport_country"`
	PassportPlaceBirth      string            `json:"passport_place_birth"`
	MedicalNotes            string            `json:"medical_notes"`
	MedicalCategories       []MedicalCategory `json:"medical_categories"`
	Notes                   string            `json:"notes"`
	CreatedAt               time.Time         `json:"created_at"`
	CreatedBy               CreatedBy         `json:"created_by"`
	PreferredName           string            `json:"preferred_name"`
	Pronouns                Pronouns          `json:"pronouns"`
	StudentLocation         *string           `json:"student_location"`
	StudentLocationFrom     *string           `json:"student_loc_from"`
	StudentLocationTo       *string           `json:"student_loc_to"`
	Locations               []Location        `json:"locations"`
	StudentPhotoFlag        bool              `json:"student_photo_flag"`
	PortalRegisteredFlag    bool              `json:"portal_registered_flag"`
	PortalPrivacyAcceptFlag bool              `json:"portal_privacy_accept_flag"`
	PortalPrivacyAcceptDate string            `json:"portal_privacy_accept_date"`
	StudentSupportFlag      *string           `json:"student_support_flag"`
	StudentTest             StudentTest       `json:"student_test"`
}

// Address represents an eBECAS address.
type Address struct {
	Line1    string  `json:"line1"`
	Line2    string  `json:"line2"`
	Line3    string  `json:"line3"`
	Suburb   string  `json:"suburb"`
	State    string  `json:"state"`
	Postcode string  `json:"postcode"`
	Country  Country `json:"country"`
}

// Country represents an eBECAS country.
type Country struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// Gender represents an eBECAS gender.
type Gender struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// Title represents an eBECAS title.
type Title struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// Language represents an eBECAS language.
type Language struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// MedicalCategory represents an eBECAS medical category.
type MedicalCategory struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// Pronouns represents a person's pronouns.
type Pronouns struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// Location represents an eBECAS location.
type Location struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// StudentTest represents a student's test details.
type StudentTest struct {
	TestType  TestType `json:"test_type"`
	TestScore string   `json:"test_score"`
	TestDate  string   `json:"test_date"`
}

// TestType represents a student's test type.
type TestType struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// StudentFile represents a file attached to an eBECAS student.
type StudentFile struct {
	ID           int64            `json:"id"`
	Name         string           `json:"name"`
	DocumentName string           `json:"document_name"`
	Description  string           `json:"description"`
	FileSize     int64            `json:"file_size"`
	DocumentType DocumentType     `json:"document_type"`
	CreatedAt    time.Time        `json:"created_at"`
	CreatedBy    CreatedBy        `json:"created_by"`
	Diary        *string          `json:"diary"`
	Type         StudentFileType  `json:"type"`
	Category     DocumentCategory `json:"category"`
	Object       DocumentObject   `json:"object"`
}

// DocumentType represents an eBECAS document type.
type DocumentType struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// StudentFileType represents an eBECAS student file type.
type StudentFileType struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// DocumentCategory represents an eBECAS document category.
type DocumentCategory struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// DocumentObject represents the eBECAS object associated with a document.
type DocumentObject struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// StudentContact represents student's contact.
type StudentContact struct {
	ID          int64       `json:"id"`
	ContactType ContactType `json:"contact_type"`
	Name        string      `json:"name"`
	ContactInfo string      `json:"contact_info"`
	Address     Address     `json:"address"`
	Mobile      string      `json:"mobile"`
	Phone       string      `json:"phone"`
	Email       string      `json:"email"`
	CreatedAt   time.Time   `json:"created_at"`
	CreatedBy   CreatedBy   `json:"created_by"`
}

// ContactType represents eBECAS contact type.
type ContactType struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

// GetStudentContacts retrieves student's contacts.
func (c *Client) GetStudentContacts(ctx context.Context, studentID int64) ([]StudentContact, int, error) {
	if studentID <= 0 {
		return nil, http.StatusBadRequest, fmt.Errorf("student ID must be greater than zero")
	}

	requestURL := fmt.Sprintf("%s/students/%d/contacts", c.baseURL, studentID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(
			"create get contacts request for student %d: %w",
			studentID,
			err,
		)
	}

	data, statusCode, err := c.do(req)
	if err != nil {
		return nil, statusCode, fmt.Errorf("get contacts for student %d, status %d: %w", studentID, statusCode, err)
	}

	if statusCode != http.StatusOK {
		return nil, statusCode,
			fmt.Errorf(
				"get contacts for student %d returned status %d: %s",
				studentID,
				statusCode,
				strings.TrimSpace(string(data)),
			)
	}

	contacts := make([]StudentContact, 0)
	if err := json.Unmarshal(data, &contacts); err != nil {
		return nil, statusCode,
			fmt.Errorf("decode contacts response for student %d: %w", studentID, err)
	}

	return contacts, statusCode, nil
}

type CreateStudentContactInput struct {
	ContactType int64  `json:"contact_type"`
	Name        string `json:"name"`
	ContactInfo string `json:"contact_info"`
	Mobile      string `json:"mobile"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     struct {
		Line1    string `json:"line1"`
		Line2    string `json:"line2"`
		Line3    string `json:"line3"`
		State    string `json:"state"`
		Postcode string `json:"postcode"`
		Country  int64  `json:"country"`
	} `json:"address"`
}

// CreateStudentContact creates student's contact.
func (c *Client) CreateStudentContact(
	ctx context.Context,
	studentID int64,
	input CreateStudentContactInput,
) (StudentContact, int, error) {
	var contact StudentContact

	if studentID <= 0 {
		return contact, http.StatusBadRequest, fmt.Errorf("student ID must be greater than zero")
	}

	if input.ContactType <= 0 {
		return contact, http.StatusBadRequest, fmt.Errorf("contact type must be greater than zero")
	}

	if input.Name == "" {
		return contact, http.StatusBadRequest, fmt.Errorf("contact name is required")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return contact, http.StatusInternalServerError, fmt.Errorf(
			"encode contact request for student %d: %w",
			studentID,
			err,
		)
	}

	requestURL := fmt.Sprintf("%s/students/%d/contacts", c.baseURL, studentID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewReader(body))
	if err != nil {
		return contact, http.StatusInternalServerError, fmt.Errorf(
			"create contact request for student %d: %w",
			studentID,
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")

	data, statusCode, err := c.do(req)
	if err != nil {
		return contact, statusCode, fmt.Errorf(
			"create contact for student %d, status %d: %w",
			studentID,
			statusCode,
			err,
		)
	}

	if statusCode != http.StatusOK {
		return contact, statusCode,
			fmt.Errorf(
				"create contact for student %d returned status %d: %s",
				studentID,
				statusCode,
				strings.TrimSpace(string(data)),
			)
	}

	if err := json.Unmarshal(data, &contact); err != nil {
		return contact, statusCode, fmt.Errorf(
			"decode contact response for student %d: %w",
			studentID,
			err,
		)
	}

	return contact, statusCode, nil
}

type UpdateStudentTestInput struct {
	TestType  int64  `json:"test_type"`
	TestScore string `json:"test_score"`
	TestDate  string `json:"test_date"`
}

// UpdateStudentTest updates student's pre-entry test.
func (c *Client) UpdateStudentTest(
	ctx context.Context,
	studentID int64,
	input UpdateStudentTestInput,
) (Student, int, error) {
	var student Student

	if studentID <= 0 {
		return student, http.StatusBadRequest, fmt.Errorf("student ID must be greater than zero")
	}

	if input.TestType <= 0 {
		return student, http.StatusBadRequest, fmt.Errorf("test type must be greater than zero")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return student, http.StatusInternalServerError, fmt.Errorf(
			"encode test request for student %d: %w",
			studentID,
			err,
		)
	}

	requestURL := fmt.Sprintf("%s/students/%d/student-test", c.baseURL, studentID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, requestURL, bytes.NewReader(body))
	if err != nil {
		return student, http.StatusInternalServerError, fmt.Errorf(
			"create update test request for student %d: %w",
			studentID,
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")

	data, statusCode, err := c.do(req)
	if err != nil {
		return student, statusCode, fmt.Errorf(
			"update test for student %d, status %d: %w",
			studentID,
			statusCode,
			err,
		)
	}

	if statusCode != http.StatusOK {
		return student, statusCode,
			fmt.Errorf(
				"update test for student %d returned status %d: %s",
				studentID,
				statusCode,
				strings.TrimSpace(string(data)),
			)
	}

	if err := json.Unmarshal(data, &student); err != nil {
		return student, statusCode, fmt.Errorf(
			"decode test response for student %d: %w",
			studentID,
			err,
		)
	}

	return student, statusCode, nil
}

type SearchStudentsInput struct {
	Q          string         `json:"q"`
	Fields     []StudentField `json:"fields"`
	Sort       StudentSort    `json:"sort"`
	Filter     StudentFilter  `json:"filter"`
	Pagination Pagination     `json:"pagination"`
}

type SearchStudentsResponse struct {
	Data []struct {
		ID                    int64  `json:"id"`
		StudentsFirstName     string `json:"students-first_name"`
		StudentsLastName      string `json:"students-last_name"`
		StudentsStudentNumber string `json:"students_student_number"`
	} `json:"data"`
	Meta struct {
		Pagination Pagination `json:"pagination"`
	} `json:"meta"`
}

// SearchStudents retrieves students using the provided search criteria.
func (c *Client) SearchStudents(
	ctx context.Context,
	input SearchStudentsInput,
) (SearchStudentsResponse, int, error) {
	var response SearchStudentsResponse

	body, err := json.Marshal(input)
	if err != nil {
		return response, http.StatusInternalServerError, fmt.Errorf("encode students search request: %w", err)
	}

	requestURL := fmt.Sprintf("%s/students/search", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewReader(body))
	if err != nil {
		return response, http.StatusInternalServerError, fmt.Errorf("create students search request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	data, statusCode, err := c.do(req)
	if err != nil {
		return response, statusCode, fmt.Errorf("search students, status %d: %w", statusCode, err)
	}

	if statusCode != http.StatusOK {
		return response, statusCode, fmt.Errorf(
			"search students returned status %d: %s",
			statusCode,
			strings.TrimSpace(string(data)),
		)
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return response, statusCode, fmt.Errorf("decode search students response: %w", err)
	}

	return response, statusCode, nil
}

type CreateStudentInput struct {
	StudentNumber      string       `json:"student_number"`
	FirstName          string       `json:"first_name"`
	MiddleName         string       `json:"middle_name"`
	LastName           string       `json:"last_name"`
	DOB                string       `json:"dob"`
	Gender             int64        `json:"gender"`
	Country            int64        `json:"country"`
	PreferredName      string       `json:"preferred_name"`
	Pronouns           int64        `json:"pronouns"`
	Title              int64        `json:"title"`
	OSStudent          bool         `json:"os_student"`
	OnlineStudent      bool         `json:"online_student"`
	Language           int64        `json:"language"`
	CitizenCountry     int64        `json:"citizen_country"`
	BirthCountry       int64        `json:"birth_country"`
	PassportCode       string       `json:"passport_code"`
	PassportTypeCode   string       `json:"passport_type_code"`
	PassportIssuer     string       `json:"passport_issuer"`
	PassportIssueDate  string       `json:"passport_issue_date"`
	PassportExpiryDate string       `json:"passport_expiry_date"`
	PassportCountry    int64        `json:"passport_country"`
	PassportPlaceBirth string       `json:"passport_place_birth"`
	MedicalNotes       string       `json:"medical_notes"`
	MedicalCategories  []int64      `json:"medical_categories"`
	Notes              string       `json:"notes"`
	LocalEmail         string       `json:"local_email"`
	LocalPhone         string       `json:"local_phone"`
	LocalMobile        string       `json:"local_mobile"`
	LocalAddress       AddressInput `json:"local_address"`
	OSEmail            string       `json:"os_email"`
	OSPhone            string       `json:"os_phone"`
	OSMobile           string       `json:"os_mobile"`
	OSAddress          AddressInput `json:"os_address"`
	PostalAddress      AddressInput `json:"postal_address"`
}

type AddressInput struct {
	Line1    string `json:"line1"`
	Line2    string `json:"line2"`
	Line3    string `json:"line3"`
	State    string `json:"state"`
	Postcode string `json:"postcode"`
	Country  int64  `json:"country"`
}

// CreateStudent creates a student.
func (c *Client) CreateStudent(
	ctx context.Context,
	studentID int64,
	input CreateStudentInput,
) (Student, int, error) {
	var student Student

	body, err := json.Marshal(input)
	if err != nil {
		return student, http.StatusInternalServerError, fmt.Errorf("encode student request: %w", err)
	}

	requestURL := fmt.Sprintf("%s/students", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewReader(body))
	if err != nil {
		return student, http.StatusInternalServerError, fmt.Errorf("create student request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	data, statusCode, err := c.do(req)
	if err != nil {
		return student, statusCode, fmt.Errorf("create student, status %d: %w", statusCode, err)
	}

	if statusCode != http.StatusOK {
		return student, statusCode, fmt.Errorf(
			"create student returned status %d: %s",
			statusCode,
			strings.TrimSpace(string(data)),
		)
	}

	if err := json.Unmarshal(data, &student); err != nil {
		return student, statusCode, fmt.Errorf("decode create student response: %w", err)
	}

	return student, statusCode, nil
}

// UploadStudentFileInput contains the student document upload details.
type UploadStudentFileInput struct {
	StudentID   int64
	FileName    string
	File        io.Reader
	Category    int64
	Name        string
	Description string
}

// UploadStudentFile uploads a file to a student's document store.
func (c *Client) UploadStudentFile(ctx context.Context, input UploadStudentFileInput) (int, error) {
	if input.StudentID <= 0 {
		return http.StatusBadRequest, fmt.Errorf("student ID must be greater than zero")
	}

	if strings.TrimSpace(input.FileName) == "" {
		return http.StatusBadRequest, fmt.Errorf("file name is required")
	}

	if input.File == nil {
		return http.StatusBadRequest, fmt.Errorf("file is required")
	}

	reader, writer := io.Pipe()
	defer reader.Close()

	multipartWriter := multipart.NewWriter(writer)

	requestURL := fmt.Sprintf("%s/students/%d/documents", c.baseURL, input.StudentID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, reader)
	if err != nil {
		_ = writer.CloseWithError(err)
		return http.StatusInternalServerError, fmt.Errorf(
			"create upload student file request for student %d: %w",
			input.StudentID,
			err,
		)
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	go func() {
		if err := writeStudentFileMultipart(multipartWriter, input); err != nil {
			_ = writer.CloseWithError(err)
			return
		}

		if err := multipartWriter.Close(); err != nil {
			_ = writer.CloseWithError(fmt.Errorf("close multipart writer: %w", err))
			return
		}

		_ = writer.Close()
	}()

	_, statusCode, err := c.do(req)
	if err != nil {
		return statusCode, fmt.Errorf(
			"upload file for student %d, status %d: %w",
			input.StudentID,
			statusCode,
			err,
		)
	}

	return statusCode, nil
}

func writeStudentFileMultipart(writer *multipart.Writer, input UploadStudentFileInput) error {
	filePart, err := writer.CreateFormFile("file", strings.TrimSpace(input.FileName))
	if err != nil {
		return fmt.Errorf("create multipart file field: %w", err)
	}

	if _, err := io.Copy(filePart, input.File); err != nil {
		return fmt.Errorf("copy file into multipart request: %w", err)
	}

	if input.Category > 0 {
		if err := writer.WriteField("category", strconv.FormatInt(input.Category, 10)); err != nil {
			return fmt.Errorf("write category field: %w", err)
		}
	}

	if description := strings.TrimSpace(input.Description); description != "" {
		if err := writer.WriteField("description", description); err != nil {
			return fmt.Errorf("write description field: %w", err)
		}
	}

	if name := strings.TrimSpace(input.Name); name != "" {
		if err := writer.WriteField("name", name); err != nil {
			return fmt.Errorf("write name field: %w", err)
		}
	}

	return nil
}

// GetStudentByStudentNumber retrieves a student using their student number.
func (c *Client) GetStudentByStudentNumber(ctx context.Context, studentNumber string) (Student, int, error) {
	var student Student

	studentNumber = strings.TrimSpace(studentNumber)
	if studentNumber == "" {
		return student, http.StatusBadRequest, fmt.Errorf("student number is required")
	}

	requestURL := fmt.Sprintf("%s/students/student-number/%s", c.baseURL, url.PathEscape(studentNumber))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return student, http.StatusInternalServerError, fmt.Errorf(
			"create get student by student number request for student number %q: %w",
			studentNumber,
			err,
		)
	}

	data, statusCode, err := c.do(req)
	if err != nil {
		return student, statusCode, fmt.Errorf(
			"get student by student number %q, status %d: %w",
			studentNumber,
			statusCode,
			err,
		)
	}

	if statusCode != http.StatusOK {
		return student, statusCode, fmt.Errorf(
			"get student by student number %q returned status %d: %s",
			studentNumber,
			statusCode,
			strings.TrimSpace(string(data)),
		)
	}

	if err := json.Unmarshal(data, &student); err != nil {
		return student, statusCode, fmt.Errorf(
			"decode get student by student number %q response: %w",
			studentNumber,
			err,
		)
	}

	return student, statusCode, nil
}

// GetStudent retrieves a student using their student ID.
func (c *Client) GetStudent(ctx context.Context, studentID int64) (Student, int, error) {
	var student Student

	if studentID <= 0 {
		return student, http.StatusBadRequest, fmt.Errorf("student ID must be greater than zero")
	}

	requestURL := fmt.Sprintf("%s/students/%d", c.baseURL, studentID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return student, http.StatusInternalServerError, fmt.Errorf(
			"create get student request for student %d: %w",
			studentID,
			err,
		)
	}

	data, statusCode, err := c.do(req)
	if err != nil {
		return student, statusCode, fmt.Errorf("get student %d, status %d: %w", studentID, statusCode, err)
	}

	if statusCode != http.StatusOK {
		return student, statusCode, fmt.Errorf(
			"get student %d returned status %d: %s",
			studentID,
			statusCode,
			strings.TrimSpace(string(data)),
		)
	}

	if err := json.Unmarshal(data, &student); err != nil {
		return student, statusCode, fmt.Errorf("decode get student %d response: %w", studentID, err)
	}

	return student, statusCode, nil
}
