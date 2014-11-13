#!/bin/bash

#
# [!] Run this script from the root folder of this repository
#

THIS_SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPTDIR}"

export SSH_RSA_PRIVATE_KEY=$(<testkey_ok)
export SSH_KEY_SAVE_PATH="$(PWD)/testsave"
export IS_REMOVE_OTHER_IDENTITIES="true"

bash ../step.sh
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
echo " Loaded identities:"
source ~/.bashrc
ssh-add -l
echo "-------------"