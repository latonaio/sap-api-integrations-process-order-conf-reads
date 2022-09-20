package sap_api_caller

import (
	"fmt"
	"io/ioutil"
	sap_api_output_formatter "sap-api-integrations-process-order-conf-reads/SAP_API_Output_Formatter"
	"strings"
	"sync"

	sap_api_request_client_header_setup "github.com/latonaio/sap-api-request-client-header-setup"

	"github.com/latonaio/golang-logging-library-for-sap/logger"
)

type SAPAPICaller struct {
	baseURL         string
	sapClientNumber string
	requestClient   *sap_api_request_client_header_setup.SAPRequestClient
	log             *logger.Logger
}

func NewSAPAPICaller(baseUrl, sapClientNumber string, requestClient *sap_api_request_client_header_setup.SAPRequestClient, l *logger.Logger) *SAPAPICaller {
	return &SAPAPICaller{
		baseURL:         baseUrl,
		requestClient:   requestClient,
		sapClientNumber: sapClientNumber,
		log:             l,
	}
}

func (c *SAPAPICaller) AsyncGetProcessOrderConfirmation(orderID, batch, confirmationGroup string, accepter []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(accepter))
	for _, fn := range accepter {
		switch fn {
		case "ConfByOrderID":
			func() {
				c.ConfByOrderID(orderID)
				wg.Done()
			}()
		case "MaterialMovements":
			func() {
				c.MaterialMovements(batch)
				wg.Done()
			}()
		case "BatchCharacteristic":
			func() {
				c.BatchCharacteristic(batch)
				wg.Done()
			}()
		case "ConfByOrderIDConfGroup":
			func() {
				c.ConfByOrderIDConfGroup(orderID, confirmationGroup)
				wg.Done()
			}()
		default:
			wg.Done()
		}
	}

	wg.Wait()
}

func (c *SAPAPICaller) ConfByOrderID(orderID string) {
	confByOrderIDData, err := c.callProcessOrderConfSrvAPIRequirementConfByOrderID("ProcOrdConf2", orderID)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(confByOrderIDData)
}

func (c *SAPAPICaller) callProcessOrderConfSrvAPIRequirementConfByOrderID(api, orderID string) ([]sap_api_output_formatter.Confirmation, error) {
	url := strings.Join([]string{c.baseURL, "API_PROC_ORDER_CONFIRMATION_2_SRV", api}, "/")
	param := c.getQueryWithConfByOrderID(map[string]string{}, orderID)

	resp, err := c.requestClient.Request("GET", url, param, "")
	if err != nil {
		return nil, fmt.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToConfirmation(byteArray, c.log)
	if err != nil {
		return nil, fmt.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) MaterialMovements(batch string) {
	data, err := c.callProcessOrderConfSrvAPIRequirementMaterialMovements("ProcOrdConfMatlDocItm", batch)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callProcessOrderConfSrvAPIRequirementMaterialMovements(api, batch string) ([]sap_api_output_formatter.MaterialMovements, error) {
	url := strings.Join([]string{c.baseURL, "API_PROC_ORDER_CONFIRMATION_2_SRV", api}, "/")

	param := c.getQueryWithMaterialMovements(map[string]string{}, batch)

	resp, err := c.requestClient.Request("GET", url, param, "")
	if err != nil {
		return nil, fmt.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToMaterialMovements(byteArray, c.log)
	if err != nil {
		return nil, fmt.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) BatchCharacteristic(batch string) {
	data, err := c.callProcessOrderConfSrvAPIRequirementBatchCharacteristic("ProcOrderConfBatchCharc", batch)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(data)
}

func (c *SAPAPICaller) callProcessOrderConfSrvAPIRequirementBatchCharacteristic(api, batch string) ([]sap_api_output_formatter.BatchCharacteristic, error) {
	url := strings.Join([]string{c.baseURL, "API_PROC_ORDER_CONFIRMATION_2_SRV", api}, "/")

	param := c.getQueryWithBatchCharacteristic(map[string]string{}, batch)

	resp, err := c.requestClient.Request("GET", url, param, "")
	if err != nil {
		return nil, fmt.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToBatchCharacteristic(byteArray, c.log)
	if err != nil {
		return nil, fmt.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) ConfByOrderIDConfGroup(orderID, confirmationGroup string) {
	confbyOrderIDConfGroupData, err := c.callProcessOrderConfSrvAPIRequirementConfByOrderIDConfGroup("ProcOrdConf2", orderID, confirmationGroup)
	if err != nil {
		c.log.Error(err)
		return
	}
	c.log.Info(confbyOrderIDConfGroupData)
}

func (c *SAPAPICaller) callProcessOrderConfSrvAPIRequirementConfByOrderIDConfGroup(api, orderID, confirmationGroup string) ([]sap_api_output_formatter.Confirmation, error) {
	url := strings.Join([]string{c.baseURL, "API_PROC_ORDER_CONFIRMATION_2_SRV", api}, "/")
	param := c.getQueryWithConfByOrderIDConfGroup(map[string]string{}, orderID, confirmationGroup)

	resp, err := c.requestClient.Request("GET", url, param, "")
	if err != nil {
		return nil, fmt.Errorf("API request error: %w", err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	data, err := sap_api_output_formatter.ConvertToConfirmation(byteArray, c.log)
	if err != nil {
		return nil, fmt.Errorf("convert error: %w", err)
	}
	return data, nil
}

func (c *SAPAPICaller) getQueryWithConfByOrderID(params map[string]string, orderID string) map[string]string {
	if len(params) == 0 {
		params = make(map[string]string, 1)
	}
	params["$filter"] = fmt.Sprintf("OrderID eq '%s'", orderID)
	return params
}

func (c *SAPAPICaller) getQueryWithMaterialMovements(params map[string]string, batch string) map[string]string {
	if len(params) == 0 {
		params = make(map[string]string, 1)
	}
	params["$filter"] = fmt.Sprintf("Batch eq '%s'", batch)
	return params
}

func (c *SAPAPICaller) getQueryWithBatchCharacteristic(params map[string]string, batch string) map[string]string {
	if len(params) == 0 {
		params = make(map[string]string, 1)
	}
	params["$filter"] = fmt.Sprintf("Batch eq '%s'", batch)
	return params
}

func (c *SAPAPICaller) getQueryWithConfByOrderIDConfGroup(params map[string]string, orderID, confirmationGroup string) map[string]string {
	if len(params) == 0 {
		params = make(map[string]string, 1)
	}
	params["$filter"] = fmt.Sprintf("OrderID eq '%s' and ConfirmationGroup eq '%s'", orderID, confirmationGroup)
	return params
}
