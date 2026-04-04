package user

import (
	"errors"

	"webGL-720yun/internal/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) FindUserByID(userID uint) (*model.User, error) {
	var user model.User
	if err := r.db.Preload("Role").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByUsernameOrEmail(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Preload("Role").Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *Repository) DeleteUser(userID uint) error {
	return r.db.Delete(&model.User{}, userID).Error
}

func (r *Repository) GetUserList(req *UserListRequest) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.Model(&model.User{}).Preload("Role")

	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Email != "" {
		query = query.Where("email LIKE ?", "%"+req.Email+"%")
	}
	if req.RoleID > 0 {
		query = query.Where("role_id = ?", req.RoleID)
	}
	if req.ClassID > 0 {
		query = query.Joins("JOIN students ON students.id = users.id").Where("students.class_id = ?", req.ClassID)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?", keyword, keyword, keyword)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *Repository) FindRoleByID(roleID uint) (*model.Role, error) {
	var role model.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, err
	}
	return &role, nil
}

func (r *Repository) FindStudentByID(userID uint) (*model.Student, error) {
	var student model.Student
	if err := r.db.First(&student, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("学生不存在")
		}
		return nil, err
	}
	return &student, nil
}

func (r *Repository) FindStudentByStudentID(studentID string) (*model.Student, error) {
	var student model.Student
	if err := r.db.Where("student_id = ?", studentID).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &student, nil
}

func (r *Repository) CreateStudent(student *model.Student) error {
	return r.db.Create(student).Error
}

func (r *Repository) UpdateStudent(student *model.Student) error {
	return r.db.Save(student).Error
}

func (r *Repository) CreateStudentTeacher(studentTeacher *model.StudentTeacher) error {
	return r.db.Create(studentTeacher).Error
}

func (r *Repository) DeleteStudentTeachers(studentID uint) error {
	return r.db.Where("student_id = ?", studentID).Delete(&model.StudentTeacher{}).Error
}

func (r *Repository) FindClassByName(name string) (*model.Class, error) {
	var class model.Class
	if err := r.db.Where("name = ?", name).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &class, nil
}

func (r *Repository) DeleteUserSessions(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.UserSession{}).Error
}

func (r *Repository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *Repository) Begin() *gorm.DB {
	return r.db.Begin()
}
