#! /bin/bash

if ! which oc > /dev/null; then
    echo "'oc' command not found"
    echo 'install and try again - https://docs.okd.io/latest/cli_reference/openshift_cli/getting-started-cli.html'
    exit
fi

# require oc client version 4.11 or greater
veroc () {
    requiredVer="4.11"
    ocVer=$(oc version --client | awk -F : {'print $2'} | tr -d "[:space:]")
    if [ -z "${ocVer}" ]; then
        return 0
    fi
    if [ "${ocVer}" == "${requiredVer}" ]; then
        return 1
    fi
    if [  "${ocVer}" == "`echo -e "${ocVer}\n${requiredVer}" | sort -V | head -n1`" ]; then
        return 0
    fi
    return 1
}

if veroc; then
    echo ""
    echo "This script requires 'oc' client version ${requiredVer} or greater."
    echo "Currently running ${ocVer}"
    echo ""
    echo 'upgrade and try again - https://docs.okd.io/latest/cli_reference/openshift_cli/getting-started-cli.html'
    echo ""
    exit
fi

if ! which jq > /dev/null; then
   echo "'jq' command not found"
   echo 'install the jq 1.5+ and try again - https://stedolan.github.io/jq/download/'
   exit
fi

# require jq version 1.5 or greater
verjq() {
    jqRequiredVer="jq-1.5"
    jqVer=$(jq --version)
    if [ -z "${jqVer}" ]; then
        return 0
    fi
    if [ "${jqVer}" == "${jqRequiredVer}" ]; then
        return 1
    fi
    if [  "${jqVer}" == "`echo -e "${jqVer}\n${jqRequiredVer}" | sort -V | head -n1`" ]; then
        return 0
    fi
    return 1
}

if verjq; then
    echo ""
    echo "This script requires 'jq' version ${jqRequiredVer} or greater."
    echo "Currently running ${jqVer}"
    echo ""
    echo 'upgrade and try again - https://stedolan.github.io/jq/download/'
    echo ""
    exit
fi

ocuser=$(oc whoami)
echo "Logged in as ${ocuser}"

if [ $(oc auth can-i get csv) != "yes" ]; then
    echo "user cannot get csv"
    echo "'oc login ...' with a user that has admin rights to get csv and try again"
    exit
fi

installns=$(oc get csv dbaas-operator.v0.2.0 --ignore-not-found -o template --template '{{index .metadata.annotations "olm.operatorNamespace"}}')

if [ $(oc auth can-i update deployment --subresource=scale -n ${installns}) != "yes" ]; then
    echo "user cannot scale deployments in namespace ${installns}"
    echo "'oc login ...' with a user that has admin rights to scale deployments in namespace ${installns} and try again"
    exit
fi

INS_CRDS="crdbdbaasinstance crunchybridgeinstance mongodbatlasinstance rdsinstance dbaasinstance"

for insCRD in $INS_CRDS; do
    if [ $(oc auth can-i get $insCRD --all-namespaces) == "yes" ]; then
        echo -n "."
    else
        echo "user cannot get ${insCRD}"
        echo "'oc login ...' with a user that has admin rights to get ${insCRD} and try again"
        exit
    fi
done

for insCRD in $INS_CRDS; do
    if [ $(oc auth can-i patch $insCRD --subresource=status --all-namespaces) == "yes" ]; then
        echo -n "."
    else
        echo "user cannot patch ${insCRD} status subresource"
        echo "'oc login ...' with a user that has admin rights to patch ${insCRD} status subresource and try again"
        exit
    fi
done

echo ""
echo "Stop RHODA"
oc scale --replicas=0 deployment dbaas-operator-controller-manager -n ${installns}
oc scale --replicas=0 deployment ccapi-k8s-operator-controller-manager -n ${installns}
oc scale --replicas=0 deployment crunchy-bridge-operator-controller-manager -n ${installns}
oc scale --replicas=0 deployment mongodb-atlas-operator -n ${installns}
oc scale --replicas=0 deployment rds-dbaas-operator-controller-manager -n ${installns}

echo ""
echo "Upgrading RHODA 0.2.0 dbaasinstance CRs"
SUPPORTED_PHASES="[\"Unknown\",\"Pending\",\"Creating\",\"Updating\",\"Deleting\",\"Deleted\",\"Ready\",\"Error\",\"Failed\"]"
for insCRD in $INS_CRDS; do
    dbaasinstances=$(oc get $insCRD --all-namespaces -o=json | jq --argjson PHASES "$SUPPORTED_PHASES" '.items[] | select(.status.phase as $p | $PHASES | all(.!=$p)) | "\(.metadata.name),\(.metadata.namespace)"')
    if [ ! -z "${dbaasinstances}" ]; then
        saveIFS="$IFS"
        for ins in $dbaasinstances; do
            nsn=${ins:1:-1}
            IFS=, read -r name namespace <<< $nsn
            oc patch $insCRD $name -n $namespace --subresource='status' --type='merge' -p '{"status":{"phase":"Unknown"}}'
        done
        IFS="$saveIFS"  # Set IFS back to normal!
    fi
done