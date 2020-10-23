# d7024e-lab-grp-1
[Trello Backlog](https://trello.com/b/JBu8iKev/d7024e)
# Commands & usage
## Unit test
To run unit tests for the implementation run the following command from `d7024e-lab-grp-1/src/d7024e/`
```
go test -cover
```
## Application
### Building Docker image
To build the Docker image with the name `ubuntu` run the following command from the root directory `d7024e-lab-grp-1/`
```
sudo docker build -t ubuntu -f Dockerfile . 
```
### Spinup containers
To spinup `50` container using the created image `ubuntu` run the following command from the root directory `d7024e-lab-grp-1/`
```
sh spinup.sh 50 ubuntu
```
This command will also clear all previous containers on the network by running the `clear.sh` script before spinning up new containers.
### IMPORTANT!
For some unknown reason all containers are not continuing execution after being started every time. This can be checked by running the command in section **Check for available containers and get their reference name** below. When running the command the status and some other information for each launched container will be displayed in the `STATUS` field. If the spinup command in section **Spinup containers** was previously run every created container should have a status showing the time they have been up for. If the `STATUS` field instead says the container has exited the spinup command in section **Spinup containers** should be run and again and this process should be repeated untill all containers are running.
### Check for available containers and get their reference name
To check what containers are available on the network by running the following command from the root directory ``d7024e-lab-grp-1/`
```
sudo docker ps -a
```
Note that the `NAMES` field will be the reference to use when attaching to a cointainer see section **Attach to a cointer to allow use of CLI**.
### Attach to a container to allow use of the Kademlia CLI
To attach to a container with the reference name `u1` use the following command from the root directory `d7024e-lab-grp-1/`
```
sudo docker attach u1
```
### Detach from a container without shutting down the container
To detach from a container without shutting down the container use the hotkey command `CTRL` + `P` + `Q`.
For exiting and ternminating the cointainer see section **Kademlia CLI command exit**.
### Using the Kademlia CLI
Once attached to a running container the Kademlia CLI should be available. Get a list of available commands in the Kademlia CLI by using the following command
```
help
```
#### Kademlia CLI command put
While in the Kademlia CLI the put command is available and used by running the following command
```
put filename content
```
where `filename` is the value that gets hashed and used as a key and `content` is the content that will be stored and retrieved when using the related key.
#### Kademlia CLI command get
While in the Kademlia CLI the put command is available and used by running the following command
```
get filename
```
where `filename` is the unhashed key (should be same as the entered value `filename` when using the command put to store data).
#### Kademlia CLI command exit
While in the Kademlia CLI the exit command is available and used by running the following command
```
exit
```
This will exit and terminate the container. For exiting the container without ternminating see section [**Detach from a container without shutting down the container**](#building-docker-image).
