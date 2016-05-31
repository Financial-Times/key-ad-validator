# Privision README

## Required ENV variables:

| Name 			| Example value 	|
|---------------	|---------------	|
|`AWS_DEFAULT_REGION` 	|eu-west-1    		|
|`AWS_SECRET_ACCESS_KEY`|asdasdasd		|
|`AWS_ACCESS_KEY_ID`	|asdasdasd		|
|`TF_VAR_keys_uri`	|test 			|
|`TF_VAR_ldap_server`	|test			|
|`TF_VAR_ldap_port`	|test			|
|`TF_VAR_ldap_user`	|test			|
|`TF_VAR_ldap_password`	|test			|

# Get the current state:

```sh
terraform remote config -backend=s3 -backend-config="bucket=ft.terraform.state-files" -backend-config="key=coco-key-ad-valid-svc/terraform.tfstate" -backend-config="region=eu-west-1"
```

# Plan

```sh
terraform plan
```

# Apply

```sh
terraform apply
```

