package contact

import (
	"database/sql"
	"os"
	"testing"

	"github.com/codersgyan/camp/internal/database"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	code := m.Run()

	testDB.Close()
	os.Exit(code)
}

func TestContactRepositoryCreate(t *testing.T) {
	testDB = setupDB(t)
	repo := NewRepository(testDB)

	contact := &Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@codersgyan.com",
	}

	_, err := repo.CreateContactOrUpsertTags(contact)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	var count int
	testDB.QueryRow("SELECT COUNT(*) FROM contacts WHERE email = $1",
		contact.Email).Scan(&count)

	if count != 1 {
		t.Errorf("Expected 1 contact, got %d", count)
	}
}

func TestContactRepositoryCreateWithTags(t *testing.T) {
	testDB = setupDB(t)
	repo := NewRepository(testDB)

	contact := &Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@codersgyan.com",
		Tags:      []Tag{{Text: "purchase:golang"}, {Text: "subscribed:platform"}},
	}

	_, err := repo.CreateContactOrUpsertTags(contact)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	var count int
	testDB.QueryRow("SELECT COUNT(*) FROM contacts WHERE email = $1",
		contact.Email).Scan(&count)

	if count != 1 {
		t.Errorf("Expected 1 contact, got %d", count)
	}

	var tagsCount int
	testDB.QueryRow("SELECT COUNT(*) FROM contact_tag").Scan(&tagsCount)

	if tagsCount != 2 {
		t.Errorf("Expected 2 tags, got %d", tagsCount)
	}
}

func TestContactRepositoryUpsertWithTagsIfEmailExists(t *testing.T) {
	testDB = setupDB(t)
	repo := NewRepository(testDB)

	contact1 := &Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@codersgyan.com",
		Tags:      []Tag{{Text: "purchase:golang"}, {Text: "subscribed:platform"}},
	}

	_, err := repo.CreateContactOrUpsertTags(contact1)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	contact2 := &Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@codersgyan.com",
		Tags:      []Tag{{Text: "joined:annual"}},
	}

	_, err = repo.CreateContactOrUpsertTags(contact2)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	var count int
	testDB.QueryRow("SELECT COUNT(*) FROM contacts WHERE email = $1",
		contact1.Email).Scan(&count)

	if count != 1 {
		t.Errorf("Expected 1 contact, got %d", count)
	}

	var tagsCount int
	testDB.QueryRow("SELECT COUNT(*) FROM contact_tag").Scan(&tagsCount)

	if tagsCount != 3 {
		t.Errorf("Expected 3 tags, got %d", tagsCount)
	}
}

func TestContactRepositoryThrowErrorIfTagsNotSent(t *testing.T) {
	testDB = setupDB(t)
	repo := NewRepository(testDB)

	contact1 := &Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@codersgyan.com",
		Tags:      []Tag{{Text: "purchase:golang"}, {Text: "subscribed:platform"}},
	}

	_, err := repo.CreateContactOrUpsertTags(contact1)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	// here we are not sending the tags
	contact2 := &Contact{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@codersgyan.com",
	}

	_, err = repo.CreateContactOrUpsertTags(contact2)
	if err == nil {
		t.Fatalf("Expected error when tags are not sent, but got nil")
	}

	var count int
	testDB.QueryRow("SELECT COUNT(*) FROM contacts WHERE email = $1",
		contact1.Email).Scan(&count)

	if count != 1 {
		t.Errorf("Expected 1 contact, got %d", count)
	}

	var tagsCount int

	testDB.QueryRow("SELECT COUNT(*) FROM contact_tag").Scan(&tagsCount)

	if tagsCount != 2 {
		t.Errorf("Expected 2 tags, got %d", tagsCount)
	}
}

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err := database.RunMigration(db); err != nil {
		t.Fatal(err)
	}

	return db
}
