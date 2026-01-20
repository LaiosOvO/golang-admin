package dict

import (
	"fmt"
	"gorm.io/gorm"
)

// DictType 字典类型
type DictType struct {
	ID       int    `gorm:"primarykey" json:"id"`
	Name     string `gorm:"size:100;not null" json:"name"`
	Type     string `gorm:"size:100;not null;uniqueIndex" json:"type"`
	Status   int    `gorm:"default:1" json:"status"` // 0-禁用 1-启用
	Remark   string `gorm:"size:500" json:"remark"`
	CreateBy uint   `json:"createBy"`
	UpdateBy uint   `json:"updateBy"`
	CreateAt uint   `json:"createAt"`
	UpdateAt uint   `json:"updateAt"`
}

// DictData 字典数据
type DictData struct {
	ID       int    `gorm:"primarykey" json:"id"`
	DictSort int    `gorm:"default:0" json:"dictSort"`
	Label    string `gorm:"size:100;not null" json:"label"`
	Value    string `gorm:"size:100;not null" json:"value"`
	DictType string `gorm:"size:100;not null;index" json:"dictType"`
	Status   int    `gorm:"default:1" json:"status"` // 0-禁用 1-启用
	Remark   string `gorm:"size:500" json:"remark"`
	CreateBy uint   `json:"createBy"`
	UpdateBy uint   `json:"updateBy"`
	CreateAt uint   `json:"createAt"`
	UpdateAt uint   `json:"updateAt"`
}

// TableName 设置表名
func (DictType) TableName() string {
	return "system_dict_type"
}

func (DictData) TableName() string {
	return "system_dict_data"
}

// Service 字典服务
type Service struct {
	db *gorm.DB
}

// NewService 创建字典服务
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// GetDictTypes 获取字典类型列表
func (s *Service) GetDictTypes(page, pageSize int, name, dictType string, status int) ([]DictType, int64, error) {
	var dictTypes []DictType
	var total int64

	query := s.db.Model(&DictType{})

	// 添加查询条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if dictType != "" {
		query = query.Where("type LIKE ?", "%"+dictType+"%")
	}
	if status != -1 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&dictTypes).Error
	if err != nil {
		return nil, 0, err
	}

	return dictTypes, total, nil
}

// GetDictTypeByID 根据ID获取字典类型
func (s *Service) GetDictTypeByID(id int) (*DictType, error) {
	var dictType DictType
	err := s.db.First(&dictType, id).Error
	if err != nil {
		return nil, err
	}
	return &dictType, nil
}

// GetDictTypeByType 根据类型获取字典类型
func (s *Service) GetDictTypeByType(dictType string) (*DictType, error) {
	var dt DictType
	err := s.db.Where("type = ?", dictType).First(&dt).Error
	if err != nil {
		return nil, err
	}
	return &dt, nil
}

// CreateDictType 创建字典类型
func (s *Service) CreateDictType(dictType *DictType) error {
	// 检查类型是否已存在
	var count int64
	err := s.db.Model(&DictType{}).Where("type = ?", dictType.Type).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("字典类型 %s 已存在", dictType.Type)
	}

	return s.db.Create(dictType).Error
}

// UpdateDictType 更新字典类型
func (s *Service) UpdateDictType(dictType *DictType) error {
	// 检查类型是否被其他记录使用
	if dictType.Type != "" {
		var count int64
		err := s.db.Model(&DictType{}).Where("type = ? AND id != ?", dictType.Type, dictType.ID).Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("字典类型 %s 已存在", dictType.Type)
		}
	}

	return s.db.Save(dictType).Error
}

// DeleteDictType 删除字典类型
func (s *Service) DeleteDictType(id int) error {
	// 检查是否有关联的字典数据
	var count int64
	err := s.db.Model(&DictData{}).Where("dict_type = (SELECT type FROM system_dict_type WHERE id = ?)", id).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("该字典类型下存在字典数据，无法删除")
	}

	return s.db.Delete(&DictType{}, id).Error
}

// GetDictDataList 获取字典数据列表
func (s *Service) GetDictDataList(page, pageSize int, dictType string, label string, status int) ([]DictData, int64, error) {
	var dictData []DictData
	var total int64

	query := s.db.Model(&DictData{})

	// 添加查询条件
	if dictType != "" {
		query = query.Where("dict_type = ?", dictType)
	}
	if label != "" {
		query = query.Where("label LIKE ?", "%"+label+"%")
	}
	if status != -1 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("dict_sort ASC, id DESC").Find(&dictData).Error
	if err != nil {
		return nil, 0, err
	}

	return dictData, total, nil
}

