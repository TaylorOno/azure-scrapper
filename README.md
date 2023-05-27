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
az functionapp create -g az-scrapper-rg -n az_scrapper -s azscrapperstorage --consumption-plan-location eastus --runtime custom --disable-app-insights
```
> NOTE: using flag `--functions-version 4` results in an error `The parameter WEBSITE_CONTENTSHARE has an invalid value.` there does not appear to be any documentation about additional parameters required at this time.
