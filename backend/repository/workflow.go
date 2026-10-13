package repository

import (
	"backend/exceptions"
	"backend/model"
	"errors"

	"github.com/jinzhu/gorm"
)

func CreateWorkflow(userId uint, name string, description string) (*model.Workflow, *exceptions.Error) {
	w := model.Workflow{
		UserId:      userId,
		Name:        name,
		Description: description,
		Clicks:      0,
		Content:     "{\"content\":{\"nodes\":[{\"id\":\"1\",\"type\":\"start\",\"dimensions\":{\"width\":250,\"height\":80},\"computedPosition\":{\"x\":100,\"y\":100,\"z\":0},\"handleBounds\":{\"source\":[{\"id\":null,\"position\":\"right\",\"nodeId\":\"1\",\"type\":\"source\",\"x\":245,\"y\":35,\"width\":10,\"height\":10}],\"target\":[]},\"selected\":false,\"dragging\":false,\"resizing\":false,\"initialized\":false,\"isParent\":false,\"position\":{\"x\":100,\"y\":100},\"data\":{\"name\":\"开始\",\"input\":[]},\"events\":{}},{\"id\":\"2\",\"type\":\"end\",\"dimensions\":{\"width\":250,\"height\":80},\"computedPosition\":{\"x\":700,\"y\":100,\"z\":0},\"handleBounds\":{\"source\":[],\"target\":[{\"id\":null,\"position\":\"left\",\"nodeId\":\"2\",\"type\":\"target\",\"x\":-5,\"y\":35,\"width\":10,\"height\":10}]},\"selected\":false,\"dragging\":false,\"resizing\":false,\"initialized\":false,\"isParent\":false,\"position\":{\"x\":700,\"y\":100},\"data\":{\"name\":\"结束\",\"output\":[]},\"events\":{}}],\"edges\":[]}}",
	}

	if err := db.Create(&w).Error; err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	return &w, nil
}

func DeleteWorkflow(userId, workflowId uint) *exceptions.Error {
	var w model.Workflow
	if err := db.First(&w, workflowId).Error; err != nil {
		return exceptions.InternalServerError(err)
	}

	if w.UserId != userId {
		return exceptions.ErrPermissionDenied.WithDetail("You are not authorized to delete this workflow")
	}

	if err := db.Delete(&w).Error; err != nil {
		return exceptions.InternalServerError(err)
	}
	return nil
}

func DeleteCorrespondingRecordsAssociatedWithWorkflow(workflowId uint) *exceptions.Error {
	if err1 := db.Where("workflow_id = ?", workflowId).Delete(&model.FavoriteWorkflowRecord{}).Error; err1 != nil {
		return exceptions.InternalServerError(err1)
	}

	if err2 := db.Where("workflow_id = ?", workflowId).Delete(&model.LikeWorkflowRecord{}).Error; err2 != nil {
		return exceptions.InternalServerError(err2)
	}
	return nil
}

func GetWorkflowList(userId uint) ([]*model.Workflow, *exceptions.Error) {
	var workflows []*model.Workflow
	if err := db.Preload("BelongsTo").Where("user_id = ?", userId).Order("updated_at DESC").Find(&workflows).Error; err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	return workflows, nil
}

func GetAllPublicWorkflows(page, pageSize uint) ([]*model.Workflow, *exceptions.Error) {
	var workflows []*model.Workflow
	err := db.Table("workflows").
		Select(`
			workflows.id, workflows.created_at, workflows.updated_at, workflows.deleted_at,
			workflows.user_id, workflows.name, workflows.description, workflows.cover, workflows.clicks, workflows.status, workflows.content,
			COUNT(DISTINCT l.id) AS likes,
			COUNT(DISTINCT f.id) AS collects
		`).
		Joins("LEFT JOIN like_workflow_records l ON l.workflow_id = workflows.id AND l.deleted_at IS NULL").
		Joins("LEFT JOIN favorite_workflow_records f ON f.workflow_id = workflows.id AND f.deleted_at IS NULL").
		Where("status = ?", model.StatusPublic).
		Group("id").
		Order("clicks DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Preload("BelongsTo").
		Find(&workflows).Error
	// Note that `db.Raw` cannot use Preload method and I don't know why...
	//err := db.Raw(`
	//		SELECT
	//			w.id, w.created_at, w.updated_at, w.deleted_at,
	//			w.user_id, w.name, w.description, w.cover, w.clicks, w.status, w.content,
	//			COUNT(l.id) AS likes,
	//			COUNT(f.id) AS collects
	//		FROM workflows w
	//		LEFT JOIN like_workflow_records l ON l.workflow_id = w.id AND l.deleted_at IS NULL
	//		LEFT JOIN favorite_workflow_records f ON f.workflow_id = w.id AND f.deleted_at IS NULL
	//		WHERE w.deleted_at IS NULL AND w.status = ?
	//		GROUP BY w.id
	//		ORDER BY w.created_at DESC
	//		LIMIT ? OFFSET ?
	//	`, model.StatusPublic, pageSize, (page-1)*pageSize).
	//	Scan(&workflows).
	//	Preload("BelongsTo").Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // FIXME: nil or empty slice?
		}
		return nil, exceptions.InternalServerError(err)
	}
	return workflows, nil
}

func GetPublicWorkflowsWithSearchKeywords(page, pageSize uint, keyword string) ([]*model.Workflow, *exceptions.Error) {
	var workflows []*model.Workflow
	err := db.Where("status = ? AND name LIKE ?", model.StatusPublic, "%"+keyword+"%").
		Order("clicks DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Preload("BelongsTo").
		Find(&workflows).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, exceptions.InternalServerError(err)
	}
	return workflows, nil
}

