if [ -z "$USE_PKGINST" ]; then
    if [ $USING_PACMAN = 1 ]; then
        if [ -x "$(command -v sudo)" ]; then
            USE_PKGINST='sudo pacman -S --noconfirm '
        else
            USE_PKGINST='pacman -S --noconfirm '
        fi
    elif [ -x "$(command -v dpkg)" ]; then
        USE_PKGINST='sudo apt-get install -y '
    elif [ -x "$(command -v dnf)" ]; then
        USE_PKGINST='sudo dnf install -y '
    elif [[ -x "$(command -v pkcon)" && -d /run/systemd/system ]]; then
        USE_PKGINST='pkcon install -y '
    fi
fi

if [ -v USE_PKGINST ]; then
    echo "INFO: package installer = $USE_PKGINST"
else
    echo
    echo "What's your package manager? such as:"
    echo '- `sudo pacman -S --noconfirm `'
    echo '- `sudo apt-get install -y `'
    echo '- `sudo dnf install -y `'
    echo '- ...'
    read -p 'Input the command: ' USE_PKGINST
fi
