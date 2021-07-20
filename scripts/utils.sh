source $shdir/scripts/const.sh

function cecho {
    timestamp=`date --rfc-3339=seconds`
    RED="\033[1;31m"
    GREEN="\033[0;32m"
    YELLOW="\033[0;33m"
    BLUE="\033[1;34m"
    DARKGRAY="\033[0;90m"
    NC="\033[0m"
    ERROR="\033[1;31m[ERROR]\t"
    SUCCESS="\033[0;32m[SUCCESS]\t"
    WARNING="\033[0;33m[WARNING]\t"
    INFO="\033[1;34m[INFO]\t"
    str="${@:2:${#@}}"
    printf "[$timestamp] ${!1} ${str}  ${NC}\n"
}

function print_error {
  cecho "ERROR" "${FUNCNAME[2]}"
}

function print_success {
  cecho "SUCCESS" "${FUNCNAME[2]}"
}

function command {
    cecho "INFO" "$1"
    $1 >&$bdir/log.txt
    cat $bdir/log.txt
    res=$?
    if [ $res -ne 0 ]; then
      print_error
      # exit
    fi
    print_success
}