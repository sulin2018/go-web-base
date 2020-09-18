package models

import (
	"time"

	logs "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          uint       `gorm:"primaryKey" uri:"id" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
	Username    string     `gorm:"type:varchar(50);not null;unique" description:"用户名" json:"username"`
	Password    string     `gorm:"type:varchar(100)" json:"password,omitempty"`
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
	Description string `gorm:"type:varchar(30)" json:"description"`

	Users    []*User  `gorm:"many2many:user_permission" json:"users"`
	UserIds  []uint   `gorm:"-" json:"user_ids"`
	Groups   []*Group `gorm:"many2many:group_permission" json:"groups"`
	GroupIds []uint   `gorm:"-" json:"group_ids"`
}

type Group struct {
	ID          uint   `gorm:"primaryKey" uri:"id" json:"id"`
	Name        string `gorm:"type:varchar(30);not null" json:"name"`
	Description string `gorm:"type:varchar(30)" json:"description"`

	Permissions   []*Permission `gorm:"many2many:group_permission" json:"permissions"`
	PermissionIds []uint        `gorm:"-" json:"permission_ids"`
	Users         []*User       `gorm:"many2many:user_group" json:"users"`
	UserIds       []uint        `gorm:"-" json:"user_ids"`
}

func (s *User) EncryptPassword() {
	hash, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
	if err != nil {
		logs.Error("encrypt password error:", err)
	}
	s.Password = string(hash)
}

func (s *User) SetPassword(tempPw string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(tempPw), bcrypt.DefaultCost)
	if err != nil {
		logs.Error("set password error:", err)
	}
	db.Model(s).Update("Password", string(hash))
}

func (s *User) CheckPassword() bool {
	if s.Username == "" {
		return false
	}
	var userPassword []string
	db.Table("user").Where("username = ?", s.Username).Pluck("password", &userPassword)
	if len(userPassword) > 0 {
		err := bcrypt.CompareHashAndPassword([]byte(userPassword[0]), []byte(s.Password))
		if err != nil {
			logs.Warn("password error:", s.Username)
			return false
		} else {
			logs.Info("password ok:", s.Username)
			return true
		}
	}
	return false
}

func (s *User) GetAllAssociationIds() {
	var groupIds []uint
	db.Table("user_group").Where("user_id = ?", s.ID).Pluck("group_id", &groupIds)
	s.GroupIds = groupIds

	var permissionIds []uint
	db.Table("user_permission").Where("user_id = ?", s.ID).Pluck("permission_id", &permissionIds)
	s.PermissionIds = permissionIds
}

func (s *User) GetAllAssociations() {
	db.Preload("Permissions").Preload("Groups").First(s)
}

func (s *User) Create() {
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

	db.Create(s)
}

func (s *User) Delete() {
	if s.ID == 0 {
		return
	}
	// clear relations
	db.Model(s).Association("Groups").Clear()
	db.Model(s).Association("Permissions").Clear()
	db.Delete(s)
}

func (s *User) Update() {
	// replace relations
	var groups []*Group
	var permissions []*Permission

	if s.GroupIds != nil {
		for _, groupId := range s.GroupIds {
			group := Group{ID: groupId}
			groups = append(groups, &group)
		}
		db.Model(s).Association("Groups").Replace(groups)
	}

	if s.PermissionIds != nil {
		for _, permId := range s.PermissionIds {
			perm := Permission{ID: permId}
			permissions = append(permissions, &perm)
		}
		db.Model(s).Association("Permissions").Replace(permissions)
	}

	db.Model(s).Updates(s)
}

func (s *Group) Page(results interface{}, page uint) {
	if page == 0 {
		db.Find(results)
	} else {
		db.Scopes(DBPage(page)).Find(results)
	}
}

func (s *Group) Detail() {
	db.First(s)
	s.GetAllAssociationIds()
}

func (s *Group) CreateOrUpdate() {
	if Exist(s) {
		s.Update()
	} else {
		s.Create()
	}
}

func (s *Group) GetAllAssociationIds() {
	var userIds []uint
	db.Table("user_group").Where("group_id = ?", s.ID).Pluck("user_id", &userIds)
	s.UserIds = userIds

	var permissionIds []uint
	db.Table("group_permission").Where("group_id = ?", s.ID).Pluck("permission_id", &permissionIds)
	s.PermissionIds = permissionIds
}

func (s *Group) GetAllAssociations() {
	db.Preload("Permissions").Preload("Users").First(s)
}

func (s *Group) AppendPermission(obj interface{}) {
	db.Model(s).Association("Permissions").Append(obj)
}

func (s *Group) Create() {
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

	db.Create(s)
}

func (s *Group) Delete() {
	if s.ID == 0 {
		return
	}
	// clear relations
	db.Model(s).Association("Users").Clear()
	db.Model(s).Association("Permissions").Clear()
	db.Delete(s)
}

func (s *Group) Update() {
	// replace relations
	var users []*User
	var permissions []*Permission

	if s.UserIds != nil {
		for _, uerId := range s.UserIds {
			user := User{ID: uerId}
			users = append(users, &user)
		}
		db.Model(s).Association("Users").Replace(users)
	}

	if s.PermissionIds != nil {
		for _, permId := range s.PermissionIds {
			perm := Permission{ID: permId}
			permissions = append(permissions, &perm)
		}
		db.Model(s).Association("Permissions").Replace(permissions)
	}

	db.Model(s).Updates(s)
}

func (s *Permission) GetAllAssociationIds() {
	var groupIds []uint
	db.Table("group_permission").Where("permission_id = ?", s.ID).Pluck("group_id", &groupIds)
	s.GroupIds = groupIds

	var userIds []uint
	db.Table("user_permission").Where("permission_id = ?", s.ID).Pluck("user_id", &userIds)
	s.UserIds = userIds
}

func (s *Permission) GetAllAssociations() {
	db.Preload("Users").Preload("Groups").First(s)
}

func (s *Permission) Create() {
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

	db.Create(s)
}

func (s *Permission) Delete() {
	if s.ID == 0 {
		return
	}
	// clear relations
	db.Model(s).Association("Groups").Clear()
	db.Model(s).Association("Users").Clear()
	db.Delete(s)
}

func (s *Permission) Update() {
	// replace relations
	var groups []*Group
	var users []*User

	if s.GroupIds != nil {
		for _, groupId := range s.GroupIds {
			group := Group{ID: groupId}
			groups = append(groups, &group)
		}
		db.Model(s).Association("Groups").Replace(groups)
	}

	if s.UserIds != nil {
		for _, userId := range s.UserIds {
			user := User{ID: userId}
			users = append(users, &user)
		}
		db.Model(s).Association("Users").Replace(users)
	}

	db.Model(s).Updates(s)
}
