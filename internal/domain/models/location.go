package models

import "time"

type Province struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Code      string    `json:"code" gorm:"unique;not null"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    
    Cities    []City    `json:"cities,omitempty"`
}

type City struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    ProvinceID uint      `json:"province_id"`
    Province   Province  `json:"province" gorm:"foreignKey:ProvinceID"`
    Name       string    `json:"name" gorm:"not null"`
    Code       string    `json:"code" gorm:"unique;not null"`
    Type       string    `json:"type" gorm:"comment:'kabupaten or kota'"`
    CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    
    Districts  []District `json:"districts,omitempty"`
}

type District struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    CityID    uint      `json:"city_id"`
    City      City      `json:"city" gorm:"foreignKey:CityID"`
    Name      string    `json:"name" gorm:"not null"`
    Code      string    `json:"code" gorm:"unique;not null"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Address struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    UserID      uint      `json:"user_id"`
    User        User      `json:"-" gorm:"foreignKey:UserID"`
    Name        string    `json:"name" gorm:"not null"`
    PhoneNumber string    `json:"phone_number" gorm:"not null"`
    ProvinceID  uint      `json:"province_id"`
    Province    Province  `json:"province" gorm:"foreignKey:ProvinceID"`
    CityID      uint      `json:"city_id"`
    City        City      `json:"city" gorm:"foreignKey:CityID"`
    DistrictID  uint      `json:"district_id"`
    District    District  `json:"district" gorm:"foreignKey:DistrictID"`
    PostalCode  string    `json:"postal_code" gorm:"not null"`
    DetailAddress string    `json:"detail_address" gorm:"not null"`
    IsDefault   bool      `json:"is_default" gorm:"default:false"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   *time.Time `json:"-" gorm:"index"`
}

type Courier struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Code      string    `json:"code" gorm:"unique;not null"`
    Name      string    `json:"name" gorm:"not null"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}