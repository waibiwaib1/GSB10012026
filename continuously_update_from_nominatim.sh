#!/bin/bash -e

# For Nominatim < 3.7 set this to the Nominatim build directory.
# For newer versions, this must be the project directory of your import.
: ${NOMINATIM_DIR:=.}
: ${PHOTON_JAR:=photon.jar}
: ${PHOTON_DB_NAME:=nominatim}
: ${PHOTON_DB_USER:=nominatim}
: ${PHOTON_DB_PASSWORD:=}

while true
do
    starttime=`date +%s`

    # First consume and index updates in the Nominatim database. Photon
    # consumes the completed rows through triggers installed with
    # -nominatim-update-init-for.

    # For Nominatim versions < 3.7 use the following:
    # (cd $NOMINATIM_DIR && ./utils/update.php --import-osmosis)
    nominatim replication --project-dir $NOMINATIM_DIR --once

    # Now copy the indexed changes into Photon's database.
    java -jar $PHOTON_JAR -database $PHOTON_DB_NAME -user $PHOTON_DB_USER -password $PHOTON_DB_PASSWORD -nominatim-update

    # Sleep a bit if updates take less than a minute.
    # If you consume hourly or daily diffs adapt the period accordingly.
    endtime=`date +%s`
    elapsed=$((endtime - starttime))
    if [[ $elapsed -lt 60 ]]
    then
        sleepy=$((60 - $elapsed))
        echo "Sleeping for ${sleepy}s..."
        sleep $sleepy
    fi
done
