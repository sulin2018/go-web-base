package models

import (
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          uint       `gorm:"primaryKey" uri:"id" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
	Username    string     `gorm:"type:varchar(50);not null;unique" description:"用户名" json:"username"`
	Password    string     `gorm:"type:varchar(100)" json:"password,omitempty"` // omitempty作用: 该字段为空时, 对象转json将忽略该字段 -字符: 可彻底忽略
	ChineseName string     `gorm:"type:varchar(25)" description:"中文名" json:"chinese_name"`
	// UserDN      string     `gorm:"type:varchar(100)" description:"LDAP DN" json:"user_dn"`
	Active    bool   `gorm:"default:1" json:"active"`
	Superuser bool   `gorm:"default:0" json:"superuser"`
	Phone     string `gorm:"type:varchar(20)" description:"手机" json:"phone"`

	Permissions   []*Permission `gorm:"many2many:user_permission" json:"permissions"`
	PermissionIds []uint        `gorm:"-" json:"permission_ids"`
	Groups        []*Group      `gorm:"many2many:user_group" json:"groups"`
	GroupIds      []uint        `gorm:"-" json:"group_ids"`
}

type Permission struct {
	ID          uint   `gorm:"primaryKey" uri:"id" json:"id"`
	Name        string `gorm:"type:varchar(30);not null;unique" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	Users    []*User  `gorm:"many2many:user_permission" json:"users"`
	UserIds  []uint   `gorm:"-" json:"user_ids"`
	Groups   []*Group `gorm:"many2many:group_permission" json:"groups"`
	GroupIds []uint   `gorm:"-" json:"group_ids"`
}

