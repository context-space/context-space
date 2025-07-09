package hubspot

import (
	"context"
	"fmt"
	"net/http"
)

// Define constants for API paths used by handlers.
const (
	endpointSearchContacts  = "/crm/v3/objects/contacts/search"
	endpointSearchCompanies = "/crm/v3/objects/companies/search"
	endpointCreateContact   = "/crm/v3/objects/contacts"
	endpointUpdateContact   = "/crm/v3/objects/contacts/{contactId}"
	endpointDeleteContact   = "/crm/v3/objects/contacts/{contactId}"
	endpointCreateDeal      = "/crm/v3/objects/deals"
	endpointDeleteDeal      = "/crm/v3/objects/deals/{dealId}"
	endpointLogCall         = "/crm/v3/objects/calls"
	endpointSendEmail       = "/marketing/v3/transactional/single-email/send"
	endpointCreateTask      = "/crm/v3/objects/tasks"
	endpointCreateTicket    = "/crm/v3/objects/tickets"
)

// Define constants for operation IDs used by handlers.
const (
	operationIDSearchContacts  = "search_contacts"
	operationIDSearchCompanies = "search_companies"
	operationIDCreateContact   = "create_contact"
	operationIDUpdateContact   = "update_contact"
	operationIDDeleteContact   = "delete_contact"
	operationIDCreateDeal      = "create_deal"
	operationIDDeleteDeal      = "delete_deal"
	operationIDLogCall         = "log_call"
	operationIDCreateTask      = "create_task"
	operationIDCreateTicket    = "create_ticket"
)

// Define a struct for each operation's parameters based on config.
//

// SearchContactsParams defines parameters for the Search Contacts operation.
type SearchContactsParams struct {
	// Define fields based on search_contacts parameters in config

	Filtergroups []interface{} `mapstructure:"filterGroups" validate:"omitempty"` // Filter groups to apply to the search. Each group's filters are ANDed, groups are ORed.

	Query string `mapstructure:"query" validate:"omitempty"` // Keyword string for full-text search across contact properties.

	Properties []string `mapstructure:"properties" validate:"omitempty"` // List of contact property names to return.

	Limit int `mapstructure:"limit" validate:"omitempty"` // Maximum number of results to return.

	After string `mapstructure:"after" validate:"omitempty"` // Pagination cursor to get the next page of results.

}

// SearchCompaniesParams defines parameters for the Search Companies operation.
type SearchCompaniesParams struct {
	// Define fields based on search_companies parameters in config

	Filtergroups []interface{} `mapstructure:"filterGroups" validate:"omitempty"` // Filter groups for company properties.

	Query string `mapstructure:"query" validate:"omitempty"` // Keyword string for full-text search.

	Properties []string `mapstructure:"properties" validate:"omitempty"` // List of company property names to return.

	Limit int `mapstructure:"limit" validate:"omitempty"` // Maximum number of results.

	After string `mapstructure:"after" validate:"omitempty"` // Pagination cursor.

}

// CreateContactParams defines parameters for the Create Contact operation.
type CreateContactParams struct {
	// Define fields based on create_contact parameters in config

	Properties map[string]interface{} `mapstructure:"properties" validate:"required"` // Object containing contact properties (e.g., email, firstname, lastname).

}

// UpdateContactParams defines parameters for the Update Contact operation.
type UpdateContactParams struct {
	// Define fields based on update_contact parameters in config

	Contactid string `mapstructure:"contactId" validate:"required"` // The ID of the contact to update.

	Properties map[string]interface{} `mapstructure:"properties" validate:"required"` // Object containing contact properties to update.

}

// DeleteContactParams defines parameters for the Delete Contact operation.
type DeleteContactParams struct {
	// Define fields based on delete_contact parameters in config

	Contactid string `mapstructure:"contactId" validate:"required"` // The ID of the contact to delete.

}

// CreateDealParams defines parameters for the Create Deal operation.
type CreateDealParams struct {
	// Define fields based on create_deal parameters in config

	Properties map[string]interface{} `mapstructure:"properties" validate:"required"` // Object containing deal properties (e.g., dealname, amount, pipeline, dealstage).

}

// DeleteDealParams defines parameters for the Delete Deal operation.
type DeleteDealParams struct {
	// Define fields based on delete_deal parameters in config

	Dealid string `mapstructure:"dealId" validate:"required"` // The ID of the deal to delete.

}

