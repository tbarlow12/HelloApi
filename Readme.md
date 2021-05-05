# Deploy the HelloApi to Azure VM

This documentation breaks down the steps to deploy the .net framework app (HelloApi) to azure VM.

## Steps Break Down

### Step 1: Deploy Azure Resource using terraform script

```text
terraform workspace new dev
terraform init
terraform plan --out plan.out
terraform apply plan.out
```

### Step 2: Set up Bastion access on vm from Azure portal

[Instruction here](https://docs.microsoft.com/en-us/azure/bastion/quickstart-host-portal)

[Todo Terraform](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/bastion_host)

### Step 3: Set up IIS server within VM

Use Bastion access into vm, and set up the IIS server inside vm.

[Instruction here](https://github.com/aspnet/Tooling/blob/AspNetVMs/docs/create-asp-net-vm-with-webdeploy.md#install-iis-web-server-plus-web-management-service-and-aspnet-46)

### Step 4: Download and set up Web Deploy 3.6

[Download here](https://www.microsoft.com/en-us/download/details.aspx?id=43717)

[Instruction here](https://github.com/aspnet/Tooling/blob/AspNetVMs/docs/create-asp-net-vm-with-webdeploy.md#install-web-deploy-36)

### Step 5: Set up Windows Firewall access

Inside the window vm, open Windows Defender Firewall Advanced settings.
Add inbound rule and outbound rule for port 8172 and port 80 under TCP portocal.

[Firewall port doc](https://docs.microsoft.com/en-us/aspnet/web-forms/overview/deployment/configuring-server-environments-for-web-deployment/configuring-a-web-server-for-web-deploy-publishing-web-deploy-handler#configure-firewall-exceptions)

### Step 6: Publish .net framework app through Visual Studio

[Instruction here](https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-dev-start-howto-vm-dotnet?view=azs-2008#deploy-and-run-the-app)
