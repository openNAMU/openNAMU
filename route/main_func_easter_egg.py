from .tool.func import *

def main_func_easter_egg():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ip_or_user(ip) == 0:
            curs.execute(db_change('select name from user_set where id = ? and name = ?'), [ip, 'get_🥚'])
            if not curs.fetchall():
                curs.execute(db_change('insert into user_set (name, id, data) values ("get_🥚", ?, "Y")'), [ip])
                conn.commit()
    
        if platform.system() == 'Linux':
            if platform.machine() in ["AMD64", "x86_64"]:
                data = os.popen(os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".amd64.bin")).read()
            else:
                data = os.popen(os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".arm64.bin")).read()
        else:
            if platform.machine() in ["AMD64", "x86_64"]:
                data = os.popen(os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".amd64.exe")).read()
            else:
                data = os.popen(os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".arm64.exe")).read()

        return easy_minify(flask.render_template(skin_check(),
            imp = ['Easter Egg', wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = data,
            menu = 0
        ))