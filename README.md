# Compile
```
GOARM=6 GOARCH=arm GOOS=linux go build
```

# Raspbery pi setup
Enable auto login
```
sudo raspi-config
```
Select “Boot Options” then “Desktop/CLI” then “Console Autologin”

## **Startup with SYSTEMD**

The best method (that I've found) to running a go program on a Raspberry Pi at startup is to use the **systemd** files. **systemd** provides a standard process for controlling what programs run when a Linux system boots up. 

Note that **systemd** is available only from the Jessie versions of Raspbian OS.

### **1. Create A Unit File**

Create a service file at the following location:
```
sudo touch /lib/systemd/system/simpsons.service
```

Edit the file to look like this:

```
[Unit]
 Description=Simpsons Shuffler
 After=multi-user.target

 [Service]
 Type=idle
 ExecStart=/home/pi/simpsons

 [Install]
 WantedBy=multi-user.target
```

This defines a new service called “Simpsons Shuffler” and we are requesting that it is launched once the multi-user environment is available. The “ExecStart” parameter is used to specify the command we want to run. The “Type” is set to “idle” to ensure that the ExecStart command is run only when everything else has loaded. Note that the paths are absolute.

The permission on the unit file needs to be set to 644 :

```
sudo chmod 644 /lib/systemd/system/simpsons.service
```

### **2. Configure systemd**

Now the unit file has been defined we can tell systemd to start it during the boot sequence :

```
sudo systemctl daemon-reload
sudo systemctl enable simpsons.service
```

Reboot the Pi and your custom service should run:

```
sudo reboot
```



## update app if neded
```
sudo systemctl stop simpsons.service
 ```

 copy file

 reboot