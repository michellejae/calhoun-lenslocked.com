package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func NewServices(dialect, connectionInfo string) (*Services, error) {
	//ToDo: config this
	db, err := gorm.Open(dialect, connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryService(db),
		Image:   NewImageService(),
		db:      db,
	}, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

// closes the DB connection
func (s *Services) Close() error {
	return s.db.Close()
}

// drops all tables and rebuilds them
func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()

}

// Automigrate will attempt to automatiaclly migrate the all table
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error

}
