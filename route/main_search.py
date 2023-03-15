from .tool.func import *

def main_search():
    with get_db_connect() as conn:
        return redirect('/search/' + url_pas(flask.request.form.get('search', 'test')))
