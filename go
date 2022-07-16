#!/bin/bash
set -e

# Constants
ZPLG_LIST=${ZPLG_LIST:-$HOME/.zshrc.pluglist}
ZPLG_LOADER=${ZPLG_LOADER:-$HOME/.zshrc.plugloader}
ZSH="${ZSH:-$HOME/.oh-my-zsh}"
ZSH_CUSTOM="$ZSH/custom"

# Environment
DIR=$(dirname "$0")

if [ -x "$(command -v lsb_release)" ]; then
    DISTRO=$(lsb_release -si)
elif [ -x /msys2.exe ]; then
    MSYS2=1
fi

USING_PACMAN=$(
    [ ! -x "$(command -v pacman)" ]
    echo $?
)

function ask_common_aliases() {
    [ -v USE_COMMON_ALIASES ] && return
    echo
    echo "Document: https://github.com/ohmyzsh/ohmyzsh/tree/master/plugins/common-aliases"
    read -p 'Use ohmyzsh common-aliases plugin (y/N)? ' r
    [[ "$r" =~ ^(Y|y|)$ ]] && USE_COMMON_ALIASES=1
}

function ask_recommend_theme() {
    [ -v USE_RECOMMEND_THEME ] && return
    [ -f ~/.zshrc ] && grep -q '^ZSH_THEME="steeef"$' ~/.zshrc && return
    if [[ ! -f ~/.zshrc || "$(grep '^ZSH_THEME=' ~/.zshrc | cut -d'"' -f2)" != "steeef" ]]; then
        echo
        read -p 'Use ohmyzsh steeef theme (Y/n)? ' r
        [[ "$r" =~ ^(Y|y|)$ ]] && USE_RECOMMEND_THEME=1
    fi
}

function ask_syntax_highlight() {
    [ -v USE_SYNTAX_HIGHLIGHT ] && return
    [ -f ~/.zshrc ] && grep -q '^source .*zsh-syntax-highlighting.zsh' ~/.zshrc && return
    echo
    read -p 'Install zsh-syntax-highlighting (Y/n)? ' r
    [[ "$r" =~ ^(Y|y|)$ ]] && USE_SYNTAX_HIGHLIGHT=1
}

function ask_auto_suggest() {
    [ -v USE_AUTO_SUGGEST ] && return
    [ -f ~/.zshrc ] && grep -q '^source .*zsh-autosuggestions.zsh' ~/.zshrc && return
    echo
    read -p 'Install zsh-autosuggestions (Y/n)? ' r
    [[ "$r" =~ ^(Y|y|)$ ]] && USE_AUTO_SUGGEST=1
}

function ask_pacman_pkgfile() {
    [ -v USE_PACMAN_PKGFILE ] && return
    [[ ! -z "$(pacman -Qs zsh)" ]] && grep -q '^source /usr/share/doc/pkgfile/command-not-found.zsh$' ~/.zshrc && return
    echo
    read -p 'Install `pkgfile` to provide command-not-found advice (Y/n)? ' r
    [[ "$r" =~ ^(Y|y|)$ ]] && USE_PACMAN_PKGFILE=1
}

source "$DIR/_findpkgmgr.bash"
ask_common_aliases
ask_recommend_theme
if [ -v DISTRO ]; then
    ask_syntax_highlight
    ask_auto_suggest
fi
[ $USING_PACMAN = 1 ] && ask_pacman_pkgfile

[ ! -x "$(command -v zsh)" ] && $USE_PKGINST zsh
test -d "$ZSH" || NOLOG=1 RUNZSH=no "$DIR/instomz"
echo >>~/.zshrc

cat "$DIR/_recommend.zsh" >>~/.zshrc

# No Update:
sed -i "/disable automatic updates/ s/^#[ ]*//" ~/.zshrc
echo 'Manually update: omz update'

if [ -v USE_RECOMMEND_THEME ]; then
    cp $ZSH/themes/steeef.zsh-theme $ZSH_CUSTOM/themes
    sed -e '/^PROMPT=\$/i local exit_code="%(?,,C:%{$fg[red]%}%?%{$reset_color%})"' -i $ZSH_CUSTOM/themes/steeef.zsh-theme
    sed -e "/^PROMPT=\\$'$/{n;s/$/ [%*] \$exit_code/}" -i $ZSH_CUSTOM/themes/steeef.zsh-theme
    if [ -v DISTRO ]; then
        sed -e "/^PROMPT=\\$'$/{n;s/%m/&(\$(lsb_release -si))/}" -i $ZSH_CUSTOM/themes/steeef.zsh-theme
    fi
    sed -i 's/^ZSH_THEME=".*"/ZSH_THEME="steeef"/' ~/.zshrc
fi

# Plugins:
TMPL_ZPLG_LIST="$DIR/_omzpluglist"
TMPL_ZPLG_LOADER="$DIR/_omzplugloader.zsh"

[[ -f "$ZPLG_LIST" ]] || cat "$TMPL_ZPLG_LIST" >"$ZPLG_LIST"
[ -x "$(command -v systemctl)" ] && echo systemd >>"$ZPLG_LIST"
if [ -x "$(command -v dpkg)" ]; then
    if [ "$DISTRO" = Ubuntu ]; then
        echo ubuntu >>"$ZPLG_LIST"
    else
        echo debian >>"$ZPLG_LIST"
    fi
    echo 'unalias acs &>/dev/null' >>~/.zshrc
    echo "alias acse='apt-cache search'" >>~/.zshrc
fi
[ -x "$(command -v dnf)" ] && echo dnf >>"$ZPLG_LIST"
[ "$USING_PACMAN" = 1 ] && echo archlinux >>"$ZPLG_LIST"

cat "$TMPL_ZPLG_LOADER" >"$ZPLG_LOADER"
sed -i "s#plugins=(git)\$#source \"$ZPLG_LOADER\"#" ~/.zshrc
echo "INFO: Plugins loader at: $ZPLG_LOADER"
echo "INFO: Plugins list at: $ZPLG_LIST"
echo

if [ "$USE_SYNTAX_HIGHLIGHT" = 1 ]; then
    $USE_PKGINST zsh-syntax-highlighting
    if [ $USING_PACMAN = 1 ]; then
        echo 'source /usr/share/zsh/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh' >>~/.zshrc
    else
        echo 'source /usr/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh' >>~/.zshrc
    fi
fi

if [ "$USE_AUTO_SUGGEST" = 1 ]; then
    $USE_PKGINST zsh-autosuggestions
    if [ $USING_PACMAN = 1 ]; then
        echo 'source /usr/share/zsh/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh' >>~/.zshrc
    else
        echo 'source /usr/share/zsh-autosuggestions/zsh-autosuggestions.zsh' >>~/.zshrc
    fi
    cat <<END >>~/.zshrc
pasteinit(){OLD_SELF_INSERT=\${\${(s.:.)widgets[self-insert]}[2,3]};zle -N self-insert url-quote-magic}
pastefinish(){zle -N self-insert \$OLD_SELF_INSERT}
zstyle :bracketed-paste-magic paste-init pasteinit
zstyle :bracketed-paste-magic paste-finish pastefinish
END
fi

if [ "$USE_PACMAN_PKGFILE" = 1 ]; then
    $USE_PKGINST pkgfile
    echo 'source /usr/share/doc/pkgfile/command-not-found.zsh' | tee -a ~/.zshrc >/dev/null
    echo "==> Run 'sudo pkgfile -u'"
    if [ -v DISTRO ]; then
        sudo pkgfile -u
    elif [ -v MSYS2 ]; then
        pkgfile -u
    fi
fi

if [ -v MSYS2 ]; then
    echo 'alias sudo=""' | tee -a ~/.zshrc >/dev/null
fi

[ "$USE_COMMON_ALIASES" = 1 ] && echo common-aliases >>"$ZPLG_LIST"

echo >>~/.zshrc

[[ "$NOLOG" = 1 ]] || echo "$0" >>~/deployworkenv.log
