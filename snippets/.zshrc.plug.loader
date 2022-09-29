source ~/.oh-my-zsh/lib/clipboard.zsh
detect-clipboard && plugins+=(copybuffer copyfile copypath)

for it in rsync tmux docker rustc kate node yarn; do
    [ -n "$(whence "$it")" ] && plugins+=("$it")
done

[ -n "$(whence fzf)" ] && plugins+=(fzf zsh-interactive-cd)

[ -n "$(whence go)" ] && plugins+=(golang)

[ -n "$(whence pip)" ] || [ -x "$(whence pip3)" ] && plugins+=(pip)
