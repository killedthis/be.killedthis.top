#!/bin/bash

# shellcheck disable=SC1091
source /home/www/secret/r53
source /home/www/secret/db

readarray -t SERVICES < <(mysql -u"$DBUSERNAME" -p"$DBPASSWORD" -h dbserver -D funny -e "SELECT serviceprovider FROM killedthis" | rev | rev | tail -n +2 | sort -u | uniq)

printf "Service Providers Found %s \n" "${SERVICES[@]}"
echo -e

for SERVICE in "${SERVICES[@]}"; do
	hostnamelist="${hostnamelist} ${SERVICE,}.killedthis.top"
done

#create catchall for all hostnames on :80
printf "server {\n\tlisten 80;\n\tserver_name%s;\n\n\treturn 301 https://\$server_name\$request_uri;\n\n\tlocation ~ /\.(?!well-known).* {\n\t\tdeny all;\n\t}\n\n}\n\n" "$hostnamelist" >/home/www/hosts/dynamicservers.conf

for hostname in "${SERVICES[@]}"; do
	(
		cd /home/www/golang/r53b || exit
		/home/www/golang/r53b/r53b.go "${hostname,}"
	)


	printf "server {\n\tlisten 443 quic;\n\tlisten 443 ssl;\n\tlisten [::]:443 quic;\n\tlisten [::]:443 ssl;\n\thttp2 on;\n\n\t
	server_name %s.killedthis.top;\n\tinclude /home/www/hosts/common/dynamicvhost.conf;\n\n\tlocation / {\n\t\troot /home/www/sites/killedthis.top/;\n\t\tindex \$server_name.html;\n\t}\n\n\tlocation ~ /\.(?!well-known).* {\n\t\tdeny all;\n\t}\n\n}\n" "${hostname,}" >>/home/www/hosts/dynamicservers.conf
	echo -e
done

if nginx -t; then
	printf "We can reload"
	systemctl reload nginx
fi
