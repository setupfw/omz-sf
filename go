#!/bin/sh
set -e
cd $(dirname $0)

# Constants
ZSH=${ZSH:-$HOME/.oh-my-zsh}
ZSH_CUSTOM="$ZSH/custom"
ZSH_PLUGLOADER=${ZSH_PLUGLOADER:-$HOME/.zshrc.plug.loader}
ZSH_PLUGLOADER_LIST=${ZSH_PLUGLOADER_LIST:-$HOME/.zshrc.plug.list}

source lib/prompt-options
source lib/find-pkgmgr

if [ ! -x "$(command -v zsh)" ]; then $INSTPKG zsh; fi
if [ ! -x "$(command -v git)" ]; then $INSTPKG git; fi

if [ ! -d "$ZSH" ]; then
   if [ ! -x "$(command -v curl)" ] && [ ! -x "$(command -v wget)" ]; then
      $INSTPKG wget
   fi

   INSTALLER_URL=${INSTALLER_URL:-https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh}
   SCRIPT_CODE="$([ -x "$(command -v wget)" ] && wget "$INSTALLER_URL" -O - || curl -fsSL "$INSTALLER_URL")"
   RUNZSH=no sh -c "$SCRIPT_CODE"
fi

echo >>~/.zshrc
cat src/snippet-after-omz.zsh >>~/.zshrc

if [ -x /msys2.exe ]; then
   ismsys=1
fi

if [ "$DISABLE_OMZ_AUTOUPDATE" = 1 ]; then
   sed -i "/disable automatic updates/ s/^#[ ]*//" ~/.zshrc
fi

if [ "$USE_RECOMMEND_THEME" = 1 ]; then
   cp $ZSH/themes/steeef.zsh-theme $ZSH_CUSTOM/themes
   sed -i \
      -e '/^PROMPT=\$/i local exit_code="%(?,,C:%{$fg[red]%}%?%{$reset_color%})"' \
      -e "/^PROMPT=\\$'$/{n;s/$/ [%*] \$exit_code/}" \
      $ZSH_CUSTOM/themes/steeef.zsh-theme
   if [ -x "$(command -v lsb_release)" ]; then
      sed -e "/^PROMPT=\\$'$/{n;s/%m/&(\$(lsb_release -si))/}" -i $ZSH_CUSTOM/themes/steeef.zsh-theme
   fi
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

   if [ "$USE_COMMON_ALIASES" = 1 ]; then
      echo common-aliases >>"$ZSH_PLUGLOADER_LIST"
   fi
fi

sed -i "s#plugins=(git)\$#source \"$ZSH_PLUGLOADER\"#" ~/.zshrc
echo "• Plugins dynamic loader: $ZSH_PLUGLOADER"
echo "• Plugins static list: $ZSH_PLUGLOADER_LIST"
echo

if [ "$USE_SYNTAX_HIGHLIGHT" = 1 ]; then
   $INSTPKG zsh-syntax-highlighting
   dirs='/usr/share/zsh/plugins/zsh-syntax-highlighting
/usr/share/zsh-syntax-highlighting'

   for dir in $dir; do
      if [ -d "$dir" ]; then
         echo 'source $dir/zsh-syntax-highlighting.zsh' >>~/.zshrc
         break
      fi
   done
fi

if [ "$USE_AUTO_SUGGEST" = 1 ]; then
   $INSTPKG zsh-autosuggestions
   dirs='/usr/share/zsh/plugins/zsh-autosuggestions
/usr/share/zsh-autosuggestions'

   for dir in $dir; do
      if [ -d "$dir" ]; then
         echo 'source $dir/zsh-autosuggestions.zsh' >>~/.zshrc
         break
      fi
   done

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
   elif [ -n "$ismsys" ]; then
      pkgfile -u
   fi
fi

if [ -n "$ismsys" ]; then
   echo "alias sudo=''" | tee -a ~/.zshrc >/dev/null
fi

echo 'source ~/.zshrc.local' >>~/.zshrc
touch ~/.zshrc.local

if [ -n "$USE_LOCALIZE" ]; then
   cat src/snippet-localize.zsh >>~/.zshrc.local
fi

if [ "$RUNZSH" != no ]; then
   exec zsh -l
fi