type Group struct {
	ID          uint   `gorm:"primaryKey" uri:"id" json:"id"`
	Name        string `gorm:"type:varchar(30);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	Permissions   []*Permission `gorm:"many2many:group_permission" json:"permissions"`
	PermissionIds []uint        `gorm:"-" json:"permission_ids"`
	Users         []*User       `gorm:"many2many:user_group" json:"users"`
	UserIds       []uint        `gorm:"-" json:"user_ids"`
}

// func (s *User) AfterFind(tx *gorm.DB) (err error) {
// 	if s.Password != "" {
// 		s.Password = ""
// 	}
// 	return
// }

func (s *User) EncryptPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("encrypt password error:", err)
		return err
	}
	s.Password = string(hash)
	return nil
}

func (s *User) SetPassword(tempPw string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(tempPw), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("set password error:", err)
		return err
	}
	db.Model(s).Update("Password", string(hash))
	return nil
}

func (s *User) CheckPassword() bool {
	if s.Username == "" {
		return false
	}
	var userPassword string
	row := db.Table("user").Where("username = ?", s.Username).Select("password").Row()
	err := row.Scan(&userPassword)
	if err == nil {
		err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(s.Password))
		if err != nil {
			logrus.Warn("password error:", s.Username)
			return false
		} else {
			logrus.Info("password ok:", s.Username)
			return true
		}
	} else {
		logrus.Warn("get user password error: ", s.Username)
		logrus.Error(err)
	}
	return false
}

func (s *User) LoadPermAssociationIds() error {
	var permissionIds []uint
	result := db.Table("user_permission").Where("user_id = ?", s.ID).Pluck("permission_id", &permissionIds)
	if result.Error != nil {
		logrus.Error(result.Error)
		logrus.Error(result.Error)
		return result.Error
	}
	s.PermissionIds = permissionIds
	return nil
}

func (s *User) LoadGroupAssociationIds() error {
	var groupIds []uint
	result := db.Table("user_group").Where("user_id = ?", s.ID).Pluck("group_id", &groupIds)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	s.GroupIds = groupIds
	return nil
}

func (s *User) LoadAllAssociationIds() error {
	var groupIds []uint
	result := db.Table("user_group").Where("user_id = ?", s.ID).Pluck("group_id", &groupIds)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	s.GroupIds = groupIds

	var permissionIds []uint
	result = db.Table("user_permission").Where("user_id = ?", s.ID).Pluck("permission_id", &permissionIds)
	if result.Error != nil {
		logrus.Error(result.Error)
		logrus.Error(result.Error)
		return result.Error
	}
	s.PermissionIds = permissionIds
	return nil
}

func (s *User) LoadAllAssociations() error {
	result := db.Preload("Groups").Preload("Permissions").First(s)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (s *User) Create() error {
	// add relations
	var groups []*Group
	var permissions []*Permission

	for _, groupId := range s.GroupIds {
		group := Group{ID: groupId}
		groups = append(groups, &group)
	}
	s.Groups = groups

	for _, permId := range s.PermissionIds {
		perm := Permission{ID: permId}
		permissions = append(permissions, &perm)
	}
	s.Permissions = permissions

	result := db.Create(s)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (s *User) Delete() error {
	if s.ID == 0 {
		return nil
	}
	// clear relations
	db.Model(s).Association("Groups").Clear()
	db.Model(s).Association("Permissions").Clear()
	result := db.Delete(s)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (s *User) Update() error {
	// replace relations
	var groups []*Group
	var permissions []*Permission

	if s.GroupIds != nil {
		for _, groupId := range s.GroupIds {
			group := Group{ID: groupId}
			groups = append(groups, &group)
		}
		result := db.Model(s).Association("Groups").Replace(groups)
		if result.Error != nil {
			logrus.Error(result.Error)
			return result.Error
		}
	}

	if s.PermissionIds != nil {
		for _, permId := range s.PermissionIds {
			perm := Permission{ID: permId}
			permissions = append(permissions, &perm)
		}
		result := db.Model(s).Association("Permissions").Replace(permissions)
		if result.Error != nil {
			logrus.Error(result.Error)
			return result.Error
		}
	}

	result := db.Model(s).Updates(s)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (s *Group) LoadAllAssociationIds() error {
	var userIds []uint
	result := db.Table("user_group").Where("group_id = ?", s.ID).Pluck("user_id", &userIds)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	s.UserIds = userIds

	var permissionIds []uint
	result = db.Table("group_permission").Where("group_id = ?", s.ID).Pluck("permission_id", &permissionIds)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	s.PermissionIds = permissionIds
	return nil
}

func (s *Group) LoadAllAssociations() error {
	result := db.Preload("Users").Preload("Permissions").First(s)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (s *Group) Create() error {
	// add relations
	var users []*User
	var permissions []*Permission

	for _, uerId := range s.UserIds {
		user := User{ID: uerId}
		users = append(users, &user)
	}
	s.Users = users

	for _, permId := range s.PermissionIds {
		perm := Permission{ID: permId}
		permissions = append(permissions, &perm)
	}
	s.Permissions = permissions

	result := db.Create(s)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (s *Group) Delete() error {
	if s.ID == 0 {
		return nil
	}
	// clear relations
	db.Model(s).Association("Users").Clear()
	db.Model(s).Association("Permissions").Clear()
	result := db.Delete(s)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (s *Group) Update() error {
	// replace relations
	var users []*User
	var permissions []*Permission

	if s.UserIds != nil {
		for _, uerId := range s.UserIds {
			user := User{ID: uerId}
			users = append(users, &user)
		}
		result := db.Model(s).Association("Users").Replace(users)
		if result.Error != nil {
			logrus.Error(result.Error)
			return result.Error
		}
	}

	if s.PermissionIds != nil {
		for _, permId := range s.PermissionIds {
			perm := Permission{ID: permId}
			permissions = append(permissions, &perm)
		}
		result := db.Model(s).Association("Permissions").Replace(permissions)
		if result.Error != nil {
			logrus.Error(result.Error)
			return result.Error
		}
	}

	result := db.Model(s).Updates(s)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (s *Permission) LoadGroupAssociationIds() error {
	var groupIds []uint
	result := db.Table("group_permission").Where("permission_id = ?", s.ID).Pluck("group_id", &groupIds)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	s.GroupIds = groupIds
	return nil
}

func (s *Permission) LoadAllAssociationIds() error {
	var groupIds []uint
	result := db.Table("group_permission").Where("permission_id = ?", s.ID).Pluck("group_id", &groupIds)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	s.GroupIds = groupIds

	var userIds []uint
	result = db.Table("user_permission").Where("permission_id = ?", s.ID).Pluck("user_id", &userIds)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	s.UserIds = userIds
	return nil
}

func (s *Permission) LoadAllAssociations() error {
	result := db.Preload("Users").Preload("Permissions").First(s)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (s *Permission) Create() error {
	// add relations
	var groups []*Group
	var users []*User

	for _, groupId := range s.GroupIds {
		group := Group{ID: groupId}
		groups = append(groups, &group)
	}
	s.Groups = groups

	for _, userId := range s.UserIds {
		user := User{ID: userId}
		users = append(users, &user)
	}
	s.Users = users

	result := db.Create(s)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (s *Permission) Delete() error {
	if s.ID == 0 {
		return nil
	}
	// clear relations
	db.Model(s).Association("Groups").Clear()
	db.Model(s).Association("Users").Clear()
	result := db.Delete(s)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}

func (s *Permission) Update() error {
	// replace relations
	var groups []*Group
	var users []*User

	if s.GroupIds != nil {
		for _, groupId := range s.GroupIds {
			group := Group{ID: groupId}
			groups = append(groups, &group)
		}
		result := db.Model(s).Association("Groups").Replace(groups)
		if result.Error != nil {
			logrus.Error(result.Error)
			return result.Error
		}
	}

	if s.UserIds != nil {
		for _, userId := range s.UserIds {
			user := User{ID: userId}
			users = append(users, &user)
		}
		result := db.Model(s).Association("Users").Replace(users)
		if result.Error != nil {
			logrus.Error(result.Error)
			return result.Error
		}
	}

	result := db.Model(s).Updates(s)
	if result.Error != nil {
		logrus.Error(result.Error)
		return result.Error
	}
	return nil
}
