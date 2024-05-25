package ambulance_hr

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/patrikkuchar/ambulance-hr-webapi/internal/db_service"
	"log"
	"net/http"
)

// CreateUser - Create new user
func (this *implHRManagementAPI) CreateUser(ctx *gin.Context) {

	// Get db service from context
	db, ok := getDbServiceFromContext(ctx)
	if !ok {
		return
	}

	userEntry := UserEntry{}
	err := ctx.BindJSON(&userEntry)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	// Check if user with email already exists
	emailExists, emailErr := userWithEmailExists(ctx, db, userEntry.Email)
	if emailErr != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to check if user with email already exists",
				"error":   emailErr.Error(),
			})
		return
	}
	if emailExists {
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "User with email already exists",
				"error":   "User with email already exists",
			})
		return

	}

	// Convert UserEntry to UserDto
	user := UserDto{
		Id:               uuid.New().String(),
		Name:             userEntry.Name,
		Role:             userEntry.Role,
		Phone:            userEntry.Phone,
		Email:            userEntry.Email,
		Department:       userEntry.Department,
		PersonalDocument: []PersonalDocument{},
	}

	err = db.CreateDocument(ctx, user.Id, &user)

	switch err {
	case nil:
		usersResponse, listErr := getUserListWithDbContext(ctx, db)
		if listErr != nil {
			ctx.JSON(
				http.StatusInternalServerError,
				gin.H{
					"status":  "Internal Server Error",
					"message": "Failed to get users from database",
					"error":   listErr.Error(),
				})
			return
		}

		ctx.JSON(
			http.StatusCreated,
			usersResponse,
		)
	case db_service.ErrConflict:
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Ambulance already exists",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create ambulance in database",
				"error":   err.Error(),
			},
		)
	}
}

// DeleteUser - Delete user by id
func (this *implHRManagementAPI) DeleteUser(ctx *gin.Context) {

	// Get db service from context
	db, ok := getDbServiceFromContext(ctx)
	if !ok {
		return
	}

	// Get user id from path
	userId := ctx.Param("userId")

	// Delete user from db
	err := db.DeleteDocument(ctx, userId)

	switch err {
	case nil:
		ctx.JSON(
			http.StatusNoContent,
			gin.H{},
		)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "User not found",
				"error":   err.Error(),
			})
	default:
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to delete user from database",
				"error":   err.Error(),
			})
	}
	return
}

// GetUser - Get user by id
func (this *implHRManagementAPI) GetUser(ctx *gin.Context) {

	// Get db service from context
	db, ok := getDbServiceFromContext(ctx)
	if !ok {
		return
	}

	// Get user id from path
	userId := ctx.Param("userId")

	// Get user from db
	user, err := db.FindDocument(ctx, userId)

	switch err {
	case nil:
		ctx.JSON(
			http.StatusOK,
			user,
		)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "User not found",
				"error":   err.Error(),
			})
	default:
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to get user from database",
				"error":   err.Error(),
			})
	}
	return
}

// GetUsers - Get all users
func (this *implHRManagementAPI) GetUsers(ctx *gin.Context) {
	users, err := getUserList(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to get users from database",
				"error":   err.Error(),
			})
		return
	}

	ctx.JSON(
		http.StatusOK,
		users,
	)
}

// UpdateUserDepartment - Update user department
func (this *implHRManagementAPI) UpdateUserDepartment(ctx *gin.Context) {

	// Get db service from context
	db, ok := getDbServiceFromContext(ctx)
	if !ok {
		return
	}

	// Get user id from path
	userId := ctx.Param("userId")

	body := DepartmentDto{}
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	// Get user by id
	user, userErr := getUserById(ctx, db, userId)
	if userErr != nil {
		return
	}

	// Update user
	user.Department = body.Department
	err = db.UpdateDocument(ctx, userId, user)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to update user",
				"error":   err.Error(),
			})
		return
	}

	ctx.JSON(
		http.StatusOK,
		user,
	)
}

// function that returns users []UserList
func getUserList(ctx *gin.Context) ([]UserList, error) {
	// Get db service from context
	db, ok := getDbServiceFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get db service from context")
	}

	return getUserListWithDbContext(ctx, db)
}

func getUserListWithDbContext(ctx *gin.Context, db db_service.DbService[UserDto]) ([]UserList, error) {
	// Get users from db
	users, err := getUsers(ctx, db)
	if err != nil {
		return nil, err
	}

	// Map []UserDto to []UserList
	userList := make([]UserList, len(users))
	for i, user := range users {
		userList[i] = UserList{
			Id:         user.Id,
			Name:       user.Name,
			Role:       user.Role,
			Department: user.Department,
		}
	}

	return userList, nil
}

// function that returns users []UserDto from db
func getUsers(ctx *gin.Context, db db_service.DbService[UserDto]) ([]*UserDto, error) {
	// Get all users from db
	users, err := db.GetAllDocuments(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func getDbServiceFromContext(ctx *gin.Context) (db_service.DbService[UserDto], bool) {
	// get db service from context
	value, exists := ctx.Get("db_service")
	if exists {
		log.Printf("db_service retrieved from context: %v", value)
	} else {
		log.Println("db_service not found in context")
	}
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return nil, false
	}

	db, ok := value.(db_service.DbService[UserDto])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return nil, false
	}

	return db, true
}

func userWithEmailExists(ctx *gin.Context, db db_service.DbService[UserDto], email string) (bool, error) {
	_, err := db.FindDocumentByEmail(ctx, email)

	switch err {
	case nil:
		return true, nil
	case db_service.ErrNotFound:
		return false, nil
	default:
		return false, err
	}

}