// LogCallParams defines parameters for the Log Call operation.
type LogCallParams struct {
	// Define fields based on log_call parameters in config

	Properties map[string]interface{} `mapstructure:"properties" validate:"required"` // Object containing call properties (e.g., hs_call_body, hs_call_duration, hs_callOutcome).

	Associations []interface{} `mapstructure:"associations" validate:"omitempty"` // List of associations to other CRM records (e.g., contact, company, deal).

}

// CreateTaskParams defines parameters for the Create Task operation.
type CreateTaskParams struct {
	// Define fields based on create_task parameters in config

	Properties map[string]interface{} `mapstructure:"properties" validate:"required"` // Object containing task properties (e.g., hs_task_subject, hs_task_body, hs_task_due_date).

	Associations []interface{} `mapstructure:"associations" validate:"omitempty"` // List of associations to other CRM records.

}

// CreateTicketParams defines parameters for the Create Ticket operation.
type CreateTicketParams struct {
	// Define fields based on create_ticket parameters in config

	Properties map[string]interface{} `mapstructure:"properties" validate:"required"` // Object containing ticket properties (e.g., subject, content, hs_pipeline, hs_pipeline_stage).

	Associations []interface{} `mapstructure:"associations" validate:"omitempty"` // List of associations to other CRM records (e.g., contact, company).

}

// OperationHandler defines the function signature for handling a specific API operation.
type OperationHandler func(ctx context.Context, params interface{}) (map[string]interface{}, error)

// OperationDefinition combines parameter schema and handler.
type OperationDefinition struct {
	Schema                interface{}      // Parameter schema (struct pointer)
	Handler               OperationHandler // Operation handler function
	PermissionIdentifiers []string         // List of internal permission identifiers required
}

// Operations maps operation IDs to their definitions.
type Operations map[string]OperationDefinition

// RegisterOperation registers the parameter schema and handler.
func (a *HubspotAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler, requiredPerms []string) {
	a.BaseAdapter.RegisterOperation(operationID, schema) // Register schema for validation
	if a.operations == nil {
		a.operations = make(Operations)
	}
	a.operations[operationID] = OperationDefinition{
		Schema:                schema,
		Handler:               handler,
		PermissionIdentifiers: requiredPerms,
	}
}

// registerOperations is called by the adapter constructor to register all supported operations.
func (a *HubspotAdapter) registerOperations() {

	a.RegisterOperation(
		operationIDSearchContacts,
		&SearchContactsParams{},
		handleSearchContacts,
		[]string{"crm.objects.contacts.read"},
	)

	a.RegisterOperation(
		operationIDSearchCompanies,
		&SearchCompaniesParams{},
		handleSearchCompanies,
		[]string{"crm.objects.companies.read"},
	)

	a.RegisterOperation(
		operationIDCreateContact,
		&CreateContactParams{},
		handleCreateContact,
		[]string{"crm.objects.contacts.write"},
	)

	a.RegisterOperation(
		operationIDUpdateContact,
		&UpdateContactParams{},
		handleUpdateContact,
		[]string{"crm.objects.contacts.write"},
	)

	a.RegisterOperation(
		operationIDDeleteContact,
		&DeleteContactParams{},
		handleDeleteContact,
		[]string{"crm.objects.contacts.write"},
	)

	a.RegisterOperation(
		operationIDCreateDeal,
		&CreateDealParams{},
		handleCreateDeal,
		[]string{"crm.objects.deals.write"},
	)

	a.RegisterOperation(
		operationIDDeleteDeal,
		&DeleteDealParams{},
		handleDeleteDeal,
		[]string{"crm.objects.deals.write"},
	)

	a.RegisterOperation(
		operationIDLogCall,
		&LogCallParams{},
		handleLogCall,
		[]string{"crm.objects.calls.write"},
	)

	a.RegisterOperation(
		operationIDCreateTask,
		&CreateTaskParams{},
		handleCreateTask,
		[]string{"crm.objects.tasks.write"},
	)

	a.RegisterOperation(
		operationIDCreateTicket,
		&CreateTicketParams{},
		handleCreateTicket,
		[]string{"crm.objects.tickets.write"},
	)
}

// handleSearchContacts constructs parameters for the REST adapter for the Search Contacts operation.
func handleSearchContacts(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*SearchContactsParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for search_contacts")
	}

	requestBody := make(map[string]interface{})
	if p.Filtergroups != nil {
		requestBody["filterGroups"] = p.Filtergroups
	}
	if p.Query != "" {
		requestBody["query"] = p.Query
	}
	if len(p.Properties) > 0 {
		requestBody["properties"] = p.Properties
	}
	if p.Limit > 0 {
		requestBody["limit"] = p.Limit
	}
	if p.After != "" {
		requestBody["after"] = p.After
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointSearchContacts,
	}
	if len(requestBody) > 0 {
		restParams["body"] = requestBody
	}

	return restParams, nil
}

