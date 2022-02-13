if [ $1 = "help" ]
then
        echo "Commands:"
        echo "1 - help (shows this help menu)"
        echo "2 - stop (stops the Xonia server)"
        echo "3 - start (starts the Xonia server)"
        echo "4 - edit-env (edits server's env vars, specify editor with the second parameter, like so: edit-server-env <editor, eg vim>)"
        echo "5 - replace-server-binary (replaces server binary, use the command like so: replace-server-binary <binary filename and path, eg ./new_server>)"
        echo "6 - restart (restarts the server, use the command after replace-server-binary)"
        echo "7 - status (shows server status)"
        echo "8 - deploy (deploys specified version of server, use it like this: deploy <version, eg 2.0.3>)"
fi
if [ $1 = "stop" ]
then
        sudo systemctl stop xonia-server
        echo "done"
fi
if [ $1 = "start" ]
then
        sudo systemctl start xonia-server
        echo "done"
fi
if [ $1 = "edit-env" ]
then
        sudo $2 ~/.local/bin/xonia-server/.env
fi
if [ $1 = "replace-server-binary" ]
then
        sudo cp -f $2 ~/.local/bin/xonia-server/server-binary
        echo "done"
fi
if [ $1 = "restart" ]
then
        sudo systemctl restart xonia-server
        echo "done"
fi
if [ $1 = "status" ]
then
        sudo systemctl status xonia-server
fi
if [ $1 = "deploy" ]
then
        echo "Downloading binary archive..."
        wget "https://github.com/Xoniaapp/Server/releases/download/$2/xoniaapp-$2-linux-amd64.tar.gz"
        echo "Extracting archive..."
        tar xf xoniaapp-*.tar.gz
        echo "Installing binary from archive..."
        mxs replace-server-binary xoniaapp
        echo "Restarting server..."
        mxs restart
        echo "Cleaning up junk..."
        rm -f xoniaap* .env.example LICENSE
        echo "Deployment finished!"
        echo "Do 'mxs status' to check server status just in case if the server has crashed."
fi
