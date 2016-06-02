# Infrastructure provision README

## Required ENV variables:

### AWS keys to provision an instance

```sh
export AWS_DEFAULT_REGION=eu-west-1
export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

### Application

```sh
export TF_VAR_keys_uri=https://raw.githubusercontent.com/Financial-Times/up-ssh-keys/master/authorized_keys
export TF_VAR_ldap_server=ldap.server.com
export TF_VAR_ldap_port=389
export TF_VAR_ldap_user=tyler.durden
export TF_VAR_ldap_password=uniquesnowflake
```

Sync terraform state locally

```sh
terraform remote config -backend=s3 -backend-config="bucket=ft.terraform.state-files" -backend-config="key=coco-key-ad-valid-svc/terraform.tfstate" -backend-config="region=eu-west-1"
```

Plan - displayes the execution plan + changes if there any detected

```sh
terraform plan
```

Apply - apply the plan, will provision/decom/update resources according to the plan

```sh
terraform apply
```
