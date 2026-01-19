package system

import (
	"errors"
	"gin-admin-pro/internal/dao/system"

	"gorm.io/gorm"
)

// MenuService 菜单服务层
type MenuService struct {
	menuDAO *system.MenuDAO
}

// NewMenuService 创建菜单服务实例
func NewMenuService(menuDAO *system.MenuDAO) *MenuService {
	return &MenuService{menuDAO: menuDAO}
}

// GetList 获取菜单列表
func (s *MenuService) GetList(req *system.MenuListReq) ([]system.MenuResp, error) {
	return s.menuDAO.GetList(req)
}

// GetByID 根据ID获取菜单详情
func (s *MenuService) GetByID(id uint) (*system.MenuDetailResp, error) {
	if id == 0 {
		return nil, errors.New("菜单ID不能为空")
	}

	menu, err := s.menuDAO.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("菜单不存在")
		}
		return nil, err
	}

	return menu, nil
}

// Create 创建菜单
func (s *MenuService) Create(req *system.CreateMenuReq, createBy uint) (uint, error) {
	// 参数验证
	if req.Name == "" {
		return 0, errors.New("菜单名称不能为空")
	}
	if req.Type < 1 || req.Type > 3 {
		return 0, errors.New("菜单类型不正确")
	}

	// 检查同级下菜单名称是否重复
	exists, err := s.menuDAO.CheckNameExists(req.Name, req.ParentID, nil)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("同级下已存在相同名称的菜单")
	}

	// 如果没有指定排序，自动获取最大排序值+1
	if req.Sort == 0 {
		maxSort, err := s.menuDAO.GetMaxSort(req.ParentID)
		if err != nil {
			return 0, err
		}
		req.Sort = maxSort + 1
	}

	return s.menuDAO.Create(req, createBy)
}

// Update 更新菜单
func (s *MenuService) Update(req *system.UpdateMenuReq, updateBy uint) error {
	if req.ID == 0 {
		return errors.New("菜单ID不能为空")
	}

	// 获取原菜单信息
	menu, err := s.menuDAO.GetByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("菜单不存在")
		}
		return err
	}

	// 检查是否将菜单设置为自己的子菜单
	if req.ParentID != nil && *req.ParentID == req.ID {
		return errors.New("不能将菜单设置为自己的子菜单")
	}

	// 检查是否造成循环引用
	if req.ParentID != nil && *req.ParentID != menu.ParentID {
		if err := s.checkCircularReference(req.ID, *req.ParentID); err != nil {
			return err
		}
	}

	// 如果修改了名称或父级，检查重复
	parentID := menu.ParentID
	if req.ParentID != nil {
		parentID = *req.ParentID
	}

	if (req.Name != "" && req.Name != menu.Name) || (req.ParentID != nil && *req.ParentID != menu.ParentID) {
		name := req.Name
		if name == "" {
			name = menu.Name
		}

		exists, err := s.menuDAO.CheckNameExists(name, parentID, &req.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("同级下已存在相同名称的菜单")
		}
	}

	return s.menuDAO.Update(req, updateBy)
}

// Delete 删除菜单
func (s *MenuService) Delete(id uint) error {
	if id == 0 {
		return errors.New("菜单ID不能为空")
	}

	err := s.menuDAO.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("菜单不存在或存在子菜单，无法删除")
		}
		return err
	}

	return nil
}

// GetAllSimpleList 获取所有菜单简单列表（用于角色授权）
func (s *MenuService) GetAllSimpleList() ([]system.MenuSimpleResp, error) {
	return s.menuDAO.GetAllSimpleList()
}

// GetUserMenus 获取用户菜单列表
func (s *MenuService) GetUserMenus(userID uint) ([]system.MenuResp, error) {
	if userID == 0 {
		return nil, errors.New("用户ID不能为空")
	}

	return s.menuDAO.GetUserMenus(userID)
}

// checkCircularReference 检查循环引用
func (s *MenuService) checkCircularReference(menuID, parentID uint) error {
	// 获取父级菜单链
	parentChain, err := s.menuDAO.GetParentChain(parentID)
	if err != nil {
		return err
	}

	// 检查当前菜单是否在父级菜单链中
	for _, parent := range parentChain {
		if parent.ID == menuID {
			return errors.New("不能将菜单设置为自己的子菜单，会造成循环引用")
		}
	}

	return nil
}
