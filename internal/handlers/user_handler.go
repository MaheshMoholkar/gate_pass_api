package handlers

import (
	"database/sql"

	"gate_pass_api/internal/database"
	"gate_pass_api/internal/types"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	store *database.Store
}

func NewUserHandler(store *database.Store) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) HandleSignUp(ctx *fiber.Ctx) error {
	var user *types.User
	err := ctx.BodyParser(&user)
	if err != nil {
		return err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	passwordHash := string(hashedBytes)

	query := `INSERT INTO users (mobile_no, password_hash) VALUES (@mobile_no, @password_hash)`

	_, err = h.store.DB.Exec(query, sql.Named("mobile_no", user.MobileNo), sql.Named("password_hash", passwordHash))
	if err != nil {
		return err
	}

	return ctx.SendStatus(200)
}

func (h *UserHandler) HandleGetUser(ctx *fiber.Ctx) error {
	mobile_no, ok := ctx.Locals("mobile_no").(string)
	if !ok || mobile_no == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	query := `SELECT id, mobile_no, password_hash FROM users WHERE mobile_no = @p1`
	var user types.User

	err := h.store.DB.QueryRow(query, sql.Named("p1", mobile_no)).Scan(&user.ID, &user.MobileNo, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "User not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return ctx.JSON(user)
}

func (h *UserHandler) HandleUpdateUser(ctx *fiber.Ctx) error {
	mobile_no, ok := ctx.Locals("mobile_no").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	var user types.User
	err := ctx.BodyParser(&user)
	if err != nil {
		return err
	}

	query := `UPDATE users SET mobile_no = @new_mobile_no WHERE mobile_no = @mobile_no`

	_, err = h.store.DB.Exec(query,
		sql.Named("new_mobile_no", user.MobileNo),
		sql.Named("mobile_no", mobile_no))

	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) HandleUpdatePassword(ctx *fiber.Ctx) error {
	mobile_no, ok := ctx.Locals("mobile_no").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	var request struct {
		NewPassword string `json:"new_password"`
	}
	err := ctx.BodyParser(&request)
	if err != nil {
		return err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	passwordHash := string(hashedBytes)

	query := `UPDATE users SET password_hash = @password_hash WHERE mobile_no = @mobile_no`

	_, err = h.store.DB.Exec(query,
		sql.Named("password_hash", passwordHash),
		sql.Named("mobile_no", mobile_no))

	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) HandleDeleteUser(ctx *fiber.Ctx) error {
	mobile_no, ok := ctx.Locals("mobile_no").(string)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}

	query := `DELETE FROM users WHERE mobile_no = @mobile_no`

	_, err := h.store.DB.Exec(query, sql.Named("mobile_no", mobile_no))
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
