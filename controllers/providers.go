package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strconv"
	"strings"
)

func (r *DBaaSInventoryReconciler) getDBaaSProviders(ctx *context.Context, namespacedName types.NamespacedName) (v1alpha1.DBaaSProviderList, error) {
	logger := log.FromContext(*ctx, "dbaassinventory", namespacedName)

	var cmList v1.ConfigMapList
	opts := []client.ListOption{
		client.InNamespace("dbaas-operator"),
		client.MatchingLabels{"related-to": "dbaas-operator"},
		client.MatchingLabels{"type": "dbaas-provider-registration"},
	}

	if err := r.List(*ctx, &cmList, opts...); err != nil {
		logger.Error(err, "Error reading ConfigMaps for the configured DBaaS Providers")
		return v1alpha1.DBaaSProviderList{}, err
	}

	providers := make([]v1alpha1.DBaaSProvider, len(cmList.Items))
	for i, cm := range cmList.Items {
		var provider v1alpha1.DBaaSProvider
		if providerName, exists := cm.Data["provider"]; exists {
			provider.Provider = v1alpha1.DatabaseProvider{Name: providerName}
		} else {
			return v1alpha1.DBaaSProviderList{}, errors.New("provider name is missing of the configured DBaaS provider")
		}
		if inventoryKind, exists := cm.Data["inventory_kind"]; exists {
			provider.InventoryKind = inventoryKind
		} else {
			return v1alpha1.DBaaSProviderList{}, fmt.Errorf("inventory kind is missing of the configured DBaaS provider %s", provider.Provider.Name)
		}
		if connectionKind, exists := cm.Data["connection_kind"]; exists {
			provider.ConnectionKind = connectionKind
		} else {
			return v1alpha1.DBaaSProviderList{}, fmt.Errorf("connection kind is missing of the configured DBaaS provider %s", provider.Provider.Name)
		}
		if credentialsFields, exists := cm.Data["credentials_fields"]; exists {
			if authenticationFields, err := parseCredentialFields(credentialsFields, provider.Provider.Name); err != nil {
				return v1alpha1.DBaaSProviderList{}, err
			} else {
				provider.AuthenticationFields = authenticationFields
			}
		} else {
			return v1alpha1.DBaaSProviderList{}, fmt.Errorf("credential fields are missing of the configured DBaaS provider %s", provider.Provider.Name)
		}
		providers[i] = provider
	}

	return v1alpha1.DBaaSProviderList{Items: providers}, nil
}

func parseCredentialFields(credentialsFields string, providerName string) ([]v1alpha1.AuthenticationField, error) {
	fieldProperties := strings.Split(credentialsFields, "\n")
	fieldMap := map[string]v1alpha1.AuthenticationField{}
	for _, fieldProperty := range fieldProperties {
		if field, property, value, err := parseCredentialField(fieldProperty, providerName); err != nil {
			return nil, err
		} else {
			authenticationField, exists := fieldMap[field]
			if !exists {
				fieldMap[field] = authenticationField
			}
			if err := setAuthenticationFieldProperty(&authenticationField, property, value, providerName); err != nil {
				return nil, err
			}
		}
	}
	fields := make([]v1alpha1.AuthenticationField, 0, len(fieldMap))
	for _, field := range fieldMap {
		fields = append(fields, field)
	}
	return fields, nil
}

func parseCredentialField(fieldString string, providerName string) (field string, property string, value string, err error) {
	fieldToken := strings.SplitN(fieldString, ".", 2)
	if len(fieldToken) != 2 {
		err = fmt.Errorf("invalid credential field of DBaaS provider %s: %s", providerName, fieldString)
		return
	}
	field = strings.TrimSpace(fieldToken[0])

	propertyValue := strings.SplitN(fieldToken[1], ":", 2)
	if len(propertyValue) != 2 {
		err = fmt.Errorf("invalid property of credential field %s of DBaaS provider %s: %s", field, providerName, fieldString)
		return
	}
	property = strings.TrimSpace(propertyValue[0])
	value = strings.TrimSpace(propertyValue[1])
	return
}

func setAuthenticationFieldProperty(authenticationField *v1alpha1.AuthenticationField, property string, value string, providerName string) (err error) {
	if property == "name" {
		authenticationField.Name = value
	} else if property == "masked" {
		if authenticationField.Masked, err = strconv.ParseBool(value); err != nil {
			err = fmt.Errorf("invalid credential field property %s of DBaaS provider %s: %v", property, providerName, err)
		}
	} else {
		err = fmt.Errorf("invalid credential field property %s of DBaaS provider %s", property, providerName)
	}
	return
}
