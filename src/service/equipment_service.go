package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/E-cercise/E-cercise/src/config"
	"github.com/E-cercise/E-cercise/src/data/dto"
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
	GetRecommendEquipmentData(user *model.User) (*response.EquipmentsResponse, error)
	AddEquipment(req request.EquipmentPostRequest, context context.Context) error
	GetEquipmentDetail(eqID uuid.UUID) (*response.EquipmentDetailResponse, error)
	UpdateEquipment(eqID uuid.UUID, context context.Context, req request.EquipmentPutRequest) error
	DeleteEquipment(eqID uuid.UUID, context context.Context) error
	GetAllEquipmentCategories() (*response.CategoriesResponse, error)
	GetAllEquipmentsDetail(eqIDs []uuid.UUID) (*response.EquipmentDetailComparisonResponse, error)
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
	var muscleGroup []string
	if q.MuscleGroup != "" {
		muscleGroup = strings.Split(q.MuscleGroup, ",")
	}
	equipments, err := s.equipmentRepo.FindEquipmentList(q.Q, muscleGroup, paginator, q.Category, int(q.MinBudget), int(q.MaxBudget))

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
			ID:              equipment.ID,
			Name:            equipment.Name,
			Price:           price,
			ImagePath:       imagePath,
			MuscleGroupUsed: helper.GetMuscleGroupIDFromEquipment(equipment),
		}
		resp.Equipments = append(resp.Equipments, eq)
	}
	return &resp, nil
}

