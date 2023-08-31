# upgraded-config-mang
upgraded-config-mang is a config management tool similar to Ansible. All the executions happen on ssh.

Note: This is repo is tested only on 1.Mac 2.Linux, Also this can support only executions on remote node which are Ubuntu Os-Release

## Inventory File
The tools takes a inventory file in json format which needs specific key, an example [inventory.json](./inventory.json) will look as follows
```
{
    "inventory":{
        "hosts": ["3.80.105.244","3.80.61.136"],
        "username": "root",
        "password": "xxxxxx"
    }
}

```
Note: This tool is tested only to test against username and password not with keys.

inventory: is a the parent where all hosts have same username and password.
hosts: is an array of hosts that can be defined.
username: username to login with
password: password to login with.

## Tasks File
The tool takes a tasks file where we defined all the tasks to be run, there is action key which determines what type of task to be run and each action will have its own keys.
An example [task.json](./tasks.json) file is as follows
```
{
    "tasks": [
        { "action": "packageInstall", "packageNames": ["apache2", "php", "libapache2-mod-php"]},
        { "action": "removeFile", "path": "/var/www/html/index.html" },
        { "action": "remoteFile", "localPath": "index.php", "remotePath": "/var/www/html/index.php", "owner":"root", "group": "root", "perm": "644", "depends":{"serviceName":"apache2"}},
        { "action": "serviceStart", "serviceName": "apache2"}
      ]
}
```
In the above example, there are actions for each type of task and their corresponding keys. Use next section to understand what each task does.

## Supporting Tasks
1. #### packageInstall
Installs a package on Ubuntu, The task can install multiple packages or single package.
Expects that all the repos are configured.
```
{ "action": "packageInstall", "packageNames": ["apache2", "php", "libapache2-mod-php"]},
```

2. #### removeFile
Removes a file from a path
```
{ "action": "removeFile", "path": "/var/www/html/index.html" },
```

3. #### remoteFile
Copys a file from local to remote server, It also ensures all the permissions are set properly
`"depends":{"serviceName":"apache2"}` is an optional field, when provided it will restart the service defined, It ensures that the restart happens only if there is a change in the file 
```        
{ "action": "remoteFile", "localPath": "index.php", "remotePath": "/var/www/html/index.php", "owner":"root", "group": "root", "perm": "644", "depends":{"serviceName":"apache2"}},
```

4. #### serviceStart
This task will Start a service, if the service is already running, it will skip it
```
{ "action": "serviceStart", "serviceName": "apache2"}
```

5. #### serviceStop
This task will stop a service, if the service is already stopped, it will skip it
```
{ "action": "serviceStop", "serviceName": "apache2"}
```

6. #### serviceRestart
This task will restart a service
```
{ "action": "serviceRestart", "serviceName": "apache2"}
```

7. #### packageRemove
Removes a package on Ubuntu, The task can remove multiple packages or single package.
Expects that all the repos are configured.
```
{ "action": "packageRemove", "packageNames": ["apache2"]},
```   

## Building the Binary

### Pre-requisites
* GoLang
* Make

Running bootstrap.sh will install `Go`, and any dependencies for building the package will be managed via `make`

### Build
This tool uses make to build the binary and spits out binary for linux and darwin
```make binary```

## Usage 
The tool takes 2 arguments, inventory and tasks

```
  -inventory string
    	inventory file path, must be a JSON, EG: ./inventory.json
  -tasks string
    	Task file path, must be a JSON EG: ./tasks.json
```

Example to run the Binary:
```
./upgraded-config-mang-<OSNAME> -inventory ./inventory.json -tasks ./tasks.json
```
Example to run from Go:
```
go run main.go -inventory ./inventory.json -tasks ./tasks.json
```






