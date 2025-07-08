from .tool.func import *

async def main_func_easter_egg():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ip_or_user(ip) == 0:
            curs.execute(db_change('select name from user_set where id = ? and name = ?'), [ip, 'get_🥚'])
            if not curs.fetchall():
                curs.execute(db_change('insert into user_set (name, id, data) values ("get_🥚", ?, "Y")'), [ip])
    
        data = ''

        return easy_minify(flask.render_template(await skin_check(),
            imp = ['Easter Egg', await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = data,
            menu = 0
        ))