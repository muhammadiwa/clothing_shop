// package repository

// import (
// 	"database/sql"
// 	"fmt"

// 	"clothing-shop-api/internal/domain/interfaces"
// 	"clothing-shop-api/internal/domain/models"
// )

// type categoryRepository struct {
// 	db *sql.DB
// }

// func NewCategoryRepository(db *sql.DB) interfaces.CategoryRepository {
// 	return &categoryRepository{
// 		db: db,
// 	}
// }

// func (r *categoryRepository) Create(category *models.Category) error {
// 	query := `
//         INSERT INTO categories (
//             name, slug, description, parent_id,
//             created_at, updated_at
//         ) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
//     `
// 	result, err := r.db.Exec(query,
// 		category.Name, category.Slug,
// 		category.Description, category.ParentID,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to create category: %v", err)
// 	}

// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return fmt.Errorf("failed to get last insert ID: %v", err)
// 	}

// 	category.ID = uint(id)
// 	return nil
// }

// func (r *categoryRepository) FindByID(id uint) (*models.Category, error) {
// 	query := `
//         SELECT
//             c1.id, c1.name, c1.slug, c1.description, c1.parent_id,
//             c1.created_at, c1.updated_at,
//             c2.id, c2.name, c2.slug, c2.description
//         FROM categories c1
//         LEFT JOIN categories c2 ON c1.parent_id = c2.id
//         WHERE c1.id = ? AND c1.deleted_at IS NULL
//     `

// 	category := &models.Category{}
// 	var parent models.Category
// 	var parentID *uint
// 	var parentName, parentSlug, parentDesc *string

// 	err := r.db.QueryRow(query, id).Scan(
// 		&category.ID, &category.Name, &category.Slug,
// 		&category.Description, &parentID,
// 		&category.CreatedAt, &category.UpdatedAt,
// 		&parent.ID, &parentName, &parentSlug, &parentDesc,
// 	)

// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to find category: %v", err)
// 	}

// 	if parentID != nil {
// 		category.ParentID = parentID
// 		parent.Name = *parentName
// 		parent.Slug = *parentSlug
// 		parent.Description = *parentDesc
// 		category.Parent = &parent
// 	}

// 	// Get children categories
// 	children, err := r.findChildren(id)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get children categories: %v", err)
// 	}
// 	category.Children = children

// 	return category, nil
// }

// func (r *categoryRepository) FindAll() ([]*models.Category, error) {
// 	query := `
//         SELECT
//             id, name, slug, description, parent_id,
//             created_at, updated_at
//         FROM categories
//         WHERE deleted_at IS NULL
//         ORDER BY
//             CASE WHEN parent_id IS NULL THEN 0 ELSE 1 END,
//             name ASC
//     `

// 	rows, err := r.db.Query(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query categories: %v", err)
// 	}
// 	defer rows.Close()

// 	categories := make([]*models.Category, 0)
// 	for rows.Next() {
// 		category := &models.Category{}
// 		err := rows.Scan(
// 			&category.ID, &category.Name, &category.Slug,
// 			&category.Description, &category.ParentID,
// 			&category.CreatedAt, &category.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan category: %v", err)
// 		}
// 		categories = append(categories, category)
// 	}

// 	// Build tree structure
// 	categoryMap := make(map[uint]*models.Category)
// 	var rootCategories []*models.Category

// 	// First pass: Create map of all categories
// 	for _, category := range categories {
// 		categoryMap[category.ID] = category
// 	}

// 	// Second pass: Build tree structure
// 	for _, category := range categories {
// 		if category.ParentID == nil {
// 			rootCategories = append(rootCategories, category)
// 		} else {
// 			parent := categoryMap[*category.ParentID]
// 			if parent != nil {
// 				parent.Children = append(parent.Children, *category)
// 			}
// 		}
// 	}

// 	return rootCategories, nil
// }

// func (r *categoryRepository) Update(category *models.Category) error {
// 	query := `
//         UPDATE categories
//         SET
//             name = ?,
//             slug = ?,
//             description = ?,
//             parent_id = ?,
//             updated_at = CURRENT_TIMESTAMP
//         WHERE id = ? AND deleted_at IS NULL
//     `

// 	result, err := r.db.Exec(query,
// 		category.Name, category.Slug,
// 		category.Description, category.ParentID,
// 		category.ID,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to update category: %v", err)
// 	}

// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get affected rows: %v", err)
// 	}
// 	if rows == 0 {
// 		return fmt.Errorf("category not found")
// 	}

// 	return nil
// }

// func (r *categoryRepository) Delete(id uint) error {
// 	// First check if category has any children
// 	var childCount int
// 	err := r.db.QueryRow("SELECT COUNT(*) FROM categories WHERE parent_id = ? AND deleted_at IS NULL", id).Scan(&childCount)
// 	if err != nil {
// 		return fmt.Errorf("failed to check for children categories: %v", err)
// 	}
// 	if childCount > 0 {
// 		return fmt.Errorf("cannot delete category with children")
// 	}

// 	// Check if category has any products
// 	var productCount int
// 	err = r.db.QueryRow("SELECT COUNT(*) FROM products WHERE category_id = ? AND deleted_at IS NULL", id).Scan(&productCount)
// 	if err != nil {
// 		return fmt.Errorf("failed to check for products: %v", err)
// 	}
// 	if productCount > 0 {
// 		return fmt.Errorf("cannot delete category with products")
// 	}

// 	// Soft delete the category
// 	query := `
//         UPDATE categories
//         SET
//             deleted_at = CURRENT_TIMESTAMP
//         WHERE id = ? AND deleted_at IS NULL
//     `
// 	result, err := r.db.Exec(query, id)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete category: %v", err)
// 	}

// 	rows, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("failed to get affected rows: %v", err)
// 	}
// 	if rows == 0 {
// 		return fmt.Errorf("category not found")
// 	}

// 	return nil
// }

// // Helper methods
// func (r *categoryRepository) findChildren(parentID uint) ([]models.Category, error) {
// 	query := `
//         SELECT
//             id, name, slug, description, parent_id,
//             created_at, updated_at
//         FROM categories
//         WHERE parent_id = ? AND deleted_at IS NULL
//         ORDER BY name ASC
//     `

// 	rows, err := r.db.Query(query, parentID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	children := make([]models.Category, 0)
// 	for rows.Next() {
// 		var child models.Category
// 		err := rows.Scan(
// 			&child.ID, &child.Name, &child.Slug,
// 			&child.Description, &child.ParentID,
// 			&child.CreatedAt, &child.UpdatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		children = append(children, child)
// 	}

// 	return children, nil
// }
