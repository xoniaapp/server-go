# Functions

# This is for getting started ascii art.
function gettingStarted() {
echo "
 ██████╗ ███████╗████████╗████████╗██╗███╗   ██╗ ██████╗     ███████╗████████╗ █████╗ ██████╗ ████████╗███████╗██████╗          
██╔════╝ ██╔════╝╚══██╔══╝╚══██╔══╝██║████╗  ██║██╔════╝     ██╔════╝╚══██╔══╝██╔══██╗██╔══██╗╚══██╔══╝██╔════╝██╔══██╗         
██║  ███╗█████╗     ██║      ██║   ██║██╔██╗ ██║██║  ███╗    ███████╗   ██║   ███████║██████╔╝   ██║   █████╗  ██║  ██║         
██║   ██║██╔══╝     ██║      ██║   ██║██║╚██╗██║██║   ██║    ╚════██║   ██║   ██╔══██║██╔══██╗   ██║   ██╔══╝  ██║  ██║         
╚██████╔╝███████╗   ██║      ██║   ██║██║ ╚████║╚██████╔╝    ███████║   ██║   ██║  ██║██║  ██║   ██║   ███████╗██████╔╝██╗██╗██╗
 ╚═════╝ ╚══════╝   ╚═╝      ╚═╝   ╚═╝╚═╝  ╚═══╝ ╚═════╝     ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═════╝ ╚═╝╚═╝╚═╝
"
}

function finished(){
    echo "    
███████╗██╗███╗   ██╗██╗███████╗██╗  ██╗███████╗██████╗ ██╗
██╔════╝██║████╗  ██║██║██╔════╝██║  ██║██╔════╝██╔══██╗██║
█████╗  ██║██╔██╗ ██║██║███████╗███████║█████╗  ██║  ██║██║
██╔══╝  ██║██║╚██╗██║██║╚════██║██╔══██║██╔══╝  ██║  ██║╚═╝
██║     ██║██║ ╚████║██║███████║██║  ██║███████╗██████╔╝██╗
╚═╝     ╚═╝╚═╝  ╚═══╝╚═╝╚══════╝╚═╝  ╚═╝╚══════╝╚═════╝ ╚═╝
"
}

### = Main = ###

URL="https://github.com/Xoniaapp/Server/releases/download/$2/xoniaapp-$2-linux-amd64.tar.gz"
VERSION="1.0.0"
# Help Command

if [ $1 = "help" ]
then
    echo "
 ██████╗██╗     ██╗
██╔════╝██║     ██║
██║     ██║     ██║
██║     ██║     ██║
╚██████╗███████╗██║
 ╚═════╝╚══════╝╚═╝
    "
    echo "Available Commands"
    echo "---------------------------------------------------------------------"
    echo "help              : Shows this menu"
    echo "install <version> : Installs the required packages."
    echo "uninstall         : Uninstalls the packages that came with install"
    echo "----------------------------------------------------------------------"
    echo "Version $VERSION"
fi
if [ $1 = "install" ]
then
    gettingStarted # For the ASCII Art.
    echo "Updating Repository..."
    sudo apt update -y
    echo "Upgrading System..."
    sudo apt upgrade -y
    echo "Installing required packages..."
    sudo apt install curl git wget nano htop -y
    echo "Installing Nodejs & NPM..."
    curl -fsSL https://deb.nodesource.com/setup_17.x | sudo -E bash -
    sudo apt-get install -y nodejs
    echo "Installing PM2..."
    sudo npm i -g pm2
    echo "Deleting xoniaapp-api-server directory..." # Delete if there was one
    rm -rf ./xoniaapp-server
    echo "Making new dirctory..."
    mkdir ./xoniaapp-server
    echo "Changing directory to /xoniaapp-server"
    cd ./xoniaapp-server
    echo "Detected Version $2 | Downloading from CI/CD Server..."
    wget "$URL"
    echo "Unzipping binary..."
    tar xvf xoniaapp*.tar.gz
    echo "Clearning unless files..."
    sudo rm -rf deploy.sh LICENSE .env.example xoniaapp-$2-linux-amd64.tar.gz
    echo "Editing .env file..."
    nano .env
    finished
fi
if [ $1="start" ]
then
    echo "Starting..."
    pm2 start ./xoniaapp-server/xoniaapp
fi
if [ $1="stop" ]
then
    echo "Starting..."
    pm2 stop xoniaapp
fi