// GetDictDataByID 根据ID获取字典数据
func (s *Service) GetDictDataByID(id int) (*DictData, error) {
	var dictData DictData
	err := s.db.First(&dictData, id).Error
	if err != nil {
		return nil, err
	}
	return &dictData, nil
}

// GetDictDataByType 根据字典类型获取所有启用的字典数据
func (s *Service) GetDictDataByType(dictType string) ([]DictData, error) {
	var dictData []DictData
	err := s.db.Where("dict_type = ? AND status = 1", dictType).
		Order("dict_sort ASC, id ASC").
		Find(&dictData).Error
	if err != nil {
		return nil, err
	}
	return dictData, nil
}

// CreateDictData 创建字典数据
func (s *Service) CreateDictData(dictData *DictData) error {
	// 检查字典类型是否存在
	var dictType DictType
	err := s.db.Where("type = ? AND status = 1", dictData.DictType).First(&dictType).Error
	if err != nil {
		return fmt.Errorf("字典类型 %s 不存在或已禁用", dictData.DictType)
	}

	return s.db.Create(dictData).Error
}

// UpdateDictData 更新字典数据
func (s *Service) UpdateDictData(dictData *DictData) error {
	// 检查字典类型是否存在
	if dictData.DictType != "" {
		var dictType DictType
		err := s.db.Where("type = ? AND status = 1", dictData.DictType).First(&dictType).Error
		if err != nil {
			return fmt.Errorf("字典类型 %s 不存在或已禁用", dictData.DictType)
		}
	}

	return s.db.Save(dictData).Error
}

// DeleteDictData 删除字典数据
func (s *Service) DeleteDictData(id int) error {
	return s.db.Delete(&DictData{}, id).Error
}

// GetDictValueByLabel 根据字典类型和标签获取值
func (s *Service) GetDictValueByLabel(dictType, label string) (string, error) {
	var dictData DictData
	err := s.db.Where("dict_type = ? AND label = ? AND status = 1", dictType, label).
		First(&dictData).Error
	if err != nil {
		return "", err
	}
	return dictData.Value, nil
}

// GetDictLabelByValue 根据字典类型和值获取标签
func (s *Service) GetDictLabelByValue(dictType, value string) (string, error) {
	var dictData DictData
	err := s.db.Where("dict_type = ? AND value = ? AND status = 1", dictType, value).
		First(&dictData).Error
	if err != nil {
		return "", err
	}
	return dictData.Label, nil
}

// RefreshDictCache 刷新字典缓存
func (s *Service) RefreshDictCache() error {
	// 这里可以实现缓存刷新逻辑
	// 比如清除Redis中的字典缓存
	return nil
}

// ExportDictType 导出字典类型
func (s *Service) ExportDictType(ids []int) ([]DictType, error) {
	var dictTypes []DictType
	query := s.db.Model(&DictType{})
	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}
	err := query.Find(&dictTypes).Error
	return dictTypes, err
}

// ExportDictData 导出字典数据
func (s *Service) ExportDictData(dictType string, ids []int) ([]DictData, error) {
	var dictData []DictData
	query := s.db.Model(&DictData{})
	if dictType != "" {
		query = query.Where("dict_type = ?", dictType)
	}
	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}
	err := query.Find(&dictData).Error
	return dictData, err
}

// GetDictTypesSimple 获取字典类型简单列表（用于下拉选择）
func (s *Service) GetDictTypesSimple() ([]DictType, error) {
	var dictTypes []DictType
	err := s.db.Model(&DictType{}).
		Select("id", "name", "type").
		Where("status = 1").
		Find(&dictTypes).Error
	return dictTypes, err
}

// GetDictDataSimple 根据类型获取字典数据简单列表（用于下拉选择）
func (s *Service) GetDictDataSimple(dictType string) ([]DictData, error) {
	var dictData []DictData
	err := s.db.Model(&DictData{}).
		Select("label", "value").
		Where("dict_type = ? AND status = 1", dictType).
		Order("dict_sort ASC, id ASC").
		Find(&dictData).Error
	return dictData, err
}

// ValidateDictType 验证字典类型
func (s *Service) ValidateDictType(dictType string) error {
	var count int64
	err := s.db.Model(&DictType{}).Where("type = ? AND status = 1", dictType).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("字典类型 %s 不存在或已禁用", dictType)
	}
	return nil
}
