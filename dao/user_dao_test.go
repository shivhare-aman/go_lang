package dao

import (
	"golang/models"
	"log"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

// SetupTestDB initializes an in-memory SQLite database and migrates the schema.
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Note{}, &models.CreditCard{})
	if err != nil {
		t.Fatalf("Failed to migrate database schema: %v", err)
	}

	return db
}

func TestUserDao_Create(t *testing.T) {
	db := SetupTestDB(t)
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get DB from GORM: %v", err)
		}
		sqlDB.Close()
	}()

	userDao := NewUserDao(db)

	t.Run("Success", func(t *testing.T) {
		user := &models.User{
			Username: "johndoe",
			Password: "password123",
		}

		err := userDao.Create(user)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)

		// Verify in database
		var fetchedUser models.User
		err = db.First(&fetchedUser, user.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, user.Username, fetchedUser.Username)
		assert.Equal(t, user.Password, fetchedUser.Password)
	})

	t.Run("Duplicate Username", func(t *testing.T) {
		user1 := &models.User{
			Username: "janedoe",
			Password: "password1",
		}
		err := userDao.Create(user1)
		assert.NoError(t, err)

		user2 := &models.User{
			Username: "janedoe", // Duplicate username
			Password: "password2",
		}
		err = userDao.Create(user2)
		assert.Error(t, err)
	})
}

func TestUserDao_GetByID(t *testing.T) {
	db := SetupTestDB(t)
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get DB from GORM: %v", err)
		}
		sqlDB.Close()
	}()

	userDao := NewUserDao(db)

	// Insert a user with associated notes and credit card
	user := &models.User{
		Username: "alice",
		Password: "password123",
		Notes: []models.Note{
			{Name: "Note1", Content: "Alice's first note"},
			{Name: "Note2", Content: "Alice's second note"},
		},
		CreditCard: models.CreditCard{
			Number: "4111111111111111",
		},
	}

	err := userDao.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	t.Run("Found", func(t *testing.T) {
		fetchedUser, err := userDao.GetByID(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.Username, fetchedUser.Username)
		assert.Equal(t, user.Password, fetchedUser.Password)
		assert.Len(t, fetchedUser.Notes, 2)
		assert.Equal(t, user.CreditCard.Number, fetchedUser.CreditCard.Number)
	})

	t.Run("Not Found", func(t *testing.T) {
		fetchedUser, err := userDao.GetByID(999) // Non-existent ID
		assert.Error(t, err)
		assert.Nil(t, fetchedUser)
	})
}

func TestUserDao_GetAll(t *testing.T) {
	db := SetupTestDB(t)
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get DB from GORM: %v", err)
		}
		sqlDB.Close()
	}()

	userDao := NewUserDao(db)

	// Insert multiple users
	users := []models.User{
		{Username: "user1", Password: "pass1"},
		{Username: "user2", Password: "pass2"},
		{Username: "jane_doe", Password: "password"},
		{Username: "john_smith", Password: "password"},
		{Username: "jane_smith", Password: "password"},
	}

	for i := range users {
		err := userDao.Create(&users[i])
		assert.NoError(t, err)
	}

	t.Run("GetAll with Pagination and Search - Success", func(t *testing.T) {
		page := 1
		pageSize := 2
		searchQuery := "jane"

		fetchedUsers, err := userDao.GetAll((page-1)*pageSize, pageSize, searchQuery)
		assert.NoError(t, err)
		assert.Len(t, fetchedUsers, 2)
		for _, user := range fetchedUsers {
			assert.Contains(t, user.Username, "jane")
		}
	})

	t.Run("GetAll without Search - Success", func(t *testing.T) {
		page := 2
		pageSize := 2
		searchQuery := ""

		fetchedUsers, err := userDao.GetAll((page-1)*pageSize, pageSize, searchQuery)
		assert.NoError(t, err)
		assert.Len(t, fetchedUsers, 2)
	})

	t.Run("GetAll with No Results", func(t *testing.T) {
		page := 10
		pageSize := 5
		searchQuery := "nonexistent"

		fetchedUsers, err := userDao.GetAll((page-1)*pageSize, pageSize, searchQuery)
		assert.NoError(t, err)
		assert.Len(t, fetchedUsers, 0)
	})

	t.Run("GetAll with Case-Insensitive Search", func(t *testing.T) {
		page := 1
		pageSize := 10
		searchQuery := "JANe"

		fetchedUsers, err := userDao.GetAll((page-1)*pageSize, pageSize, searchQuery)
		assert.NoError(t, err)
		assert.Len(t, fetchedUsers, 2)
		for _, user := range fetchedUsers {
			assert.Contains(t, user.Username, "jane")
		}
	})
}

func TestUserDao_Update(t *testing.T) {
	db := SetupTestDB(t)
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get DB from GORM: %v", err)
		}
		sqlDB.Close()
	}()

	userDao := NewUserDao(db)

	// Insert a user
	user := &models.User{
		Username: "bob",
		Password: "password123",
	}
	err := userDao.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	t.Run("Success", func(t *testing.T) {
		// Update user details
		user.Username = "bobby"
		user.Password = "newpassword"

		err := userDao.Update(user)
		assert.NoError(t, err)

		// Verify updates
		var updatedUser models.User
		err = db.First(&updatedUser, user.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "bobby", updatedUser.Username)
		assert.Equal(t, "newpassword", updatedUser.Password)
	})

	t.Run("Update Non-Existent User", func(t *testing.T) {
		nonExistentUser := &models.User{
			ID:       999,
			Username: "ghost",
			Password: "ghostpassword",
		}

		err := userDao.Update(nonExistentUser)
		assert.NoError(t, err) // GORM's Save will create or update; depends on implementation

		// Verify that the user was created
		var fetchedUser models.User
		err = db.First(&fetchedUser, nonExistentUser.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "ghost", fetchedUser.Username)
		assert.Equal(t, "ghostpassword", fetchedUser.Password)
	})
}

func TestUserDao_Delete(t *testing.T) {
	db := SetupTestDB(t)
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get DB from GORM: %v", err)
		}
		sqlDB.Close()
	}()

	userDao := NewUserDao(db)

	// Insert a user
	user := &models.User{
		Username: "charlie",
		Password: "password123",
	}
	err := userDao.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	t.Run("Success", func(t *testing.T) {
		err := userDao.Delete(user.ID)
		assert.NoError(t, err)

		// Verify deletion
		var deletedUser models.User
		err = db.First(&deletedUser, user.ID).Error
		assert.Error(t, err) // Expect an error because the user should not exist
	})

	t.Run("Delete Non-Existent User", func(t *testing.T) {
		err := userDao.Delete(999) // Non-existent ID
		assert.NoError(t, err)     // No error should be returned for non-existent users
	})
}
