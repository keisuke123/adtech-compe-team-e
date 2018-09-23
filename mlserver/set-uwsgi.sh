source /home/mizutaninikkou/miniconda3/bin/activate predict
pkill -9 uwsgi
uwsgi --ini uwsgi.ini
