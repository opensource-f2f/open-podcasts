#!/bin/sh

while [ $# -gt 0 ]; do
  case "$1" in
    --server=*)
      server="${1#*=}"
      ;;
    --showFile=*)
      showFile="${1#*=}"
      ;;
    --itemsPattern=*)
      itemsPattern="${1#*=}"
      ;;
    --output=*)
      output="${1#*=}"
      ;;
    *)
      printf "***************************\n"
      printf "* Error: Invalid argument.*\n"
      printf "***************************\n"
      exit 1
  esac
  shift
done

if [ "$showFile" == "" ] || [ "$itemsPattern" == "" ] || [ "$output" == "" ]
then
  echo "Flags --showFile, --itemsPattern, --output are required."
  exit 1
fi

yaml-rss --server ${server} --show-file ${showFile} --items-pattern "${itemsPattern}" > ${output}
