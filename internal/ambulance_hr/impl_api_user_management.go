package ambulance_hr

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/patrikkuchar/ambulance-hr-webapi/internal/db_service"
	"net/http"
)

// AddPersonalDocument - Add personal document to user
func (this *implUserManagementAPI) AddPersonalDocument(ctx *gin.Context) {
	// Get user id from path
	userId := ctx.Param("userId")

	// Parse the request body
	body := PersonalDocumentEntry{}
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

	// Get db service from context
	db, ok := getDbServiceFromContext(ctx)
	if !ok {
		return
	}

	// Get user by id
	user, userErr := getUserById(ctx, db, userId)
	if userErr != nil {
		return
	}

	// Map PersonalDocumentEntry to PersonalDocument
	personalDocument := PersonalDocument{
		Id:      uuid.New().String(),
		Name:    body.Name,
		Content: body.Content,
	}

	// Add personal document to user
	user.PersonalDocument = append(user.PersonalDocument, personalDocument)

	// Update user
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

	ctx.JSON(http.StatusOK, personalDocument)
}

// LoginUser - Login user
func (this *implUserManagementAPI) LoginUser(ctx *gin.Context) {
	// Parse the request body
	body := LoginEntry{}
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

	// Get db service from context
	db, ok := getDbServiceFromContext(ctx)
	if !ok {
		return
	}

	// Find user by email
	user, err := db.FindDocumentByEmail(ctx, body.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the found user
	ctx.JSON(http.StatusOK, user)
}

// UpdatePersonalDocument - Update personal document of user
func (this *implUserManagementAPI) UpdatePersonalDocument(ctx *gin.Context) {
	// Get user id from path
	userId := ctx.Param("userId")

	// Parse the request body
	body := PersonalDocument{}
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

	// Get db service from context
	db, ok := getDbServiceFromContext(ctx)
	if !ok {
		return
	}

	// Get user by id
	user, userErr := getUserById(ctx, db, userId)
	if userErr != nil {
		return
	}

	// Find personal document by id
	var personalDocumentFromDocument *PersonalDocument
	for i, doc := range user.PersonalDocument {
		if doc.Id == body.Id {
			personalDocumentFromDocument = &user.PersonalDocument[i]
			break
		}
	}

	// Update personal document
	if personalDocumentFromDocument != nil {
		personalDocumentFromDocument.Name = body.Name
		personalDocumentFromDocument.Content = body.Content
	} else {
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Personal document not found",
			})
		return
	}

	// Update user
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

	ctx.JSON(http.StatusOK, personalDocumentFromDocument)
}

func getUserById(ctx *gin.Context, db db_service.DbService[UserDto], id string) (*UserDto, error) {
	user, err := db.FindDocument(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return nil, err
	}
	return user, nil
}
