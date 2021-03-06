try:
    import os
    import subprocess
    import shutil
    import stat
    import json
    import sys
    import platform
    import tarfile
    import hashlib

    PYTHON2 = 2
    PYTHON3 = 3

    pyversion = None
    installpath = "/usr/bin"

    # Detect OS and architecture
    if sys.version_info[0] < 3:
        pyversion = PYTHON2
        import urllib2 as urllib
    else:
        pyversion = PYTHON3
        import urllib.request as urllib

    if platform.machine() != "x86_64":
        if platform.machine() != "":
            print("System architecture unsupported, build from source instead")
            exit()
        elif pyversion == PYTHON3:
            q = input(
                "Unable to detect system architecture, are you using x86_64 (AMD64) [y/N]: ")
            if q.lower() != "y":
                print("System architecture unsupported, build from source instead")
                exit(1)
        else:
            q = raw_input(
                "Unable to detect system architecture, are you using x86_64 (AMD64) [y/N]: ")
            if q.lower() != "y":
                print("System architecture unsupported, build from source instead")
                exit(1)

    if platform.system() != "Linux" and platform.system() != "Darwin":
        if platform.system() == "Windows":
            print(
                "Windows is not supported by the Python installer, use pre-compiled binaries instead")
        else:
            print("OS unsupported, build from source instead")
        exit(1)

    if platform.system() == "Darwin":
        installpath = "/usr/local/bin"

    # Check for existance and write permissions
    if not os.path.isdir(installpath):
        print("Unable to locate " + installpath)
        exit(1)
    if not os.path.isdir("/tmp"):
        print("Unable to locate /tmp")
        exit(1)

    if not os.access(installpath, os.W_OK):
        print("Insufficient permissions, please elevate")
        exit(1)
    if not os.access("/tmp", os.W_OK):
        print("Insufficient permissions, please elevate")
        exit(1)

    # Get latest release
    print("Gathering system info ...")
    if os.path.isfile(installpath + "/emailbomber"):
        if pyversion == PYTHON2:
            q = raw_input(
                "emailbomber is already installed, reinstall/update [Y/n]: ")
        else:
            q = str(
                input("emailbomber is already installed, reinstall/update [Y/n]: "))

        if q.lower() == "n":
            exit(0)

        print("Removing " + installpath + "/emailbomber...")

        try:
            os.remove(installpath + "/emailbomber")
        except OSError:
            print("Insufficient permissions, please elevate")
            exit(1)

        print("Removed " + installpath + "/emailbomber...")

    print("Getting release information from https://api.github.com/repos/amblar/emailbomber/releases...")
    response = urllib.urlopen(
        "https://api.github.com/repos/amblar/emailbomber/releases")
    releases = json.loads(response.read())

    for release in releases:
        print("Getting asset list from " + release["assets_url"] + "...")
        response = urllib.urlopen(release["assets_url"])
        assets = json.loads(response.read())

        for asset in assets:
            if platform.system().lower() in asset["name"] and "tar.gz" in asset["name"]:
                print("Downloading " + asset["browser_download_url"] + "...")
                if(pyversion == PYTHON2):
                    response = urllib.urlopen(asset["browser_download_url"])
                    rawfile = response.read()

                    with open("/tmp/emailbomber.tar.gz", "wb") as FILE:
                        FILE.write(rawfile)

                else:
                    urllib.urlretrieve(
                        asset["browser_download_url"], "/tmp/emailbomber.tar.gz")
                break

        if os.path.isfile("/tmp/emailbomber.tar.gz"):
            break

    if not os.path.isfile("/tmp/emailbomber.tar.gz"):
        print("There aren't any compatible pre-compiled binaries available for your system, build from source instead")
        exit(1)

    print("Extracting " + asset["name"] + "...")
    filename = asset["name"].replace(".tar.gz", "")
    for asset in assets:
        if asset["name"] == filename + ".sha256":
            print("Getting SHA256 checksum from " +
                  asset["browser_download_url"] + "...")
            if pyversion == PYTHON2:
                response = urllib.urlopen(asset["browser_download_url"])
                rawfile = response.read()

                with open("/tmp/emailbomber.sha256", "wb") as FILE:
                    FILE.write(rawfile)

            else:
                urllib.urlretrieve(
                    asset["browser_download_url"], "/tmp/emailbomber.sha256")
            break

    if os.path.isfile("/tmp/emailbomber.sha256"):
        with open("/tmp/emailbomber.sha256", "r") as checksum:
            checksum = checksum.read().split(" ")[0]
        print("Verifying checksum " + checksum + "...")

        checksumFILE = hashlib.sha256()
        with open("/tmp/emailbomber.tar.gz", "rb") as FILE:
            for chunk in iter(lambda: FILE.read(4096), b""):
                checksumFILE.update(chunk)

        if checksum.strip() != checksumFILE.hexdigest().strip():
            print("Integrity check failed, invalid checksum")
            print("Calculated checksum: " + checksumFILE.hexdigest().strip())
            print("Provided checksum: " + checksum.strip())

            if os.path.isfile("/tmp/emailbomber"):
                os.remove("/tmp/emailbomber")
            if os.path.isdir("/tmp/emailbomber.sha256"):
                shutil.rmtree("/tmp/emailbomber.sha256")
            if os.path.isfile("/tmp/emailbomber.tar.gz"):
                os.remove("/tmp/emailbomber.tar.gz")
            exit(1)

        print("Integrity check successful")
    else:
        print("Can't find checksum for file, skipping integrity check...")

    print("Extracting " + asset["name"] + "...")
    os.mkdir("/tmp/emailbomber")
    tar = tarfile.TarFile.open("/tmp/emailbomber.tar.gz", "r:gz")
    tar.extractall("/tmp/emailbomber")
    tar.close()

    print("Creating files in " + installpath + "...")
    shutil.move("/tmp/emailbomber/" + filename + "/emailbomber", installpath + "/emailbomber")

    shutil.rmtree("/tmp/emailbomber")
    os.remove("/tmp/emailbomber.tar.gz")
    os.remove("/tmp/emailbomber.sha256")
    print("Installation successful, use emailbomber command to start")

except KeyboardInterrupt:
    print("\nKeyboard interrupt detected")
    print("Cleaning up...")

    if os.path.isdir("/tmp/emailbomber"):
        shutil.rmtree("/tmp/emailbomber")
    if os.path.isfile("/tmp/emailbomber.tar.gz"):
        os.remove("/tmp/emailbomber.tar.gz")
    if os.path.isfile("/tmp/emailbomber.sha256"):
        os.remove("/tmp/emailbomber.sha256")