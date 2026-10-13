package repository

import (
	dbPkg "backend/db"
	"backend/exceptions"
	"backend/model"
)

func CreateKnowledgeFile(file *model.KnowledgeFile) *exceptions.Error {
	result := dbPkg.DB.Create(file)
	if result.Error != nil {
		return exceptions.InternalServerError(result.Error)
	}
	return nil
}

func GetKnowledgeFilesByUserId(userId uint) ([]*model.KnowledgeFile, *exceptions.Error) {
	var files []*model.KnowledgeFile
	result := dbPkg.DB.Where("user_id = ?", userId).Find(&files)
	if result.Error != nil {
		return nil, exceptions.InternalServerError(result.Error)
	}
	return files, nil
}

func GetKnowledgeFileById(fileId uint) (*model.KnowledgeFile, *exceptions.Error) {
	var file model.KnowledgeFile
	
	result := dbPkg.DB.First(&file, fileId)
	if result.Error != nil {
		return nil, exceptions.InternalServerError(result.Error)
	}
	return &file, nil
}

func DeleteKnowledgeFile(fileId uint) *exceptions.Error {
	result := dbPkg.DB.Where("id = ?", fileId).Delete(&model.KnowledgeFile{})
	if result.Error != nil {
		return exceptions.InternalServerError(result.Error)
	}
	return nil
} 