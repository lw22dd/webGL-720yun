package model

import "time"

type Class struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null;comment:班级名称" json:"name"`
	Description string    `gorm:"type:varchar(255);comment:班级描述" json:"description"`
	TeacherID   uint      `gorm:"not null;comment:班主任ID" json:"teacher_id"`
	Teacher     User      `gorm:"foreignKey:TeacherID" json:"teacher"`
	Students    []Student `gorm:"foreignKey:ClassID" json:"students"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Class) TableName() string {
	return "sys_classes"
}

type StudentTeacher struct {
	StudentID uint `gorm:"primaryKey;comment:学生ID" json:"student_id"`
	TeacherID uint `gorm:"primaryKey;comment:教师ID" json:"teacher_id"`
}

func (StudentTeacher) TableName() string {
	return "sys_student_teachers"
}