func (s *equipmentService) GetRecommendEquipmentData(user *model.User) (*response.EquipmentsResponse, error) {
	var resp response.EquipmentsResponse

	payload := dto.RecommenderRequest{
		UserType:   strings.ToLower(string(user.Experience)),
		Gender:     strings.ToLower(string(user.Gender)),
		Age:        user.Age,
		Weight:     user.Weight,
		Height:     user.Height,
		Goal:       strings.ToLower(user.Goal.Name),
		Experience: strings.ToLower(string(user.Experience)),
	}

	for _, pref := range user.UserPreferences {
		group := helper.GetTagGroup(pref.Tag.Name)
		payload.Preferences = append(payload.Preferences, dto.RecommenderPreference{
			Tag:   pref.Tag.Name,
			Group: group,
		})
	}

	recommenderURL := fmt.Sprintf("%s/recommend", config.RecommendationServiceBaseUrl)
	res, err := helper.PostJSON(recommenderURL, payload)
	if err != nil {
		logger.Log.WithError(err).Error("failed to call recommender service")
		return nil, err
	}

	var recommended dto.RecommendedResponseDTO

	if err := json.Unmarshal(res, &recommended); err != nil {
		logger.Log.WithError(err).Error("failed to parse recommender response")
		return nil, err
	}

	for _, item := range recommended {
		option, err := s.equipmentRepo.FindOptionByID(item.OptionID)
		if err != nil {
			logger.Log.WithError(err).Error("failed to fetch option ID", "optID", item.OptionID)
			return nil, err
		}

		primaryImg := helper.FindPrimaryImage(*option)
		var imagePath string
		if primaryImg != nil {
			imagePath = primaryImg.CloudinaryPath
		}

		newRecEq := response.Equipment{
			ID:               option.EquipmentID,
			Name:             helper.AbbreviateEquipmentName(item.EquipmentName, option.Name),
			Price:            option.Price,
			ImagePath:        imagePath,
			RemainingProduct: &item.RemainingProduct,
		}
		resp.Equipments = append(resp.Equipments, newRecEq)
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
			logger.Log.Error("error ", r)
			tx.Rollback()
		}
	}()

	equipmentID := uuid.New()

	newEquipment := model.Equipment{
		ID:          equipmentID,
		Name:        req.Name,
		Brand:       req.Brand,
		Category:    req.Category,
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
			imgID, err := uuid.Parse(img.ID)
			if err != nil {
				tx.Rollback()
				logger.Log.WithError(err).Error("error parsing image id")
				return err
			}

			err = s.imageService.ArchiveImage(tx, context, imgID, equipmentID, optID, img.IsPrimary)
			if err != nil {
				tx.Rollback()
				logger.Log.WithError(err).Error("error archiving image", imgID)
				return err
			}
		}
	}

	if req.Features != nil {
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

	if req.MuscleGroupUsed != nil {
		if err := s.muscleGroupRepo.UpdateGroups(s.db, req.MuscleGroupUsed, eqID); err != nil {
			logger.Log.WithError(err).Error("Failed to update muscle groups")
			return err
		}
		logger.Log.Infof("✅ Muscle groups updated for equipment ID: %v", eqID)
	}

	if req.Option != nil && req.Option.Deleted != nil {
		var opts []uuid.UUID
		for _, opt := range req.Option.Deleted {
			optID := uuid.MustParse(opt)
			if err := s.imageService.DeleteImagesByOptionID(s.db, context, optID); err != nil {
				logger.Log.WithError(err).Error("Error deleting images for option", "optionID", optID)
				return err
			}
			opts = append(opts, optID)
		}
		if err := s.equipmentRepo.DeleteEquipmentOption(s.db, opts); err != nil {
			logger.Log.WithError(err).Error("Error deleting equipment options")
			return err
		}
	}

	equipment, err := s.equipmentRepo.FindByID(eqID)
	if err != nil {
		logger.Log.WithError(err).Error("error during find equipment by ID")
		return err
	}

	if req.Option != nil && req.Option.Updated != nil {
		for _, updateOption := range req.Option.Updated {
			optID := uuid.MustParse(updateOption.ID)
			updatedOpt := model.EquipmentOption{
				ID:                optID,
				EquipmentID:       equipment.ID,
				Name:              updateOption.Name,
				Weight:            updateOption.Weight,
				Price:             updateOption.Price,
				RemainingProducts: updateOption.Available,
			}
			if err := s.equipmentRepo.SaveEquipmentOption(s.db, updatedOpt); err != nil {
				logger.Log.WithError(err).Error("error saving equipment options")
				return err
			}
			if updateOption.Images != nil {
				for _, deletedID := range updateOption.Images.DeletedID {
					delID := uuid.MustParse(deletedID.ID)
					if err := s.imageService.DeleteImage(s.db, context, delID); err != nil {
						logger.Log.WithError(err).Error("error deleting image", "imgID", delID)
						return err
					}
				}
				for _, uploadID := range updateOption.Images.UploadID {
					imgID := uuid.MustParse(uploadID.ID)
					err = s.imageService.ArchiveImage(s.db, context, imgID, equipment.ID, optID, uploadID.IsPrimary)
					if err != nil {
						logger.Log.WithError(err).Error("error archiving image", imgID)
						return err
					}
				}
			}
		}
	}

	if req.Feature != nil && req.Feature.Deleted != nil {
		var ids []uuid.UUID
		for _, feat := range req.Feature.Deleted {
			id := uuid.MustParse(feat)
			ids = append(ids, id)
		}
		if err := s.equipmentRepo.DeleteEquipmentFeature(s.db, ids); err != nil {
			logger.Log.WithError(err).Error("Error deleting equipment features")
			return err
		}
		logger.Log.Infof("Deleted %d equipment features", len(ids))
	}

	// 🧹 4. Delete Additional Fields (Attributes)
	if req.AdditionalField != nil && req.AdditionalField.Deleted != nil {
		var ids []uuid.UUID
		for _, attrID := range req.AdditionalField.Deleted {
			id := uuid.MustParse(attrID)
			ids = append(ids, id)
		}
		if err := s.equipmentRepo.DeletesAttributes(s.db, ids); err != nil {
			logger.Log.WithError(err).Error("Error deleting equipment attributes")
			return err
		}
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.Log.WithField("panic", r).Error("💥 panic occurred during UpdateEquipment")
		}
	}()

	if req.Option != nil {
		if req.Option.Created != nil {
			for _, optCreated := range req.Option.Created {
				optID := uuid.New()
				newOption := model.EquipmentOption{
					ID:                optID,
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
					err = s.imageService.ArchiveImage(tx, context, imgID, equipment.ID, optID, img.IsPrimary)
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
		if req.Feature.Created != nil {
			var feats []model.EquipmentFeature
			for _, description := range req.Feature.Created {
				feats = append(feats, model.EquipmentFeature{
					EquipmentID: equipment.ID,
					Description: description,
				})
			}
			if err = s.equipmentRepo.CreateEquipmentFeatures(tx, feats); err != nil {
				tx.Rollback()
				logger.Log.WithError(err).Error("error adding equipment feature")
				return err
			}
		}

		if req.Feature.Updated != nil {
			for _, feat := range req.Feature.Updated {
				featID := uuid.MustParse(feat.ID)
				updated := model.EquipmentFeature{
					ID:          featID,
					EquipmentID: equipment.ID,
					Description: feat.Description,
				}
				if err := s.equipmentRepo.SaveEquipmentFeature(tx, updated); err != nil {
					tx.Rollback()
					logger.Log.WithError(err).Error("error saving equipment feature")
					return err
				}
			}
		}
	}

	if req.AdditionalField != nil {
		if req.AdditionalField.Created != nil {
			var toCreate []model.Attribute
			for _, field := range req.AdditionalField.Created {
				toCreate = append(toCreate, model.Attribute{
					EquipmentID: equipment.ID,
					Key:         field.Key,
					Value:       field.Value,
				})
			}
			if err := s.equipmentRepo.AddAttributes(tx, toCreate); err != nil {
				tx.Rollback()
				logger.Log.WithError(err).Error("error adding equipment attributes")
				return err
			}
		}

		if req.AdditionalField.Updated != nil {
			for _, field := range req.AdditionalField.Updated {
				toUpdate := model.Attribute{
					ID:          uuid.MustParse(field.ID),
					EquipmentID: equipment.ID,
					Key:         field.Key,
					Value:       field.Value,
				}
				if err := s.equipmentRepo.SaveAttributes(tx, &toUpdate); err != nil {
					tx.Rollback()
					logger.Log.WithError(err).Error("cannot save attribute into equipment")
					return err
				}
			}
		}

	}

	if req.Brand != nil {
		equipment.Brand = *req.Brand
	}
	if req.Category != nil {
		equipment.Category = *req.Category
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
		logger.Log.WithError(err).Error("error updating equipment ID: ", equipment.ID)
		return err
	}

	logger.Log.Info("🧾 About to commit transaction...")
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("transaction commit failed")
		return err
	}
	logger.Log.Info("✅ Transaction committed.")

	return nil
}

func (s *equipmentService) DeleteEquipment(eqID uuid.UUID, context context.Context) error {
	tx := s.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	equipment, err := s.equipmentRepo.FindByIDTransaction(tx, eqID)

	if err != nil {
		logger.Log.WithError(err).Error("error finding equipment by id, ", "equipmentID", eqID)
		tx.Rollback()
		return err
	}

	for _, opt := range equipment.EquipmentOptions {
		for _, img := range opt.Images {
			if err := s.imageService.DeleteImage(tx, context, img.ID); err != nil {
				logger.Log.WithError(err).Error("cant delete image", "imgID", img)
				tx.Rollback()
				return err
			}
		}
	}

	if err := s.equipmentRepo.DeleteEquipment(tx, eqID); err != nil {
		tx.Rollback()
		logger.Log.WithError(err).Error("error deleting equipment ID: ", equipment.ID)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *equipmentService) GetAllEquipmentCategories() (*response.CategoriesResponse, error) {
	equipments, err := s.equipmentRepo.GetAllEquipmentCategories()
	if err != nil {
		logger.Log.WithError(err).Error("can't get all equipment categories")
		return nil, err
	}

	var resp response.CategoriesResponse

	for i, equipment := range equipments {
		category := response.Category{
			Value: int64(i + 1),
			Label: equipment.Category,
		}
		resp.Categories = append(resp.Categories, category)
	}

	return &resp, nil
}

func (s *equipmentService) GetAllEquipmentsDetail(eqIDs []uuid.UUID) (*response.EquipmentDetailComparisonResponse, error) {
	equipments, err := s.equipmentRepo.FindByIDs(eqIDs)

	if err != nil {
		return nil, err
	}

	commonKeys := helper.FindCommonAttributes(equipments)

	var filteredEquipments []response.EquipmentDetail
	for _, eq := range equipments {
		filteredData := response.EquipmentDetail{
			ID:          eq.ID,
			Name:        eq.Name,
			Brand:       eq.Brand,
			Color:       eq.Color,
			Category:    eq.Category,
			Description: eq.Description,
			Material:    eq.Material,
			Model:       eq.Model,
		}
		for _, opt := range eq.EquipmentOptions {
			option := response.Option{
				ID:        opt.ID.String(),
				Name:      opt.Name,
				Available: opt.RemainingProducts,
				Price:     opt.Price,
				Weight:    opt.Weight,
			}
			for _, img := range opt.Images {
				image := response.Image{
					ID:        img.ID.String(),
					Url:       img.CloudinaryPath,
					IsPrimary: img.IsPrimary,
				}
				option.Images = append(option.Images, image)
			}
			filteredData.Option = append(filteredData.Option, option)
		}

		var additionalAttributes []response.AdditionalField
		for _, key := range commonKeys {
			for _, attr := range eq.Attribute {
				if attr.Key == key {
					attribute := response.AdditionalField{
						ID:    attr.ID.String(),
						Key:   attr.Key,
						Value: attr.Value,
					}
					additionalAttributes = append(additionalAttributes, attribute)
					// filteredData[key] = attr.Value
				}
			}
		}
		filteredData.AdditionalField = additionalAttributes
		filteredEquipments = append(filteredEquipments, filteredData)
	}

	equipmentMap := make(map[uuid.UUID]response.EquipmentDetail)
	for _, eq := range filteredEquipments {
		equipmentMap[eq.ID] = eq
	}

	var sortedEquipments []response.EquipmentDetail
	for _, id := range eqIDs {
		if eq, exists := equipmentMap[id]; exists {
			sortedEquipments = append(sortedEquipments, eq)
		}
	}

	resp := response.EquipmentDetailComparisonResponse{
		Equipments: sortedEquipments,
	}

	return &resp, nil
}
