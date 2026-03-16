from .tool.func import *

async def login_logout():
    with get_db_connect() as conn:
        return_url = flask.request.args.get('return', '')
        if not return_url.startswith('/') or return_url.startswith('//') or '\\' in return_url:
            return_url = ''

        flask.session.pop('state', None)
        flask.session.pop('id', None)

        if return_url != '':
            return redirect(conn, return_url)
        else:
            return redirect(conn, '/user')