// handleSearchCompanies constructs parameters for the REST adapter for the Search Companies operation.
func handleSearchCompanies(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*SearchCompaniesParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for search_companies")
	}

	requestBody := make(map[string]interface{})
	if p.Filtergroups != nil {
		requestBody["filterGroups"] = p.Filtergroups
	}
	if p.Query != "" {
		requestBody["query"] = p.Query
	}
	if len(p.Properties) > 0 {
		requestBody["properties"] = p.Properties
	}
	if p.Limit > 0 {
		requestBody["limit"] = p.Limit
	}
	if p.After != "" {
		requestBody["after"] = p.After
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointSearchCompanies,
	}
	if len(requestBody) > 0 {
		restParams["body"] = requestBody
	}

	return restParams, nil
}

// handleCreateContact constructs parameters for the REST adapter for the Create Contact operation.
func handleCreateContact(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*CreateContactParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for create_contact")
	}

	requestBody := map[string]interface{}{
		"properties": p.Properties,
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointCreateContact,
		"body":   requestBody,
	}

	return restParams, nil
}

// handleUpdateContact constructs parameters for the REST adapter for the Update Contact operation.
func handleUpdateContact(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*UpdateContactParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for update_contact")
	}

	pathParams := map[string]string{
		"contactId": p.Contactid,
	}
	requestBody := map[string]interface{}{
		"properties": p.Properties,
	}

	restParams := map[string]interface{}{
		"method":      http.MethodPatch,
		"path":        endpointUpdateContact,
		"path_params": pathParams,
		"body":        requestBody,
	}

	return restParams, nil
}

// handleDeleteContact constructs parameters for the REST adapter for the Delete Contact operation.
func handleDeleteContact(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*DeleteContactParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for delete_contact")
	}

	pathParams := map[string]string{
		"contactId": p.Contactid,
	}

	restParams := map[string]interface{}{
		"method":      http.MethodDelete,
		"path":        endpointDeleteContact,
		"path_params": pathParams,
	}

	return restParams, nil
}

// handleCreateDeal constructs parameters for the REST adapter for the Create Deal operation.
func handleCreateDeal(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*CreateDealParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for create_deal")
	}

	requestBody := map[string]interface{}{
		"properties": p.Properties,
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointCreateDeal,
		"body":   requestBody,
	}

	return restParams, nil
}

// handleDeleteDeal constructs parameters for the REST adapter for the Delete Deal operation.
func handleDeleteDeal(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*DeleteDealParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for delete_deal")
	}

	pathParams := map[string]string{
		"dealId": p.Dealid,
	}

	restParams := map[string]interface{}{
		"method":      http.MethodDelete,
		"path":        endpointDeleteDeal,
		"path_params": pathParams,
	}

	return restParams, nil
}

// handleLogCall constructs parameters for the REST adapter for the Log Call operation.
func handleLogCall(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*LogCallParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for log_call")
	}

	requestBody := map[string]interface{}{
		"properties": p.Properties,
	}
	if p.Associations != nil {
		requestBody["associations"] = p.Associations
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointLogCall,
		"body":   requestBody,
	}

	return restParams, nil
}

// handleCreateTask constructs parameters for the REST adapter for the Create Task operation.
func handleCreateTask(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*CreateTaskParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for create_task")
	}

	requestBody := map[string]interface{}{
		"properties": p.Properties,
	}
	if p.Associations != nil {
		requestBody["associations"] = p.Associations
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointCreateTask,
		"body":   requestBody,
	}

	return restParams, nil
}

// handleCreateTicket constructs parameters for the REST adapter for the Create Ticket operation.
func handleCreateTicket(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	p, ok := params.(*CreateTicketParams)
	if !ok || p == nil {
		return nil, fmt.Errorf("invalid or missing parameters for create_ticket")
	}

	requestBody := map[string]interface{}{
		"properties": p.Properties,
	}
	if p.Associations != nil {
		requestBody["associations"] = p.Associations
	}

	restParams := map[string]interface{}{
		"method": http.MethodPost,
		"path":   endpointCreateTicket,
		"body":   requestBody,
	}

	return restParams, nil
}
