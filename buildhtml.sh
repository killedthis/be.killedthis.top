#!/bin/bash
source /home/www/secret/db

readarray -t SERVICES < <(mysql -u$DBUSERNAME -p$DBPASSWORD -h dbserver -D funny -e "SELECT serviceprovider FROM killedthis" | rev | rev | uniq | tail -n +2)

for SERVICE in "${SERVICES[@]}"; do
    LINKSLIST="${LINKSLIST} <a href=\"https://${SERVICE}.killedthis.top\">${SERVICE}</a> | "

    readarray -t TITLES < <(mysql -u$DBUSERNAME -p$DBPASSWORD -h dbserver -D funny -e "SELECT title FROM killedthis WHERE serviceprovider='$SERVICE'" | rev | rev | uniq | tail -n +2)
    readarray -t DATEOFDEATH < <(mysql -u$DBUSERNAME -p$DBPASSWORD -h dbserver -D funny -e "SELECT date FROM killedthis WHERE serviceprovider='$SERVICE'" | rev | rev | tail -n +2)

    TITLELIST=""
    for INDEX in "${!TITLES[@]}"; do
        TITLELIST="${TITLELIST}  <div class=\"kt-item\"><img src=\"{movie_thumbnail}\" alt=\"placeholder image\"><div class=\"uk-text-secondary kt-item-title\">${TITLES[$INDEX]}</div><div class=\"kt-item-date\">${DATEOFDEATH[$INDEX]::-3}</div></div>"
        # printf "%s\n" "${INDEX} ${TITLES[$INDEX]} ${DATEOFDEATH[$INDEX]}"
    done

    cp /home/www/sites/killedthis.top/site.template /home/www/sites/killedthis.top/"${SERVICE,}".killedthis.top.html
    sed -i "s|#template#|$TITLELIST|g" "/home/www/sites/killedthis.top/${SERVICE,}.killedthis.top.html"
    sed -i "s|{website_title}|${SERVICE}|g" "/home/www/sites/killedthis.top/${SERVICE,}.killedthis.top.html"
    sed -i "s|{Title}|${SERVICE} Killed This|g" "/home/www/sites/killedthis.top/${SERVICE,}.killedthis.top.html"

done

printf "document.write('\\\nLinks: %s\\\n ');\n" "$LINKSLIST" >/home/www/sites/killedthis.top/js/loadlinks.js
