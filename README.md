# azure-scrapper
Example Azure function custom handlers with go.

## Prerequisites Utilities
### az - Azure CLI
Used to log in to azure and provision the required resource group, storage, and function app

### func - Azure Function CLI
Used to test function locally simulating Azure runtime environment as well as deploy the application bundle.

### upx - Binary compression utility
Used to compress the function binary to reduce upload time at the cost of slower function start up.

## Getting started
### Creating the functionapp

Creating a resource group
```bash
az group create --name az-scrapper-rg --location eastus
```

Creating a storage account
```bash
az storage account create --name azscrapperstorage --location eastus --resource-group az-scrapper-rg --sku Standard_LRS --allow-blob-public-access false
```

Creating a functionapp
```bash
az functionapp create -g az-scrapper-rg -n az-scrapper -s azscrapperstorage --consumption-plan-location eastus --runtime custom --functions-version 4 --disable-app-insights
```
> NOTE: functionapp name must contain only uppercase, lowercase numbers and dashes, or you will get the exception `The parameter WEBSITE_CONTENTSHARE has an invalid value.`
