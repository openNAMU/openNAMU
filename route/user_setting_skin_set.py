from .tool.func import *

async def user_setting_skin_set():
    with get_db_connect() as conn:
        curs = conn.cursor()

        data = flask.make_response(await re_error(conn, 5))

        curs.execute(db_change("select data from other where name = 'language'"))
        main_data = curs.fetchall()

        data.set_cookie('language', main_data[0][0])

        curs.execute(db_change('select data from user_set where name = "lang" and id = ?'), [ip_check()])
        user_data = curs.fetchall()
        data.set_cookie('user_language', user_data[0][0] if user_data else main_data[0][0])

        return data