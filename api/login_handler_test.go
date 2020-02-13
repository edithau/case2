package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/InVisionApp/interview-test/api/dal"
	"github.com/InVisionApp/interview-test/api/session"
	"github.com/InVisionApp/rye"
	jwtgo "github.com/dgrijalva/jwt-go"
)

var _ = Describe("Login handler tests", func() {
	var (
		request       *http.Request
		response      *httptest.ResponseRecorder
		api           *API
		handlers      []rye.Handler
		mwHandler     *rye.MWHandler
		testEmail     string
		testPassword  string
		loginRequest  io.Reader
		loginResponse *LoginResponse
		errorResponse *rye.JSONStatus
	)

	BeforeEach(func() {
		mwHandler = rye.NewMWHandler(rye.Config{})
		api = &API{MWHandler: mwHandler, apiDAL: dal.NewMockDAL()}
		testEmail = "user@example.com"
		testPassword = "electric cat festival"
		loginRequest = strings.NewReader(fmt.Sprintf(`{"email": %q, "password": %q}`, testEmail, testPassword))
	})

	Describe("loginHandler", func() {
		JustBeforeEach(func() {
			handlers = []rye.Handler{
				mwCallingService(callingService),
				api.loginHandler,
			}
			response = httptest.NewRecorder()
		})

		Context("Happy Path", func() {
			It("Should succeed without errors", func() {
				now := time.Now()
				then := now.Add(jwtCookieExpiresDuration + time.Second)
				request = httptest.NewRequest("POST", "/api/login", loginRequest)
				request.Header.Set("Calling-Service", callingService)

				api.MWHandler.Handle(handlers).ServeHTTP(response, request)

				By("Checking the status code")
				Expect(response.Code).To(Equal(http.StatusOK))

				By("Checking the Login Response")
				Expect(json.Unmarshal(response.Body.Bytes(), &loginResponse)).To(BeNil())
				Expect(loginResponse.Status).To(Equal("ok"))
				Expect(loginResponse.Message).To(ContainSubstring("Login successful and session created"))
				Expect(loginResponse.JWT).To(BeEmpty())

				By("Checking the headers")
				Expect(response.Header()).ToNot(BeNil())
				cookies := response.Result().Cookies()
				Expect(len(cookies)).To(Equal(1))
				Expect(cookies[0].Name).To(Equal(jwtCookieName))
				Expect(cookies[0].Expires.After(now)).To(BeTrue())
				Expect(cookies[0].Expires.Before(then)).To(BeTrue())

				By("Checking the JWT")
				jwt := cookies[0].Value
				token, err := jwtgo.ParseWithClaims(jwt, &jwtgo.MapClaims{},
					func(token *jwtgo.Token) (interface{}, error) {
						return []byte(jwtSecret), nil
					})
				Expect(err).To(BeNil())
				claims, ok := token.Claims.(*jwtgo.MapClaims)
				Expect(ok).To(BeTrue())
				Expect(claims).ToNot(BeNil())
				Expect(claims.Valid()).To(BeNil())
				Expect(claims.VerifyIssuer(jwtIssuer, true)).To(BeTrue())
				Expect(*claims).Should(HaveKeyWithValue(jwtSessionName, session.FakeSessionID))
			})

			It("Should succeed with JWT in body", func() {
				request = httptest.NewRequest("POST", "/api/login?jwtbody", loginRequest)
				request.Header.Set("Calling-Service", callingService)

				api.MWHandler.Handle(handlers).ServeHTTP(response, request)

				By("Checking the status code")
				Expect(response.Code).To(Equal(http.StatusOK))

				By("Checking the Login Response")
				Expect(json.Unmarshal(response.Body.Bytes(), &loginResponse)).To(BeNil())
				Expect(loginResponse.Status).To(Equal("ok"))
				Expect(loginResponse.Message).To(ContainSubstring("Login successful and session created"))
				Expect(loginResponse.JWT).ToNot(BeEmpty())

				By("Checking the headers")
				Expect(response.Header()).ToNot(BeNil())
				Expect(len(response.Result().Cookies())).To(BeZero())
			})

			It("Should succeed with JWT in body and header", func() {
				now := time.Now()
				then := now.Add(jwtCookieExpiresDuration + time.Second)
				request = httptest.NewRequest("POST", "/api/login?jwtboth", loginRequest)
				request.Header.Set("Calling-Service", callingService)

				api.MWHandler.Handle(handlers).ServeHTTP(response, request)

				By("Checking the status code")
				Expect(response.Code).To(Equal(http.StatusOK))

				By("Checking the Login Response")
				Expect(json.Unmarshal(response.Body.Bytes(), &loginResponse)).To(BeNil())
				Expect(loginResponse.Status).To(Equal("ok"))
				Expect(loginResponse.Message).To(ContainSubstring("Login successful and session created"))
				Expect(loginResponse.JWT).ToNot(BeEmpty())

				By("Checking the headers")
				Expect(response.Header()).ToNot(BeNil())
				cookies := response.Result().Cookies()
				Expect(len(cookies)).To(Equal(1))
				Expect(cookies[0].Name).To(Equal(jwtCookieName))
				Expect(cookies[0].Expires.After(now)).To(BeTrue())
				Expect(cookies[0].Expires.Before(then)).To(BeTrue())

				By("Checking the JWT")
				jwt := cookies[0].Value
				token, err := jwtgo.ParseWithClaims(jwt, &jwtgo.MapClaims{},
					func(token *jwtgo.Token) (interface{}, error) {
						return []byte(jwtSecret), nil
					})
				Expect(err).To(BeNil())
				claims, ok := token.Claims.(*jwtgo.MapClaims)
				Expect(ok).To(BeTrue())
				Expect(claims).ToNot(BeNil())
				Expect(claims.Valid()).To(BeNil())
				Expect(claims.VerifyIssuer(jwtIssuer, true)).To(BeTrue())
				then = now.Add(jwtExpiresAtDuration - time.Second)
				Expect(claims.VerifyExpiresAt(then.Unix(), true))
				Expect(*claims).Should(HaveKeyWithValue(jwtSessionName, session.FakeSessionID))
			})
		})

		Context("Sad Path", func() {
			It("Should fail with HTTP status 400 - no credentials", func() {
				loginRequest = strings.NewReader(`{}`)
				request = httptest.NewRequest("POST", "/api/login", loginRequest)
				request.Header.Set("Calling-Service", callingService)

				api.MWHandler.Handle(handlers).ServeHTTP(response, request)

				By("Checking the status code")
				Expect(response.Code).To(Equal(http.StatusBadRequest))

				By("Checking the Login Response")
				Expect(json.Unmarshal(response.Body.Bytes(), &errorResponse)).To(BeNil())
				Expect(errorResponse.Status).To(Equal("error"))

				By("Checking the headers")
				Expect(response.Header()).ToNot(BeNil())
				Expect(len(response.Result().Cookies())).To(BeZero())
			})

			It("Should fail with HTTP status 400 - no email", func() {
				loginRequest = strings.NewReader(`{"password":"password"}`)
				request = httptest.NewRequest("POST", "/api/login", loginRequest)
				request.Header.Set("Calling-Service", callingService)

				api.MWHandler.Handle(handlers).ServeHTTP(response, request)

				By("Checking the status code")
				Expect(response.Code).To(Equal(http.StatusBadRequest))

				By("Checking the Login Response")
				Expect(json.Unmarshal(response.Body.Bytes(), &errorResponse)).To(BeNil())
				Expect(errorResponse.Status).To(Equal("error"))

				By("Checking the headers")
				Expect(response.Header()).ToNot(BeNil())
				Expect(len(response.Result().Cookies())).To(BeZero())
			})

			It("Should fail with HTTP status 400 - invalid email", func() {
				loginRequest = strings.NewReader(`{"email":"email","password":"password"}`)
				request = httptest.NewRequest("POST", "/api/login", loginRequest)
				request.Header.Set("Calling-Service", callingService)

				api.MWHandler.Handle(handlers).ServeHTTP(response, request)

				By("Checking the status code")
				Expect(response.Code).To(Equal(http.StatusBadRequest))

				By("Checking the Login Response")
				Expect(json.Unmarshal(response.Body.Bytes(), &errorResponse)).To(BeNil())
				Expect(errorResponse.Status).To(Equal("error"))

				By("Checking the headers")
				Expect(response.Header()).ToNot(BeNil())
				Expect(len(response.Result().Cookies())).To(BeZero())
			})

			It("Should fail with HTTP status 400 - no password", func() {
				loginRequest = strings.NewReader(`{"email":"pineapple@example.com"}`)
				request = httptest.NewRequest("POST", "/api/login", loginRequest)
				request.Header.Set("Calling-Service", callingService)

				api.MWHandler.Handle(handlers).ServeHTTP(response, request)

				By("Checking the status code")
				Expect(response.Code).To(Equal(http.StatusBadRequest))

				By("Checking the Login Response")
				Expect(json.Unmarshal(response.Body.Bytes(), &errorResponse)).To(BeNil())
				Expect(errorResponse.Status).To(Equal("error"))
				Expect(errorResponse.Message).ToNot(BeEmpty())

				By("Checking the headers")
				Expect(response.Header()).ToNot(BeNil())
				Expect(len(response.Result().Cookies())).To(BeZero())
			})

			It("Should fail with HTTP status 401 - invalid creds", func() {
				loginRequest = strings.NewReader(`{"email":"bad@example.com", "password":"bad password"}`)
				request = httptest.NewRequest("POST", "/api/login", loginRequest)
				request.Header.Set("Calling-Service", callingService)

				api.MWHandler.Handle(handlers).ServeHTTP(response, request)

				By("Checking the status code")
				Expect(response.Code).To(Equal(http.StatusUnauthorized))

				By("Checking the Login Response")
				Expect(json.Unmarshal(response.Body.Bytes(), &errorResponse)).To(BeNil())
				Expect(errorResponse.Status).To(Equal("error"))
				Expect(errorResponse.Message).To(ContainSubstring("Unauthorized: invalid credentials"))

				By("Checking the headers")
				Expect(response.Header()).ToNot(BeNil())
				Expect(len(response.Result().Cookies())).To(BeZero())
			})

			It("Should fail with HTTP status 403 - no caller set", func() {
				request = httptest.NewRequest("POST", "/api/login", loginRequest)

				api.MWHandler.Handle(handlers).ServeHTTP(response, request)

				By("Checking the status code")
				Expect(response.Code).To(Equal(http.StatusForbidden))

				By("Checking the Login Response")
				Expect(json.Unmarshal(response.Body.Bytes(), &errorResponse)).To(BeNil())
				Expect(errorResponse.Status).To(Equal("error"))
				Expect(errorResponse.Message).To(ContainSubstring("Caller not whitelisted"))

				By("Checking the headers")
				Expect(response.Header()).ToNot(BeNil())
				Expect(len(response.Result().Cookies())).To(BeZero())
			})

			It("Should fail with HTTP status 403 - invalid caller", func() {
				request = httptest.NewRequest("POST", "/api/login", loginRequest)
				request.Header.Set("Calling-Service", "whatever")

				api.MWHandler.Handle(handlers).ServeHTTP(response, request)

				By("Checking the status code")
				Expect(response.Code).To(Equal(http.StatusForbidden))

				By("Checking the Login Response")
				Expect(json.Unmarshal(response.Body.Bytes(), &errorResponse)).To(BeNil())
				Expect(errorResponse.Status).To(Equal("error"))
				Expect(errorResponse.Message).To(ContainSubstring("Caller not whitelisted"))

				By("Checking the headers")
				Expect(response.Header()).ToNot(BeNil())
				Expect(len(response.Result().Cookies())).To(BeZero())
			})

			It("Should fail with HTTP status 500 - database error", func() {
				request = httptest.NewRequest("POST", "/api/login", loginRequest)
				request.Header.Set("Calling-Service", callingService)

				api.apiDAL = dal.NewFailDAL()
				api.MWHandler.Handle(handlers).ServeHTTP(response, request)

				By("Checking the status code")
				Expect(response.Code).To(Equal(http.StatusInternalServerError))

				By("Checking the Login Response")
				Expect(json.Unmarshal(response.Body.Bytes(), &errorResponse)).To(BeNil())
				Expect(errorResponse.Status).To(Equal("error"))
				Expect(errorResponse.Message).To(ContainSubstring("database error"))

				By("Checking the headers")
				Expect(response.Header()).ToNot(BeNil())
				Expect(len(response.Result().Cookies())).To(BeZero())
			})
		})
	})
})
