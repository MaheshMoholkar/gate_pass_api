package handlers

import (
	"database/sql"
	"gate_pass_api/internal/database"
	"gate_pass_api/internal/types"
	"time"

	"github.com/gofiber/fiber/v2"
)

type VisitorHandler struct {
	store *database.Store
}

func NewVisitorHandler(store *database.Store) *VisitorHandler {
	return &VisitorHandler{
		store: store,
	}
}

func (h *VisitorHandler) HandleGetVisitors(ctx *fiber.Ctx) error {
	rows, err := h.store.DB.Query(`SELECT id, name, purpose, date, address, vehicle_no, mobile_no, image, appointment, in, out FROM visitor_form`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var visitors []types.Visitor
	for rows.Next() {
		var visitor types.Visitor
		var image sql.NullString
		var address sql.NullString
		var appointment sql.NullString
		var vehicleNo sql.NullInt32
		var inTime sql.NullTime
		var outTime sql.NullTime

		err := rows.Scan(
			&visitor.ID,
			&visitor.Name,
			&visitor.Purpose,
			&visitor.Date,
			&address,
			&vehicleNo,
			&visitor.MobileNo,
			&image,
			&appointment,
			&inTime,
			&outTime,
		)
		if err != nil {
			return err
		}

		if address.Valid {
			visitor.Address = address.String
		}

		if image.Valid {
			visitor.Image = image.String
		}

		if appointment.Valid {
			visitor.Appointment = appointment.String
		}

		if vehicleNo.Valid {
			visitor.VehicleNo = int(vehicleNo.Int32)
		}

		if inTime.Valid {
			visitor.In = inTime.Time
		}

		if outTime.Valid {
			visitor.Out = outTime.Time
		}

		visitors = append(visitors, visitor)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return ctx.JSON(visitors)
}

func (h *VisitorHandler) HandleVisitorEntryForm(ctx *fiber.Ctx) error {
	var form types.Visitor
	err := ctx.BodyParser(&form)
	if err != nil {
		return err
	}

	query := `INSERT INTO visitor_form (name, purpose, date, address, vehicle_no, mobile_no, image, appointment, in) 
              VALUES (@name, @purpose, @date, @address, @vehicle_no, @mobile_no, @image, @appointment, @in)`

	_, err = h.store.DB.Exec(query,
		sql.Named("name", form.Name),
		sql.Named("purpose", form.Purpose),
		sql.Named("date", form.Date),
		sql.Named("address", form.Address),
		sql.Named("vehicle_no", form.VehicleNo),
		sql.Named("mobile_no", form.MobileNo),
		sql.Named("image", form.Image),
		sql.Named("appointment", form.Appointment),
		sql.Named("in", time.Now()))

	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *VisitorHandler) HandleVisitorExit(ctx *fiber.Ctx) error {
	mobileNo := ctx.Params("mobile_no")
	if mobileNo == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Mobile number is required")
	}

	query := `UPDATE visitor_form SET out = @out WHERE mobile_no = @mobile_no AND out IS NULL`

	_, err := h.store.DB.Exec(query,
		sql.Named("out", time.Now()),
		sql.Named("mobile_no", mobileNo))

	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusOK)
}
