from .tool.func import *

async def main_view_file(data = ''):
    with get_db_connect() as conn:
        if data == 'robots.txt':
            curs = conn.cursor()

            curs.execute(db_change("select data from other where name = 'robot_default'"))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                return flask.Response(get_default_robots_txt(conn), mimetype = 'text/plain')
            else:
                curs.execute(db_change("select data from other where name = 'robot'"))
                db_data = curs.fetchall()
                if db_data:
                    return flask.Response(db_data[0][0], mimetype = 'text/plain')
                else:
                    return flask.Response(get_default_robots_txt(conn), mimetype = 'text/plain')
        elif os.path.exists(data):
            if re.search(r'\.txt$', data, flags = re.I):
                return flask.send_from_directory('./', data, mimetype = 'text/plain')
            else:
                return flask.send_from_directory('./', data, mimetype = 'text/xml')
        else:
            return ''