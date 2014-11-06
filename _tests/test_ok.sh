#!/bin/bash

#
# [!] Run this script from the root folder of this repository
#

THIS_SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPTDIR}"

SSH_RSA_PUBLIC_KEY=$(<testkey_ok) SSH_KEY_SAVE_PATH="$(PWD)/testsave" bash ../step.sh
if [ $? -ne 0 ] ; then
	echo
	echo "-------------"
	echo "[!!!] Failed, should return with 0 (ok)!"
	echo "-------------"
	exit 1
fi

echo
echo "-------------"
echo "Test: OK"
echo "-------------"