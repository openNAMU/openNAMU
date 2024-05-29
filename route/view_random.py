from .tool.func import *

from .go_api_w_random import api_w_random

def view_random():
    with get_db_connect() as conn:
        data = json.loads(api_w_random().data)["data"]
        
        return redirect(conn, '/w/' + url_pas(data))