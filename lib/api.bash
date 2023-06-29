appendLineIfNotExist() {
   if [ -z "$2" ]; then
      echo >>"$1"
   elif ! grep -qF "$2" "$1"; then
      echo "$2" >>"$1"
   fi
}

trimFinalNewlines() {
   local tmp=$(mktemp)

   awk '/^$/{n=n RS}; /./{printf "%s",n; n=""; print}' $1 >$tmp
   rm $1
   mv $tmp $1
}

if [[ "$OSTYPE" == darwin* ]]; then
   [ ! -x "$(command -v gsed)" ] && brew install gsed || :
   alias sed=gsed
fi
