# upgraded-config-mang
upgraded-config-mang is a config management tool similar to Chef.

Note: This is repo is tested only on 2.Ubuntu-Linux, Also this can support only executions on remote node which are Ubuntu Os-Release

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
perm needs to have a padding 0, so a perm will look like `0644 || 0777`
`"depends":{"serviceName":"apache2"}` is an optional field, when provided it will restart the service defined, It ensures that the restart happens only if there is a change in the file 
```        
{ "action": "remoteFile", "localPath": "index.php", "remotePath": "/var/www/html/index.php", "owner":"root", "group": "root", "perm": "0644", "depends":{"serviceName":"apache2"}},
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
The tool takes 1 arguments, tasks

```
  -tasks string
    	Task file path, must be a JSON EG: ./tasks.json
```

Example to run the Binary:
```
./upgraded-config-mang-<OSNAME> -tasks ./tasks.json
```
Example to run from Go:
```
go run main.go -tasks ./tasks.json
```

## Download bootstrap.sh
ssh to the remote machine where this tool need to run and download the bootstrap.sh
```
curl -o ./bootstrap.sh https://raw.githubusercontent.com/reddydinesh427/upgraded-config-mang/main/bootstrap.sh
chmod +x bootstrap.sh 
./bootstrap.sh # this will install the necessary components and run the upgrade-config-mang tool
```
A sample output is added at [sample_output_bootstrap.txt](./sample_output_bootstrap.txt)





