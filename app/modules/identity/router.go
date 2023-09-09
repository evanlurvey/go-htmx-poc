package identity

import (
	"htmx-poc/app"
	"htmx-poc/app/csrf"
	"htmx-poc/app/forms"
	"htmx-poc/app/logging"
	"htmx-poc/app/template"
	"htmx-poc/validation"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewRouter(templates template.TemplateEngine, db *DB, form forms.Service) *Router {
	return &Router{
		templates: templates,
		db:        db,
		form:      form,
	}
}

type Router struct {
	templates template.TemplateEngine
	db        *DB
	form      forms.Service
}

func (r *Router) Setup(rtr fiber.Router) {
	rtr = rtr.Group("/identity")
	rtr.Get("/login", r.LoginGET)
	rtr.Post("/login", r.LoginPOST)
	rtr.Get("/create-account", r.CreateAccountGET)
	rtr.Post("/create-account", r.CreateAccountPOST)
	rtr.Get("/account-created", r.AccountCreatedGET)
}

func (r *Router) LoginGET(c *fiber.Ctx) error {
	return r.templates.Render(c, "pages/identity/form.html", map[string]any{
		"form": LoginForm,
	})
}

func (r *Router) LoginPOST(c *fiber.Ctx) error {
	ctx := c.UserContext()
	c.Set("HX-Trigger", "auth")
	l := logging.FromContext(ctx)
	var formData LoginFormData
	if err := r.form.Parse(c, &formData); err != nil {
		return err
	}

	req := formData.LoginRequest
	if ve := validation.ValidateStruct(req); ve != nil {
		form := LoginForm.GenerateFields(req, ve)
		return r.templates.Render(c, "pages/identity/form.html", fiber.Map{
			"form": form,
		})
	}

	// check recent attempts with this email.
	// won't leak account info because it will report it as locked even if the email doesn't exist
	const accountLocked = "account locked, too many attempts"
	if attempts, err := r.db.getRecentAttempts(ctx, req.Email, time.Minute*10); err != nil {
		l.Error("request to db failed", zap.Error(err))
		// we don't want the user to know we are busted.
		form := LoginForm.GenerateFields(req)
		form.Error = accountLocked
		return r.templates.Render(c, "pages/identity/form.html", fiber.Map{
			"form": form,
		})
	} else if attempts.Unsuccessful() > 10 {
		l.Info("account locked", zap.String("email", req.Email))
		form := LoginForm.GenerateFields(req)
		form.Error = accountLocked
		return r.templates.Render(c, "pages/identity/form.html", fiber.Map{
			"form": form,
		})
	}

	const genericError = "invalid email or password"

	user, err := r.db.getUserByEmail(ctx, req.Email)
	if err != nil {
		l.Error("request to db failed", zap.Error(err))
		// we don't want the user to know we are busted.
		form := LoginForm.GenerateFields(req)
		form.Error = genericError
		return r.templates.Render(c, "pages/identity/form.html", fiber.Map{
			"form": form,
		})
	} else if user.ID == "" {
		r.db.createLoginAttempt(ctx, loginAttempt{
			id:      app.NewID(),
			at:      time.Now(),
			email:   req.Email,
			outcome: loginOutcome_invalidEmail,
		})
		form := LoginForm.GenerateFields(req)
		form.Error = genericError
		return r.templates.Render(c, "pages/identity/form.html", fiber.Map{
			"form": form,
		})
	}

	if !user.password.VerifyPassword(req.Password) {
		r.db.createLoginAttempt(ctx, loginAttempt{
			id:      app.NewID(),
			at:      time.Now(),
			user_id: user.ID,
			email:   req.Email,
			outcome: loginOutcome_invalidPassword,
		})
		form := LoginForm.GenerateFields(req)
		form.Error = genericError
		return r.templates.Render(c, "pages/identity/form.html", fiber.Map{
			"form": form,
		})
	}

	// TODO: Create session

	r.db.createLoginAttempt(ctx, loginAttempt{
		id:      app.NewID(),
		at:      time.Now(),
		user_id: user.ID,
		email:   req.Email,
		outcome: loginOutcome_success,
	})

	session := newAuthenticatedSession(user.User)
	r.db.storeSession(ctx, session)

	c.Cookie(&fiber.Cookie{
		Name:     sessionCookie,
		Value:    session.token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		Secure:   true,
		HTTPOnly: true, // still sent with ajax calls
		SameSite: "Lax",
	})

	return c.Redirect("/", 303)
}

func (r *Router) CreateAccountGET(c *fiber.Ctx) error {
	return r.templates.Render(c, "pages/identity/form.html", map[string]any{
		"form": CreateAccountForm,
	})
}

func (r *Router) CreateAccountPOST(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var formData CreateAccountFormData
	if err := r.form.Parse(c, &formData); err != nil {
		return err
	}

	req := formData.CreateAccountRequest
	if ve := validation.ValidateStruct(req); ve != nil {
		form := CreateAccountForm.GenerateFields(req, ve)
		return r.templates.Render(c, "pages/identity/form.html", fiber.Map{
			"form": form,
		})
	}

	// check for existing account
	if user, err := r.db.getUserByEmail(ctx, req.Email); err != nil {
		return r.templates.Render(c, "pages/sorry.html", fiber.Map{})
	} else if user.ID != "" { // give them the same message to not leak accounts
		// TODO: Send email letting them know they had an account
		return c.Redirect("/identity/account-created", 303)
	}

	// create account
	err := r.db.storeUser(user{
		User: User{
			ID:        app.NewID(),
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
		},
		password: newPHC(req.Password),
	})
	if err != nil {
		return r.templates.Render(c, "pages/sorry.html", fiber.Map{})
	}

	return c.Redirect("/identity/account-created", 303)
}

func (r *Router) AccountCreatedGET(c *fiber.Ctx) error {
	return r.templates.Render(c, "pages/identity/account-created.html", fiber.Map{})
}
