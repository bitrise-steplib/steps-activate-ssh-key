#!/bin/bash
set -e
#
# [!] Run this script from the root folder of this repository
#

THIS_SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "${THIS_SCRIPTDIR}"

export ssh_rsa_private_key=$(<testkey_wrong)
export ssh_key_save_path="$(pwd)/testsave"
export is_remove_other_identities="true"

set +e
#bash ../step.sh
go run ../main.go ../util.go
if [ $? -eq 0 ] ; then
	echo
	echo "-------------"
	echo "[!!!] Failed, should return with !=0 (error)!"
	echo "-------------"
	exit 1
fi

echo
echo "-------------"
echo "Test: OK"
echo "-------------"
