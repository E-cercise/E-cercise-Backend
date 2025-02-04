package service

import (
	"context"
	"fmt"
	"github.com/E-cercise/E-cercise/src/data/request"
	"github.com/E-cercise/E-cercise/src/data/response"
	"github.com/E-cercise/E-cercise/src/helper"
	"github.com/E-cercise/E-cercise/src/logger"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/E-cercise/E-cercise/src/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

type EquipmentService interface {
	GetEquipmentData(q request.EquipmentListRequest, paginatior *helper.Paginator) (*response.EquipmentsResponse, error)
	AddEquipment(req request.EquipmentPostRequest, context context.Context) error
	GetEquipmentDetail(eqID uuid.UUID) (*response.EquipmentDetailResponse, error)
	UpdateEquipment(eqID uuid.UUID, context context.Context, req request.EquipmentPutRequest) error
}

type equipmentService struct {
	db              *gorm.DB
	equipmentRepo   repository.EquipmentRepository
	muscleGroupRepo repository.MuscleGroupRepository
	imageService    ImageService
}

func NewEquipmentService(db *gorm.DB, equipmentRepo repository.EquipmentRepository, muscleGroupRepo repository.MuscleGroupRepository, imageService ImageService) EquipmentService {
	return &equipmentService{db: db, equipmentRepo: equipmentRepo, muscleGroupRepo: muscleGroupRepo, imageService: imageService}
}

func (s *equipmentService) GetEquipmentData(q request.EquipmentListRequest, paginator *helper.Paginator) (*response.EquipmentsResponse, error) {
	muscleGroup := strings.Split(q.MuscleGroup, ",")
	equipments, err := s.equipmentRepo.FindEquipmentList(q.Q, muscleGroup, paginator)

	if err != nil {
		logger.Log.WithError(err).Error("error during find all equipments")
		return nil, err
	}

	var resp response.EquipmentsResponse

	for _, equipment := range equipments {
		primaryImage := helper.FindPrimaryImageFromEquipment(equipment)
		var imagePath string
		if primaryImage == nil {
			newName := strings.ReplaceAll(equipment.Name, " ", "+")
			imagePath = fmt.Sprintf("https://placehold.co/600x400?text=%v/png", newName)
		} else {
			imagePath = primaryImage.CloudinaryPath
		}

		price := findEquipmentMinimumPrice(equipment)

		eq := response.Equipment{
			ID:        equipment.ID,
			Name:      equipment.Name,
			Price:     price,
			ImagePath: imagePath,
		}
		resp.Equipments = append(resp.Equipments, eq)

	}
	return &resp, nil

}

func findEquipmentMinimumPrice(equipment model.Equipment) float64 {
	minimumPrice := equipment.EquipmentOptions[0].Price
	for _, option := range equipment.EquipmentOptions {
		if option.Price < minimumPrice {
			minimumPrice = option.Price
		}
	}
	return minimumPrice
}

func (s *equipmentService) AddEquipment(req request.EquipmentPostRequest, context context.Context) error {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	equipmentID := uuid.New()

	newEquipment := model.Equipment{
		ID:          equipmentID,
		Name:        req.Name,
		Brand:       req.Band,
		Description: req.Description,
		Model:       req.Model,
		Color:       req.Color,
		Material:    req.Material,
	}

	err := s.equipmentRepo.CreateEquipment(tx, newEquipment)
	if err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error creating equipment")
		return err
	}

	for _, option := range req.Options {
		optID := uuid.New()

		newOption := model.EquipmentOption{
			ID:                optID,
			EquipmentID:       equipmentID,
			Name:              option.Name,
			Weight:            option.Weight,
			Price:             option.Price,
			RemainingProducts: option.Available,
		}

		if err := s.equipmentRepo.AddEquipmentOption(tx, newOption); err != nil {
			tx.Rollback()
			logger.Log.WithError(err).Error("error adding equipment options", newOption)
			return err
		}

		for _, img := range option.Images {
			imgID := uuid.MustParse(img.ID)
			err = s.imageService.ArchiveImage(tx, context, imgID, equipmentID, optID, img.IsPrimary)
			if err != nil {
				tx.Rollback()
				logger.Log.WithError(err).Error("error archiving image", imgID)
				return err
			}
		}
	}

	var feats []model.EquipmentFeature

	for _, featStr := range req.Features {
		feat := model.EquipmentFeature{
			EquipmentID: equipmentID,
			Description: featStr,
		}
		feats = append(feats, feat)
	}

	if err = s.equipmentRepo.CreateEquipmentFeatures(tx, feats); err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error adding equipment feature")
		return err
	}

	if len(req.AdditionalFields) > 0 {
		var atts []model.Attribute

		for _, field := range req.AdditionalFields {
			newAttribute := model.Attribute{
				EquipmentID: equipmentID,
				Key:         field.Key,
				Value:       field.Value,
			}
			atts = append(atts, newAttribute)
		}

		if err := s.equipmentRepo.AddAttributes(tx, atts); err != nil {
			tx.Rollback()
			logger.Log.WithError(err).Error("error cant add attribute into equipment", equipmentID)
			return err
		}
	}

	for _, groupID := range req.MuscleGroupUsed {
		if _, err := s.muscleGroupRepo.FindByID(tx, groupID); err != nil {
			tx.Rollback()
			logger.Log.WithError(err).Error("error finding muscleGroupID", groupID)
			return err
		}
	}

	if err := s.muscleGroupRepo.AddGroup(tx, req.MuscleGroupUsed, equipmentID); err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error adding muscle Group to equipment")
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *equipmentService) GetEquipmentDetail(eqID uuid.UUID) (*response.EquipmentDetailResponse, error) {

	equipment, err := s.equipmentRepo.FindByID(eqID)
	if err != nil {
		logger.Log.WithError(err).Error("cant find equipment id:", eqID)
		return nil, err
	}

	resp := response.FormatEquipmentDetailResponse(equipment)

	return resp, nil
}

