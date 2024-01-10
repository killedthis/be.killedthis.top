#!/bin/bash
# shellcheck disable=SC1091
source /home/www/secret/db

#setup vars for paths


# DEFIMAGE=$(cat sites/killedthis.top/img/cbimage.jpg | base64)
readarray -t SERVICES < <(mysql -u"${DBUSERNAME}" -p"${DBPASSWORD}" -h dbserver -D funny -e "SELECT serviceprovider FROM killedthis" | rev | rev | tail -n +2 | sort -u | uniq)

for SERVICE in "${SERVICES[@]}"; do
	LINKSLIST="${LINKSLIST} <a href=\"https://${SERVICE}.killedthis.top\">${SERVICE}</a> | "
	MENUSTR="${MENUSTR} <li><a href=\"https://${SERVICE}.killedthis.top\">${SERVICE}</a></li>"

	readarray -t TITLES < <(mysql -u"${DBUSERNAME}" -p"${DBPASSWORD}" -h dbserver -D funny -e "SELECT title FROM killedthis WHERE serviceprovider='${SERVICE}'" | rev | rev | uniq | tail -n +2)
	readarray -t YEARS < <(mysql -u"${DBUSERNAME}" -p"${DBPASSWORD}" -h dbserver -D funny -e "SELECT date FROM killedthis WHERE serviceprovider='${SERVICE}'" | rev | rev | tail -n +2 | cut -c -4 | sort -s | uniq)
	readarray -t MONTHS < <(mysql -u"${DBUSERNAME}" -p"${DBPASSWORD}" -h dbserver -D funny -e "SELECT date FROM killedthis WHERE serviceprovider='${SERVICE}'" | rev | rev | tail -n +2 | cut -c 6-7 | sort -s | uniq)
	readarray -t DATEOFDEATH < <(mysql -u"${DBUSERNAME}" -p"${DBPASSWORD}" -h dbserver -D funny -e "SELECT date FROM killedthis WHERE serviceprovider='${SERVICE}'" | rev | rev | tail -n +2)
	readarray -t TMDBID < <(mysql -u"${DBUSERNAME}" -p"${DBPASSWORD}" -h dbserver -D funny -e "SELECT tmdbid FROM killedthis WHERE serviceprovider='${SERVICE}'" | rev | rev | tail -n +2)

	YEARFILTER=""
	for INDEX in "${!YEARS[@]}"; do
		YEARFILTER="${YEARFILTER} <li uk-filter-control=\"filter: [year='${YEARS[${INDEX}]}']; group: year\"><a href=\"#\">${YEARS[${INDEX}]}</a></li>"
	done

	MONTHFILTER=""
	for INDEX in "${!MONTHS[@]}"; do
		MONTHFILTER="${MONTHFILTER} <li uk-filter-control=\"filter: [month='${MONTHS[${INDEX}]}']; group: month\"><a href=\"#\">${MONTHS[${INDEX}]}</a></li>"
	done

	TITLELIST=""
	for INDEX in "${!TITLES[@]}"; do
		if [[ ${TMDBID[${INDEX}]} == "NULL" ]]; then
			IMAGENAME="img/cbimage.jpg"
			(
				cd /home/www/golang/tmdb/ || exit
				/home/www/golang/tmdb/tmdb.go -s "${TITLES[${INDEX}]}"
			)
		else
			IMAGENAME="img/posters/${TMDBID[${INDEX}]}.webp"
			LARGEIMAGENAME="img/posters/${TMDBID[${INDEX}]}w780.webp"
			(
				cd /home/www/golang/tmdb/ || exit
				/home/www/golang/tmdb/tmdb.go -t="${TMDBID[${INDEX}]}"
			)
		fi

		TITLELIST="${TITLELIST} <li class=\"kt-item\" year=\"${DATEOFDEATH[${INDEX}]::-6}\" month=\"${DATEOFDEATH[${INDEX}]:5:2}\" tmdbid=\"${TMDBID[${INDEX}]}\"><a class=\"uk-link-reset\" href=\"#modal-full${INDEX}\" uk-toggle><img src=\"${IMAGENAME}\" loading=\"lazy\" class=\"kt-thumbnail\" alt=\"No movie thumbnail\"></a><div class=\"uk-text-secondary kt-item-title\">${TITLES[${INDEX}]}</div><div class=\"kt-item-date\">${DATEOFDEATH[${INDEX}]::-3}</div>    <div id=\"modal-full${INDEX}\" class=\"uk-modal-full\" uk-modal> <div class=\"uk-modal-dialog\"> <button class=\"uk-modal-close-full uk-close-large\" type=\"button\" uk-close></button> <div class=\"uk-grid-collapse uk-child-width-1-2@s uk-flex-middle\" uk-grid> <div class=\"uk-background-cover\" style=\"background-image: url('${LARGEIMAGENAME}');\" uk-height-viewport></div> <div class=\"uk-padding-large\"> <h1>${TITLES[${INDEX}]}</h1> <h2>${DATEOFDEATH[${INDEX}]::-6} - ${DATEOFDEATH[${INDEX}]:5:2}</h2> <p>foobar</p> </div> </div> </div> </div>  </li>"
	done

	cp /home/www/sites/killedthis.top/site.template /home/www/sites/killedthis.top/"${SERVICE,}".killedthis.top.html

	#insert content into template
	sed -i "s|{template}|${TITLELIST}|g" "/home/www/sites/killedthis.top/${SERVICE,}.killedthis.top.html"
	sed -i "s|{year_filter}|${YEARFILTER}|g" "/home/www/sites/killedthis.top/${SERVICE,}.killedthis.top.html"
	sed -i "s|{month_filter}|${MONTHFILTER}|g" "/home/www/sites/killedthis.top/${SERVICE,}.killedthis.top.html"

	#replace website title with service provider name
	sed -i "s|{website_title}|${SERVICE} killed the following shows|g" "/home/www/sites/killedthis.top/${SERVICE,}.killedthis.top.html"
done

cp sites/killedthis.top/root.template sites/killedthis.top/root.html
cp sites/killedthis.top/menu.template sites/killedthis.top/menu.js
sed -i "s|{menu}|${MENUSTR}|g" "/home/www/sites/killedthis.top/menu.js"
sed -i "s/{website_title}/KilledThis/g" "/home/www/sites/killedthis.top/root.html"
