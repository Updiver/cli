# Dumper CLI

Dumper CLI is an utility to dump repositories from github and bitbucket.

## Intro

`dumper-cli` allows you to quickly dump repositories from github or bitbucket without going over every repository and cloning that manually.

**Example**:

_Bitbucket_

```
./dumper-cli dump bitbucket -u USERNAME -d DESTINATION_FOLDER -t APP_PASSWORD
```

_Github_

```
./dumper-cli dump github -u USERNAME -d DESTINATION_FOLDER -t PERSONAL_ACCESS_TOKEN
```


## Get executable

Download executable from Relase page for your platform.

## Prerequisite

__Github__

To get personal access token for github account you need to:

- Go to Settings 
- Open Developer Settings (bottom left-hand side) 
- Personal access tokens and generate new one (classic)
- Select scopes (we need these three as we dump all repositories based on projects/user/repo):
  
  - repo
  - user
  - project

__Bitbucket__

To get application password for bitbucket account you need to:

- Go to Personal Bitbucket Settings
- App passwords menu item
- Create app password
- For Permissions:

  - Account -> Read
  - Workspace membership -> Read
  - Projects -> Read
  - Repositories -> Automatically should be Read

To get username you need to:

- Go to Personal Bitbucket Settings
- Account settings
- Bitbucket profile settings -> copy _Username_