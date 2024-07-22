package handlers

import (
	"database/sql"
	"gate_pass_api/internal/database"
	"gate_pass_api/internal/types"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthParams struct {
	MobileNo string `json:"mobile_no"`
	Password string `json:"password"`
}

type AuthHandler struct {
	store *database.Store
}

func NewAuthHandler(store *database.Store) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request payload",
		})
	}

	var user types.User
	query := `SELECT id, mobile_no, password_hash FROM users WHERE mobile_no = @mobile_no`
	err := h.store.DB.QueryRow(query, sql.Named("mobile_no", params.MobileNo)).Scan(&user.ID, &user.MobileNo, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid mobile_no or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Compare the stored hash with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid mobile_no or password",
		})
	}

	token, err := GenerateJWT(user.MobileNo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"token":  token,
	})
}

func GenerateJWT(mobile_no string) (string, error) {
	claims := jwt.MapClaims{
		"mobile_no": mobile_no,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
