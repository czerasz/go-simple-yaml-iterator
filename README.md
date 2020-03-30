# Simple YAML Iterator

Iterates over items in a specified YAML file and executes specified command.
Item properties are passed to the command as environment variables.

## Usage

```bash
$ simple-yaml-iterator -h
Usage of simple-yaml-iterator:
  -command string
    	Command to execute for each item
  -dry-run
    	Only pretend to run
  -file string
    	Path to YAML file
  -root string
    	Root key inside YAML
  -version
    	Show application version
```

### Example

Given a `config.yml`:

```yaml
aws:
  - aws_profile: dev
    aws_default_region: eu-central-1
  - envaws_profile: prod
    aws_default_region: us-east-1
```

By runnig:

```bash
simple-yaml-iterator -file config.yml -root aws -command 'aws ec2 describe-instances' 
```

The following commands will be executed:

```bash
AWS_DEFAULT_REGION=eu-central-1 AWS_PROFILE=dev aws ec2 describe-instances
ENVAWS_PROFILE=prod AWS_DEFAULT_REGION=us-east-1 aws ec2 describe-instances
```

### Dry Run

Test it without executing commands (dry run mode):

```bash
simple-yaml-iterator -dry-run -file config.yml -root aws -command 'aws ec2 describe-instances' 
```

The output:

```bash
AWS_DEFAULT_REGION=eu-central-1 AWS_PROFILE=dev aws ec2 describe-instances
ENVAWS_PROFILE=prod AWS_DEFAULT_REGION=us-east-1 aws ec2 describe-instances
```