func GetPopularWorkflows(pageSize uint) ([]*model.Workflow, *exceptions.Error) {
	var workflows []*model.Workflow
	err := db.Order("clicks DESC").Limit(pageSize).Preload("BelongsTo").Find(&workflows).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, exceptions.InternalServerError(err)
	}
	return workflows, nil
}

func GetWorkflow(userId, workflowId uint) (*model.Workflow, *exceptions.Error) {
	var w model.Workflow
	if err := db.Preload("BelongsTo").First(&w, workflowId).Error; err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	if err := w.ParseContent(); err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	w.Clicks += 1
	db.Save(&w)

	if checkLikeWorkflowStatus(userId, workflowId) {
		w.LikeStatus = true
	}

	if checkWorkflowFavoriteStatus(userId, workflowId) {
		w.CollectStatus = true
	}

	likeCnt, err := getWorkflowLikeCount(workflowId)
	if err != nil {
		return nil, err
	}
	w.Likes = likeCnt

	favouriteCnt, err := getWorkflowFavouriteCount(workflowId)
	if err != nil {
		return nil, err
	}
	w.Collects = favouriteCnt

	return &w, nil
}

func UpdateWorkflow(workflow *model.Workflow) (*model.Workflow, *exceptions.Error) {
	if err := db.Save(workflow).Error; err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	return workflow, nil
}

// checkWorkflowFavoriteStatus returns true if the user favorites the workflow, false otherwise.
func checkWorkflowFavoriteStatus(userId, workflowId uint) bool {
	var record model.FavoriteWorkflowRecord
	err := db.Where("user_id = ? AND workflow_id = ?", userId, workflowId).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		panic(err) // FIXME
	}
	return true
}

func FavoriteWorkflow(userId, workflowId uint) (bool, *exceptions.Error) {
	var user model.User
	if err := db.First(&user, userId).Error; err != nil {
		return false, exceptions.InternalServerError(err)
	}

	var workflow model.Workflow
	if err := db.First(&workflow, workflowId).Error; err != nil {
		return false, exceptions.InternalServerError(err)
	}

	if checkWorkflowFavoriteStatus(userId, workflowId) {
		db.Where("workflow_id = ? AND user_id = ?", workflowId, userId).Delete(&model.FavoriteWorkflowRecord{})
		return false, nil
	}

	var newRecord model.FavoriteWorkflowRecord
	newRecord.UserId = userId
	newRecord.WorkflowId = workflowId
	if err := db.Create(&newRecord).Error; err != nil {
		return false, exceptions.InternalServerError(err)
	}

	return true, nil
}

func GetFavoriteWorkflowList(userId uint) ([]*model.FavoriteWorkflowRecord, *exceptions.Error) {
	var favoriteWorkflows []*model.FavoriteWorkflowRecord
	err := db.Where("user_id = ?", userId).Preload("CorrespondingWorkflow").Find(&favoriteWorkflows).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, exceptions.InternalServerError(err)
	}
	return favoriteWorkflows, nil
}

func getUserFavouriteWorkflowCount(userId uint) (uint, *exceptions.Error) {
	var count uint
	err := db.Model(&model.FavoriteWorkflowRecord{}).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, exceptions.InternalServerError(err)
	}
	return count, nil
}

func getUserPublishedWorkflowCount(userId uint) (uint, *exceptions.Error) {
	var count uint
	err := db.Model(&model.Workflow{}).Where("user_id = ? AND status = ?", userId, model.StatusPublic).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, exceptions.InternalServerError(err)
	}
	return count, nil
}

func getUserFollowersCount(userId uint) (uint, *exceptions.Error) {
	// TODO
	return 0, nil
}

func getWorkflowFavouriteCount(workflowId uint) (uint, *exceptions.Error) {
	var count uint
	err := db.Model(&model.FavoriteWorkflowRecord{}).Where("workflow_id = ?", workflowId).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, exceptions.InternalServerError(err)
	}
	return count, nil
}

func checkLikeWorkflowStatus(userId, workflowId uint) bool {
	var record model.LikeWorkflowRecord
	err := db.Where("user_id = ? AND workflow_id = ?", userId, workflowId).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		panic(err)
	}
	return true
}

func LikeWorkflow(userId, workflowId uint) (bool, *exceptions.Error) {
	var user model.User
	if err := db.First(&user, userId).Error; err != nil {
		return false, exceptions.InternalServerError(err)
	}

	var workflow model.Workflow
	if err := db.First(&workflow, workflowId).Error; err != nil {
		return false, exceptions.InternalServerError(err)
	}

	if checkLikeWorkflowStatus(userId, workflowId) {
		db.Where("workflow_id = ? AND user_id = ?", workflowId, userId).Delete(&model.LikeWorkflowRecord{})
		return false, nil
	}

	var newRecord model.LikeWorkflowRecord
	newRecord.UserId = userId
	newRecord.WorkflowId = workflowId
	if err := db.Create(&newRecord).Error; err != nil {
		return false, exceptions.InternalServerError(err)
	}

	return true, nil
}

func getWorkflowLikeCount(workflowId uint) (uint, *exceptions.Error) {
	var count uint
	err := db.Model(&model.LikeWorkflowRecord{}).Where("workflow_id = ?", workflowId).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, exceptions.InternalServerError(err)
	}
	return count, nil
}
