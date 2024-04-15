package user_test

import (
	"bytes"
	"encoding/json"
	"github.com/cristalhq/jwt/v5"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	. "github.com/onsi/ginkgo" //nolint:revive
	. "github.com/onsi/gomega" //nolint:revive
	"gitlab.sch.ocrv.com.rzd/blockchain/trainees/quiz-svc/internal/enums"
	"gitlab.sch.ocrv.com.rzd/blockchain/trainees/quiz-svc/internal/models"
	"gitlab.sch.ocrv.com.rzd/blockchain/trainees/quiz-svc/internal/services"
	"gitlab.sch.ocrv.com.rzd/blockchain/trainees/quiz-svc/internal/services/mocks"
	"gitlab.sch.ocrv.com.rzd/blockchain/trainees/quiz-svc/internal/web/controllers/apiv1/user"
	"gitlab.sch.ocrv.com.rzd/blockchain/trainees/quiz-svc/internal/web/helpers"
	"net/http"
	"net/http/httptest"
	"time"
)

const (
	URIPrefix = "/api/v1/users"
)

var (
	commonUserID, _ = uuid.NewV4()
	superUserID, _  = uuid.NewV4()
	pagCfg          = &helpers.PaginationConfig{
		DefaultLimit:       10,
		DefaultOrderColumn: "created_at",
		DefaultOrder:       ">",
	}
	pag = &helpers.Pagination{
		PageURI:     URIPrefix + "/" + commonUserID.String() + "/attachedQuizzes",
		Limit:       10,
		Offset:      0,
		Order:       ">",
		OrderColumn: "created_at",
	}
	tokenTTL = time.Duration(10)
	token    = &jwt.Token{}
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
			pagCfg,
		)

		userCtrl = user.NewController(userService)

		commonUser = &models.User{
			UID:      commonUserID,
			Email:    "user@mail.ru",
			FullName: "Zubenko Michael Petrovich",
			Password: "12345",
			AttachedQuizzes: []*models.UserAttachedQuiz{
				{
					UserID: commonUserID,
					QuizID: 1,
					Quiz: &models.Quiz{
						Title:       "gg",
						Description: "efaf",
					},
				},
				{
					UserID: commonUserID,
					QuizID: 2,
					Quiz: &models.Quiz{
						Title:       "gg1",
						Description: "efadasf",
					},
					CompletedAt: time.Now(),
				},
			},
			Role: enums.Common,
		}
		superUser = &models.User{
			UID:      superUserID,
			Email:    "superUser@mail.ru",
			Role:     enums.SuperUser,
			Password: "123457",
		}
	})

	Describe("GetUser()", func() {
		Context("user is trying to get himself", func() {
			It("should return user by id", func() {
				userRepo.Mock.On("UserExists", commonUserID).Twice().Return(true, nil)
				userRepo.Mock.On("GetUserByID", commonUserID).Twice().Return(commonUser, nil)

				ctx.AddParam("userID", commonUserID.String())
				ctx.Set("callerID", commonUserID.String())
				ctx.Request = httptest.NewRequest(http.MethodGet, URIPrefix+"/me", nil)

				userCtrl.GetMe(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
			})
		})

		Context("superuser is trying to get a user", func() {
			It("should return status ok", func() {
				userRepo.Mock.On("UserExists", superUserID).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
				userRepo.Mock.On("UserExists", commonUserID).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByID", commonUserID).Once().Return(commonUser, nil)

				ctx.AddParam("userID", commonUserID.String())
				ctx.Set("callerID", superUserID.String())
				ctx.Request = httptest.NewRequest(http.MethodGet, URIPrefix+"/"+commonUserID.String(), nil)

				userCtrl.GetUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
			})
		})

		Context("superuser is trying to get a non-existent user", func() {
			It("should return not found with text 'user not found'", func() {
				nonExistentUserID, _ := uuid.NewV4()

				userRepo.Mock.On("UserExists", superUserID).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
				userRepo.Mock.On("UserExists", nonExistentUserID).Once().Return(false, nil)

				ctx.AddParam("userID", nonExistentUserID.String())
				ctx.Set("callerID", superUserID.String())
				ctx.Request = httptest.NewRequest(http.MethodGet, URIPrefix+"/"+nonExistentUserID.String(), nil)

				userCtrl.GetUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("common user is trying to get an other user", func() {
			It("should return status forbidden", func() {
				userRepo.Mock.
					On("UserExists", commonUserID).
					Once().
					Return(true, nil)
				userRepo.Mock.
					On("GetUserByID", commonUserID).
					Once().
					Return(commonUser, nil)

				ctx.AddParam("userID", commonUserID.String())
				ctx.Set("callerID", commonUserID.String())
				ctx.Request = httptest.NewRequest(http.MethodGet, URIPrefix+"/"+superUserID.String(), nil)

				userCtrl.GetUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("CreateUser()", func() {
		It("should create user and return status ok", func() {
			commonUser.AttachedQuizzes = []*models.UserAttachedQuiz{}

			userJSON, err := json.Marshal(commonUser)
			Expect(err).To(BeNil())

			userRepo.Mock.On("UserEmailExists", commonUser.Email).Once().Return(false, nil)
			authManager.Mock.On("HashString", commonUser.Password).Return(commonUser.Password, nil)
			userRepo.Mock.On("CreateUser", commonUser).Once().Return(nil)
			authManager.Mock.On("GetTokenTTL").Return(tokenTTL)
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
		Context("superuser updates a user", func() {
			It("should update existing user and return status ok", func() {
				newEmail := "userNewEmail@mail.ru" //nolint:goconst

				upd := helpers.Upd{
					Param: "Email",
					Value: newEmail,
				}

				updJSON, err := json.Marshal(upd)
				Expect(err).To(BeNil())

				userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
				userRepo.Mock.On("UserExists", commonUserID).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByID", commonUserID).Once().Return(commonUser, nil)
				userRepo.Mock.On("UpdateUser", commonUser, "Email", newEmail).Once().Return(nil)

				ctx.Set("callerID", superUserID.String())
				ctx.AddParam("userID", commonUserID.String())
				ctx.Request = httptest.NewRequest(http.MethodPatch, URIPrefix+"/"+commonUserID.String(), bytes.NewBuffer(updJSON))
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.UpdateUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
			})
		})

		Context("user who's superuser wants to update does not exist", func() {
			It("should fail and return status not found", func() {
				nonExistentUserID, _ := uuid.NewV4() //nolint:goconst

				newEmail := "someNewEmail@mail.ru" //nolint:goconst
				upd := helpers.Upd{
					Param: "Email",
					Value: newEmail,
				}

				updJSON, err := json.Marshal(upd)
				Expect(err).To(BeNil())

				userRepo.Mock.On("UserExists", superUserID).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
				userRepo.Mock.On("UserExists", nonExistentUserID).Once().Return(false, nil)

				ctx.Set("callerID", superUserID.String())
				ctx.AddParam("userID", nonExistentUserID.String())
				ctx.Request = httptest.NewRequest(http.MethodPatch, URIPrefix+"/"+commonUserID.String(), bytes.NewBuffer(updJSON))
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.UpdateUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("common user is trying to update an other user", func() {
			It("should return status forbidden", func() {
				newEmail := "userNewEmail@mail.ru" //nolint:goconst

				upd := helpers.Upd{
					Param: "Email",
					Value: newEmail,
				}

				updJSON, err := json.Marshal(upd)
				Expect(err).To(BeNil())

				userRepo.Mock.
					On("UserExists", commonUserID).
					Once().
					Return(true, nil)
				userRepo.Mock.
					On("GetUserByID", commonUserID).
					Once().
					Return(commonUser, nil)

				ctx.Set("callerID", commonUserID.String())
				ctx.AddParam("userID", commonUserID.String())
				ctx.Request = httptest.NewRequest(http.MethodPatch, URIPrefix+"/"+superUserID.String(), bytes.NewBuffer(updJSON))
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.UpdateUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	Describe("DeleteUser()", func() {
		It("should delete existing user and return status ok", func() {
			userRepo.Mock.On("UserExists", superUserID).Once().Return(true, nil)
			userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
			userRepo.Mock.On("UserExists", commonUserID).Once().Return(true, nil)
			userRepo.Mock.On("GetUserByID", commonUserID).Once().Return(commonUser, nil)
			userRepo.Mock.On("DeleteUser", commonUser).Once().Return(nil)

			ctx.AddParam("userID", commonUserID.String())
			ctx.Set("callerID", superUserID.String())
			ctx.Request = httptest.NewRequest(http.MethodDelete, URIPrefix+"/"+commonUserID.String(), nil)
			ctx.Request.Header.Set("Content-Type", "application/json")

			userCtrl.DeleteUser(ctx)
			Expect(w.Result().StatusCode).To(Equal(http.StatusOK))
		})

		Context("user who's superuser wants to delete does not exist", func() {
			It("should fail and return status not found", func() {
				nonExistentUserID, _ := uuid.NewV4() //nolint:goconst

				userRepo.Mock.On("UserExists", superUserID).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByID", superUserID).Once().Return(superUser, nil)
				userRepo.Mock.On("UserExists", nonExistentUserID).Once().Return(false, nil)

				ctx.AddParam("userID", nonExistentUserID.String())
				ctx.Set("callerID", superUserID.String())
				ctx.Request = httptest.NewRequest(http.MethodPatch, URIPrefix+"/"+nonExistentUserID.String(), nil)
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.DeleteUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("user without privileges wants to delete a user", func() {
			It("should return status not found", func() {
				userRepo.Mock.
					On("UserExists", commonUserID).
					Once().
					Return(true, nil)
				userRepo.Mock.On("GetUserByID", commonUserID).
					Once().
					Return(commonUser, nil)

				ctx.AddParam("userID", commonUserID.String())
				ctx.Set("callerID", commonUserID.String())
				ctx.Request = httptest.NewRequest(http.MethodPatch, URIPrefix+"/"+superUserID.String(), nil)
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.DeleteUser(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Describe("Login()", func() {
			It("should create jwt token and return him with status ok", func() {
				userJSON, err := json.Marshal(commonUser)
				Expect(err).To(BeNil())

				userRepo.Mock.On("UserEmailExists", commonUser.Email).Once().Return(true, nil)
				userRepo.Mock.On("GetUserByEmail", commonUser.Email).Once().Return(commonUser, nil)
				authManager.Mock.On("HashString", commonUser.Password).Return(commonUser.Password, nil)
				authManager.Mock.On("GetTokenTTL").Return(tokenTTL)
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
	})

	Describe("GetUserAttachedQuizzes()", func() {
		Context("no filters", func() {
			It("should return status ok and all quizzes that attached to user", func() {
				userRepo.Mock.On("UserExists", commonUserID).Return(true, nil)
				userRepo.Mock.On("GetUserByID", superUserID).Return(superUser, nil)
				userRepo.Mock.
					On(
						"GetUserAttachedQuizzes",
						commonUserID,
						&models.AttachedQuizFilter{},
						pag,
						true,
					).
					Return(commonUser.AttachedQuizzes, nil, nil)
				userRepo.Mock.On("GetLastUserAttachedQuizID", commonUserID).Once().Return(uint(len(commonUser.AttachedQuizzes)), nil)

				ctx.Set("callerID", superUserID.String())
				ctx.AddParam("userID", commonUserID.String())
				ctx.Request = httptest.NewRequest(http.MethodGet, URIPrefix+"/"+commonUserID.String()+"/attachedQuizzes", nil)
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.GetUserAttachedQuizzes(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))

				var res user.AttachedQuizzesResponse
				err := json.NewDecoder(w.Result().Body).Decode(&res)
				Expect(err).To(BeNil())

				Expect(len(res.Data)).To(Equal(len(commonUser.AttachedQuizzes)))
			})
		})

		Context("with only completed filter", func() {
			It("should return status ok and all the quizzes that user has completed", func() {
				numberOfCompletedQuizzes := 1
				completedQuizID := uint(2)
				filters := "?onlyCompleted=true"

				userRepo.Mock.On("UserExists", commonUserID).Return(true, nil)
				userRepo.Mock.On("GetUserByID", superUserID).Return(superUser, nil)
				userRepo.Mock.
					On(
						"GetUserAttachedQuizzes",
						commonUserID,
						&models.AttachedQuizFilter{OnlyCompleted: true},
						pag,
						true,
					).
					Return(commonUser.AttachedQuizzes[1:], nil, nil)
				userRepo.Mock.On("GetLastUserAttachedQuizID", commonUserID).Once().Return(uint(len(commonUser.AttachedQuizzes)), nil)

				ctx.Set("callerID", superUserID.String())
				ctx.AddParam("userID", commonUserID.String())
				ctx.Request = httptest.NewRequest(http.MethodGet, URIPrefix+"/"+commonUserID.String()+filters, nil)
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.GetUserAttachedQuizzes(ctx)
				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))

				var res user.AttachedQuizzesResponse
				err := json.NewDecoder(w.Result().Body).Decode(&res)
				Expect(err).To(BeNil())

				Expect(len(res.Data)).To(Equal(numberOfCompletedQuizzes))
				Expect(res.Data[0].QuizID).To(Equal(completedQuizID))
			})
		})

		Context("with only not completed filter", func() {
			It("should return status ok and all the quizzes that user has not completed yet", func() {
				numberOfNotCompletedQuizzes := 1
				notCompletedQuizID := uint(1)
				filters := "?onlyNotCompleted=true"

				userRepo.Mock.On("UserExists", commonUserID).Return(true, nil)
				userRepo.Mock.On("GetUserByID", superUserID).Return(superUser, nil)
				userRepo.Mock.
					On(
						"GetUserAttachedQuizzes",
						commonUserID,
						&models.AttachedQuizFilter{OnlyNotCompleted: true},
						pag,
						true,
					).
					Return(commonUser.AttachedQuizzes[:1], nil, nil)
				userRepo.Mock.On("GetLastUserAttachedQuizID", commonUserID).Once().Return(uint(len(commonUser.AttachedQuizzes)), nil)

				ctx.Set("callerID", superUserID.String())
				ctx.AddParam("userID", commonUserID.String())
				ctx.Request = httptest.NewRequest(http.MethodGet, URIPrefix+"/"+commonUserID.String()+filters, nil)
				ctx.Request.Header.Set("Content-Type", "application/json")

				userCtrl.GetUserAttachedQuizzes(ctx)

				Expect(w.Result().StatusCode).To(Equal(http.StatusOK))

				var res user.AttachedQuizzesResponse
				err := json.NewDecoder(w.Result().Body).Decode(&res)
				Expect(err).To(BeNil())

				Expect(len(res.Data)).To(Equal(numberOfNotCompletedQuizzes))
				Expect(res.Data[0].QuizID).To(Equal(notCompletedQuizID))
			})
		})
	})
})