func (s *equipmentService) UpdateEquipment(eqID uuid.UUID, context context.Context, req request.EquipmentPutRequest) error {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	equipment, err := s.equipmentRepo.FindByIDTransaction(tx, eqID)
	if err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error during find equipment by ID")
		return err
	}

	if req.MuscleGroupUsed != nil {
		if err := s.muscleGroupRepo.UpdateGroups(tx, req.MuscleGroupUsed, equipment.ID); err != nil {
			tx.Rollback()
			logger.Log.WithError(err).Error("error updated muscle Group to equipment")
			return err
		}
	}

	if req.Option != nil {
		// created option
		for _, optCreated := range req.Option.Created {
			newOption := model.EquipmentOption{
				EquipmentID:       equipment.ID,
				Name:              optCreated.Name,
				Weight:            optCreated.Weight,
				Price:             optCreated.Price,
				RemainingProducts: optCreated.Available,
			}

			if err := s.equipmentRepo.AddEquipmentOption(tx, newOption); err != nil {
				tx.Rollback()
				logger.Log.WithError(err).Error("error adding equipment options", newOption)
				return err
			}

			for _, img := range optCreated.Images {
				imgID := uuid.MustParse(img.ID)
				err = s.imageService.ArchiveImage(tx, context, imgID, equipment.ID, newOption.ID, img.IsPrimary)
				if err != nil {
					tx.Rollback()
					logger.Log.WithError(err).Error("error archiving image", imgID)
					return err
				}
			}

		}

		//deleted option
		if req.Option.Deleted != nil {
			var opts []uuid.UUID

			for _, opt := range req.Option.Deleted {
				opts = append(opts, uuid.MustParse(opt))
			}

			if err := s.equipmentRepo.DeleteEquipmentOption(tx, opts); err != nil {
				tx.Rollback()
				logger.Log.WithError(err).Error("error deleting equipment options")
				return err
			}
		}

		// updated option
		for _, updateOption := range req.Option.Updated {
			toUpdatedID := uuid.MustParse(updateOption.ID)
			toUpdatedGroup := model.EquipmentOption{
				ID:                toUpdatedID,
				EquipmentID:       equipment.ID,
				Name:              updateOption.Name,
				Weight:            updateOption.Weight,
				Price:             updateOption.Price,
				RemainingProducts: updateOption.Available,
			}

			if err := s.equipmentRepo.SaveEquipmentOption(tx, toUpdatedGroup); err != nil {
				tx.Rollback()
				logger.Log.WithError(err).Error("error saving equipment options")
				return err
			}

			if updateOption.Images != nil {
				for _, deletedID := range updateOption.Images.DeletedID {
					deletedUUID := uuid.MustParse(deletedID.ID)
					if err := s.imageService.DeleteImage(tx, context, deletedUUID); err != nil {
						logger.Log.WithError(err).Error("error deleting image", "imgID", deletedUUID)
						return err
					}
				}

				for _, uploadID := range updateOption.Images.UploadID {
					imgID := uuid.MustParse(uploadID.ID)
					optID := uuid.MustParse(updateOption.ID)
					err = s.imageService.ArchiveImage(tx, context, imgID, equipment.ID, optID, uploadID.IsPrimary)
					if err != nil {
						tx.Rollback()
						logger.Log.WithError(err).Error("error archiving image", imgID)
						return err
					}
				}
			}

		}
	}

	if req.Feature != nil {
		var feats []model.EquipmentFeature

		for _, description := range req.Feature.Created {
			feat := model.EquipmentFeature{
				EquipmentID: equipment.ID,
				Description: description,
			}
			feats = append(feats, feat)
		}

		if err = s.equipmentRepo.CreateEquipmentFeatures(tx, feats); err != nil {
			tx.Rollback()
			logger.Log.WithError(err).Error("error adding equipment feature")
			return err
		}

		if req.Feature.Deleted != nil {
			var feats []uuid.UUID

			for _, feat := range req.Feature.Deleted {
				feats = append(feats, uuid.MustParse(feat))
			}

			if err := s.equipmentRepo.DeleteEquipmentFeature(tx, feats); err != nil {
				tx.Rollback()
				logger.Log.WithError(err).Error("error deleting equipment Features")
				return err
			}
		}

		for _, updateFeature := range req.Feature.Updated {
			toUpdatedID := uuid.MustParse(updateFeature.ID)
			toUpdatedFeat := model.EquipmentFeature{
				ID:          toUpdatedID,
				EquipmentID: equipment.ID,
				Description: updateFeature.Description,
			}

			if err := s.equipmentRepo.SaveEquipmentFeature(tx, toUpdatedFeat); err != nil {
				tx.Rollback()
				logger.Log.WithError(err).Error("error saving equipment feature")
				return err
			}
		}
	}

	if req.AdditionalField != nil {

		var toCreateAtts []model.Attribute
		for _, field := range req.AdditionalField.Created {
			newAttribute := model.Attribute{
				EquipmentID: equipment.ID,
				Key:         field.Key,
				Value:       field.Value,
			}
			toCreateAtts = append(toCreateAtts, newAttribute)
		}

		if err := s.equipmentRepo.AddAttributes(tx, toCreateAtts); err != nil {
			tx.Rollback()
			logger.Log.WithError(err).Error("error cant add attribute into equipment", equipment.ID)
			return err
		}

		for _, field := range req.AdditionalField.Updated {
			toUpdateAttribute := model.Attribute{
				ID:          equipment.ID,
				EquipmentID: equipment.ID,
				Key:         field.Key,
				Value:       field.Value,
			}

			if err := s.equipmentRepo.SaveAttributes(tx, &toUpdateAttribute); err != nil {
				tx.Rollback()
				logger.Log.WithError(err).Error("cannot save attribute into equipment")
				return err
			}
		}

		var toDeletedAttributesID []uuid.UUID
		for _, attrID := range req.AdditionalField.Deleted {
			toDeletedAttributesID = append(toDeletedAttributesID, uuid.MustParse(attrID))
		}

		if err = s.equipmentRepo.DeletesAttributes(tx, toDeletedAttributesID); err != nil {
			tx.Rollback()
			logger.Log.WithError(err).Error("error deleting equipment attribute")
			return err
		}

	}

	if req.Brand != nil {
		equipment.Brand = *req.Brand
	}

	if req.Color != nil {
		equipment.Color = *req.Color
	}

	if req.Material != nil {
		equipment.Material = *req.Material
	}

	if req.Model != nil {
		equipment.Model = *req.Model
	}

	if req.Name != nil {
		equipment.Name = *req.Name
	}

	if req.Description != nil {
		equipment.Description = *req.Description
	}

	if err := s.equipmentRepo.SaveEquipment(tx, equipment); err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error update equipment ID: ", equipment.ID)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
