package application

import (
	"context"
	"errors"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/provideradapter/domain"
	"github.com/context-space/context-space/backend/internal/shared/apierrors"
	contractAdapter "github.com/context-space/context-space/backend/internal/shared/contract/provideradapter"
	contractProvider "github.com/context-space/context-space/backend/internal/shared/contract/providercore"
	contractTranslation "github.com/context-space/context-space/backend/internal/shared/contract/providertranslation"
	"github.com/context-space/context-space/backend/internal/shared/serviceerrors"
	"github.com/context-space/context-space/backend/internal/shared/types"
	"github.com/context-space/context-space/backend/internal/shared/utils"
	"golang.org/x/text/language"
)

type ProviderAdapterService struct {
	providerCoreACL        domain.ProvidercoreAcl
	providerTranslationACL domain.ProviderTranslationAcl
	adapterRepo            domain.ProviderAdapterConfigRepository
	obs                    *observability.ObservabilityProvider
}

func NewProviderAdapterService(
	providerCoreACL domain.ProvidercoreAcl,
	providerTranslationACL domain.ProviderTranslationAcl,
	adapterRepo domain.ProviderAdapterConfigRepository,
	obs *observability.ObservabilityProvider,
) *ProviderAdapterService {
	return &ProviderAdapterService{
		providerCoreACL:        providerCoreACL,
		providerTranslationACL: providerTranslationACL,
		adapterRepo:            adapterRepo,
		obs:                    obs,
	}
}

func (s *ProviderAdapterService) GetProviderAdapterByIdentifier(ctx context.Context, identifier string, preferredLang language.Tag) (*contractAdapter.ProviderAdapterInfoDTO, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "ProviderAdapterService.GetProviderAdapterByIdentifier")
	defer span.End()

	// Get provider basic information from ProviderCore with translation support
	provider, err := s.providerCoreACL.GetProvidercoreDataWithoutTranslation(ctx, identifier)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}
	if provider == nil {
		return nil, apierrors.NewNotFoundError("", err)
	}

	// Get adapter configuration (including permissions)
	adapterConfig, err := s.adapterRepo.GetByIdentifier(ctx, identifier)
	if err != nil {
		return nil, apierrors.NewInternalError("", err)
	}

	providerTranslation, err := s.providerTranslationACL.GetProviderTranslation(ctx, identifier, preferredLang)
	if err != nil && !errors.Is(err, serviceerrors.ErrTranslationNotFound) {
		return nil, apierrors.NewInternalError("", err)
	}

	// Convert to ProviderAdapterDTO
	if providerTranslation == nil {
		return mapProviderAdapterInfoDTOWithoutTranslation(provider, adapterConfig), nil
	}

	return mapProviderAdapterInfoDTO(providerTranslation, provider, adapterConfig), nil
}

func mapProviderAdapterInfoDTO(
	providerTranslation *contractTranslation.ProviderTranslationDTO,
	provider *contractProvider.ProviderDTO,
	adapterConfig *domain.ProviderAdapterConfig,
) *contractAdapter.ProviderAdapterInfoDTO {
	adapterInfoDTO := &contractAdapter.ProviderAdapterInfoDTO{
		Identifier:  provider.Identifier,
		Name:        providerTranslation.Name,
		Description: providerTranslation.Description,
		AuthType:    provider.AuthType,
		Status:      provider.Status,
		IconURL:     provider.IconURL,
		Categories:  providerTranslation.Categories,
		Operations:  make([]contractAdapter.OperationDTO, 0, len(provider.Operations)),
		Permissions: adapterConfig.Permissions,
	}

	// Map provider translation operations and parameters to adapter operations and parameters
	transOpMap := make(map[string]contractTranslation.OperationDTO)
	transParamMap := make(map[string]contractTranslation.ParameterDTO)
	for _, transOp := range providerTranslation.Operations {
		transOpMap[transOp.Identifier] = transOp
		for _, transParam := range transOp.Parameters {
			transParamMap[utils.StringsBuilder(transOp.Identifier, ":", transParam.Name)] = transParam
		}
	}
	for _, op := range provider.Operations {
		opDTO := contractAdapter.OperationDTO{
			Identifier:          op.Identifier,
			Name:                op.Name,
			Description:         op.Description,
			Category:            op.Category,
			RequiredPermissions: op.RequiredPermissions,
			Parameters:          make([]contractAdapter.ParameterDTO, 0, len(op.Parameters)),
		}
		if transOp, ok := transOpMap[op.Identifier]; ok {
			opDTO.Name = transOp.Name
			opDTO.Description = transOp.Description
		}
		for _, param := range op.Parameters {
			paramDTO := contractAdapter.ParameterDTO{
				Name:        param.Name,
				Description: param.Description,
				Type:        string(param.Type),
				Required:    param.Required,
				Enum:        param.Enum,
				Default:     param.Default,
			}
			transParam, ok := transParamMap[utils.StringsBuilder(op.Identifier, ":", param.Name)]
			if ok {
				paramDTO.Description = transParam.Description
			}
			opDTO.Parameters = append(opDTO.Parameters, paramDTO)
		}
		adapterInfoDTO.Operations = append(adapterInfoDTO.Operations, opDTO)
	}

	// Map provider permissions to adapter permissions
	transPermMap := make(map[string]types.Permission)
	for _, perm := range providerTranslation.Permissions {
		transPermMap[perm.Identifier] = perm
	}
	for i, perm := range adapterConfig.Permissions {
		if transPerm, ok := transPermMap[perm.Identifier]; ok {
			adapterInfoDTO.Permissions[i] = transPerm
		}
	}

	return adapterInfoDTO
}

func mapProviderAdapterInfoDTOWithoutTranslation(
	provider *contractProvider.ProviderDTO,
	adapterConfig *domain.ProviderAdapterConfig,
) *contractAdapter.ProviderAdapterInfoDTO {
	adapterInfoDTO := &contractAdapter.ProviderAdapterInfoDTO{
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    provider.AuthType,
		Status:      provider.Status,
		IconURL:     provider.IconURL,
		Categories:  provider.Categories,
		Operations:  make([]contractAdapter.OperationDTO, 0, len(provider.Operations)),
		Permissions: adapterConfig.Permissions,
	}
	for _, op := range provider.Operations {
		opDTO := contractAdapter.OperationDTO{
			Identifier:          op.Identifier,
			Name:                op.Name,
			Description:         op.Description,
			Category:            op.Category,
			RequiredPermissions: op.RequiredPermissions,
			Parameters:          make([]contractAdapter.ParameterDTO, 0, len(op.Parameters)),
		}
		for _, param := range op.Parameters {
			opDTO.Parameters = append(opDTO.Parameters, contractAdapter.ParameterDTO{
				Name:        param.Name,
				Description: param.Description,
				Type:        string(param.Type),
				Required:    param.Required,
				Enum:        param.Enum,
				Default:     param.Default,
			})
		}
		adapterInfoDTO.Operations = append(adapterInfoDTO.Operations, opDTO)
	}
	return adapterInfoDTO
}
