from .tool.func import *

from .go_api_w_random import api_w_random

async def view_random():
    with get_db_connect() as conn:
        data = (await api_w_random())["data"]
        
        return redirect(conn, '/w/' + url_pas(data))