package application

import (
	"github.com/context-space/context-space/backend/internal/providercore/domain"
	contractProvider "github.com/context-space/context-space/backend/internal/shared/contract/providercore"
	contractTranslation "github.com/context-space/context-space/backend/internal/shared/contract/providertranslation"
	"github.com/context-space/context-space/backend/internal/shared/utils"
)

func ProviderToDTONoTranslation(provider *domain.Provider, needOperations bool) *contractProvider.ProviderDTO {
	providerDTO := &contractProvider.ProviderDTO{
		ID:          provider.ID,
		Identifier:  provider.Identifier,
		Name:        provider.Name,
		Description: provider.Description,
		AuthType:    string(provider.AuthType),
		Status:      string(provider.Status),
		Tags:        provider.Tags,
		IconURL:     provider.IconURL,
		Categories:  provider.Categories,
		Embedding:   provider.Embedding,
	}

	if needOperations {
		providerDTO.Operations = make([]contractProvider.OperationDTO, 0, len(provider.Operations))
		for _, operation := range provider.Operations {
			providerDTO.Operations = append(providerDTO.Operations, OperationToDTO(&operation))
		}
	}

	return providerDTO
}

func ProviderToDTO(provider *domain.Provider, translation *contractTranslation.ProviderTranslationDTO, needOperations bool) *contractProvider.ProviderDTO {
	providerDTO := &contractProvider.ProviderDTO{
		ID:          provider.ID,
		Identifier:  provider.Identifier,
		Name:        translation.Name,
		Description: translation.Description,
		AuthType:    string(provider.AuthType),
		Status:      string(provider.Status),
		Tags:        provider.Tags,
		IconURL:     provider.IconURL,
		Categories:  provider.Categories,
		Embedding:   provider.Embedding,
	}

	if needOperations {
		opTranslationMap := make(map[string]contractTranslation.OperationDTO)
		opParamTranslationMap := make(map[string]contractTranslation.ParameterDTO)
		for _, op := range translation.Operations {
			opTranslationMap[op.Identifier] = op
			for _, param := range op.Parameters {
				opParamTranslationMap[utils.StringsBuilder(op.Identifier, ":", param.Name)] = param
			}
		}

		translatedOps := make([]contractProvider.OperationDTO, 0, len(provider.Operations))
		// Apply operation translations if found, otherwise keep original from provider
		for _, op := range provider.Operations {
			opDTO := OperationToDTO(&op)
			if opTranslation, exists := opTranslationMap[op.Identifier]; exists {
				// Use translation
				opDTO.Name = opTranslation.Name
				opDTO.Description = opTranslation.Description
				for i, param := range opDTO.Parameters {
					if paramTranslation, exists := opParamTranslationMap[utils.StringsBuilder(op.Identifier, ":", param.Name)]; exists {
						opDTO.Parameters[i].Name = paramTranslation.Name
						opDTO.Parameters[i].Description = paramTranslation.Description
					}
				}
				translatedOps = append(translatedOps, opDTO)
			} else {
				translatedOps = append(translatedOps, opDTO)
			}
		}
		providerDTO.Operations = translatedOps
	}

	return providerDTO
}

func OperationToDTO(operation *domain.Operation) contractProvider.OperationDTO {
	operationDTO := contractProvider.OperationDTO{
		ID:                  operation.ID,
		Identifier:          operation.Identifier,
		Name:                operation.Name,
		Description:         operation.Description,
		Category:            operation.Category,
		RequiredPermissions: operation.RequiredPermissions,
		Parameters:          make([]contractProvider.ParameterDTO, 0, len(operation.Parameters)),
	}
	for _, param := range operation.Parameters {
		operationDTO.Parameters = append(operationDTO.Parameters, contractProvider.ParameterDTO{
			Name:        param.Name,
			Description: param.Description,
			Type:        string(param.Type),
			Required:    param.Required,
			Enum:        param.Enum,
			Default:     param.Default,
		})
	}
	return operationDTO
}
