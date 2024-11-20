package routes

import (
	"portfolio/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(publicRoutes, authenticatedRoutes *gin.RouterGroup) {
	publicRoutes.POST("/login", controllers.Login())
}

func LayoutRoutes(publicRoutes, authenticatedRoutes *gin.RouterGroup) {
	publicRoutes.GET("/manage-layout", controllers.ManageLayout())

	authenticatedRoutes.POST("/manage-layout", controllers.ManageLayout())
	authenticatedRoutes.PUT("/manage-layout", controllers.ManageLayout())
}

func CertificateRoutes(publicRoutes, authenticatedRoutes *gin.RouterGroup) {
	publicRoutes.GET("/certificate/get-all", controllers.GetAllCertificates())
	publicRoutes.GET("/certificate/get-one/:id", controllers.GetOneCertificate())

	authenticatedRoutes.POST("/certificate/create", controllers.CreateCertificate())
	authenticatedRoutes.PUT("/certificate/update/:id", controllers.UpdateCertificate())
	authenticatedRoutes.DELETE("/certificate/delete/:id", controllers.DeleteCertificate())
}

func ServiceRoutes(publicRoutes, authenticatedRoutes *gin.RouterGroup) {
	publicRoutes.GET("/service/get-all", controllers.GetAllServices())
	publicRoutes.GET("/service/get-one/:id", controllers.GetOneService())

	authenticatedRoutes.POST("/service/create", controllers.CreateService())
	authenticatedRoutes.PUT("/service/update/:id", controllers.UpdateService())
	authenticatedRoutes.DELETE("/service/delete/:id", controllers.DeleteService())
}

func ProjectRoutes(publicRoutes, authenticatedRoutes *gin.RouterGroup) {
	publicRoutes.GET("/project/get-all", controllers.GetAllProjects())
	publicRoutes.GET("/project/get-one/:id", controllers.GetOneProject())

	authenticatedRoutes.POST("/project/create", controllers.CreateProject())
	authenticatedRoutes.PUT("/project/update/:id", controllers.UpdateProject())
	authenticatedRoutes.DELETE("/project/delete/:id", controllers.DeleteProject())
}

func EmailRoutes(publicRoutes, authenticatedRoutes *gin.RouterGroup) {
	publicRoutes.POST("/email/create", controllers.CreateEmail())

	authenticatedRoutes.GET("/email/get-all", controllers.GetAllEmails())
	authenticatedRoutes.GET("/email/get-one/:id", controllers.GetOneEmail())
	authenticatedRoutes.DELETE("/email/delete/:id", controllers.DeleteEmail())
}

// Expense App
func UserRoutes(publicRoutes, expenseRoutes *gin.RouterGroup, expenseAdminRoutes *gin.RouterGroup) {
	publicRoutes.POST("/expense/register", controllers.RegisterUser())
	publicRoutes.POST("/expense/login", controllers.LoginUser())

	expenseRoutes.GET("/me", controllers.GetCurrentUser())
	expenseRoutes.PUT("/update-user-info", controllers.UpdateUserInfo())
	expenseRoutes.PUT("/update-user-password", controllers.UpdateUserPassword())

	expenseAdminRoutes.GET("/get-all-users", controllers.GetAllUsers())
	expenseAdminRoutes.PUT("/update-user-role", controllers.UpdateUserRole())
	expenseAdminRoutes.DELETE("/delete-user/:id", controllers.DeleteUser())
}

func ExpenseCategoryRoutes(expenseRoutes *gin.RouterGroup) {
	expenseRoutes.POST("/create-category", controllers.CreateExpenseCategory())
	expenseRoutes.PUT("/update-category/:id", controllers.UpdateExpenseCategory())
	expenseRoutes.DELETE("/delete-category/:id", controllers.DeleteExpenseCategory())
	expenseRoutes.GET("/get-all-categories", controllers.GetAllExpenseCategories())
	expenseRoutes.GET("/get-one-category/:id", controllers.GetOneExpenseCategory())
}

func ExpenseItemRoutes(expenseRoutes *gin.RouterGroup) {
	expenseRoutes.POST("/create-item", controllers.CreateExpenseItem())
	expenseRoutes.PUT("/update-item/:id", controllers.UpdateExpenseItem())
	expenseRoutes.DELETE("/delete-item/:id", controllers.DeleteExpenseItem())
	expenseRoutes.GET("/get-all-incomes", controllers.GetAllIncomes())
	expenseRoutes.GET("/get-all-outcomes", controllers.GetAllOutcomes())
}
