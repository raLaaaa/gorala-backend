package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/raLaaaa/gorala/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseService struct{}

func (d DatabaseService) CreateUser(user *(models.User)) error {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})
	if err = db.Create(&user).Error; err != nil {
		return err
	}

	return err
}

func (d DatabaseService) CreateTask(task *(models.Task)) error {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Task{})
	db.Create(&task)

	return err
}

func (d DatabaseService) UpdateTask(task *(models.Task)) error {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Task{})
	db.Save(&task)

	return err
}

func (d DatabaseService) DeleteTask(id uint64) error {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.Delete(&models.Task{}, id)

	return err
}

func (d DatabaseService) FindAllTasksByUserID(idRequestor uint) ([]models.Task, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Task{})
	requestor, errRequestor := d.FindUserByID(idRequestor)

	if errRequestor != nil {
		fmt.Println(errRequestor)
	}

	tasks := []models.Task{}
	db.Model(&requestor).Association("AllTasks").Find(&tasks)

	return tasks, err
}

func (d DatabaseService) FindAllTasksOfDateByUserID(idRequestor uint, date time.Time) ([]models.Task, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Task{})
	requestor, errRequestor := d.FindUserByID(idRequestor)

	if errRequestor != nil {
		fmt.Println(errRequestor)
	}

	tasks := []models.Task{}
	db.Model(&requestor).Where("execution_date = ?", date).Association("AllTasks").Find(&tasks)

	return tasks, err
}

func (d DatabaseService) FindAllTasksOfDateInRange(idRequestor uint, start time.Time, end time.Time) ([]models.Task, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Task{})
	requestor, errRequestor := d.FindUserByID(idRequestor)

	if errRequestor != nil {
		fmt.Println(errRequestor)
	}

	tasks := []models.Task{}
	db.Model(&requestor).Where("execution_date BETWEEN ? AND ?", start, end).Association("AllTasks").Find(&tasks)

	return tasks, err
}

func (d DatabaseService) FindUserByID(id uint) (*models.User, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})

	var user = models.User{}
	db.First(&user, id)

	return &user, err
}

func (d DatabaseService) FindTaskByID(id uint64) (*models.Task, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})

	var task = models.Task{}
	db.First(&task, id)

	return &task, err
}

//Find user from email
func (d DatabaseService) FindByEmail(email string) (*models.User, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{})

	var user = models.User{}

	if err := db.Where("email = ?", &email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, err
}

func (d DatabaseService) CreateConfirmationToken(user *models.User) (*models.ConfirmationToken, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	tokenUUID := uuid.New()

	db.AutoMigrate(&models.ConfirmationToken{})
	token := models.ConfirmationToken{
		Token:     tokenUUID.String(),
		UserID:    user.ID,
		Activated: false,
	}

	db.Create(&token)

	return &token, err
}

func (d DatabaseService) ResolveConfirmationToken(token string) (bool, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.ConfirmationToken{})

	tokenObj := models.ConfirmationToken{}

	if err := db.First(&tokenObj, "token = ?", token).Error; err != nil {
		return false, errors.New("invalid token")
	}

	if tokenObj.Activated {
		return false, errors.New("token already redeemed")
	}

	user, err := d.FindUserByID(tokenObj.UserID)

	if err != nil {
		return false, errors.New("could not find user")
	}

	user.Accepted = true
	db.Save(&user)

	fmt.Println(user.ID)

	tokenObj.Activated = true
	db.Save(tokenObj)

	return true, err
}
