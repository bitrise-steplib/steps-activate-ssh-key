#!/bin/bash

#
# [!] Run this script from the root folder of this repository
#

THIS_SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPTDIR}"

export ssh_rsa_private_key=$(<testkey_ok)
export ssh_key_save_path="$(PWD)/testsave"
export is_remove_other_identities="true"

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