package system

import (
	"errors"
	"gin-admin-pro/internal/dao/system"
	usermodel "gin-admin-pro/internal/model/system"

	"gorm.io/gorm"
)

// DeptService 部门服务层
type DeptService struct {
	deptDAO *system.DeptDAO
}

// NewDeptService 创建部门服务实例
func NewDeptService(deptDAO *system.DeptDAO) *DeptService {
	return &DeptService{deptDAO: deptDAO}
}

// GetList 获取部门列表
func (s *DeptService) GetList(req *system.DeptListReq) ([]system.DeptResp, error) {
	return s.deptDAO.GetList(req)
}

// GetByID 根据ID获取部门详情
func (s *DeptService) GetByID(id uint) (*system.DeptDetailResp, error) {
	if id == 0 {
		return nil, errors.New("部门ID不能为空")
	}

	dept, err := s.deptDAO.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("部门不存在")
		}
		return nil, err
	}

	return dept, nil
}

// Create 创建部门
func (s *DeptService) Create(req *system.CreateDeptReq, createBy uint) (uint, error) {
	// 参数验证
	if req.Name == "" {
		return 0, errors.New("部门名称不能为空")
	}

	// 检查同级下部门名称是否重复
	exists, err := s.deptDAO.CheckNameExists(req.Name, req.ParentID, nil)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("同级下已存在相同名称的部门")
	}

	// 如果没有指定排序，自动获取最大排序值+1
	if req.Sort == 0 {
		maxSort, err := s.deptDAO.GetMaxSort(req.ParentID)
		if err != nil {
			return 0, err
		}
		req.Sort = maxSort + 1
	}

	// 验证负责人是否存在
	if req.Leader != 0 {
		// 这里需要查询用户表，但为了不产生循环依赖，可以在DAO层处理
		// 或者通过其他方式验证
	}

	return s.deptDAO.Create(req, createBy)
}

// Update 更新部门
func (s *DeptService) Update(req *system.UpdateDeptReq, updateBy uint) error {
	if req.ID == 0 {
		return errors.New("部门ID不能为空")
	}

	// 获取原部门信息
	dept, err := s.deptDAO.GetByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("部门不存在")
		}
		return err
	}

	// 检查是否将部门设置为自己的子部门
	if req.ParentID != nil && *req.ParentID == req.ID {
		return errors.New("不能将部门设置为自己的子部门")
	}

	// 检查是否造成循环引用
	if req.ParentID != nil && *req.ParentID != dept.ParentID {
		if err := s.checkCircularReference(req.ID, *req.ParentID); err != nil {
			return err
		}
	}

	// 如果修改了名称或父级，检查重复
	parentID := dept.ParentID
	if req.ParentID != nil {
		parentID = *req.ParentID
	}

	if (req.Name != "" && req.Name != dept.Name) || (req.ParentID != nil && *req.ParentID != dept.ParentID) {
		name := req.Name
		if name == "" {
			name = dept.Name
		}

		exists, err := s.deptDAO.CheckNameExists(name, parentID, &req.ID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("同级下已存在相同名称的部门")
		}
	}

	// 验证负责人是否存在
	if req.Leader != nil && *req.Leader != 0 {
		// 需要验证用户是否存在
	}

	return s.deptDAO.Update(req, updateBy)
}

// Delete 删除部门
func (s *DeptService) Delete(id uint) error {
	if id == 0 {
		return errors.New("部门ID不能为空")
	}

	err := s.deptDAO.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("部门不存在或存在子部门，无法删除")
		}
		return err
	}

	return nil
}

// GetAllSimpleList 获取所有部门简单列表
func (s *DeptService) GetAllSimpleList() ([]system.DeptSimpleResp, error) {
	return s.deptDAO.GetAllSimpleList()
}

// GetUsersByDept 获取部门下的用户列表
func (s *DeptService) GetUsersByDept(deptID uint) ([]usermodel.User, error) {
	if deptID == 0 {
		return nil, errors.New("部门ID不能为空")
	}

	return s.deptDAO.GetUsersByDept(deptID)
}

// checkCircularReference 检查循环引用
func (s *DeptService) checkCircularReference(deptID, parentID uint) error {
	// 获取父级部门链
	parentChain, err := s.deptDAO.GetParentChain(parentID)
	if err != nil {
		return err
	}

	// 检查当前部门是否在父级部门链中
	for _, parent := range parentChain {
		if parent.ID == deptID {
			return errors.New("不能将部门设置为自己的子部门，会造成循环引用")
		}
	}

	return nil
}
