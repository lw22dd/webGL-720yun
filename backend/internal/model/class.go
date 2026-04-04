package model

import "time"

type Class struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	TeacherID   uint      `gorm:"not null" json:"teacher_id"`
	Teacher     User      `gorm:"foreignKey:TeacherID" json:"teacher"`
	Students    []Student `gorm:"foreignKey:ClassID" json:"students"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Class) TableName() string {
	return "sys_classes"
}

type StudentTeacher struct {
	StudentID uint `gorm:"primaryKey" json:"student_id"`
	TeacherID uint `gorm:"primaryKey" json:"teacher_id"`
}

func (StudentTeacher) TableName() string {
	return "sys_student_teachers"
}
