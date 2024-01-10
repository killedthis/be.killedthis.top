source /home/www/secret/r53
lego -d killedthis.top -d *.killedthis.top --email="admin@ear.pm" --dns route53 run

lego -d earentir.dev -d *.earentir.dev --email="admin@ear.pm" --dns route53 run

echo "key in .lego/certificates"
