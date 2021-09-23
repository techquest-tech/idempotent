package gorm

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Idempotent struct {
	ObjKey    string `gorm:"primaryKey"`
	CreatedAt time.Time
}

type IdempotentGormService struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func (s *IdempotentGormService) key2str(key interface{}) string {
	str, ok := key.(string)
	if !ok {
		str = fmt.Sprintf("%s", key)
		s.Logger.Info("convert key to string.", zap.String("result", str), zap.Any("origin", key))
	}
	return str
}

func (s *IdempotentGormService) Duplicated(key interface{}) (bool, error) {
	item := &Idempotent{}
	str := s.key2str(key)

	err := s.DB.First(item, str)
	if err != nil || item.ObjKey == "" {
		return false, nil
	}
	s.Logger.Debug("key found.", zap.Time("createdAt", item.CreatedAt), zap.String("key", str))
	return true, nil
}

func (s *IdempotentGormService) Save(key interface{}) error {
	item := &Idempotent{
		ObjKey: s.key2str(key),
	}
	err := s.DB.Save(item).Error
	if err != nil {
		s.Logger.Error("save Idempotent item failed.", zap.Any("error", err))
		return err
	}
	s.Logger.Info("save idempontent item done", zap.String("key", item.ObjKey))
	return nil
}

func (s *IdempotentGormService) AllKeys() ([]interface{}, error) {
	items := make([]Idempotent, 1000)
	err := s.DB.Limit(1000).Order("CreatedAt desc").Find(&items).Error
	if err != nil {
		s.Logger.Error("get last 1000 keys failed.", zap.Any("error", err))
		return nil, err
	}
	keys := make([]interface{}, len(items))
	for index, item := range items {
		keys[index] = item.ObjKey
	}
	return keys, nil
}
