from .tool.func import *

async def main_search():
    with get_db_connect() as conn:
        return redirect(conn, '/search/' + url_pas(flask.request.form.get('search', 'test')))