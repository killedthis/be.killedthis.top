#!/bin/bash

# shellcheck source=/home/www/secret/r53
source /home/www/secret/r53
source /home/www/secret/db

hostnames=$(mysql -u$DBUSERNAME -p$DBPASSWORD -h dbserver -D funny -e "SELECT serviceprovider FROM killedthis" | rev | rev | uniq | tail -n +2)

printf "Hostnames Found %s \n" "${hostnames[@]}"

hostnamelist=""
for hostname in ${hostnames[@]}; do
	hostnamelist="${hostnamelist} ${hostname,}.killedthis.top"
done

#create catchall for all hostnames on :80
rm /home/www/hosts/dynamicservers.conf
printf "server {\n\tlisten 80;\n\tserver_name%s;\n\n\treturn 301 https://\$server_name\$request_uri;\n}\n\n" "$hostnamelist" >>/home/www/hosts/dynamicservers.conf

for hostname in ${hostnames[@]}; do
	/home/www//r53b/r53b "${hostname,}"
	printf "server {\n\tlisten 443 ssl http2;\n\tserver_name %s.killedthis.top;\n\tinclude /home/www/hosts/common/dynamicvhost.conf;\n\n\tlocation / {\n\t\troot /home/www/sites/killedthis.top/;\n\t\tindex \$server_name.html;\n\t}\n}\n\n" "${hostname,}" >>/home/www/hosts/dynamicservers.conf
done

nginx -t

if [ $? = 0 ]; then
	printf "We can reload"
	systemctl reload nginx
fi
