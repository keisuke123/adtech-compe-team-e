source /home/ozilikepop/miniconda3/bin/activate predictor
pkill -9 uwsgi
uwsgi --ini uwsgi.ini
