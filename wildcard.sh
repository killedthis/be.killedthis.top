source /home/www/secret/r53
lego -d killedthis.top -d *.killedthis.top --email="admin@ear.pm" --dns route53 --path /etc/lego/certificates/ run

echo "key in /root/.lego/certificates"
