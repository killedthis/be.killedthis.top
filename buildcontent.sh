#!/bin/bash
# shellcheck disable=SC1091
source /home/www/secret/db

# DEFIMAGE=$(cat sites/killedthis.top/img/cbimage.jpg | base64)
readarray -t SERVICES < <(mysql -u"$DBUSERNAME" -p"$DBPASSWORD" -h dbserver -D funny -e "SELECT serviceprovider FROM killedthis" | rev | rev | tail -n +2 | sort -u | uniq)

for SERVICE in "${SERVICES[@]}"; do
    LINKSLIST="${LINKSLIST} <a href=\"https://${SERVICE}.killedthis.top\">${SERVICE}</a> | "
    MENUSTR="${MENUSTR} <li><a href=\"https://${SERVICE}.killedthis.top\">${SERVICE}</a></li>"

    readarray -t TITLES < <(mysql -u"$DBUSERNAME" -p"$DBPASSWORD" -h dbserver -D funny -e "SELECT title FROM killedthis WHERE serviceprovider='$SERVICE'" | rev | rev | uniq | tail -n +2)
    readarray -t YEARS < <(mysql -u"$DBUSERNAME" -p"$DBPASSWORD" -h dbserver -D funny -e "SELECT date FROM killedthis WHERE serviceprovider='$SERVICE'" | rev | rev | tail -n +2 | cut -c -4 | sort -s | uniq)
    readarray -t MONTHS < <(mysql -u"$DBUSERNAME" -p"$DBPASSWORD" -h dbserver -D funny -e "SELECT date FROM killedthis WHERE serviceprovider='$SERVICE'" | rev | rev | tail -n +2 | cut -c 6-7 | sort -s | uniq)
    readarray -t DATEOFDEATH < <(mysql -u"$DBUSERNAME" -p"$DBPASSWORD" -h dbserver -D funny -e "SELECT date FROM killedthis WHERE serviceprovider='$SERVICE'" | rev | rev | tail -n +2)
    readarray -t TMDBID < <(mysql -u"$DBUSERNAME" -p"$DBPASSWORD" -h dbserver -D funny -e "SELECT tmdbid FROM killedthis WHERE serviceprovider='$SERVICE'" | rev | rev | tail -n +2)

    YEARFILTER=""
    for INDEX in "${!YEARS[@]}"; do
        YEARFILTER="$YEARFILTER <li uk-filter-control=\"filter: [year='${YEARS[$INDEX]}']; group: year\"><a href=\"#\">${YEARS[$INDEX]}</a></li>"
    done

    MONTHFILTER=""
    for INDEX in "${!MONTHS[@]}"; do
        MONTHFILTER="$MONTHFILTER <li uk-filter-control=\"filter: [month='${MONTHS[$INDEX]}']; group: month\"><a href=\"#\">${MONTHS[$INDEX]}</a></li>"
    done

    TITLELIST=""
    for INDEX in "${!TITLES[@]}"; do
        if [ "${TMDBID[$INDEX]}" == "NULL" ]; then
            IMAGENAME="img/cbimage.jpg"
            /home/www/tmdb/tmdb -s ${TITLES[$INDEX]}
        else
            IMAGENAME="img/posters/${TMDBID[$INDEX]}.jpg"

            if [ -f "$IMAGENAME" ]; then
                /home/www/tmdb/tmdb -t ${TMDBID[$INDEX]}
            fi

        fi

        TITLELIST="${TITLELIST} <li class=\"kt-item\" year=\"${DATEOFDEATH[$INDEX]::-6}\" month=\"${DATEOFDEATH[$INDEX]:5:2}\" tmdbid=\"${TMDBID[$INDEX]}\"><img src=\"$IMAGENAME\" class=\"kt-thumbnail\" alt=\"No movie thumbnail\"><div class=\"uk-text-secondary kt-item-title\">${TITLES[$INDEX]}</div><div class=\"kt-item-date\">${DATEOFDEATH[$INDEX]::-3}</div></li>"
    done

    cp /home/www/sites/killedthis.top/site.template /home/www/sites/killedthis.top/"${SERVICE,}".killedthis.top.html

    #insert content into template
    sed -i "s|{template}|$TITLELIST|g" "/home/www/sites/killedthis.top/${SERVICE,}.killedthis.top.html"
    sed -i "s|{year_filter}|$YEARFILTER|g" "/home/www/sites/killedthis.top/${SERVICE,}.killedthis.top.html"
    sed -i "s|{month_filter}|$MONTHFILTER|g" "/home/www/sites/killedthis.top/${SERVICE,}.killedthis.top.html"

    #replace website title with service provider name
    sed -i "s|{website_title}|${SERVICE}|g" "/home/www/sites/killedthis.top/${SERVICE,}.killedthis.top.html"
    sed -i "s/{website_title}/KilledThis/g" "/home/www/sites/killedthis.top/root.html"
done

cp sites/killedthis.top/root.template sites/killedthis.top/root.html
sed -i "s|{menu}|$MENUSTR|g" "/home/www/sites/killedthis.top/root.html"
