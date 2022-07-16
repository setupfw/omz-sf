#!/bin/bash
set -e
cd $(dirname $0)

# Constants
ZSH="${ZSH:-$HOME/.oh-my-zsh}"
ZSH_CUSTOM="$ZSH/custom"
ZSH_PLUGLOADER=${ZSH_PLUGLOAD:-$HOME/.zshrc.plugloader.zsh}
ZSH_PLUGLOADER_LIST=${ZSH_PLUGLOADER_LIST:-$HOME/.zshrc.plugloader-list.txt}

source lib/prompt-options
source lib/find-pkgmgr

[ ! -x "$(command -v zsh)" ] && $INSTPKG zsh
[ ! -x "$(command -v git)" ] && $INSTPKG git

if [ ! -d "$ZSH" ]; then
   [ ! -x "$(command -v curl)" ] && [ ! -x "$(command -v wget)" ] && $INSTPKG wget

   INSTALLER_URL=${INSTALLER_URL:-https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh}
   SCRIPT_CODE="$([ -x "$(command -v wget)" ] && wget "$INSTALLER_URL" -O - || curl -fsSL "$INSTALLER_URL")"
   RUNZSH=no sh -c "$SCRIPT_CODE"
fi

echo >>~/.zshrc
cat src/snippet-after-omz.zsh >>~/.zshrc

if [ -x /msys2.exe ]; then
   ismsys=1
fi

[ "$DISABLE_OMZ_AUTOUPDATE" = 1 ] &&
   sed -i "/disable automatic updates/ s/^#[ ]*//" ~/.zshrc

if [ "$USE_RECOMMEND_THEME" = 1 ]; then
   cp $ZSH/themes/steeef.zsh-theme $ZSH_CUSTOM/themes
   sed -i \
      -e '/^PROMPT=\$/i local exit_code="%(?,,C:%{$fg[red]%}%?%{$reset_color%})"' \
      -e "/^PROMPT=\\$'$/{n;s/$/ [%*] \$exit_code/}" \
      $ZSH_CUSTOM/themes/steeef.zsh-theme
   [ -x "$(command -v lsb_release)" ] && sed -e "/^PROMPT=\\$'$/{n;s/%m/&(\$(lsb_release -si))/}" -i $ZSH_CUSTOM/themes/steeef.zsh-theme
   sed -i 's/^ZSH_THEME=".*"/ZSH_THEME="steeef"/' ~/.zshrc
fi

if [ "$USE_PLUGLOADER" = 1 ]; then
   cat <<END >$ZSH_PLUGLOADER
plugins=()
while read -r line; do
    if [[ "\$line" != '#'* ]]; then
        read -A list <<<"\$line"
        plugins+=("\${list[@]}")
    fi
done <$ZSH_PLUGLOADER_LIST
END

   if [ "$USE_RECOMMEND_PLUGINS" = 1 ]; then
      [ ! -f "$ZSH_PLUGLOADER_LIST" ] &&
         cat src/recommended-plugins.txt >$ZSH_PLUGLOADER_LIST

      src/recommended-plugins-tweaker $ZSH_PLUGLOADER $ZSH_PLUGLOADER_LIST
   else
      touch $ZSH_PLUGLOADER_LIST
   fi

   [ "$USE_COMMON_ALIASES" = 1 ] && echo common-aliases >>"$ZSH_PLUGLOADER_LIST"
fi

sed -i "s#plugins=(git)\$#source \"$ZSH_PLUGLOADER\"#" ~/.zshrc
echo "INFO: Plugins loader at: $ZSH_PLUGLOADER"
echo "INFO: Plugins list at: $ZSH_PLUGLOADER_LIST"
echo

if [ "$USE_SYNTAX_HIGHLIGHT" = 1 ]; then
   $INSTPKG zsh-syntax-highlighting
   if [ -x "$(command -v pacman)" ]; then
      echo 'source /usr/share/zsh/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh' >>~/.zshrc
   else
      echo 'source /usr/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh' >>~/.zshrc
   fi
fi

if [ "$USE_AUTO_SUGGEST" = 1 ]; then
   $INSTPKG zsh-autosuggestions
   if [ -x "$(command -v pacman)" ]; then
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
   $INSTPKG pkgfile
   echo 'source /usr/share/doc/pkgfile/command-not-found.zsh' | tee -a ~/.zshrc >/dev/null
   echo "==> Run 'sudo pkgfile -u'"
   if [ -x "$(command -v lsb_release)" ]; then
      sudo pkgfile -u
   elif [ -v ismsys ]; then
      pkgfile -u
   fi
fi

if [ -v ismsys ]; then
   echo "alias sudo=''" | tee -a ~/.zshrc >/dev/null
fi

echo 'source ~/.zshrc.local' >>~/.zshrc
touch ~/.zshrc.local
