# Disclaimer
We do not condone, nor take any responsibility or liability whatsoever, under any circumstances, for any damages inflicted on property without explicit permission from the owner.
  
# Installation
There are several ways to install Emailbomber. Pre-compiled binaries are available for Windows, Linux and Darwin with AMD64 architecture. On Darwin and Linux, a Python installer is available.
  
**Python Installer**  
For Darwin and Linux, a Python installer is available. Simply run:  
`$ sh -c "$(curl https://raw.githubusercontent.com/amblar/master/install.sh -s)"`  
The correct binary will automatically be installed into `/usr/bin` or `/usr/local/bin` on MacOS.  
  
**Pre-compiled Binaries**  
Download the [latest release](https://github.com/amblar/emailbomber/releases/latest) from the releases tab which matches your system. Then copy the binary to `/usr/bin` or `/usr/local/bin` on MacOS, after that you will be able to run it using `emailbomber` in your terminal. Example:  
`$ tar -xvzf emailbomber-v0.1.0-linux-amd64.tar.gz`  
`$ sudo cp kryer-v0.1.0-linux-amd64/emailbomber /usr/bin/emailbomber`  
You can now run it:  
`$ emailbomber --help`  
  
On Windows, you can place the executable in a new directory in Program Files, for example, and then add it to your environment variables. 
  
To set your environment variables open Control Panel > System and Security > System > Advanced System Settings > Environment Variables
  
Now select path and click edit, then click browse and select the containing directory of the executable. Press OK and you should be able to run it using the `emailbomber` command in the command prompt.
  
**Building from Source**  
If pre-compiled binaries are not available for your system or you don't want to use them for other reasons, you can build Emailbomber yourself from source. To do so you will need a working Go environment. 
  
Start by cloning the repository into `YourGopath/src/github.com/amblar/emailbomber`. You can then build and install it using `$ sudo make install`, unless you do not have a `/usr/bin` directory (will install to `/usr/local/bin` on MacOS due to SIP).
  
If that is the case you can build using `$ make build` or `$ go build` in the Kryer directory.  
  
# Usage
To run Emailbomber:  
`$ emailbomber -a 50 -t 20 -n YourName -e youremail@gmail.com -r theiremail@gmail.com -P "yourpassword" -c yourtextfile.txt -m 2`  
