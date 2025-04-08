// package repository

// import (
// 	"clothing-shop-api/internal/domain/models"
// 	"clothing-shop-api/internal/domain/services"
// 	"database/sql"
// )

// type paymentRepository struct {
// 	db *sql.DB
// }

// func NewPaymentRepository(db *sql.DB) services.PaymentRepository {
// 	return &paymentRepository{db: db}
// }

// func (r *paymentRepository) Create(payment *models.Payment) error {
// 	query := `
//         INSERT INTO payments (
//             order_id, payment_method, payment_channel, amount, status,
//             transaction_id, payment_token, va_number, payment_url, expiry_time
//         ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
//     `

// 	result, err := r.db.Exec(
// 		query,
// 		payment.OrderID, payment.PaymentMethod, payment.PaymentChannel,
// 		payment.Amount, payment.Status, payment.TransactionID,
// 		payment.PaymentToken, payment.VANumber, payment.PaymentURL,
// 		payment.ExpiryTime,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return err
// 	}

// 	payment.ID = uint(id)
// 	return nil
// }

// func (r *paymentRepository) FindByID(id uint) (*models.Payment, error) {
// 	query := `
//         SELECT id, order_id, payment_method, payment_channel, amount, status,
//                transaction_id, payment_token, va_number, payment_url, expiry_time,
//                paid_at, created_at, updated_at
//         FROM payments
//         WHERE id = ?
//     `

// 	payment := &models.Payment{}
// 	err := r.db.QueryRow(query, id).Scan(
// 		&payment.ID, &payment.OrderID, &payment.PaymentMethod,
// 		&payment.PaymentChannel, &payment.Amount, &payment.Status,
// 		&payment.TransactionID, &payment.PaymentToken, &payment.VANumber,
// 		&payment.PaymentURL, &payment.ExpiryTime, &payment.PaidAt,
// 		&payment.CreatedAt, &payment.UpdatedAt,
// 	)

// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return payment, nil
// }

// func (r *paymentRepository) FindByOrderID(orderID uint) (*models.Payment, error) {
// 	query := `
//         SELECT
//             id, order_id, payment_method, payment_channel, amount, status,
//             transaction_id, payment_token, va_number, payment_url, expiry_time,
//             paid_at, created_at, updated_at
//         FROM payments
//         WHERE order_id = ?
//     `

// 	payment := &models.Payment{}
// 	err := r.db.QueryRow(query, orderID).Scan(
// 		&payment.ID, &payment.OrderID, &payment.PaymentMethod,
// 		&payment.PaymentChannel, &payment.Amount, &payment.Status,
// 		&payment.TransactionID, &payment.PaymentToken, &payment.VANumber,
// 		&payment.PaymentURL, &payment.ExpiryTime, &payment.PaidAt,
// 		&payment.CreatedAt, &payment.UpdatedAt,
// 	)

// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return payment, nil
// }

// func (r *paymentRepository) FindByTransactionID(transactionID string) (*models.Payment, error) {
// 	query := `
//         SELECT
//             id, order_id, payment_method, payment_channel, amount, status,
//             transaction_id, payment_token, va_number, payment_url, expiry_time,
//             paid_at, created_at, updated_at
//         FROM payments
//         WHERE transaction_id = ?
//     `

// 	payment := &models.Payment{}
// 	err := r.db.QueryRow(query, transactionID).Scan(
// 		&payment.ID, &payment.OrderID, &payment.PaymentMethod,
// 		&payment.PaymentChannel, &payment.Amount, &payment.Status,
// 		&payment.TransactionID, &payment.PaymentToken, &payment.VANumber,
// 		&payment.PaymentURL, &payment.ExpiryTime, &payment.PaidAt,
// 		&payment.CreatedAt, &payment.UpdatedAt,
// 	)

// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return payment, nil
// }

// func (r *paymentRepository) Update(payment *models.Payment) error {
// 	query := `
//         UPDATE payments
//         SET
//             payment_method = ?,
//             payment_channel = ?,
//             amount = ?,
//             status = ?,
//             transaction_id = ?,
//             payment_token = ?,
//             va_number = ?,
//             payment_url = ?,
//             expiry_time = ?,
//             paid_at = ?,
//             updated_at = CURRENT_TIMESTAMP
//         WHERE id = ?
//     `

// 	_, err := r.db.Exec(
// 		query,
// 		payment.PaymentMethod,
// 		payment.PaymentChannel,
// 		payment.Amount,
// 		payment.Status,
// 		payment.TransactionID,
// 		payment.PaymentToken,
// 		payment.VANumber,
// 		payment.PaymentURL,
// 		payment.ExpiryTime,
// 		payment.PaidAt,
// 		payment.ID,
// 	)

// 	return err
// }

// func (r *paymentRepository) Delete(id uint) error {
// 	query := "DELETE FROM payments WHERE id = ?"
// 	_, err := r.db.Exec(query, id)
// 	return err
// }
