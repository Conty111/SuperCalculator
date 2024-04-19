package user_test

import (
	"bytes"
	"encoding/json"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/enums"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/interfaces/mocks"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/services"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/controllers/apiv1/user"
	"github.com/cristalhq/jwt/v5"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2" //nolint:revive
	. "github.com/onsi/gomega"    //nolint:revive
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

const (
	URIPrefix = "/api/v1/users"
)

var (
	commonUserID uint = 100
	superUserID  uint = 777
	tokenTTL          = time.Duration(10)
	token             = &jwt.Token{}
)

var _ = Describe("User/Controller", func() {
	var (
		t                     GinkgoTInterface
		userCtrl              *user.Controller
		ctx                   *gin.Context
		w                     *httptest.ResponseRecorder
		userRepo              *mocks.UserManager
		userService           *services.UserService
		authManager           *mocks.AuthManager
		commonUser, superUser *models.User
	)

	BeforeEach(func() {
		gin.SetMode(gin.ReleaseMode)

		w = httptest.NewRecorder()
		ctx, _ = gin.CreateTestContext(w)

		t = GinkgoT()

		userRepo = mocks.NewUserManager(t)
		authManager = mocks.NewAuthManager(t)
		userService = services.NewUserService(
			userRepo,
			authManager,
		)

		userCtrl = user.NewController(userService)

		commonUser = &models.User{
			Email:    "user@mail.ru",
			Username: "Zubenko Michael Petrovich",
			Password: "12345",
			Role:     enums.Common,
		}
		superUser = &models.User{
			Email:    "superUser@mail.ru",
			Role:     enums.Admin,
			Password: "123457",
		}
		commonUser.ID = commonUserID
		superUser.ID = superUserID
	})

	Describe("GetUser()", func() {
		//Context("user is trying to get himself", func() {
		//	It("should return user by id", func() {
		//		userRepo.Mock.On("UserExists", commonUserID).Twice().Return(true, nil)
		//		userRepo.Mock.On("GetUserByID", commonUserID).Twice().Return(commonUser, nil)
		//
		//		ctx.AddParam("userID", strconv.Itoa(commonUserID))
		//		ctx.Set("callerID", commonUserID)
		//		ctx.Request = httptest.NewRequest(http.MethodGet, URIPrefix+"/me", nil)
		//
		//		userCtrl.GetMe(ctx)
		//		Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
		//	})
		//})

		Context("admin is trying to get a user", func() {
			It("should return status ok", func() {
				//userRepo.Mock.On("UserExists", superUserID).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
				userRepo.Mock.On("UserExists", commonUserID).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByID", commonUserID).Once().Return(commonUser, nil)

				ctx.AddParam("userID", strconv.Itoa(int(commonUserID)))
				ctx.Set("callerID", superUserID)
				ctx.Request = httptest.NewRequest(http.MethodGet, URIPrefix+"/"+strconv.Itoa(int(commonUserID)), nil)

				userCtrl.GetUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
			})
		})

		Context("admin is trying to get a non-existent user", func() {
			It("should return not found with text 'user not found'", func() {
				var nonExistentUserID uint = 12345

				//userRepo.Mock.On("UserExists", superUserID).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
				userRepo.Mock.On("UserExists", nonExistentUserID).Once().Return(false, nil)

				ctx.AddParam("userID", strconv.Itoa(int(nonExistentUserID)))
				ctx.Set("callerID", superUserID)
				ctx.Request = httptest.NewRequest(http.MethodGet, URIPrefix+"/"+strconv.Itoa(int(nonExistentUserID)), nil)

				userCtrl.GetUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("common user is trying to get an other user", func() {
			It("should return status forbidden", func() {
				//userRepo.Mock.
				//	On("UserExists", superUserID).
				//	Once().
				//	Return(true, nil)
				userRepo.Mock.
					On("GetUserByID", commonUserID).
					Once().
					Return(commonUser, nil)

				ctx.AddParam("userID", strconv.Itoa(int(superUserID)))
				ctx.Set("callerID", commonUserID)
				ctx.Request = httptest.NewRequest(http.MethodGet, URIPrefix+"/"+strconv.Itoa(int(superUserID)), nil)

				userCtrl.GetUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusForbidden))
			})
		})
	})

	Describe("CreateUser()", func() {
		It("should create user and return status ok", func() {
			userJSON, err := json.Marshal(commonUser)
			Expect(err).To(BeNil())

			userRepo.Mock.On("UserEmailExists", commonUser.Email).Once().Return(false, nil)
			authManager.Mock.On("HashString", commonUser.Password).Return(commonUser.Password, nil)
			userRepo.Mock.On("CreateUser", commonUser).Once().Return(nil)
			//authManager.Mock.On("GetTokenTTL").Return(tokenTTL)
			authManager.Mock.On("BuildToken", commonUserID).Return(token, nil)

			ctx.Request = httptest.NewRequest(http.MethodPost, URIPrefix+"/create", bytes.NewBuffer(userJSON))
			ctx.Request.Header.Set("Content-Type", "application/json")

			userCtrl.CreateUser(ctx)
			Expect(w.Result().StatusCode).To(Equal(http.StatusCreated))
		})

		Context("no email provided", func() {
			It("should return bad request", func() {
				commonUser.Email = ""

				userJSON, err := json.Marshal(commonUser)
				Expect(err).To(BeNil())

				ctx.Request = httptest.NewRequest(http.MethodPost, URIPrefix+"/create", bytes.NewBuffer(userJSON))
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.CreateUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusBadRequest))
			})
		})
	})

	Describe("UpdateUser()", func() {
		Context("admin updates a user", func() {
			It("should update existing user and return status ok", func() {
				newUsername := "HeyYo" //nolint:goconst

				upd := user.Upd{
					Param: "Username",
					Value: newUsername,
				}

				updJSON, err := json.Marshal(upd)
				Expect(err).To(BeNil())

				userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
				userRepo.Mock.On("UserExists", commonUserID).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByID", commonUserID).Once().Return(commonUser, nil)
				userRepo.Mock.On("UpdateUser", commonUser, "Username", newUsername).Once().Return(nil)

				ctx.Set("callerID", superUserID)
				ctx.AddParam("userID", strconv.Itoa(int(commonUserID)))
				ctx.Request = httptest.NewRequest(http.MethodPatch, URIPrefix+"/"+strconv.Itoa(int(commonUserID)), bytes.NewBuffer(updJSON))
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.UpdateUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
			})
		})

		Context("user who's admin wants to update does not exist", func() {
			It("should fail and return status not found", func() {
				var nonExistentUserID uint = 12345 //nolint:goconst

				newEmail := "someNewEmail@mail.ru" //nolint:goconst
				upd := user.Upd{
					Param: "Email",
					Value: newEmail,
				}

				updJSON, err := json.Marshal(upd)
				Expect(err).To(BeNil())

				//userRepo.Mock.On("UserExists", superUserID).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
				userRepo.Mock.On("UserExists", nonExistentUserID).Once().Return(false, nil)

				ctx.Set("callerID", superUserID)
				ctx.AddParam("userID", strconv.Itoa(int(nonExistentUserID)))
				ctx.Request = httptest.NewRequest(http.MethodPatch, URIPrefix+"/"+strconv.Itoa(int(nonExistentUserID)), bytes.NewBuffer(updJSON))
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.UpdateUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("common user is trying to update an other user", func() {
			It("should return status forbidden", func() {
				newEmail := "userNewEmail@mail.ru" //nolint:goconst

				upd := user.Upd{
					Param: "Email",
					Value: newEmail,
				}

				updJSON, err := json.Marshal(upd)
				Expect(err).To(BeNil())

				//userRepo.Mock.
				//	On("UserExists", commonUserID).
				//	Once().
				//	Return(true, nil)
				userRepo.Mock.
					On("GetUserByID", commonUserID).
					Once().
					Return(commonUser, nil)

				ctx.Set("callerID", commonUserID)
				ctx.AddParam("userID", strconv.Itoa(int(commonUserID)))
				ctx.Request = httptest.NewRequest(http.MethodPatch, URIPrefix+"/"+strconv.Itoa(int(superUserID)), bytes.NewBuffer(updJSON))
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.UpdateUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusForbidden))
			})
		})

		Context("admin trying to update not allowed parameter", func() {
			It("should return status forbidden", func() {
				newPasswd := "123456789" //nolint:goconst

				upd := user.Upd{
					Param: "Password",
					Value: newPasswd,
				}

				updJSON, err := json.Marshal(upd)
				Expect(err).To(BeNil())

				//userRepo.Mock.On("UpdateUser", commonUser, "Password", newPasswd).Once().Return(nil)
				userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
				//userRepo.Mock.On("GetUserByID", commonUserID).Once().Return(commonUser, nil)
				userRepo.Mock.On("UserExists", commonUserID).Once().Return(true, nil)

				ctx.Set("callerID", superUserID)
				ctx.AddParam("userID", strconv.Itoa(int(commonUserID)))
				ctx.Request = httptest.NewRequest(http.MethodPatch, URIPrefix+"/"+strconv.Itoa(int(superUserID)), bytes.NewBuffer(updJSON))
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.UpdateUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusForbidden))
			})
		})
	})

	Describe("DeleteUser()", func() {
		It("should delete existing user and return status ok", func() {
			//userRepo.Mock.On("UserExists", superUserID).Once().Return(true, nil)
			userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
			userRepo.Mock.On("UserExists", commonUserID).Once().Return(true, nil)
			userRepo.Mock.On("GetUserByID", commonUserID).Once().Return(commonUser, nil)
			userRepo.Mock.On("DeleteUser", commonUser).Once().Return(nil)

			ctx.AddParam("userID", strconv.Itoa(int(commonUserID)))
			ctx.Set("callerID", superUserID)
			ctx.Request = httptest.NewRequest(http.MethodDelete, URIPrefix+"/"+strconv.Itoa(int(commonUserID)), nil)
			ctx.Request.Header.Set("Content-Type", "application/json")

			userCtrl.DeleteUser(ctx)
			Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
		})

		Context("user who's admin wants to delete does not exist", func() {
			It("should fail and return status not found", func() {
				var nonExistentUserID uint = 12345 //nolint:goconst

				//userRepo.Mock.On("UserExists", superUserID).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
				userRepo.Mock.On("UserExists", nonExistentUserID).Once().Return(false, nil)

				ctx.AddParam("userID", strconv.Itoa(int(nonExistentUserID)))
				ctx.Set("callerID", superUserID)
				ctx.Request = httptest.NewRequest(http.MethodPatch, URIPrefix+"/"+strconv.Itoa(int(nonExistentUserID)), nil)
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.DeleteUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("user without privileges wants to delete a user", func() {
			It("should return status not found", func() {
				//userRepo.Mock.
				//	On("UserExists", commonUserID).
				//	Once().
				//	Return(true, nil)
				userRepo.Mock.On("GetUserByID", commonUserID).
					Once().
					Return(commonUser, nil)

				ctx.AddParam("userID", strconv.Itoa(int(commonUserID)))
				ctx.Set("callerID", commonUserID)
				ctx.Request = httptest.NewRequest(http.MethodPatch, URIPrefix+"/"+strconv.Itoa(int(superUserID)), nil)
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.DeleteUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusForbidden))
			})
		})
	})
	Describe("Login()", func() {
		It("should create jwt token and return him with status ok", func() {
			userJSON, err := json.Marshal(commonUser)
			Expect(err).To(BeNil())

			userRepo.Mock.On("UserEmailExists", commonUser.Email).Once().Return(true, nil)
			userRepo.Mock.On("GetUserByEmail", commonUser.Email).Once().Return(commonUser, nil)
			authManager.Mock.On("HashString", commonUser.Password).Return(commonUser.Password, nil)
			//authManager.Mock.On("GetTokenTTL").Return(tokenTTL)
			authManager.Mock.On("BuildToken", commonUserID).Return(token, nil)

			ctx.Request = httptest.NewRequest(http.MethodPost, URIPrefix+"/login", bytes.NewBuffer(userJSON))
			ctx.Request.Header.Set("Content-Type", "application/json")

			userCtrl.Login(ctx)
			Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
		})

		Context("no email provided", func() {
			It("should return bad request", func() {
				commonUser.Email = ""

				userJSON, err := json.Marshal(commonUser)
				Expect(err).To(BeNil())

				ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewBuffer(userJSON))
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.Login(ctx)

				Expect(w.Result().StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("no password provided", func() {
			It("should return bad request", func() {
				commonUser.Password = ""

				userJSON, err := json.Marshal(commonUser)
				Expect(err).To(BeNil())

				ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewBuffer(userJSON))
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.Login(ctx)

				Expect(w.Result().StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("incorrect password provided", func() {
			It("should return status forbidden", func() {
				var creds struct {
					Email    string `json:"email"`
					Password string `json:"password"`
				}
				creds.Password = "1234"
				creds.Email = commonUser.Email

				userJSON, err := json.Marshal(creds)
				Expect(err).To(BeNil())

				userRepo.Mock.On("UserEmailExists", commonUser.Email).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByEmail", commonUser.Email).Once().Return(commonUser, nil)
				authManager.Mock.On("HashString", creds.Password).Once().Return(creds.Password, nil)

				ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewBuffer(userJSON))
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.Login(ctx)

				Expect(w.Result().StatusCode).To(Equal(http.StatusForbidden))
			})
		})

		Context("user with provided credentials does not exist", func() {
			It("should return not found", func() {
				userJSON, err := json.Marshal(commonUser)
				Expect(err).To(BeNil())

				userRepo.Mock.On("UserEmailExists", commonUser.Email).Once().Return(false, nil)

				ctx.Request = httptest.NewRequest(http.MethodPost, URIPrefix+"/login", bytes.NewBuffer(userJSON))
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.Login(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})
	Describe("GetAllUsers()", func() {
		Context("admin gets all users", func() {
			It("should return all existed users and status ok", func() {
				userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
				userRepo.Mock.On("GetAllUsers", superUserID).Once().Return([]*models.User{
					commonUser,
					superUser,
				}, nil)
				ctx.Set("callerID", superUserID)
				ctx.Request = httptest.NewRequest(http.MethodGet, URIPrefix, nil)
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.GetAllUsers(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
			})
		})

	})
})
