package handlers

import (
	"database/sql"
	"time"

	"gate_pass_api/internal/database"
	"gate_pass_api/internal/types"

	"github.com/gofiber/fiber/v2"
)

type StaffHandler struct {
	store *database.Store
}

func NewStaffHandler(store *database.Store) *StaffHandler {
	return &StaffHandler{
		store: store,
	}
}

func (h *StaffHandler) HandleAddStaff(ctx *fiber.Ctx) error {
	var staff types.Staff
	err := ctx.BodyParser(&staff)
	if err != nil {
		return err
	}

	query := `INSERT INTO staff (name, mobile_no, image) VALUES (@name, @mobile_no, @image)`

	_, err = h.store.DB.Exec(query,
		sql.Named("name", staff.Name),
		sql.Named("mobile_no", staff.MobileNo),
		sql.Named("image", staff.Image))
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *StaffHandler) HandleStaffEntryForm(ctx *fiber.Ctx) error {
	mobileNo := ctx.Query("mobile_no")
	if mobileNo == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("mobile_no is required")
	}

	// Retrieve staff details using mobile_no
	var staff types.Staff
	query := `SELECT id, name, mobile_no, image FROM staff WHERE mobile_no = @mobile_no`
	err := h.store.DB.QueryRow(query, sql.Named("mobile_no", mobileNo)).Scan(&staff.ID, &staff.Name, &staff.MobileNo, &staff.Image)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Staff Not Found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}

	var form types.StaffEntry
	err = ctx.BodyParser(&form)
	if err != nil {
		return err
	}

	query = `INSERT INTO staff_entry (staff_id, name, purpose, entry_time, image) VALUES (@staff_id, @name, @purpose, @entry_time, @image)`
	_, err = h.store.DB.Exec(query,
		sql.Named("staff_id", staff.ID),
		sql.Named("name", staff.Name),
		sql.Named("purpose", form.Purpose),
		sql.Named("entry_time", time.Now()),
		sql.Named("image", staff.Image))

	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *StaffHandler) HandleStaffExit(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("ID is required")
	}

	query := `UPDATE staff_entry SET exit_time = @exit_time WHERE id = @id`

	_, err := h.store.DB.Exec(query,
		sql.Named("exit_time", time.Now()),
		sql.Named("id", id))

	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *StaffHandler) HandleGetAllStaff(ctx *fiber.Ctx) error {
	query := `SELECT name, mobile_no FROM staff`

	rows, err := h.store.DB.Query(query)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}
	defer rows.Close()

	var staffList []types.Staff
	for rows.Next() {
		var staff types.Staff
		err := rows.Scan(&staff.Name, &staff.MobileNo)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Internal Server Error",
			})
		}
		staffList = append(staffList, staff)
	}

	if err := rows.Err(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}

	return ctx.JSON(staffList)
}
